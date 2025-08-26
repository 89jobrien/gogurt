package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

func DocsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./docs")
	if err != nil {
		http.Error(w, "Could not read docs directory", http.StatusInternalServerError)
		return
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"files": filenames})
}
