package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/judgeservice"
	"github.com/NCNUCodeOJ/BackendQuestionDatabase/models"
	router "github.com/NCNUCodeOJ/BackendQuestionDatabase/routers"
	"github.com/NCNUCodeOJ/BackendQuestionDatabase/styleservice"
	"github.com/gin-gonic/gin"
)

var srv *http.Server

func start() {
	models.Setup()

	judgeservice.Setup()
	styleservice.Setup()

	r := router.SetupRouter()
	if gin.Mode() == "debug" {
		srv = &http.Server{
			Addr:    "localhost:8080",
			Handler: r,
		}
	} else {
		srv = &http.Server{
			Addr:    ":8080",
			Handler: r,
		}
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()
}

func end() {
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func main() {
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	if arg == "ping" {
		resp, err := http.Get("http://localhost:8080/ping")
		if err != nil {
			os.Exit(1)
		}
		if resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		start()
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal, 1)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
		end()
	}
}
