package main

import (
	"estiam/routes"
	"estiam/dictionary"
	"fmt"
	"net/http"
)

func main() {
    dict, err := dictionary.New("dictionary/dictionary.json")
    if err != nil {
        fmt.Println("An error occurred while initializing the dictionary: ", err)
        return
    }

    router := routes.AddRoutes(dict)

    fmt.Println("Server listening on :8080")
    http.ListenAndServe(":8080", router)
}