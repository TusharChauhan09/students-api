package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TusharChauhan09/students-api/internal/config"
)

func main(){
	// 1 configration
	cfg := config.MustLoad()

	// 2 database

	// 3 router 
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
	})

	// server
	server := http.Server {
		Addr: cfg.HTTPServer.Address,
		Handler: router,
	}

	fmt.Printf("server started %s", cfg.HTTPServer.Address)

	err := server.ListenAndServe(); if err != nil{
		log.Fatal("failed to start server")
	}

}
