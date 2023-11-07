package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Prakhar2898/mongoapi/router"
)

func main() {
	fmt.Println("MongoDB API")
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":4040", router.Router()))
	fmt.Println("Listening at port 4040...")
}
