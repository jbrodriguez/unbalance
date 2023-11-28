package server

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"unbalance/common"
	"unbalance/domain"
	"unbalance/logger"
	"unbalance/services/core"
)

type Server struct {
	ctx    *domain.Context
	core   *core.Core
	engine *echo.Echo
}

func Create(ctx *domain.Context, core *core.Core) *Server {
	return &Server{
		ctx:  ctx,
		core: core,
	}
}

func (s *Server) Start() error {
	s.engine = echo.New()

	s.engine.HideBanner = true

	s.engine.Use(middleware.Recover())
	s.engine.Use(middleware.CORS())
	s.engine.Use(middleware.Gzip())
	// s.engine.Use(middleware.Logger())

	// Define a "/" endpoint to serve index.html from the embed FS
	// s.engine.GET("/*", indexHandler)

	// s.engine.GET("/assets/*", echo.WrapHandler(assetsHandler(web.Dist)))
	// s.engine.Static("/img/*", filepath.Join(s.ctx.DataDir, "img"))

	// s.engine.GET("/ws", s.wsHandler)

	api := s.engine.Group(common.APIEndpoint)
	api.GET("/config", s.getConfig)
	api.GET("/state", s.getState)
	api.GET("/storage", s.getStorage)
	api.GET("/operation", s.getOperation)
	api.GET("/history", s.getHistory)

	port := fmt.Sprintf(":%s", s.ctx.Port)
	go func() {
		err := s.engine.Start(port)
		if err != nil {
			logger.Red("unable to start http server: %s", err)
			os.Exit(1)
		}
	}()

	logger.Blue("started service server (listening http on %s) ...", port)

	return nil
}

func (s *Server) getConfig(c echo.Context) error {
	return c.JSON(200, s.core.GetConfig())
}

func (s *Server) getState(c echo.Context) error {
	return c.JSON(200, s.core.GetState())
}

func (s *Server) getStorage(c echo.Context) error {
	return c.JSON(200, s.core.GetStorage())
}

func (s *Server) getOperation(c echo.Context) error {
	return c.JSON(200, s.core.GetOperation())
}

func (s *Server) getHistory(c echo.Context) error {
	return c.JSON(200, s.core.GetHistory())
}
