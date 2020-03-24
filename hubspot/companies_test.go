package hubspot

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var responseCompanyList string = `{
	"companies": [
	  {
		"portalId": 62515,
		"additionalDomains": [
		  
		],
		"properties": {
		  "website": {
			"sourceId": null,
			"timestamp": 1457513066540,
			"versions": [
			  {
				"timestamp": 1457513066540,
				"sourceVid": [
				  
				],
				"name": "website",
				"value": "example.com",
				"source": "COMPANIES"
			  }
			],
			"value": "example.com",
			"source": "COMPANIES"
		  },
		  "name": {
			"sourceId": "name",
			"timestamp": 1464484587592,
			"versions": [
			  {
				"name": "name",
				"sourceId": "name",
				"timestamp": 1464484587592,
				"value": "Example Company",
				"source": "BIDEN",
				"sourceVid": [
				  
				]
			  }
			],
			"value": "Example Company",
			"source": "BIDEN"
		  }
		},
		"isDeleted": false,
		"companyId": 115200636
	  },
	  {
		"portalId": 62515,
		"additionalDomains": [
		  
		],
		"properties": {
		  "website": {
			"sourceId": null,
			"timestamp": 1457535205549,
			"versions": [
			  {
				"timestamp": 1457535205549,
				"sourceVid": [
				  
				],
				"name": "website",
				"value": "test.com",
				"source": "COMPANIES"
			  }
			],
			"value": "test.com",
			"source": "COMPANIES"
		  },
		  "name": {
			"sourceId": "name",
			"timestamp": 1468832771769,
			"versions": [
			  {
				"name": "name",
				"sourceId": "name",
				"timestamp": 1468832771769,
				"value": "Test Company",
				"source": "BIDEN",
				"sourceVid": [
				  
				]
			  }
			],
			"value": "Test Company",
			"source": "BIDEN"
		  }
		},
		"isDeleted": false,
		"companyId": 115279791
	  }
	],
	"has-more": true,
	"offset": 115279791
  }`

type Company struct {
	ID      int64 `hubspot:"id"`
	Deleted bool  `hubspot:"deleted"`
	Name    string
	Website string
	VAT     string `hubspot:"name=umsatzsteuerid"`
}

func TestCompanyInterfaceImpl(t *testing.T) {
	var companies ICompanies = NewCompanies(&TestRest{}, NewModel(reflect.TypeOf(Company{})))

	// just a noop to make sure go compiler doesn't babble about the var not being used
	if companies != nil {
		return
	}
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

func TestCompanyListDefault(t *testing.T) {
	rest := &TestRest{Response: readTestResponse(responseCompanyList)}

	companies := NewCompanies(rest, NewModel(reflect.TypeOf(Company{})))
	pageresponse, err := companies.List(nil, "name")
	require.NoError(t, err)

	require.Equal(t, "GET companies/v2/companies/paged?hapikey=xyz&properties=name", rest.LastRequest())
	require.True(t, pageresponse.HasMore)
	require.Equal(t, int64(115279791), pageresponse.Offset)
	require.Equal(t, 2, len(pageresponse.Data))

	company, ok := pageresponse.Data[0].(*Company)
	require.True(t, ok)
	require.Equal(t, "Example Company", company.Name)
	company, ok = pageresponse.Data[1].(*Company)
	require.True(t, ok)
	require.Equal(t, "Test Company", company.Name)
}
