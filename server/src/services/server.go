package services

import (
	"fmt"
	// "os"
	"path/filepath"

	"jbrodriguez/unbalance/server/src/common"
	"jbrodriguez/unbalance/server/src/domain"
	"jbrodriguez/unbalance/server/src/dto"
	"jbrodriguez/unbalance/server/src/lib"
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

	s.engine.HideBanner = true

	s.engine.Use(mw.Logger())
	s.engine.Use(mw.Recover())

	s.engine.Static("/", filepath.Join(location, "index.html"))
	s.engine.Static("/app", filepath.Join(location, "app"))
	s.engine.Static("/img", filepath.Join(location, "img"))

	s.engine.GET("/skt", echo.WrapHandler(websocket.Handler(s.handleWs)))

	api := s.engine.Group(apiVersion)

	api.GET("/config", s.getConfig)
	api.GET("/state", s.getState)
	api.GET("/storage", s.getStorage)
	api.GET("/operation", s.getOperation)
	api.GET("/history", s.getHistory)

	api.PUT("/config/notifyCalc", s.setNotifyCalc)
	api.PUT("/config/notifyMove", s.setNotifyMove)
	api.PUT("/config/reservedSpace", s.setReservedSpace)
	api.PUT("/config/verbosity", s.setVerbosity)
	api.PUT("/config/checkUpdate", s.setCheckUpdate)
	api.GET("/update", s.getUpdate)
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
	s.bus.Pub(msg, common.API_GET_CONFIG)

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

// func (s *Server) getStatus(c echo.Context) (err error) {
// 	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
// 	s.bus.Pub(msg, common.API_GET_STATUS)

// 	reply := <-msg.Reply
// 	status := reply.(int64)
// 	c.JSON(200, status)

// 	return nil
// }

func (s *Server) getState(c echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, common.API_GET_STATE)

	reply := <-msg.Reply
	state := reply.(*domain.State)
	c.JSON(200, state)

	return nil
}

func (s *Server) getStorage(c echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, common.API_GET_STORAGE)

	reply := <-msg.Reply
	storage := reply.(*domain.Unraid)
	c.JSON(200, storage)

	return nil
}

func (s *Server) getOperation(c echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, common.API_GET_OPERATION)

	reply := <-msg.Reply
	operation := reply.(*domain.Operation)
	c.JSON(200, operation)

	return nil
}

func (s *Server) getHistory(c echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, common.API_GET_HISTORY)

	reply := <-msg.Reply
	history := reply.(*domain.History)
	c.JSON(200, history)

	return nil
}

func (s *Server) setNotifyCalc(c echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, common.API_NOTIFY_CALC)

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
	s.bus.Pub(msg, common.API_NOTIFY_MOVE)

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
	s.bus.Pub(msg, common.API_SET_RESERVED)

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
	s.bus.Pub(msg, common.API_SET_VERBOSITY)

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) setCheckUpdate(c echo.Context) (err error) {
	var packet dto.Packet

	err = c.Bind(&packet)
	if err != nil {
		mlog.Warning("error binding: %s", err)
	}

	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, common.API_SET_CHECKUPDATE)

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) getUpdate(c echo.Context) (err error) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, common.API_GET_UPDATE)

	reply := <-msg.Reply
	resp := reply.(string)
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
	s.bus.Pub(msg, common.API_GET_TREE)

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
	s.bus.Pub(msg, common.API_LOCATE_FOLDER)

	reply := <-msg.Reply
	resp := reply.(*dto.Location)
	c.JSON(200, &resp)

	return nil
}

func (s *Server) toggleDryRun(c echo.Context) (err error) {
	msg := &pubsub.Message{Payload: nil, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, common.API_TOGGLE_DRYRUN)

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
