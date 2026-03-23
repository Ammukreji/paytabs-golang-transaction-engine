package main

import (
	"fmt"
	"log"
	"net/http"

	"paytabs-assessment/handlers"
	"paytabs-assessment/models"
	"paytabs-assessment/storage"
	"paytabs-assessment/utils"
)

func runSeed() {
	storage.InitDB()
	storage.DB.AddCard(models.Card{
		CardNumber: "4123456789012345",
		CardHolder: "John Doe",
		PinHash:    utils.HashPIN("1234"),
		Balance:    1000,
		Status:     "ACTIVE",
	})
	fmt.Println("Database seeded with sample card:")
	fmt.Println("Card Number: 4123456789012345")
	fmt.Println("Name: John Doe")
	fmt.Println("PIN: 1234")
	fmt.Println("Balance: 1000")
	fmt.Println("Status: ACTIVE")
	fmt.Println("-------------------------------------------------")
}

func main() {
	runSeed()

	http.HandleFunc("/api/transaction", handlers.HandleTransaction)
	http.HandleFunc("/api/card/", handlers.HandleCardEndpoints)

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
