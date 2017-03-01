package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wobeproject/handlers"
	"github.com/wobeproject/logger"
)

//var l logger.Logger
const port = ":8081"

func main() {
	// Parse the environment
	var env string
	flag.StringVar(&env, "env", "local", "a string var")
	flag.Parse()
	//config.Load(env)
	l := logger.NewLogger(env)

	l.Info("Server Start", map[string]interface{}{
		"environment": env,
	})
	l.Warning("Server Port", map[string]interface{}{
		"port": port,
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
