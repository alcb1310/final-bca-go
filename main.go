package main

import (
	"fmt"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/router"
)

func main() {
	r := router.NewRouter()

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r.Router)
}
