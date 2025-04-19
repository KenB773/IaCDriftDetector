// Terraform state parser
package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type TFState struct {
	Resources []TFResource `json:"resources"`
}

type TFResource struct {
	Type      string               `json:"type"`
	Name      string               `json:"name"`
	Instances []TFResourceInstance `json:"instances"`
}

type TFResourceInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// ParseTFState loads and parses a Terraform .tfstate file from disk.
func ParseTFState(filePath string) (*TFState, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var state TFState
	if err := json.Unmarshal(bytes, &state); err != nil {
		return nil, err
	}

	return &state, nil
}
