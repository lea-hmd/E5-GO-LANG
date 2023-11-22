package main

import (
	"encoding/json"
	"estiam/dictionary"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	dict, err := dictionary.New("dictionary/dictionary.json")
	if err != nil {
		fmt.Println("An error occurred while initializing the dictionary : ", err)
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/word", AddWordHandler(dict)).Methods("POST")
	router.HandleFunc("/api/word/{word}", GetWordHandler(dict)).Methods("GET")
	router.HandleFunc("/api/word/{word}", RemoveWordHandler(dict)).Methods("DELETE")

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", router)
}

func RemoveWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		word, exists := params["word"]
		if !exists {
			http.Error(w, "Missing 'word' parameter", http.StatusBadRequest)
			return
		}

		err := d.Remove(word)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error removing word '%s': %v", word, err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Word '%s' removed successfully!", word)
	}
}

func GetWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		word, exists := params["word"]
		if !exists {
			http.Error(w, "Missing 'word' parameter", http.StatusBadRequest)
			return
		}

		entry, err := d.Get(word)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting definition for word '%s': %v", word, err), http.StatusNotFound)
			return
		}

		response := map[string]string{"definition": entry.Definition}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func AddWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		word, exists := data["word"]
		if !exists {
			http.Error(w, "Missing 'word' in request payload", http.StatusBadRequest)
			return
		}

		definition, exists := data["definition"]
		if !exists {
			http.Error(w, "Missing 'definition' in request payload", http.StatusBadRequest)
			return
		}

		err := d.Add(word, definition)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding word '%s': %v", word, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Word '%s' added successfully!", word)
	}
}
