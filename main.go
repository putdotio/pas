package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/putdotio/pas/internal/analytics"
	"github.com/putdotio/pas/internal/handler"
	"github.com/putdotio/pas/internal/server"
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
)

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		return
	}

	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", config.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	analytics := analytics.New(db, config.Secret, config.User, config.Events)
	handler := handler.New(analytics)
	server := server.New(config.ListenAddress, handler)

	go server.ListenAndServe()

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
