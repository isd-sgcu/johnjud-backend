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

	"github.com/isd-sgcu/johnjud-gateway/config"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/auth"
	"github.com/isd-sgcu/johnjud-gateway/internal/healthcheck"
	"github.com/isd-sgcu/johnjud-gateway/internal/image"
	guard "github.com/isd-sgcu/johnjud-gateway/internal/middleware/auth"
	"github.com/isd-sgcu/johnjud-gateway/internal/pet"
	"github.com/isd-sgcu/johnjud-gateway/internal/router"
	"github.com/isd-sgcu/johnjud-gateway/internal/user"
	"github.com/isd-sgcu/johnjud-gateway/internal/validator"
	authProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	userProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
	petProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	v, err := validator.NewIValidator()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "validator").
			Msg("Failed to start service")
	}

	backendConn, err := grpc.Dial(conf.Service.Backend, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "johnjud-backend").
			Msg("Cannot connect to service")
	}

	authConn, err := grpc.Dial(conf.Service.Auth, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "johnjud-auth").
			Msg("Cannot connect to service")
	}

	fileConn, err := grpc.Dial(conf.Service.File, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "johnjud-file").
			Msg("Cannot connect to service")
	}

	hc := healthcheck.NewHandler()

	userClient := userProto.NewUserServiceClient(authConn)
	userService := user.NewService(userClient)
	userHandler := user.NewHandler(userService, v)

	authClient := authProto.NewAuthServiceClient(authConn)
	authService := auth.NewService(authClient)
	authHandler := auth.NewHandler(authService, userService, v)

	authGuard := guard.NewAuthGuard(authService, constant.ExcludePath, constant.AdminPath, conf.App, constant.VersionList)

	imageClient := imageProto.NewImageServiceClient(fileConn)
	imageService := image.NewService(imageClient)
	imageHandler := image.NewHandler(imageService, v, conf.App.MaxFileSize)

	petClient := petProto.NewPetServiceClient(backendConn)
	petService := pet.NewService(petClient, imageService)
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
