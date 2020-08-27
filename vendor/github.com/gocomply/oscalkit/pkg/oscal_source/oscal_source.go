package oscal_source

import (
	"fmt"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/types/oscal"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// OSCALSource is intermediary that handles IO and low-level common operations consistently for oscalkit
type OSCALSource struct {
	UserPath string
	file     *os.File
	oscal    *oscal.OSCAL
}

// Open creates new OSCALSource and load it up
func Open(path string) (*OSCALSource, error) {
	result := OSCALSource{UserPath: path}
	return &result, result.open()
}

func OpenFromReader(name string, r io.Reader) (*OSCALSource, error) {
	s := OSCALSource{UserPath: name}
	var err error
	if s.oscal, err = oscal.New(r); err != nil {
		return nil, fmt.Errorf("Cannot parse file: %v", err)
	}
	return &s, nil
}

func (s *OSCALSource) open() error {
	var err error
	path := s.UserPath
	if !filepath.IsAbs(path) {
		if path, err = filepath.Abs(path); err != nil {
			return fmt.Errorf("Cannot get absolute path: %v", err)
		}
	}
	if _, err = os.Stat(path); err != nil {
		return fmt.Errorf("Cannot stat %s, %v", path, err)
	}
	if s.file, err = os.Open(path); err != nil {
		return fmt.Errorf("Cannot open file %s: %v", path, err)
	}
	if s.oscal, err = oscal.New(s.file); err != nil {
		return fmt.Errorf("Cannot parse file: %v", err)
	}
	return nil
}

func (s *OSCALSource) OSCAL() *oscal.OSCAL {
	return s.oscal
}

func (s *OSCALSource) DocumentFormat() constants.DocumentFormat {
	if strings.HasSuffix(s.UserPath, ".xml") {
		return constants.XmlFormat
	} else if strings.HasSuffix(s.UserPath, ".json") {
		return constants.JsonFormat
	} else {
		return constants.UnknownFormat
	}
}

// Close the OSCALSource
func (s *OSCALSource) Close() {
	if s.file != nil {
		s.file.Close()
		s.file = nil
	}
}
