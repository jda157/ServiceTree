package main

import (
	"fmt"
	"log"
	"math/rand"
	"myprojects/26_concurrency/internal/app/http-handlers"
	"myprojects/26_concurrency/internal/app/service-tree"
	"net/http"
	"time"
)

const (
	setupFile = "pkg/setup/setup.json"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	sr := service_tree.New()

	if err := sr.InitFromConfig(setupFile); err != nil {
		log.Fatalf("init failed %v", err)
	}

	http.HandleFunc("/", http_handlers.HttpHandler(sr))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("unable to start server", err)
	}

}
