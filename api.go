package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Data struct {
	Escenarios map[string][]int `json:"escenarios"`
	Ciclos     map[string][]int `json:"ciclos"`
	Eventos    map[string][]int `json:"eventos"`
}

func getRandomArray(length int, max int) []int {
	rand.Seed(time.Now().UnixNano())
	array := make([]int, length)
	for i := 0; i < length; i++ {
		array[i] = rand.Intn(max)
	}
	return array
}

func createRandomJSON() Data {
	data := Data{
		Escenarios: map[string][]int{
			"1": getRandomArray(31, 2500000000),
		},
		Ciclos: map[string][]int{
			"1": getRandomArray(31, 15000),
			"2": getRandomArray(31, 15000),
			"3": getRandomArray(31, 21000),
			"4": getRandomArray(31, 21000),
			"5": getRandomArray(31, 21000),
			"6": getRandomArray(31, 21000),
			"7": getRandomArray(31, 21000),
			"8": getRandomArray(31, 21000),
		},
		Eventos: map[string][]int{
			"1": getRandomArray(4, 10),
			"2": getRandomArray(4, 20),
			"3": getRandomArray(4, 20),
			"4": getRandomArray(4, 20),
			"5": getRandomArray(4, 20),
			"6": getRandomArray(4, 20),
			"7": getRandomArray(4, 20),
			"8": getRandomArray(4, 20),
		},
	}
	return data
}

func saveJSONToFile(filename string, data Data) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(data)
}

func main() {
	http.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		data0 := createRandomJSON()
		err := saveJSONToFile("random_data.json", data0)
		if err != nil {
			panic(err)
		}
		// Read the JSON file
		file, err := os.Open("random_data.json")
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

	http.HandleFunc("/static", func(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Server is running on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
