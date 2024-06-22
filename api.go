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

func getRandomArray(length int, max int, op int) []int {
	rand.Seed(time.Now().UnixNano())
	array := make([]int, length)
	if op == 0 {
		array[0] = 0 // Ensure the first element is always 0
	}
	for i := op; i < length; i++ {
		array[i] = rand.Intn(max)
	}
	return array
}

// Function to generate random values for the array
func generateRandomArray() []int {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate random values within specified ranges
	arr := []int{
		rand.Intn(25),    // 0 to 24
		rand.Intn(60),    // 0 to 59
		rand.Intn(8) + 1, // 1 to 8 	rand.Intn(9),   // 0 to 8
		rand.Intn(121),   // 0 to 120
	}
	return arr
}

func createRandomJSON() Data {
	data := Data{
		Escenarios: map[string][]int{
			"1": getRandomArray(31, 2500000000, 1),
		},
		Ciclos: map[string][]int{
			"1": getRandomArray(31, 15000, 1),
			"2": getRandomArray(31, 15000, 1),
			"3": getRandomArray(31, 21000, 1),
			"4": getRandomArray(31, 21000, 1),
			"5": getRandomArray(31, 21000, 1),
			"6": getRandomArray(31, 21000, 1),
			"7": getRandomArray(31, 21000, 1),
			"8": getRandomArray(31, 21000, 1),
		},
		Eventos: map[string][]int{
			"1": {0, 0, 1, 0}, // Ensure the first Eventos element is {0,0,1,0}
			"2": generateRandomArray(),
			"3": generateRandomArray(),
			"4": generateRandomArray(),
			"5": generateRandomArray(),
			"6": generateRandomArray(),
			"7": generateRandomArray(),
			"8": generateRandomArray(),
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
