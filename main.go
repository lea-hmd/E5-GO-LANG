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
	router.HandleFunc("/api/words", ListWordsHandler(dict)).Methods("GET")
	router.HandleFunc("/api/word/{word}", UpdateWordHandler(dict)).Methods("PUT")

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

		fmt.Fprintf(w, "Word '%s' removed successfully !", word)
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
			http.Error(w, fmt.Sprintf("Error getting definition for word '%s' : %v", word, err), http.StatusNotFound)
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
			http.Error(w, fmt.Sprintf("Error adding word '%s' : %v", word, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Word '%s' added successfully !", word)
	}
}
func ListWordsHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		words, entries, err := d.List()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error listing words : %v", err), http.StatusInternalServerError)
			return
		}

		response := make([]map[string]string, 0, len(words))
		for _, word := range words {
			entry := entries[word]
			wordData := map[string]string{"word": word, "definition": entry.Definition}
			response = append(response, wordData)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"words": response})
	}
}

func UpdateWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
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

		newDefinition, exists := data["definition"]
		if !exists {
			http.Error(w, "Missing 'definition' in request payload", http.StatusBadRequest)
			return
		}

		err := d.Update(word, newDefinition)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating word '%s' : %v", word, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Word '%s' updated successfully !", word)
	}
}