package hubspot

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Person struct {
	ID      int64 `hubspot:"id"`
	Deleted bool  `hubspot:"deleted"`
	Name    string
	EMail   string
	Age     int `hubspot:"name=humanage"`
}

func TestCreate(t *testing.T) {
	rest := &TestRest{}
	rest.Response = map[string]interface{}{
		"vid": 61574}

	contacts := NewContacts(rest, NewModel(reflect.TypeOf(Person{})))

	person := &Person{
		Name:  "Peter",
		EMail: "peter@lack.de",
		Age:   28}

	id, err := contacts.CreateOrUpdate("peter@lack.de", person)
	require.NoError(t, err)
	require.Equal(t, "contacts/v1/contact/createOrUpdate/email/peter@lack.de?hapikey=xyz", rest.LastRequest())
	require.Equal(t, `{"properties":[{"property":"name","value":"Peter"},{"property":"email","value":"peter@lack.de"},{"property":"humanage","value":28}]}`+"\n", rest.LastBody())
	require.Equal(t, int64(61574), id)
}
