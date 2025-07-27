package wshub

import (
	"encoding/json"
	"fmt"
	"unicode"
)

// capitalize is a helper function to safely capitalize the first letter of a string.
// It's robust against empty strings.
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// capitalizeKeys recursively traverses an interface{} and capitalizes the keys of any maps it finds.
func capitalizeKeys(data interface{}) interface{} {
	// Use a type switch to handle the different types of data we might encounter.
	switch value := data.(type) {
	// If it's a map, we iterate over its keys and values.
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range value {
			// Capitalize the key and recursively process the value.
			newMap[capitalize(k)] = capitalizeKeys(v)
		}
		return newMap

	// If it's a slice, we iterate over its elements.
	case []interface{}:
		// The slice itself doesn't have keys, but its elements might.
		newSlice := make([]interface{}, len(value))
		for i, v := range value {
			// Recursively process each element in the slice.
			newSlice[i] = capitalizeKeys(v)
		}
		return newSlice

	// For any other type (string, int, bool, etc.), return it as is.
	default:
		return data
	}
}

// toCapitalizedJSON marshals any data structure (including structs) to a JSON string
// with all keys having their first letter capitalized.
func toCapitalizedJSON(payload interface{}) ([]byte, error) {
	// Step 1: Marshal the data to JSON. This respects the `json` tags on any structs.
	tempJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to perform initial marshal: %w", err)
	}

	// Step 2: Unmarshal the JSON into a generic interface{}.
	// This converts all JSON objects into map[string]interface{}, regardless of the original type.
	var genericData interface{}
	if err := json.Unmarshal(tempJSON, &genericData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into generic interface: %w", err)
	}

	// Step 3: Recursively capitalize the keys of the generic data structure.
	capitalizedData := capitalizeKeys(genericData)

	// Step 4: Marshal the final, capitalized data structure back to JSON.
	return json.MarshalIndent(capitalizedData, "", "  ")
}
