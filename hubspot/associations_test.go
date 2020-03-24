package hubspot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var responseAssociationsList string = `{
	"results": [
	  184896670
	],
	"hasMore": false,
	"offset": 184896670
  }`

func TestAssInterfaceImpl(t *testing.T) {
	var ass IAssociations = &Associations{}

	if ass != nil {
		return
	}
}

func TestAssCreate(t *testing.T) {
	rest := &TestRest{}
	api := NewAssociations(rest)
	err := api.Create(13, 29, AssociationContactToCompany)

	require.NoError(t, err)
	require.Equal(t, "PUT crm-associations/v1/associations?hapikey=xyz", rest.LastRequest())
	request, ok := rest.LastBody().(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, int64(13), request["fromObjectId"])
	require.Equal(t, int64(29), request["toObjectId"])
	require.Equal(t, "HUBSPOT_DEFINED", request["category"])
	require.Equal(t, 1, request["definitionId"])
}

func TestAssCreateBulk(t *testing.T) {
	rest := &TestRest{}
	api := NewAssociations(rest)
	err := api.CreateBulk([]*Association{
		&Association{
			From: 13,
			To:   29,
			Type: AssociationContactToCompany},
		&Association{
			From: 92,
			To:   1221,
			Type: AssociationCompanyToDeal}})

	require.NoError(t, err)
	require.Equal(t, "PUT crm-associations/v1/associations/create-batch?hapikey=xyz", rest.LastRequest())

	request, ok := rest.LastBody().([]map[string]interface{})
	require.True(t, ok)
	require.Equal(t, 2, len(request))

	require.Equal(t, int64(13), request[0]["fromObjectId"])
	require.Equal(t, int64(29), request[0]["toObjectId"])
	require.Equal(t, "HUBSPOT_DEFINED", request[0]["category"])
	require.Equal(t, 1, request[0]["definitionId"])

	require.Equal(t, int64(92), request[1]["fromObjectId"])
	require.Equal(t, int64(1221), request[1]["toObjectId"])
	require.Equal(t, "HUBSPOT_DEFINED", request[1]["category"])
	require.Equal(t, 6, request[1]["definitionId"])
}

func TestAssListDefault(t *testing.T) {
	rest := &TestRest{
		Response: readTestResponse(responseAssociationsList)}
	api := NewAssociations(rest)
	response, err := api.List(2883, AssociationCompanyToTicket, nil)
	require.NoError(t, err)
	require.Equal(t, "GET crm-associations/v1/associations/2883/HUBSPOT_DEFINED/25?hapikey=xyz", rest.LastRequest())

	require.False(t, response.HasMore)
	require.Equal(t, int64(0), response.Offset)
	require.Equal(t, 1, len(response.Data))
	require.Equal(t, int64(184896670), response.Data[0])
}

func TestAssListPage(t *testing.T) {
	rest := &TestRest{}
	api := NewAssociations(rest)
	_, err := api.List(2883, AssociationCompanyToTicket, NewPage(230, 0))
	require.NoError(t, err)
	require.Equal(t, "GET crm-associations/v1/associations/2883/HUBSPOT_DEFINED/25?hapikey=xyz&offset=230", rest.LastRequest())
}

func TestAssListPageLimit(t *testing.T) {
	rest := &TestRest{}
	api := NewAssociations(rest)
	_, err := api.List(2883, AssociationCompanyToTicket, NewPage(233, 40))
	require.NoError(t, err)
	require.Equal(t, "GET crm-associations/v1/associations/2883/HUBSPOT_DEFINED/25?hapikey=xyz&offset=233&limit=40", rest.LastRequest())
}

func TestAssDelete(t *testing.T) {
	rest := &TestRest{}
	api := NewAssociations(rest)
	err := api.Delete(50, 28, AssociationAdvisorToCompany)
	require.NoError(t, err)
	require.Equal(t, "PUT crm-associations/v1/associations/delete?hapikey=xyz", rest.LastRequest())

	response, ok := rest.LastBody().(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, int64(50), response["fromObjectId"])
	require.Equal(t, int64(28), response["toObjectId"])
	require.Equal(t, "HUBSPOT_DEFINED", response["category"])
	require.Equal(t, 33, response["definitionId"])
}

func TestAssDeleteBulk(t *testing.T) {
	rest := &TestRest{}
	api := NewAssociations(rest)
	err := api.DeleteBulk([]*Association{
		&Association{
			From: 13,
			To:   29,
			Type: AssociationContactToCompany},
		&Association{
			From: 92,
			To:   1221,
			Type: AssociationCompanyToDeal}})

	require.NoError(t, err)
	require.Equal(t, "PUT crm-associations/v1/associations/delete-batch?hapikey=xyz", rest.LastRequest())

	response, ok := rest.LastBody().([]map[string]interface{})
	require.True(t, ok)

	require.Equal(t, int64(13), response[0]["fromObjectId"])
	require.Equal(t, int64(29), response[0]["toObjectId"])
	require.Equal(t, "HUBSPOT_DEFINED", response[0]["category"])
	require.Equal(t, 1, response[0]["definitionId"])

	require.Equal(t, int64(92), response[1]["fromObjectId"])
	require.Equal(t, int64(1221), response[1]["toObjectId"])
	require.Equal(t, "HUBSPOT_DEFINED", response[1]["category"])
	require.Equal(t, 6, response[1]["definitionId"])
}
