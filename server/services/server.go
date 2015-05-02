package services

import (
	"apertoire.net/unbalance/server/dto"
	"apertoire.net/unbalance/server/model"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"os"
	"path/filepath"
)

const apiVersion string = "/api/v1"
const guiLocation string = "/usr/local/share/unbalance"

type Server struct {
	bus    *pubsub.PubSub
	config *model.Config
	engine *gin.Engine
	socket *Socket
}

func NewServer(bus *pubsub.PubSub, config *model.Config, socket *Socket) *Server {
	server := &Server{bus: bus, config: config, socket: socket}
	return server
}

func (s *Server) Start() {
	mlog.Info("Starting service Server ...")

	var path string
	if _, err := os.Stat("./index.html"); err == nil {
		path = "./"
	} else if _, err := os.Stat(filepath.Join(guiLocation, "index.html")); err == nil {
		path = guiLocation
	} else {
		slashdot, _ := filepath.Abs("./")
		mlog.Fatalf("Looked for web ui files in \n %s \n %s \n but didn\\'t find them", slashdot, guiLocation)
	}

	mlog.Info("Serving files from %s", path)

	s.engine = gin.New()

	s.engine.Use(gin.Recovery())
	// s.engine.Use(helper.Logging())
	s.engine.Use(static.Serve("/", static.LocalFile(path, true)))
	s.engine.NoRoute(static.Serve("/", static.LocalFile(path, true)))

	// websocket handler
	s.engine.GET("/ws", func(c *gin.Context) {
		s.socket.handler(c.Writer, c.Request)
	})

	api := s.engine.Group(apiVersion)
	{
		api.GET("/config", s.getConfig)
		api.PUT("/config", s.saveConfig)
		api.GET("/storage", s.getStorageInfo)
		api.POST("/storage/bestfit", s.calculateBestFit)
		api.POST("/storage/move", s.move)
	}

	mlog.Info("started listening on :6237")

	go s.engine.Run(":6237")
}

func (s *Server) Stop() {
	mlog.Info("stopped service Server ...")
}

func (s *Server) getConfig(c *gin.Context) {
	msg := &pubsub.Message{Reply: make(chan interface{})}
	s.bus.Pub(msg, "cmd.getConfig")

	reply := <-msg.Reply
	resp := reply.(*model.Config)
	c.JSON(200, &resp)
}

func (s *Server) saveConfig(c *gin.Context) {
	var config model.Config

	c.Bind(&config)

	msg := &pubsub.Message{Payload: &config, Reply: make(chan interface{})}
	s.bus.Pub(msg, "cmd.saveConfig")

	reply := <-msg.Reply
	resp := reply.(*model.Config)
	c.JSON(200, &resp)
}

func (s *Server) getStorageInfo(c *gin.Context) {
	msg := &pubsub.Message{Reply: make(chan interface{})}
	s.bus.Pub(msg, "cmd.getStorageInfo")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)
	c.JSON(200, &resp)
}

func (s *Server) calculateBestFit(c *gin.Context) {
	var bestFit dto.BestFit

	c.Bind(&bestFit)

	msg := &pubsub.Message{Payload: &bestFit, Reply: make(chan interface{})}
	s.bus.Pub(msg, "cmd.calculateBestFit")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)
	c.JSON(200, &resp)
}

func (s *Server) move(c *gin.Context) {
	msg := &pubsub.Message{Reply: make(chan interface{})}
	s.bus.Pub(msg, "cmd.move")

	reply := <-msg.Reply
	resp := reply.([]*dto.Move)

	c.JSON(200, &resp)
}
