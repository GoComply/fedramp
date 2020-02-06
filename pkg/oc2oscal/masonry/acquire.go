package masonry

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/opencontrol/compliance-masonry/pkg/cli/get/resources"
	"github.com/opencontrol/compliance-masonry/pkg/lib"
	"github.com/opencontrol/compliance-masonry/pkg/lib/common"
	"github.com/opencontrol/compliance-masonry/pkg/lib/opencontrol"
	"github.com/opencontrol/compliance-masonry/pkg/lib/opencontrol/versions/1.0.0"
)

func Open(uri string) (common.Workspace, error) {
	tempDir, err := ioutil.TempDir("/tmp", "oscal-masonry")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	repo := make([]common.RemoteSource, 1)
	repo[0] = schema.VCSEntry{
		URL:      uri,
		Revision: "master",
		Path:     ""}
	getter := resources.NewVCSAndLocalGetter(opencontrol.YAMLParser{})
	err = getter.GetRemoteResources(tempDir, "opencontrols", repo)
	if err != nil {
		return nil, err
	}
	workspace, errors := lib.LoadData(tempDir, tempDir+"/certifications/fedramp-high.yaml")
	if errors != nil {
		return nil, fmt.Errorf("%v", errors)
	}
	return workspace, nil
}
