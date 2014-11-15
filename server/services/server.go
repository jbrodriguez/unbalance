package services

import (
	"apertoire.net/unbalance/server/dto"
	"apertoire.net/unbalance/server/model"
	"apertoire.net/unbalance/server/static"
	"github.com/apertoire/mlog"
	"github.com/apertoire/pubsub"
	"github.com/gin-gonic/gin"
)

const apiVersion string = "/api/v1"

type Server struct {
	bus    *pubsub.PubSub
	engine *gin.Engine
}

func NewServer(bus *pubsub.PubSub) *Server {
	server := &Server{bus: bus}
	return server
}

func (self *Server) Start() {
	mlog.Info("starting service Server ...")

	self.engine = gin.New()

	self.engine.Use(gin.Recovery())
	// self.engine.Use(helper.Logging())

	self.engine.Use(static.Serve("./"))
	self.engine.NoRoute(static.Serve("./"))

	api := self.engine.Group(apiVersion)
	{
		api.GET("/config", self.getConfig)
		api.GET("/storage", self.getStorageInfo)
		api.POST("/storage/bestfit", self.calculateBestFit)
		api.POST("/storage/move", self.move)
	}

	mlog.Info("started listening on :6237")

	go self.engine.Run(":6237")
}

func (self *Server) Stop() {
	mlog.Info("stopped service Server ...")
}

func (self *Server) getConfig(c *gin.Context) {
	msg := &pubsub.Message(Reply: make(chan interface{}))
	self.bus.Pub(msg, "cmd.getConfig")

	reply := <= msg.Reply
	resp := reply.(*model.Config)
	c.JSON(200, &resp)
}

func (self *Server) getStorageInfo(c *gin.Context) {
	msg := &pubsub.Message{Reply: make(chan interface{})}
	self.bus.Pub(msg, "cmd.getStorageInfo")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)
	c.JSON(200, &resp)
}

func (self *Server) calculateBestFit(c *gin.Context) {
	var bestFit dto.BestFit

	c.Bind(&bestFit)

	msg := &pubsub.Message{Payload: &bestFit, Reply: make(chan interface{})}
	self.bus.Pub(msg, "cmd.calculateBestFit")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)
	c.JSON(200, &resp)
}

func (self *Server) move(c *gin.Context) {
	msg := &pubsub.Message{Reply: make(chan interface{})}
	self.bus.Pub(msg, "cmd.move")

	reply := <-msg.Reply
	resp := reply.([]*dto.Move)

	c.JSON(200, &resp)
}
