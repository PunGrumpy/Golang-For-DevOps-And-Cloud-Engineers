package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// "page":"assignment1",
// "words":["eigth","eigth","two","seven","one"],
// "percentages":{"eigth":0.5,"one":0.33,"seven":0.88,"two":0.66},
// "special":["one","two",null],
// "extraSpecial":[1,2,"3"]
type Assignment struct {
	Page         string             `json:"page"`
	Words        []string           `json:"words"`
	Percentages  map[string]float64 `json:"percentages"`
	Special      []*string          `json:"special"`
	ExtraSpecial []any              `json:"extraSpecial"`
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: assignment <url>")
		return
	}

	if _, err := url.ParseRequestURI(args[1]); err != nil {
		log.Fatalf("Invalid URL: %s\n", err)
		os.Exit(1)
	}

	response, err := http.Get(args[1])
	if err != nil {
		log.Fatalf("Error: %s\n", err)
		os.Exit(1)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
		os.Exit(1)
	}

	var assignment Assignment
	if err := json.Unmarshal(body, &assignment); err != nil {
		log.Fatalf("Error: %s\n", err)
		os.Exit(1)
	}

	percentages := make(map[string]float64)
	specials := make([]string, len(assignment.Special))
	extraSpecials := make([]any, len(assignment.ExtraSpecial))

	// Page
	fmt.Printf("Page: %v\n", assignment.Page)

	// Words
	fmt.Printf("Words: %s\n", strings.Join(assignment.Words, ", "))

	// Percentages
	for word, percentage := range assignment.Percentages {
		percentages[word] = percentage
	}
	fmt.Printf("Percentages: %v\n", percentages)

	// Special
	for i, special := range assignment.Special {
		if special == nil {
			specials[i] = "null"
			continue
		}
		specials[i] = *special
	}
	fmt.Printf("Special: %v\n", strings.Join(specials, ", "))

	// ExtraSpecial
	for i, extraSpecial := range assignment.ExtraSpecial {
		extraSpecials[i] = extraSpecial
		if _, ok := extraSpecial.(string); ok {
			extraSpecials[i] = fmt.Sprintf("\"%v\"", extraSpecial)
		}
		if _, ok := extraSpecial.(float64); ok {
			extraSpecials[i] = fmt.Sprintf("%v", extraSpecial)
		}
		if i != len(assignment.ExtraSpecial)-1 {
			extraSpecials[i] = fmt.Sprintf("%v,", extraSpecials[i])
		}
	}
	fmt.Printf("ExtraSpecial: %v\n", extraSpecials)
}
