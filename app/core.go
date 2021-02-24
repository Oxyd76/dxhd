package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dakyskye/dxhd/parser"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dakyskye/dxhd/logger"
)

type App struct {
	execName string
	ctx      context.Context
	cancel   context.CancelFunc
	cli      *kingpin.Application
	opts     options
}

type serverResponse string

const (
	shutoff serverResponse = "shutoff"
	reload  serverResponse = "reload"
)

func (a *App) Start() (err error) {
	logger.L().Debug("trying to start the server")

	// * parse config file
	// * start signal handler
	// * set up X11 connection
	// * listen for keybindings
	for {
		err = a.start()
		if err != nil {
			break
		}

		server := make(chan serverResponse, 1)
		go a.serveSignals(server)

		command := <-server
		switch command {
		case shutoff:
			a.cancel()
		case reload:
		}

		logger.L().WithField("command", command).Debug("received a command")

		break
	}

	return
}

func (a *App) start() (err error) {
	p, err := parser.New(a.opts.config)
	if err != nil {
		return
	}

	err = p.Parse()
	if err != nil {
		return
	}

	// TODO: data, err := p.Collect()
	return
}

func (a *App) serveSignals(server chan<- serverResponse) {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	logger.L().Debug("serving os signals")

	select {
	case sig := <-signals:
		switch sig {
		case syscall.SIGUSR1, syscall.SIGUSR2:
			server <- reload
		default:
			server <- shutoff
		}
	case <-a.ctx.Done():
		logger.L().WithError(a.ctx.Err()).Debug("main app context done")
		server <- shutoff
	}
}
