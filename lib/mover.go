package lib

import (
	"apertoire.net/unbalance/message"
	"github.com/golang/glog"
	"io"
	"os"
	"path/filepath"
)

type PassThru struct {
	io.Reader
	Progress *message.Progress
	Ch       chan *message.Progress
}

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if err == nil {
		pt.Progress.CurrentCopied += uint64(n)
		pt.Ch <- pt.Progress
	}
}

type Mover struct {
	Src      string
	Dst      string
	err      error
	Progress *message.Progress

	progressCh chan *message.Progress
	doneCh     chan bool
}

func (self *Mover) visit(path string, info os.FileInfo, err error) (e error) {
	if info.IsDir() {

	} else {
		out := "\nPath: " + path + "\nSource: " + self.Dst

		Progress.CurrentFile = path
		Progress.CurrentSize = info.Size()

		self.ch <- self.Progress
	}

	return nil
}

func (self *Mover) Copy() (chan *message.Progress, chan bool) {
	self.progressCh = make(chan *message.Progress)
	self.doneCh = make(chan bool)

	go func() {
		filepath.Walk(self.Src, self.visit)
		defer close(self.progressCh)
		defer close(self.doneCh)
	}()

	return self.progressCh, self.doneCh
}
