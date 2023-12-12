package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"unbalance/common"
	"unbalance/domain"
	"unbalance/logger"
	"unbalance/services/core"
	web "unbalance/ui"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	ctx           *domain.Context
	core          *core.Core
	engine        *echo.Echo
	ws            *websocket.Conn
	broadcastChan chan any
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
	s.engine.Use(middleware.Logger())

	// Define a "/" endpoint to serve index.html from the embed FS
	s.engine.GET("/*", indexHandler)

	s.engine.GET("/assets/*", echo.WrapHandler(assetsHandler(web.Dist)))
	// s.engine.Static("/img/*", filepath.Join(s.ctx.DataDir, "img"))

	s.engine.GET("/ws", s.wsHandler)

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

func indexHandler(c echo.Context) error {
	data, err := web.Dist.ReadFile("dist/index.html")
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "text/html", data)
}

func assetsHandler(content embed.FS) http.Handler {
	fsys, err := fs.Sub(content, "dist")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
}

func (s *Server) wsHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.Red("unable to upgrade websocket: %s", err)
		return err
	}
	defer conn.Close()

	s.ws = conn

	return s.wsRead()
}

func (s *Server) wsRead() (err error) {
	for {
		var packet domain.Packet
		err = s.ws.ReadJSON(&packet)
		if err != nil {
			logger.Red("unable to read websocket message: %s", err)
			return err
		}

		s.ctx.Hub.Pub(packet.Payload, packet.Topic)
	}
}

func (s *Server) wsWrite(packet *domain.Packet) (err error) {
	if (s.ws == nil) || (s.ws.RemoteAddr() == nil) {
		return
	}
	err = s.ws.WriteJSON(packet)
	return
}

func (s *Server) onBroadcast() {
	for msg := range s.broadcastChan {
		message := msg.(*domain.Packet)

		packet := &domain.Packet{
			Topic:   message.Topic,
			Payload: message.Payload,
		}

		err := s.wsWrite(packet)
		if err != nil {
			logger.Red("unable to write websocket message: %s", err)
		}
	}
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
