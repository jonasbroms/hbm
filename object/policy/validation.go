package policy

import (
	"fmt"
	"slices"
)

func isValidFilterKeys(filters map[string]string) error {
	validKeys := []string{
		"name",
		"user",
		"group",
		"resource-type",
		"resource-value",
		"resource-options",
		"collection",
	}

	for k := range filters {
		if !slices.Contains(validKeys, k) {
			return fmt.Errorf("%s is not a valid filter key", k)
		}
	}

	return nil
}
