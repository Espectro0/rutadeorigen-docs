package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Espectro0/rutadeorigen-docs/internal/auth"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env")
	}

	cfg, err := auth.NewClient()
	if err != nil {
		log.Fatal("Error creating WorkOS client ", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", cfg.LoginHandler)
	mux.HandleFunc("/callback", cfg.CallbackHandler)
	mux.HandleFunc("/api/me", cfg.MeHandler)
	mux.HandleFunc("/api/logout", cfg.LogoutHandler)

	mux.Handle("/", http.FileServer(http.Dir("../frontend")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
