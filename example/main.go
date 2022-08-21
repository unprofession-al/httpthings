package main

import (
	"flag"
	"fmt"
	"os"
)

type App struct {
	listener string
	static   string
}

func main() {
	var app App
	flag.StringVar(&app.listener, "listener", "0.0.0.0:8765", "ip/port to listen on")
	flag.StringVar(&app.static, "static", "", "serve given dir as http root")
	flag.Parse()

	s, err := NewServer(app.listener, app.static)
	exitOnErr(err)
	s.run()
}

func exitOnErr(errs ...error) {
	errNotNil := false
	for _, err := range errs {
		if err == nil {
			continue
		}
		errNotNil = true
		fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
	}
	if errNotNil {
		fmt.Print("\n")
		os.Exit(-1)
	}
}
