package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alcb1310/final-bca-go/internal/database"
	"github.com/alcb1310/final-bca-go/internal/router"
	_ "github.com/joho/godotenv/autoload"
)

var port = os.Getenv("PORT")

func main() {
	db := database.New()
	if db == nil {
		fmt.Fprintf(os.Stderr, "New Router: Unable to connect to database\n")
		os.Exit(1)
	}
	r := router.NewRouter(db)
	if r == nil {
		os.Exit(1)
	}

	r.GenerateRoutes()

	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r.Router); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\r\n", err)
		os.Exit(1)
	}
}
