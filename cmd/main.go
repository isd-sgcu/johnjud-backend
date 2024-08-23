package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/isd-sgcu/johnjud-backend/client/bucket"
	"github.com/isd-sgcu/johnjud-backend/config"
	"github.com/isd-sgcu/johnjud-backend/constant"
	"github.com/isd-sgcu/johnjud-backend/database"
	"github.com/isd-sgcu/johnjud-backend/internal/auth"
	"github.com/isd-sgcu/johnjud-backend/internal/auth/email"
	"github.com/isd-sgcu/johnjud-backend/internal/auth/jwt"
	"github.com/isd-sgcu/johnjud-backend/internal/auth/token"
	"github.com/isd-sgcu/johnjud-backend/internal/cache"
	"github.com/isd-sgcu/johnjud-backend/internal/healthcheck"
	"github.com/isd-sgcu/johnjud-backend/internal/image"
	guard "github.com/isd-sgcu/johnjud-backend/internal/middleware/auth"
	"github.com/isd-sgcu/johnjud-backend/internal/pet"
	"github.com/isd-sgcu/johnjud-backend/internal/router"
	"github.com/isd-sgcu/johnjud-backend/internal/user"
	"github.com/isd-sgcu/johnjud-backend/internal/utils"
	"github.com/isd-sgcu/johnjud-backend/internal/validator"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

// @title JohnJud API
// @version 1.0
// @description.markdown

// @contact.name ISD Team
// @contact.email sd.team.sgcu@gmail.com

// @schemes https http

// @securityDefinitions.apikey  AuthToken
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used

// @tag.name auth
// @tag.description.markdown

// @tag.name image
// @tag.description.markdown

// @tag.name pet
// @tag.description.markdown

// @tag.name user
// @tag.description.markdown

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "config").
			Msg("Failed to start service")
	}

	db, err := database.InitPostgresDatabase(&conf.Database, conf.App.IsDevelopment())
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth").
			Msg("Failed to init postgres connection")
	}

	cacheDb, err := database.InitRedisConnection(&conf.Redis)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth").
			Msg("Failed to init redis connection")
	}

	v, err := validator.NewIValidator()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "validator").
			Msg("Failed to start service")
	}

	hc := healthcheck.NewHandler()

	uuidUtil := utils.NewUuidUtil()
	bcryptUtil := utils.NewBcryptUtil()

	accessTokenCache := cache.NewRepository(cacheDb)
	refreshTokenCache := cache.NewRepository(cacheDb)
	resetPasswordCache := cache.NewRepository(cacheDb)

	bcryptUtils := utils.NewBcryptUtil()
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, bcryptUtils)
	userHandler := user.NewHandler(userSvc, v)

	jwtStrat := jwt.NewJwtStrategy(conf.Jwt.Secret)
	jwtUtils := jwt.NewJwtUtil()
	jwtSvc := jwt.NewService(conf.Jwt, jwtStrat, jwtUtils)
	tokenSvc := token.NewService(jwtSvc, accessTokenCache, refreshTokenCache, resetPasswordCache, uuidUtil)
	emailSvc := email.NewService(conf.Sendgrid)
	authRepo := auth.NewRepository(db)
	authSvc := auth.NewService(authRepo, userRepo, tokenSvc, emailSvc, bcryptUtil, conf.Auth)
	authHandler := auth.NewHandler(authSvc, userSvc, v)

	authGuard := guard.NewAuthGuard(authSvc, constant.ExcludePath, constant.AdminPath, conf.App, constant.VersionList)

	minioClient, err := minio.New(conf.Bucket.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Bucket.AccessKeyID, conf.Bucket.SecretAccessKey, ""),
		Secure: conf.Bucket.UseSSL,
	})
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "file").
			Msg("Failed to start Minio client")
		return
	}
	imageClient := bucket.NewClient(conf.Bucket, minioClient)
	randomUtils := utils.NewRandomUtil()
	imageRepo := image.NewRepository(db)
	imageService := image.NewService(imageClient, imageRepo, randomUtils)
	imageHandler := image.NewHandler(imageService, v, conf.App.MaxFileSize)

	petRepo := pet.NewRepository(db)
	petService := pet.NewService(petRepo, imageService)
	petHandler := pet.NewHandler(petService, imageService, v)

	r := router.NewFiberRouter(&authGuard, conf.App)

	r.GetUser("/:id", userHandler.FindOne)
	r.PutUser("", userHandler.Update)
	r.DeleteUser("/:id", userHandler.Delete)

	r.PostAuth("/signup", authHandler.Signup)
	r.PostAuth("/signin", authHandler.SignIn)
	r.PostAuth("/signout", authHandler.SignOut)
	//r.PostAuth("/me", authHandler.Validate)
	r.PostAuth("/refreshToken", authHandler.RefreshToken)
	r.PostAuth("/forgot-password", authHandler.ForgotPassword)
	r.PutAuth("/admin/reset-password", authHandler.ResetPassword)

	r.GetHealthCheck("", hc.HealthCheck)

	r.GetPet("", petHandler.FindAll)
	r.GetPet("/admin", petHandler.FindAllAdmin)
	r.GetPet("/:id", petHandler.FindOne)
	r.PostPet("", petHandler.Create)
	r.PutPet("/:id", petHandler.Update)
	r.PutPet("/:id/adopt", petHandler.Adopt)
	r.PutPet("/:id/visible", petHandler.ChangeView)
	r.DeletePet("/:id", petHandler.Delete)

	r.PostImage("", imageHandler.Upload)
	r.DeleteImage("/:id", imageHandler.Delete)

	v1 := router.NewAPIv1(r, conf.App)

	go func() {
		if err := v1.Listen(fmt.Sprintf(":%v", conf.App.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatal().
				Err(err).
				Str("service", "mgl-gateway").
				Msg("Server not close properly")
		}
	}()

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"server": func(ctx context.Context) error {
			return r.Shutdown()
		},
	})

	<-wait
}

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		sig := <-s

		log.Info().
			Str("service", "graceful shutdown").
			Msgf("got signal \"%v\" shutting down service", sig)

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Error().
				Str("service", "graceful shutdown").
				Msgf("timeout %v ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("cleaning up: %v", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Error().
						Str("service", "graceful shutdown").
						Err(err).
						Msgf("%v: clean up failed: %v", innerKey, err.Error())
					return
				}

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("%v was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
