package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cyrildever/treee/api"
	"github.com/cyrildever/treee/common/logger"
	"github.com/cyrildever/treee/config"
	"github.com/cyrildever/treee/core/index"
)

/** Usage:
 *
 *	To launch the Treee indexing engine as a micro-service:
 *	`$ ./treee -p 7001 -h localhost -init 101`
 *
 *	Stop it with Ctrl^c
 */
func main() {
	log := logger.Init("main", "application")
	conf, err := config.InitConfig(false)
	if err != nil {
		panic(err)
	}

	treee, err := index.Load(conf.IndexPath)
	if err != nil {
		log.Warn("Index doesn't exist, building one...", "error", err)
		treee, err = index.New(conf.InitPrime)
		if err != nil {
			log.Crit("Unable to instantiate new index", "error", err)
			return
		}
		log.Info("Index created", "initPrime", treee.InitPrime)
	} else {
		log.Info("Index up and running", "size", treee.Size(), "initPrime", treee.InitPrime)
	}

	index.Current = treee

	willGracefullyStopIndex()

	api.InitHTTPServer(conf)
}

// willGracefullyStopIndex ...
func willGracefullyStopIndex() {
	log := logger.Init("main", "terminating")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Goodbye ~")
		os.Exit(1)
	}()
}
