package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		// Read the JSON file
		file, err := os.Open("program.json")
		if err != nil {
			http.Error(w, "Could not open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Read file content
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Could not read file", http.StatusInternalServerError)
			return
		}

		// Write content as JSON response
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
