package main

import (
	"compressor/cmd/compressor/server"
	"compressor/internal/compressor_grpc"
	"compressor/internal/config"
	"compressor/internal/postgres"
	"compressor/pkg/logger"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfiguration("./config/config.json")
	if err != nil {
		log.Fatalf("Could not instantiate config %+s", err)
	}

	newLogger, err := logger.NewLogger(cfg.Log)
	if err != nil {
		log.Fatalf("Could not instantiate log %+s", err)
	}

	db := postgres.New(newLogger, cfg)

	defer db.Close()

	dbReady := make(chan struct{})

	go func() {
		for {
			if err := db.Session.Ping(); err != nil {
				time.Sleep(1 * time.Second)
			} else {
				newLogger.Debugf("db is ready")
				dbReady <- struct{}{}
				return
			}
		}
	}()

	urlStorage, err := postgres.NewURLStorage(db)
	if err != nil {
		newLogger.Fatalf("Could not instantiate statements %+s", err)
	}
	newServer := server.NewServer(newLogger, urlStorage)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	compressor_grpc.RegisterCompressingServiceServer(grpcServer, newServer)
	lis, err := net.Listen("tcp", cfg.Grpc.Adress)
	if err != nil {
		newLogger.Debugf("cannot listen on ", cfg.Grpc.Adress)
	}

	go func() {
		<-dbReady
		newLogger.Debugf("grpc server started")
		err := grpcServer.Serve(lis)
		if err != nil {
			newLogger.Fatalf("grpc server stopped %+s", err)
		}
	}()

	newLogger.Debugf("server started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		newLogger.Debugf("system call:%+v", oscall)
		cancel()
	}()

	<-ctx.Done()

	newLogger.Debugf("server stopped")

	grpcServer.GracefulStop()

	newLogger.Debugf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
