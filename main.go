package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
)

// Version of application
var Version string

func init() {
	if Version == "" {
		Version = "v0.0.0"
	}
}

var (
	version    = flag.Bool("version", false, "version")
	configPath = flag.String("config", "config.toml", "config file path")
	config     Config
	server     http.Server
)

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
	}

	err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/events", handleEvents)
	mux.HandleFunc("/users", handleUsers)

	server.Addr = config.ListenAddress
	server.Handler = cors.Default().Handler(mux)

	err = server.ListenAndServe()

	if err == http.ErrServerClosed {
		return
	}
	log.Fatal(err)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownTimeout := time.Duration(config.ShutdownTimeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal("shutdown error:", err)

	}
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
}
