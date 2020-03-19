package hubspot

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Company struct {
	ID      int64 `hubspot:"id"`
	Deleted bool  `hubspot:"deleted"`
	Name    string
	Website string
	VAT     string `hubspot:"name=umsatzsteuerid"`
}

func TestCreateCompany(t *testing.T) {
	rest := &TestRest{}
	rest.Response = map[string]interface{}{
		"companyId": 4886502,
		"properties": map[string]interface{}{
			"name":           map[string]interface{}{"value": "vertical GmbH"},
			"website":        map[string]interface{}{"value": "www.vertical.de"},
			"umsatzsteuerid": map[string]interface{}{"value": "123noclue"}}}

	companies := NewCompanies(rest, NewModel(reflect.TypeOf(Company{})))

	company := &Company{
		Name:    "vertical GmbH",
		Website: "www.vertical.de",
		VAT:     "123noclue"}

	created, err := companies.Create(company)
	createdcompany := created.(*Company)
	require.NoError(t, err)
	require.Equal(t, "POST companies/v2/companies?hapikey=xyz", rest.LastRequest())

	request := extractRequestProperties("name", rest.LastBody())
	require.Equal(t, "vertical GmbH", request["name"])
	require.Equal(t, "www.vertical.de", request["website"])
	require.Equal(t, "123noclue", request["umsatzsteuerid"])

	require.Equal(t, int64(4886502), createdcompany.ID)
	require.Equal(t, "vertical GmbH", createdcompany.Name)
	require.Equal(t, "www.vertical.de", createdcompany.Website)
	require.Equal(t, "123noclue", createdcompany.VAT)
}
