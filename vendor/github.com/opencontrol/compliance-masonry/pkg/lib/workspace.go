/*
 Copyright (C) 2018 OpenControl Contributors. See LICENSE.md for license.
*/

package lib

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/opencontrol/compliance-masonry/pkg/cli/clierrors"
	"github.com/opencontrol/compliance-masonry/pkg/lib/common"
	"github.com/opencontrol/compliance-masonry/pkg/lib/result"
)

// localWorkspace struct combines components, standards, and a certification data
// For more information on the opencontrol schema visit: https://github.com/opencontrol/schemas
type localWorkspace struct {
	components     *componentsMap
	standards      *standardsMap
	justifications *result.Justifications
	certification  common.Certification
}

// getKey extracts a component key from the filepath
func getKey(filePath string) string {
	_, key := filepath.Split(filePath)
	return key
}

// NewWorkspace initializes an empty OpenControl struct
func NewWorkspace() common.Workspace {
	return &localWorkspace{
		justifications: result.NewJustifications(),
		components:     newComponents(),
		standards:      newStandards(),
	}
}

// LoadData creates a new instance of OpenControl struct and loads
// the components, standards, and certification data.
func LoadData(openControlDir string, certificationPath string) (common.Workspace, []error) {
	var wg sync.WaitGroup
	ws := NewWorkspace()
	wg.Add(3)
	var componentsErrs, standardsErrs []error
	var certificationErr error
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		componentsErrs = ws.LoadComponents(filepath.Join(openControlDir, "components"))
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		standardsErrs = ws.LoadStandards(filepath.Join(openControlDir, "standards"))
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		certificationErr = ws.LoadCertification(certificationPath)
	}(&wg)
	wg.Wait()

	var errs []error
	if certificationErr != nil {
		errs = append(errs, certificationErr)
	}
	errs = append(errs, componentsErrs...)
	errs = append(errs, standardsErrs...)

	return ws, errs
}

// LoadComponents loads multiple components by searching for components in a
// given directory
func (ws *localWorkspace) LoadComponents(directory string) []error {
	var wg sync.WaitGroup
	componentsDir, err := ioutil.ReadDir(directory)
	if err != nil {
		return []error{errors.New("Error: Unable to read the directory " + directory)}
	}
	errChannel := make(chan error, len(componentsDir))
	wg.Add(len(componentsDir))
	for _, componentDir := range componentsDir {
		go func(componentDir os.FileInfo, wg *sync.WaitGroup) {
			if componentDir.IsDir() {
				componentDir := filepath.Join(directory, componentDir.Name())
				errChannel <- ws.LoadComponent(componentDir)
			}
			wg.Done()
		}(componentDir, &wg)
	}
	wg.Wait()
	close(errChannel)
	return convertErrChannelToErrorSlice(errChannel)
}

// LoadStandards loads multiple standards by searching for components in a
// given directory
func (ws *localWorkspace) LoadStandards(standardsDir string) []error {
	var wg sync.WaitGroup
	standardsFiles, err := ioutil.ReadDir(standardsDir)
	if err != nil {
		return []error{errors.New("Error: Unable to read the directory " + standardsDir)}
	}
	errChannel := make(chan error, len(standardsFiles))
	wg.Add(len(standardsFiles))
	for _, standardFile := range standardsFiles {
		go func(standardFile os.FileInfo, wg *sync.WaitGroup) {
			errChannel <- ws.LoadStandard(filepath.Join(standardsDir, standardFile.Name()))
			wg.Done()
		}(standardFile, &wg)
	}
	wg.Wait()
	close(errChannel)
	return convertErrChannelToErrorSlice(errChannel)
}

func convertErrChannelToErrorSlice(errs <-chan error) []error {
	errMessages := clierrors.NewMultiError()
	for err := range errs {
		if err != nil && len(err.Error()) > 0 {
			errMessages.Errors = append(errMessages.Errors, err)
		}
	}
	return errMessages.Errors
}

func (ws *localWorkspace) GetAllVerificationsWith(standardKey string, controlKey string) common.Verifications {
	return ws.justifications.Get(standardKey, controlKey)
}
