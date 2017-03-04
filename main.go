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

//var l logger.Logger
const port = ":8081"

func main() {
	// Parse the environment
	var env, filepath string
	flag.StringVar(&env, "env", "dev", "a string var")
	flag.StringVar(&filepath, "config", "config/config.toml", "config file path")
	flag.Parse()
	config.Load(filepath)
	cfg := config.GetApp(env)

	//.GetLog(env)
	if cfg == nil {
		log.Fatal("config not found, ", env)
	}

	// for i, cfgL := range cfg.Log {
	// 	log.Println(i + " ")
	// 	log.Print(cfgL)
	// }
	//config.GetLog(a, s)
	//log.Println(config.GetLog(env, "info"))
	log.Println("configLog "+env+"\n", cfg.Log)
	l := logger.NewLogger(cfg.Log)

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
