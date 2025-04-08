package prerun

import (
	"fmt"
)

// Define the allowed output formats
var ValidOutputFormats = []string{"dot", "puml", "mmd"}

func ValidateGraphPreRunE(outputFormat string) error {
	if !contains(ValidOutputFormats, outputFormat) {
		return fmt.Errorf("Invalid output format: %s. Allowed formats are: %v", outputFormat, ValidOutputFormats)
	}

	return nil
}

// Helper function to check if a string is in a slice of strings
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
