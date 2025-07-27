package wshub

import (
	"encoding/json"
	"reflect"
	"testing"
)

// --- Example struct that might be used in the Data field ---
type UserDetails struct {
	UserIdentifier int                    `json:"userIdentifier"`
	EmailAddress   string                 `json:"emailAddress"`
	IsActive       bool                   `json:"isActive"`
	Metadata       map[string]interface{} `json:"metadata"`
	Tags           []string               `json:"tags"`
}

func TestToCapitalizedJSON(t *testing.T) {
	// Define a struct for our table-driven tests
	testCases := []struct {
		name         string      // Name of the test case
		input        interface{} // Input to the function
		expectedJSON string      // The expected JSON output string
		expectError  bool        // Whether we expect an error
	}{
		{
			name: "Simple Struct",
			input: UserDetails{
				UserIdentifier: 101,
				EmailAddress:   "test@example.com",
				IsActive:       true,
			},
			expectedJSON: `{
				"UserIdentifier": 101,
				"EmailAddress": "test@example.com",
				"IsActive": true,
				"Metadata": null,
				"Tags": null
			}`,
		},
		{
			name: "Struct with Nested Map",
			input: UserDetails{
				UserIdentifier: 102,
				EmailAddress:   "another@example.com",
				IsActive:       false,
				Metadata: map[string]interface{}{
					"lastLogin":  "2024-01-01T12:00:00Z",
					"loginCount": 5,
				},
				Tags: []string{"beta", "tester"},
			},
			expectedJSON: `{
				"UserIdentifier": 102,
				"EmailAddress": "another@example.com",
				"IsActive": false,
				"Metadata": {
					"LastLogin": "2024-01-01T12:00:00Z",
					"LoginCount": 5
				},
				"Tags": ["beta", "tester"]
			}`,
		},
		{
			name: "Simple Map",
			input: map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
			},
			expectedJSON: `{
				"FirstName": "John",
				"LastName": "Doe"
			}`,
		},
		{
			name: "Nested Map and Slice",
			input: map[string]interface{}{
				"event": "user.created",
				"payload": map[string]interface{}{
					"userName": "jdoe",
					"roles": []interface{}{
						"editor",
						map[string]interface{}{"permissionLevel": 4},
					},
				},
			},
			expectedJSON: `{
				"Event": "user.created",
				"Payload": {
					"UserName": "jdoe",
					"Roles": [
						"editor",
						{
							"PermissionLevel": 4
						}
					]
				}
			}`,
		},
		{
			name: "Top-level Slice with Structs",
			input: []UserDetails{
				{UserIdentifier: 201, EmailAddress: "user1@test.com"},
				{UserIdentifier: 202, EmailAddress: "user2@test.com"},
			},
			expectedJSON: `[
				{
					"UserIdentifier": 201, "EmailAddress": "user1@test.com", "IsActive": false, "Metadata": null, "Tags": null
				},
				{
					"UserIdentifier": 202, "EmailAddress": "user2@test.com", "IsActive": false, "Metadata": null, "Tags": null
				}
			]`,
		},
		{
			name:         "Nil Input",
			input:        nil,
			expectedJSON: `null`,
		},
		{
			name:         "Empty map",
			input:        map[string]interface{}{},
			expectedJSON: `{}`,
		},
	}

	// --- Test Runner ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function we are testing
			actualBytes, err := toCapitalizedJSON(tc.input)

			// Check for an unexpected error
			if !tc.expectError && err != nil {
				t.Fatalf("ToCapitalizedJSON() returned an unexpected error: %v", err)
			}
			// Check for an expected error that did not occur
			if tc.expectError && err == nil {
				t.Fatalf("ToCapitalizedJSON() was expected to return an error, but it did not")
			}

			// To reliably compare JSON, we unmarshal both the actual and expected
			// results into a generic interface{} and use reflect.DeepEqual.
			// This avoids issues with whitespace and key ordering.
			var actualResult interface{}
			if err := json.Unmarshal(actualBytes, &actualResult); err != nil {
				t.Fatalf("Failed to unmarshal actual result: %v", err)
			}

			var expectedResult interface{}
			if err := json.Unmarshal([]byte(tc.expectedJSON), &expectedResult); err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			// Compare the results
			if !reflect.DeepEqual(actualResult, expectedResult) {
				// Use MarshalIndent to get a pretty-printed version for easier comparison
				prettyActual, _ := json.MarshalIndent(actualResult, "", "  ")
				prettyExpected, _ := json.MarshalIndent(expectedResult, "", "  ")
				t.Errorf("Result does not match expected.\nGot:\n%s\n\nWant:\n%s", string(prettyActual), string(prettyExpected))
			}
		})
	}
}
