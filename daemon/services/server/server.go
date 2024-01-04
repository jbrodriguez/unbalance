package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"unbalance/daemon/common"
	"unbalance/daemon/domain"
	"unbalance/daemon/logger"
	"unbalance/daemon/services/core"
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
		ctx:           ctx,
		core:          core,
		broadcastChan: ctx.Hub.Sub("socket:broadcast"),
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

	api.GET("/tree/:route", s.getTree)
	api.GET("/locate/:route", s.locate)
	api.PUT("/config/dryRun", s.toggleDryRun)
	api.PUT("/config/notifyPlan", s.setNotifyPlan)
	api.PUT("/config/notifyTransfer", s.setNotifyTransfer)
	api.PUT("/config/reservedSpace", s.setReservedSpace)

	port := fmt.Sprintf(":%s", s.ctx.Port)
	go func() {
		err := s.engine.Start(port)
		if err != nil {
			logger.Red("unable to start http server: %s", err)
			os.Exit(1)
		}
	}()

	go s.onBroadcast()

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

		logger.Green("packet %+v", packet)

		s.ctx.Hub.Pub(packet, packet.Topic)
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

		// packet := &domain.Packet{
		// 	Topic:   message.Topic,
		// 	Payload: message.Payload,
		// }

		err := s.wsWrite(message)
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

type QueryPath struct {
	Path string `json:"param:path"`
	ID   string `json:"query:id"`
}

func (s *Server) getTree(c echo.Context) error {
	param := c.Param("route")
	u, err := url.Parse(param)
	if err != nil {
		return err
	}

	path := filepath.Join("/", "mnt", path.Clean(u.Path))
	id := c.QueryParam("id")

	return c.JSON(200, s.core.GetTree(path, id))
}

func (s *Server) locate(c echo.Context) error {
	param := c.Param("route")
	u, err := url.Parse(param)
	if err != nil {
		return err
	}

	path := filepath.Join("/", "mnt", "user", path.Clean(u.Path))

	return c.JSON(200, s.core.Locate(path))
}

func (s *Server) toggleDryRun(c echo.Context) error {
	return c.JSON(200, s.core.ToggleDryRun())
}

func (s *Server) setNotifyPlan(c echo.Context) error {
	var value int
	err := c.Bind(&value)
	if err != nil {
		return err
	}

	return c.JSON(200, s.core.SetNotifyPlan(value))
}

func (s *Server) setNotifyTransfer(c echo.Context) error {
	var value int
	err := c.Bind(&value)
	if err != nil {
		return err
	}

	return c.JSON(200, s.core.SetNotifyTransfer(value))
}

func (s *Server) setReservedSpace(c echo.Context) error {
	var params struct {
		Amount uint64 `json:"amount"`
		Unit   string `json:"unit"`
	}
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	return c.JSON(200, s.core.SetReservedSpace(params.Amount, params.Unit))
}
