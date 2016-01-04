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
	"jbrodriguez/unbalance/server/net"
	// "os"
	"path/filepath"
)

const (
	API_VERSION = "/api/v1"
	CAPACITY    = 3
)

type Server struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings

	engine  *echo.Echo
	mailbox chan *pubsub.Mailbox

	pool map[*net.Connection]bool
}

func NewServer(bus *pubsub.PubSub, settings *lib.Settings) *Server {
	server := &Server{
		bus:      bus,
		settings: settings,
		pool:     make(map[*net.Connection]bool),
	}
	server.init()
	return server
}

func (s *Server) Start() {
	mlog.Info("Starting service Server ...")

	locations := []string{
		"/usr/local/emhttp/plugins/unbalance",
		"/usr/local/share/unbalance",
		".",
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
	s.engine.Static("/app", filepath.Join(location, "app"))
	s.engine.Static("/img", filepath.Join(location, "img"))

	s.engine.WebSocket("/skt", s.handleWs)

	api := s.engine.Group(API_VERSION)
	api.Put("/config/notifyCalc", s.setNotifyCalc)
	api.Put("/config/notifyMove", s.setNotifyMove)
	api.Put("/config/folder", s.addFolder)
	api.Delete("/config/folder", s.deleteFolder)
	api.Get("/config", s.getConfig)
	api.Get("/storage", s.getStorage)
	api.Post("/tree", s.getTree)
	api.Put("/config/dryRun", s.toggleDryRun)

	go s.engine.Run(":6237")

	s.mailbox = s.register(s.bus, "socket:broadcast", s.broadcast)
	go s.react()

	mlog.Info("Server started listening on :6237")
}

func (s *Server) Stop() {
	mlog.Info("stopped service Server ...")
}

func (s *Server) react() {
	for mbox := range s.mailbox {
		// mlog.Info("Core:Topic: %s", mbox.Topic)
		s.dispatch(mbox.Topic, mbox.Content)
	}
}

func (s *Server) getConfig(c *echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/get/config")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) setNotifyCalc(c *echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/config/set/notifyCalc")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) setNotifyMove(c *echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/config/set/notifyMove")

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

func (s *Server) deleteFolder(c *echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/config/delete/folder")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) getStorage(c *echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/get/storage")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) getTree(c *echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/get/tree")

	reply := <-msg.Reply
	resp := reply.(*dto.Entry)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) toggleDryRun(c *echo.Context) (err error) {
	msg := &pubsub.Message{Payload: nil, Reply: make(chan interface{}, CAPACITY)}
	s.bus.Pub(msg, "/config/toggle/dryRun")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) handleWs(c *echo.Context) (err error) {
	conn := net.NewConnection(c.Socket(), s.onMessage, s.onClose)
	s.pool[conn] = true
	err = conn.Read()
	return err
}

func (s *Server) onMessage(packet *dto.Packet) {
	s.bus.Pub(&pubsub.Message{Payload: packet.Payload}, packet.Topic)
}

func (s *Server) onClose(c *net.Connection, err error) {
	mlog.Warning("closing socket (%+v): %s", c, err)
	if _, ok := s.pool[c]; ok {
		delete(s.pool, c)
	}
}

func (s *Server) broadcast(msg *pubsub.Message) {
	packet := msg.Payload.(*dto.Packet)
	for conn, _ := range s.pool {
		conn.Write(packet)
	}
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
