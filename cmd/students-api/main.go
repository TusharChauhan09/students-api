package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TusharChauhan09/students-api/internal/config"
	"github.com/TusharChauhan09/students-api/internal/http/handlers/student"
	"github.com/TusharChauhan09/students-api/internal/storage/postgres"
	// "github.com/TusharChauhan09/students-api/internal/storage/sqlite"
)

func main(){
	// ! 1 configration
	cfg := config.MustLoad()


	// ! 2 database
	// _ ,err :=sqlite.New(cfg)
	storage , err := postgres.New(cfg)
	if err != nil{
		log.Fatal(err)
	}

	slog.Info("storage initalized", slog.String("env",cfg.Env))


	// ! 3 router 
	router := http.NewServeMux()


	router.HandleFunc("GET /api/students", student.New(storage))


	// ! 4 server
	server := http.Server {
		Addr: cfg.HTTPServer.Address,
		Handler: router,
	}

	slog.Info("server started ", slog.String("address",cfg.HTTPServer.Address))




	// intrupt shut down -> gracefull shut down
	// err := server.ListenAndServe(); if err != nil{
	// 	log.Fatal("failed to start server")
	// }

	// !  Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt , syscall.SIGINT, syscall.SIGTERM)

	go func () {
		if err := server.ListenAndServe(); err != nil{
			log.Fatal("failed to start server")
		}
	}()
	<-done

	slog.Info("shutting doen the server")
		
	ctx , cancel:=context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err !=nil {
		slog.Error("Failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
