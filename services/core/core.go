package core

import (
	"regexp"
	"time"

	"unbalance/common"
	"unbalance/domain"
	"unbalance/lib"
	"unbalance/logger"

	"github.com/teris-io/shortid"
)

const (
	certDir    = "/boot/config/ssl/certs"
	mailCmd    = "/usr/local/emhttp/webGui/scripts/notify"
	timeFormat = "Jan _2, 2006 15:04:05"
)

var (
	reFreeSpace = regexp.MustCompile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	reRsync     = regexp.MustCompile(`exit status (\d+)`)
	reProgress  = regexp.MustCompile(`(?s)^([\d,]+).*?\(.*?\)$|^([\d,]+).*?$`)
	reItems     = regexp.MustCompile(`(\d+)\s+(.*?)$`)
	reStat      = regexp.MustCompile(`[-dclpsbD]([-rwxsS]{3})([-rwxsS]{3})([-rwxtT]{3})\|(.*?)\:(.*?)\|(.*?)\|(.*)`)
)

var rsyncErrors = map[int]string{
	0:  "Success",
	1:  "Syntax or usage error",
	2:  "Protocol incompatibility",
	3:  "Errors selecting input/output files, dirs",
	4:  "Requested action not supported: an attempt was made to manipulate 64-bit files on a platform that cannot support them, or an option was specified that is supported by the client and not by the server.",
	5:  "Error starting client-server protocol",
	6:  "Daemon unable to append to log-file",
	10: "Error in socket I/O",
	11: "Error in file I/O",
	12: "Error in rsync protocol data stream",
	13: "Errors with program diagnostics",
	14: "Error in IPC code",
	20: "Received SIGUSR1 or SIGINT",
	21: "Some error returned by waitpid()",
	22: "Error allocating core memory buffers",
	23: "Partial transfer due to error",
	24: "Partial transfer due to vanished source files",
	25: "The --max-delete limit stopped deletions",
	30: "Timeout in data send/receive",
	35: "Timeout waiting for daemon connection",
	99: "Interrupted by the user",
}

type Core struct {
	ctx *domain.Context

	state *domain.State
	sid   *shortid.Shortid

	mailbox chan any
}

func Create(ctx *domain.Context) *Core {

	return &Core{
		ctx: ctx,
		state: &domain.State{
			Status: common.OpNeutral,
		},
		mailbox: ctx.Hub.Sub(common.CommandScatterPlanStart, common.CommandGatherPlanStart),
	}
}

func (c *Core) Start() error {
	err := c.sanityCheck()
	if err != nil {
		return err
	}

	unraid, err := c.getStatus()
	if err != nil {
		return err
	}

	c.state.Unraid = unraid

	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return err
	}

	c.sid = sid

	go c.mailboxHandler()

	return nil
}

func (c *Core) Stop() error {
	return nil
}

func (c *Core) GetConfig() *domain.Config {
	return &c.ctx.Config
}

func (c *Core) GetState() *domain.State {
	return c.state
}

func (c *Core) GetStorage() *domain.Unraid {
	unraid, err := c.getStatus()
	if err != nil {
		logger.Yellow("unable to get storage: %s", err)
	} else {
		c.state.Unraid = unraid
	}

	return c.state.Unraid
}

func (c *Core) GetOperation() *domain.Operation {
	return c.state.Operation
}

func (c *Core) GetHistory() *domain.History {
	c.state.History.LastChecked = time.Now()
	return c.state.History
}

func (c *Core) mailboxHandler() {
	for p := range c.mailbox {
		if c.state.Status != common.OpNeutral {
			logger.Yellow("unbalance is busy: %d", c.state.Status)
			continue
		}

		packet := p.(domain.Packet)
		switch packet.Topic {
		case common.CommandScatterPlanStart:
			var setup domain.ScatterSetup
			err := lib.Bind(packet.Payload, &setup)
			if err != nil {
				logger.Red("unable to unmarshal packet: %+v (%s)", packet.Payload, err)
				continue
			}
			go c.scatterPlanPrepare(setup)
		}
	}
}
