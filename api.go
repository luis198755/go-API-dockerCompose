package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var fases = 10

type Data struct {
	Fases      map[string][]int `json:"fases"`
	Escenarios map[string][]int `json:"escenarios"`
	Ciclos     map[string][]int `json:"ciclos"`
	Eventos    map[string][]int `json:"eventos"`
}

func getRandomArray(length int, max int, op int) []int {
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

func generateRandomArray() []int {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	arr := []int{
		rng.Intn(25),    // 0 to 24
		rng.Intn(60),    // 0 to 59
		rng.Intn(8) + 1, // 1 to 8
		rng.Intn(121),   // 0 to 120
	}
	return arr
}

func generateStructure() []int {
	fixedValues := []int{0, 10000, 375, 375, 375, 375, 375, 375, 375, 375, 3000, 20000, 375, 375, 375, 375, 375, 375, 375, 375, 3000, 30000, 375, 375, 375, 375, 375, 375, 375, 375, 3000}
	structure := make([]int, len(fixedValues))
	copy(structure, fixedValues)

	for i := 0; i < len(structure); i++ {
		if structure[i] != 0 && structure[i] != 375 && structure[i] != 3000 {
			structure[i] = GenerateRandomValue()
		}
	}

	return structure
}

func GenerateRandomValue() int {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	min := 1
	max := 120
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
			"1": {0, 0, 1, 0},
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

func randomJson(c *gin.Context) {
	data0 := createRandomJSON()
	err := saveJSONToFile("random_data.json", data0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	file, err := os.Open("random_data.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func staticJson(c *gin.Context) {
	file, err := os.Open("program.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func serveStatus(c *gin.Context) {
	c.String(http.StatusOK, "Ok")
}

func plain(c *gin.Context) {
	c.String(http.StatusOK, "<h1>Hello World</h1>")
}

func html(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, "<h1>Hello World</h1>")
}

func program(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func programsimu(c *gin.Context) {
	c.HTML(http.StatusOK, "simu.html", nil)
}

func programJsonRaw(c *gin.Context) {
	c.HTML(http.StatusOK, "json.html", nil)
}

func programJson(c *gin.Context) {
	var data Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := saveJSONToFile("program.json", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data saved successfully"})
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("*.html")

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	r.GET("/random", randomJson)
	r.GET("/static", staticJson)
	r.GET("/jsonProg", staticJson)
	r.GET("/plain", plain)
	r.GET("/html", html)
	r.GET("/status", serveStatus)
	r.GET("/programSem", program)
	r.POST("/program", programJson)
	r.GET("/json", programJsonRaw)
	r.GET("/simu", programsimu)

	fmt.Println("Server is running on port 80")
	r.Run(":80")
}