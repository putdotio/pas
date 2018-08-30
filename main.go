package main

import (
	"context"
	"database/sql"
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
	db         *sql.DB
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

	db, err = sql.Open("mysql", config.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/events", handleEvents)
	mux.HandleFunc("/api/users", handleUsers)

	server.Addr = config.ListenAddress
	server.Handler = cors.Default().Handler(mux)

	go func() {
		err = server.ListenAndServe()
		if err == http.ErrServerClosed {
			return
		}
		log.Fatal(err)
	}()

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
