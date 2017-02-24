package main

import (
	"flag"
	"net/http"

	"github.com/wobeproject/config"
	"github.com/wobeproject/handlers"
)

func main() {
	var env string
	flag.StringVar(&env, "env", "local", "a string var")
	flag.Parse()
	config.Load(env)

	http.HandleFunc("/", handlers.InputHandler)
	http.ListenAndServe(":8081", nil)

}
