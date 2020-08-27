package xml_validation

import (
	"bytes"
	"errors"
	"os/exec"
)

// Validate validates xml file against given schema
func Validate(schemaPath, inputFile string) error {
	xmllintCmd := exec.Command("xmllint", "--schema", schemaPath, inputFile, "--noout")

	xmllintCmdOutput := &bytes.Buffer{}
	xmllintCmdErr := &bytes.Buffer{}
	xmllintCmd.Stdout = xmllintCmdOutput
	xmllintCmd.Stderr = xmllintCmdErr

	if err := xmllintCmd.Run(); err != nil {
		stderr := xmllintCmdErr.String()
		if stderr == "" {
			return err
		}
		return errors.New(xmllintCmdErr.String())
	}
	return nil
}
