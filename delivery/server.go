package delivery

import (
	"fmt"

	"calibration-system.com/config"
	"calibration-system.com/delivery/controller"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/manager"
	"calibration-system.com/model"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	ucManager    manager.UsecaseManager
	engine       *gin.Engine
	authRoute    gin.IRoutes
	host         string
	tokenService authenticator.AccessToken
	cfg          config.Config
}

func (s *Server) initController() {
	controller.NewRoleController(s.engine, s.ucManager.RoleUc())
	controller.NewUserController(s.engine, s.ucManager.UserUc())
	controller.NewAuthController(s.engine, s.ucManager.AuthUc(), s.tokenService, s.cfg)
}

func (s *Server) Run() {
	s.initController()

	err := s.engine.Run(s.host)
	if err != nil {
		panic(err)
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	infra, err := manager.NewInfraManager(cfg)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Address, cfg.RedisPort),
		DB:   cfg.Db,
		// Username: "username",
	})

	tokenService := authenticator.NewTokenService(*cfg, client)

	repo := manager.NewRepoManager(infra)
	uc := manager.NewUsecaseManager(repo, cfg)

	r := gin.Default()
	r.GET("/migration", func(ctx *gin.Context) {
		infra.Migrate(
			&model.User{},
			&model.Role{},
		)
	})

	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())

	return &Server{
		ucManager:    uc,
		engine:       r,
		authRoute:    auth,
		host:         fmt.Sprintf("%s:%s", cfg.ApiHost, cfg.ApiPort),
		tokenService: tokenService,
		cfg:          *cfg,
	}
}
