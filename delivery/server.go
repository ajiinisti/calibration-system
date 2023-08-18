package delivery

import (
	"fmt"
	"net/http"

	"calibration-system.com/config"
	"calibration-system.com/delivery/controller"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/manager"
	"calibration-system.com/model"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Server struct {
	ucManager    manager.UsecaseManager
	engine       *gin.Engine
	authRoute    gin.IRoutes
	host         string
	tokenService authenticator.AccessToken
	cfg          config.Config
	websocket.Upgrader
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://talent-connect-dev.netlify.app", "http://localhost:3000", "http://167.172.84.203.nip.io", "http://167.172.84.203"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", " Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	}))

	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections (you can implement a proper origin check for production)
			return true
		},
	}

	return &Server{
		ucManager:    uc,
		engine:       r,
		authRoute:    auth,
		host:         fmt.Sprintf("%s:%s", cfg.ApiHost, cfg.ApiPort),
		tokenService: tokenService,
		cfg:          *cfg,
		Upgrader:     upgrader,
	}
}
