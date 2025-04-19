// Helper utilities
package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveToFile writes any struct to a file in JSON format
func SaveToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Contains checks if a string slice contains a value
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PrintError prints errors with a consistent prefix
func PrintError(msg string, err error) {
	fmt.Printf("❌ %s: %v\n", msg, err)
}
