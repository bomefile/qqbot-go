package main

import (
	"fmt"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	_ = service.InitQQTokenFromEnv()

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})

	// callmsgback
	http.HandleFunc("/bot/callback", service.CallbackMSg)

	baseURL := "http://localhost:8080"
	log.Println("server starting", baseURL)
	log.Println("serve", baseURL+"/", baseURL+"/hello", baseURL+"/api/count", baseURL+"/bot/callback")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
