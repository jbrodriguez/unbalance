package services

import (
	"apertoire.net/unbalance/server/dto"
	"apertoire.net/unbalance/server/model"
	"apertoire.net/unbalance/server/static"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

const apiVersion string = "/api/v1"

type Server struct {
	bus    *pubsub.PubSub
	config *model.Config
	engine *gin.Engine
}

func NewServer(bus *pubsub.PubSub, config *model.Config) *Server {
	server := &Server{bus: bus, config: config}
	return server
}

func (s *Server) Start() {
	mlog.Info("starting service Server ...")

	s.engine = gin.New()

	s.engine.Use(gin.Recovery())
	// s.engine.Use(helper.Logging())

	s.engine.Use(static.Serve("./"))
	s.engine.NoRoute(static.Serve("./"))

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
