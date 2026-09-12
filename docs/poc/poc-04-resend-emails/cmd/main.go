package main

import (
	"fmt"
	"os"

	"github.com/Espectro0/rutadeorigen-docs/internal/email"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env")
	}

	client, err := email.NewClient()
	if err != nil {
		fmt.Println("An error was expected creating Resend's client: ", err)
		os.Exit(1)
	}

	sender := os.Getenv("EMAIL_FROM")
	receiver := os.Getenv("EMAIL_TO")

	subject, body := email.BusinessInvitation("Lorem Ipsum")
	if err := email.Send(client, fmt.Sprintf("Ori - Ruta de Origen <%s>", sender), receiver, subject, body, sender); err != nil {
		fmt.Println("An error was expected: ", err)
		os.Exit(1)
	}
}
