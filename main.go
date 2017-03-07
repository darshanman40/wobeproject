package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wobeproject/config"
	"github.com/wobeproject/handlers"
	"github.com/wobeproject/logger"
)

const port = ":8081"

func main() {
	// Parse the environment
	var env, filepath string
	flag.StringVar(&env, "env", "dev", "a string var")
	flag.StringVar(&filepath, "config", "config/config.toml", "config file path")
	flag.Parse()
	config.Load(filepath)
	cfg := config.GetApp(env)

	if cfg == nil {
		log.Fatal("config not found, ", env)
	}

	l := logger.NewLogger(cfg.Log)

	l.Info("Server is starting", map[string]interface{}{
		"environment": env,
		"port":        port,
	})

	// override ctrl + C and  close logger files
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(l logger.Logger) {
		<-c
		l.Info("Exiting Server", map[string]interface{}{
			"Closing": "Logger files",
		})
		l.CloseAll()
		os.Exit(1)
	}(l)

	handlers.InitHandlers()
	http.ListenAndServe(port, nil)

}
