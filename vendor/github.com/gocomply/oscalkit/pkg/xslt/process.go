package xslt

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Transform inputFile using stylesheet
func Transform(stylesheetPath, inputFile string) (*bytes.Buffer, error) {
	xsltprocCmd := exec.Command("xsltproc", stylesheetPath, inputFile)
	xsltprocCmdOutput := &bytes.Buffer{}
	xsltprocCmdErr := &bytes.Buffer{}
	xsltprocCmd.Stdout = xsltprocCmdOutput
	xsltprocCmd.Stderr = xsltprocCmdErr

	err := xsltprocCmd.Run()
	if err != nil || xsltprocCmdErr.Len() > 0 {
		return nil, fmt.Errorf("Error running xsltproc: %v, stderr was: %s", err, xsltprocCmdErr.String())
	}
	return xsltprocCmdOutput, nil
}
