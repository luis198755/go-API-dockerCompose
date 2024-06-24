package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var fases = 10

type Data struct {
	Fases      map[string][]int `json:"fases"`
	Escenarios map[string][]int `json:"escenarios"`
	Ciclos     map[string][]int `json:"ciclos"`
	Eventos    map[string][]int `json:"eventos"`
}

func getRandomArray(length int, max int, op int) []int {
	//rand.Seed(time.Now().UnixNano())
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	array := make([]int, length)
	if op == 0 {
		array[0] = 0 // Ensure the first element is always 0
	}
	for i := op; i < length; i++ {
		array[i] = rng.Intn(max)
	}
	return array
}

// Function to generate random values for the array
func generateRandomArray() []int {
	// Seed the random number generator
	//rand.Seed(time.Now().UnixNano())
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Generate random values within specified ranges
	arr := []int{
		rng.Intn(25),    // 0 to 24
		rng.Intn(60),    // 0 to 59
		rng.Intn(8) + 1, //rand.Intn(8) + 1, // 1 to 8 	rand.Intn(9),   // 0 to 8
		rng.Intn(121),   // 0 to 120
	}
	return arr
}

func generateStructure() []int {
	// Initialize the fixed values
	fixedValues := []int{0, 10000, 375, 375, 375, 375, 375, 375, 375, 375, 3000, 20000, 375, 375, 375, 375, 375, 375, 375, 375, 3000, 30000, 375, 375, 375, 375, 375, 375, 375, 375, 3000}
	structure := make([]int, len(fixedValues))
	copy(structure, fixedValues)

	// Iterate over the structure and replace non-fixed values with random values
	for i := 0; i < len(structure); i++ {
		if structure[i] != 0 && structure[i] != 375 && structure[i] != 3000 {
			structure[i] = GenerateRandomValue() // Random value between 0 and 999 (you can change the range if needed)
		}
	}

	return structure
}

// GenerateRandomValue generates a random value from 1000 to 120000 in multiples of 1000
func GenerateRandomValue() int {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	min := 1   // Minimum multiplier (1000 * 1 = 1000)
	max := 120 // Maximum multiplier (1000 * 120 = 120000)
	randomMultiplier := rng.Intn(max-min+1) + min
	return randomMultiplier * 1000
}

func createRandomJSON() Data {
	data := Data{
		Fases: map[string][]int{
			"1": {fases},
		},
		Escenarios: map[string][]int{
			"1": getRandomArray(31, 2500000000, 1),
		},
		Ciclos: map[string][]int{
			"1": generateStructure(),
			"2": generateStructure(),
			"3": generateStructure(),
			"4": generateStructure(),
			"5": generateStructure(),
			"6": generateStructure(),
			"7": generateStructure(),
			"8": generateStructure(),
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

func randomJson(w http.ResponseWriter, r *http.Request) {
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
}

func staicJson(w http.ResponseWriter, r *http.Request) {
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
}

func serveStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

func plain(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTML Vs Plain")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "<h1>Hello World<h1>")
}

func html(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTML Vs Plain")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Hello World<h1>")
}

func program(w http.ResponseWriter, r *http.Request) {
	var fileName = "index.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Error parsing file")
		return
	}
	err = t.ExecuteTemplate(w, fileName, nil)
	if err != nil {
		fmt.Println("Error when executing template", err)
	}
}

func handleFunctions(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintln(w, "Hello World")
	case "/random":
		randomJson(w, r)
	case "/static":
		staicJson(w, r)
	case "/plain":
		plain(w, r)
	case "/html":
		html(w, r)
	case "/status":
		serveStatus(w, r)
	case "/program":
		program(w, r)
	default:
		fmt.Fprint(w, "Error")
	}
}

func main() {
	// http.HandleFunc("/random", randomJson)
	// http.HandleFunc("/static", staicJson)
	// http.HandleFunc("/status", serveStatus)
	// http.HandleFunc("/plain", plain)
	// http.HandleFunc("/html", html)

	http.HandleFunc("/", handleFunctions)

	fmt.Println("Server is running on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
