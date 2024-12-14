package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/g3ksa/lab9_otrpo/internal"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	httpAddr := os.Getenv("HTTP_PORT")
	if httpAddr == "" {
		httpAddr = "8080"
	}
	if !strings.HasPrefix(httpAddr, ":") {
		httpAddr = ":" + httpAddr
	}

	ctx := context.Background()
	redisClient := internal.NewRedisClient(ctx, redisAddr)

	hub := internal.NewHub(redisClient)
	go hub.Run(ctx)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		internal.ServeWs(hub, w, r)
	})

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	fmt.Println("Server listening on", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
