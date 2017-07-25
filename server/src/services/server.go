package services

import (
	"fmt"
	// "os"
	"path/filepath"

	"jbrodriguez/unbalance/server/src/dto"
	"jbrodriguez/unbalance/server/src/lib"
	"jbrodriguez/unbalance/server/src/model"
	"jbrodriguez/unbalance/server/src/net"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"

	"golang.org/x/net/websocket"
)

const (
	apiVersion = "/api/v1"
	capacity   = 3
)

// Server -
type Server struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	engine *echo.Echo
	actor  *actor.Actor

	pool map[*net.Connection]bool
}

// NewServer -
func NewServer(bus *pubsub.PubSub, settings *lib.Settings) *Server {
	server := &Server{
		bus:      bus,
		actor:    actor.NewActor(bus),
		settings: settings,
		pool:     make(map[*net.Connection]bool),
	}
	return server
}

// Start -
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

	s.engine.Static("/", filepath.Join(location, "index.html"))
	s.engine.Static("/app", filepath.Join(location, "app"))
	s.engine.Static("/img", filepath.Join(location, "img"))

	s.engine.GET("/skt", echo.WrapHandler(websocket.Handler(s.handleWs)))

	api := s.engine.Group(apiVersion)
	api.PUT("/config/notifyCalc", s.setNotifyCalc)
	api.PUT("/config/notifyMove", s.setNotifyMove)
	api.PUT("/config/reservedSpace", s.setReservedSpace)
	api.PUT("/config/verbosity", s.setVerbosity)
	api.GET("/config", s.getConfig)
	api.GET("/storage", s.getStorage)
	api.POST("/tree", s.getTree)
	api.POST("/locate", s.locate)
	api.PUT("/config/dryRun", s.toggleDryRun)
	api.PUT("/config/rsyncFlags", s.setRsyncFlags)

	port := fmt.Sprintf(":%s", s.settings.Port)

	go s.engine.Start(port)

	s.actor.Register("socket:broadcast", s.broadcast)
	go s.actor.React()

	mlog.Info("Server started listening on %s", port)
}

// Stop -
func (s *Server) Stop() {
	mlog.Info("stopped service Server ...")
}

func (s *Server) getConfig(c echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/config")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) setNotifyCalc(c echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/config/set/notifyCalc")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) setNotifyMove(c echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/config/set/notifyMove")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) setReservedSpace(c echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/config/set/reservedSpace")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) setVerbosity(c echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/config/set/verbosity")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) getStorage(c echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/storage")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) getTree(c echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/tree")

	reply := <-msg.Reply
	resp := reply.(*dto.Entry)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) locate(c echo.Context) (err error) {
	var packet dto.Chosen

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding packet: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/disks/locate")

	reply := <-msg.Reply
	resp := reply.([]*model.Disk)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) toggleDryRun(c echo.Context) (err error) {
	msg := &pubsub.Message{Payload: nil, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/config/toggle/dryRun")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) setRsyncFlags(c echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/config/set/rsyncFlags")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) handleWs(ws *websocket.Conn) {
	conn := net.NewConnection(ws, s.onMessage, s.onClose)
	s.pool[conn] = true
	conn.Read()
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
	for conn := range s.pool {
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
