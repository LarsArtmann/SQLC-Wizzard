package main

import (
	"fmt"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
)

func main() {
	tmpl := templates.NewHobbyTemplate()
	features := tmpl.RequiredFeatures()
	fmt.Printf("Type: %T\n", features)
	fmt.Printf("Value: %v\n", features)
	fmt.Printf("Length: %d\n", len(features))

	if len(features) > 0 {
		fmt.Printf("First element: %v (type: %T)\n", features[0], features[0])
	}
}
