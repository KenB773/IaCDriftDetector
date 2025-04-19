// Drift comparison logic
package internal

import (
	"fmt"
)

type Drift struct {
	ResourceType string
	ID           string
	Issue        string
	Severity     string
}

// CompareStateWithAWS compares the Terraform state resources with AWS live resources.
func CompareStateWithAWS(state *TFState, awsResources []FetchedResource, include []string) []Drift {
	drift := []Drift{}

	// Build maps for quick lookups
	tfMap := make(map[string]TFResourceInstance)
	for _, res := range state.Resources {
		if len(include) > 0 && !Contains(include, res.Type) {
			continue
		}
		for _, inst := range res.Instances {
			id, ok := inst.Attributes["id"].(string)
			if !ok {
				continue
			}
			tfMap[res.Type+":"+id] = inst
		}
	}

	awsMap := make(map[string]FetchedResource)
	for _, res := range awsResources {
		if len(include) > 0 && !Contains(include, res.Type) {
			continue
		}
		awsMap[res.Type+":"+res.ID] = res
	}

	// Check for missing in AWS
	for key := range tfMap {
		if _, exists := awsMap[key]; !exists {
			parts := splitKey(key)
			drift = append(drift, Drift{
				ResourceType: parts[0],
				ID:           parts[1],
				Issue:        "Resource missing in AWS",
				Severity:     "critical",
			})
		}
	}

	// Check for extra in AWS (not in tfstate)
	for key := range awsMap {
		if _, exists := tfMap[key]; !exists {
			parts := splitKey(key)
			drift = append(drift, Drift{
				ResourceType: parts[0],
				ID:           parts[1],
				Issue:        "Resource not defined in Terraform",
				Severity:     "warning",
			})
		}
	}

	return drift
}

func splitKey(key string) [2]string {
	var parts [2]string
	split := make([]string, 2)
	n, _ := fmt.Sscanf(key, "%[^:]:%s", &split[0], &split[1])
	if n == 2 {
		parts[0] = split[0]
		parts[1] = split[1]
	}
	return parts
}
