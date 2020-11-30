package uuid

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Uuindexable interface {
	SetUuid(uuid string)
}

// Generate function returns UUID (rfc 4122) compliant string suitable
// for use within OSCAL documents. OSCAL requires random based (version 4)
// UUID to be used. However it is clear that having stable UUIDs (that is
// same id for same content) is highly beneficial. This method creates
// version 4 UUIDs that are not truly random but rather content dependent.
func Refresh(element Uuindexable) error {
	element.SetUuid("")
	signature, err := marshal(element)
	if err != nil {
		return fmt.Errorf("Cannot re-calculate UUID for %v: %s", element, err)
	}
	id := uuid.NewHash(sha1.New(), uuid.Nil, signature, 4)
	element.SetUuid(id.String())
	return nil
}

func marshal(data interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(data)
	return buffer.Bytes(), err

}
