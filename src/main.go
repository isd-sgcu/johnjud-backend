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

	authHdr "github.com/isd-sgcu/johnjud-gateway/src/app/handler/auth"
	healthcheck "github.com/isd-sgcu/johnjud-gateway/src/app/handler/healthcheck"
	imageHdr "github.com/isd-sgcu/johnjud-gateway/src/app/handler/image"
	likeHdr "github.com/isd-sgcu/johnjud-gateway/src/app/handler/like"
	petHdr "github.com/isd-sgcu/johnjud-gateway/src/app/handler/pet"
	userHdr "github.com/isd-sgcu/johnjud-gateway/src/app/handler/user"
	guard "github.com/isd-sgcu/johnjud-gateway/src/app/middleware/auth"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	authSvc "github.com/isd-sgcu/johnjud-gateway/src/app/service/auth"
	imageSvc "github.com/isd-sgcu/johnjud-gateway/src/app/service/image"
	likeSvc "github.com/isd-sgcu/johnjud-gateway/src/app/service/like"
	petSvc "github.com/isd-sgcu/johnjud-gateway/src/app/service/pet"
	userSvc "github.com/isd-sgcu/johnjud-gateway/src/app/service/user"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	"github.com/isd-sgcu/johnjud-gateway/src/config"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/auth"
	authProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	userProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
	likeProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
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
	userService := userSvc.NewService(userClient)
	userHandler := userHdr.NewHandler(userService, v)

	authClient := authProto.NewAuthServiceClient(authConn)
	authService := authSvc.NewService(authClient)
	authHandler := authHdr.NewHandler(authService, userService, v)

	authGuard := guard.NewAuthGuard(authService, auth.ExcludePath, conf.App, auth.VersionList)

	imageClient := imageProto.NewImageServiceClient(fileConn)
	imageService := imageSvc.NewService(imageClient)
	imageHandler := imageHdr.NewHandler(imageService, v)

	likeClient := likeProto.NewLikeServiceClient(backendConn)
	likeService := likeSvc.NewService(likeClient)
	likeHandler := likeHdr.NewHandler(likeService, v)

	petClient := petProto.NewPetServiceClient(backendConn)
	petService := petSvc.NewService(petClient)
	petHandler := petHdr.NewHandler(petService, imageService, likeService, v)

	r := router.NewFiberRouter(&authGuard, conf.App)

	r.GetUser("/:id", userHandler.FindOne)
	r.PutUser("/", userHandler.Update)

	r.PostAuth("/signup", authHandler.Signup)
	r.PostAuth("/signin", authHandler.SignIn)
	r.PostAuth("/signout", authHandler.SignOut)
	//r.PostAuth("/me", authHandler.Validate)
	r.PostAuth("/refreshToken", authHandler.RefreshToken)

	r.GetHealthCheck("/", hc.HealthCheck)

	r.GetPet("/", petHandler.FindAll)
	r.GetPet("/:id", petHandler.FindOne)
	r.PostPet("/create", petHandler.Create)
	r.PutPet("/:id", petHandler.Update)
	r.PutPet("/:id/adopt", petHandler.Adopt)
	r.PutPet("/:id/visible", petHandler.ChangeView)
	r.DeletePet("/:id", petHandler.Delete)

	r.GetLike("/:id", likeHandler.FindByUserId)
	r.PostLike("/", likeHandler.Create)
	r.DeleteLike("/:id", likeHandler.Delete)

	r.PostImage("/create", imageHandler.Upload)
	r.GetImages("/:id", imageHandler.FindByPetId)
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
