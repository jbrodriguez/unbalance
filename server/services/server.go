package services

import (
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"jbrodriguez/unbalance/server/dto"
	"jbrodriguez/unbalance/server/lib"
	"jbrodriguez/unbalance/server/model"
	// "os"
	"path/filepath"
)

const (
	API_VERSION = "/api/v1"
	CAPACITY    = 3
)

// const guiLocation string = "/usr/local/share/unbalance"

type Server struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	engine   *echo.Echo
	// socket   *Socket
}

func NewServer(bus *pubsub.PubSub, settings *lib.Settings) *Server {
	server := &Server{
		bus:      bus,
		settings: settings,
	}
	return server
}

func (s *Server) Start() {
	mlog.Info("Starting service Server ...")

	locations := []string{
		".",
		"/usr/local/share/unbalance",
	}

	location := lib.SearchFile("index.html", locations)
	if location == "" {
		msg := ""
		for _, loc := range locations {
			msg += fmt.Sprintf("%s, ", loc)
		}
		mlog.Fatalf("Unable to find index.html. Exiting now. (searched in %s)", msg)
	}

	mlog.Info("Serving files from %s", location)

	s.engine = echo.New()

	s.engine.Use(mw.Logger())
	s.engine.Use(mw.Recover())

	s.engine.Index(filepath.Join(location, "index.html"))
	s.engine.Static("/img", filepath.Join(location, "img"))
	s.engine.WebSocket("/skt", func(c *echo.Context) (err error) {
		ws := c.Socket()
		s.bus.Pub(&pubsub.Message{Payload: ws}, "/add/connection")
		return nil
	})

	api := s.engine.Group(API_VERSION)
	api.Put("/config/folder", s.addFolder)
	api.Get("/config", s.getConfig)
	api.Get("/storage", s.getStorage)
	api.Post("/calculate", s.calculate)
	api.Post("/move", s.move)

	// s.engine = gin.New()
	// s.engine.RedirectTrailingSlash = false
	// s.engine.RedirectFixedPath = false

	// s.engine.Use(gin.Recovery())
	// // s.engine.Use(helper.Logging())
	// s.engine.Use(static.Serve("/", static.LocalFile(path, true)))

	// websocket handler
	// s.engine.GET("/ws", func(c *gin.Context) {
	// 	s.socket.handler(c.Writer, c.Request)
	// })

	// api := s.engine.Group(apiVersion)
	// {
	// 	api.GET("/config", s.getConfig)
	// 	api.PUT("/config", s.saveConfig)
	// 	api.GET("/storage", s.getStorageInfo)
	// 	api.POST("/storage/bestfit", s.calculateBestFit)
	// }

	// // s.engine.NoRoute(static.Serve("/", static.LocalFile(path, true)))
	// s.engine.NoRoute(s.noRoute)

	go s.engine.Run(":6237")

	mlog.Info("Server started listening on :6237")
}

func (s *Server) Stop() {
	mlog.Info("stopped service Server ...")
}

func (s *Server) getConfig(c *echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/get/config")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) addFolder(c *echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/config/add/folder")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

// func (s *Server) saveConfig(c *echo.Context) (err error) {
// 	var config lib.Config

// 	c.Bind(&config)

// 	msg := &pubsub.Message{Payload: &config, Reply: make(chan interface{}, CAPACITY)}
// 	s.bus.Pub(msg, "/set/config")

// 	reply := <-msg.Reply
// 	resp := reply.(*lib.Config)
// 	c.JSON(200, &resp)

// 	return nil
// }

func (s *Server) getStorage(c *echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/get/storage")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) calculate(c *echo.Context) (err error) {
	var calculate dto.Calculate

	c.Bind(&calculate)
	// mlog.Warning("Unable to bind calculate: %s", err)

	msg := &pubsub.Message{Payload: &calculate, Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/calculate")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) move(c *echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{})}
	s.bus.Pub(msg, "/move")

	reply := <-msg.Reply
	resp := reply.([]*dto.Move)

	c.JSON(200, &resp)

	return nil
}

// func (s *Server) noRoute(c *gin.Context) {
// 	var path string
// 	if _, err := os.Stat("./index.html"); err == nil {
// 		path = "./"
// 	} else if _, err := os.Stat(filepath.Join(guiLocation, "index.html")); err == nil {
// 		path = guiLocation
// 	} else {
// 		slashdot, _ := filepath.Abs("./")
// 		mlog.Fatalf("Looked for web ui files in \n %s \n %s \n but didn\\'t find them", slashdot, guiLocation)
// 	}

// 	c.File(filepath.Join(path, "index.html"))
// }
