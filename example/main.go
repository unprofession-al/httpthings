package main

import (
	"flag"
	"fmt"
	"os"
)

var app App

type App struct {
	listener string
	static   string
}

func init() {
	flag.StringVar(&app.listener, "listener", "0.0.0.0:8765", "ip/port to listen on")
	flag.StringVar(&app.static, "static", "", "serve given dir as http root")
}

func main() {
	flag.Parse()

	s := NewServer(app.listener, app.static)
	s.run()
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
