package json_validation

import (
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema"
)

// Validate validates JSON file against a specific JSON schema.
func Validate(schemaPath, inputFile string) error {
	schema, err := jsonschema.Compile(schemaPath)
	if err != nil {
		return fmt.Errorf("Error compiling OSCAL schema: %v", err)
	}

	rawFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("Error opening file: %s, %v", inputFile, err)
	}
	defer rawFile.Close()

	if err = schema.Validate(rawFile); err != nil {
		return err
	}
	return nil
}
