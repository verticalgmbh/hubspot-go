package hubspot

import "testing"

func TestAssInterfaceImpl(t *testing.T) {
	var ass IAssociations = &Associations{}

	if ass != nil {
		return
	}
}
