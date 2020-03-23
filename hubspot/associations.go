package hubspot

import (
	"fmt"

	"github.com/spf13/cast"
)

const hubspotDefinedType = "HUBSPOT_DEFINED"

// AssociationType - type of data association
type AssociationType int8

const (
	AssociationContactToCompany            AssociationType = 1
	AssociationCompanyToContact            AssociationType = 2
	AssociationDealToContact               AssociationType = 3
	AssociationContactToDeal               AssociationType = 4
	AssociationDealToCompany               AssociationType = 5
	AssociationCompanyToDeal               AssociationType = 6
	AssociationCompanyToEngagement         AssociationType = 7
	AssociationEngagementToCompany         AssociationType = 8
	AssociationContactToEngagement         AssociationType = 9
	AssociationEngagementToContact         AssociationType = 10
	AssociationDealToEngagement            AssociationType = 11
	AssociationEngagementToDeal            AssociationType = 12
	AssociationParentCompanyToChildCompany AssociationType = 13
	AssociationChildCompanyToParentCompany AssociationType = 14
	AssociationContactToTicket             AssociationType = 15
	AssociationTicketToContact             AssociationType = 16
	AssociationTicketToEngagement          AssociationType = 17
	AssociationEngagementToTicket          AssociationType = 18
	AssociationDealToLineItem              AssociationType = 19
	AssociationLineItemToDeal              AssociationType = 20
	AssociationCompanyToTicket             AssociationType = 25
	AssociationTicketToCompany             AssociationType = 26
	AssociationDealToTicket                AssociationType = 27
	AssociationTicketToDeal                AssociationType = 28
	AssociationAdvisorToCompany            AssociationType = 33
	AssociationCompanyToAdvisor            AssociationType = 34
	AssociationBoardMemberToCompany        AssociationType = 35
	AssociationCompanyToBoardMember        AssociationType = 36
	AssociationContractorToCompany         AssociationType = 37
	AssociationCompanyToContractor         AssociationType = 38
	AssociationManagerToCompany            AssociationType = 39
	AssociationCompanyToManager            AssociationType = 40
	AssociationBusinessOwnerToCompany      AssociationType = 41
	AssociationCompanyToBusinessOwner      AssociationType = 42
	AssociationPartnerToCompany            AssociationType = 43
	AssociationCompanyToPartner            AssociationType = 44
	AssociationResellerToCompany           AssociationType = 45
	AssociationCompanyToReseller           AssociationType = 46
)

// Association - association of objects in hubspot
type Association struct {
	From int64
	To   int64
	Type AssociationType
}

// IAssociations - interface for associations api
type IAssociations interface {
	Create(fromid int64, toid int64, asstype AssociationType) error
	CreateBulk(data []Association) error
	List(objectid int64, asstype AssociationType, page *Page) (*PageResponse, error)
	Delete(fromid int64, toid int64, asstype AssociationType) error
	DeleteBulk(data []Association) error
}

// Associations - hubspot associations api using rest
type Associations struct {
	rest IRestClient // client used to send requests
}

// NewAssociations - creates a new associations api
func NewAssociations(rest IRestClient) *Associations {
	return &Associations{rest: rest}
}

// Create - creates a new association of two objects in hubspot
func (api *Associations) Create(fromid int64, toid int64, asstype AssociationType) error {
	request := map[string]interface{}{
		"fromObjectId": fromid,
		"toObjectId":   toid,
		"category":     hubspotDefinedType,
		"definitionId": int(asstype)}

	_, err := api.rest.Post("crm-associations/v1/associations", request)
	return err
}

// CreateBulk - creates multiple associations in one call
func (api *Associations) CreateBulk(data []Association) error {
	var request []map[string]interface{} = make([]map[string]interface{}, len(data))
	for index, ass := range data {
		request[index] = map[string]interface{}{
			"fromObjectId": ass.From,
			"toObjectId":   ass.To,
			"category":     hubspotDefinedType,
			"definitionId": int(ass.Type)}
	}

	_, err := api.rest.Post("crm-associations/v1/associations/create-batch", request)
	return err
}

// List - lists associations of a type for an object
func (api *Associations) List(objectid int64, asstype AssociationType, page *Page) (*PageResponse, error) {
	response, err := api.rest.Get(fmt.Sprintf("crm-associations/v1/associations/%d/HUBSPOT_DEFINED/%d", objectid, int(asstype)))
	if err != nil {
		return nil, err
	}

	pr := new(PageResponse)
	pr.HasMore = cast.ToBool(response["hasMore"])
	if pr.HasMore {
		pr.Offset = cast.ToInt64(response["offset"])
	}

	contacts, ok := response["results"].([]interface{})
	if ok {
		pr.Data = contacts
	}

	return pr, nil
}

// Delete - deletes an object association
func (api *Associations) Delete(fromid int64, toid int64, asstype AssociationType) error {
	request := map[string]interface{}{
		"fromObjectId": fromid,
		"toObjectId":   toid,
		"category":     hubspotDefinedType,
		"definitionId": int(asstype)}

	_, err := api.rest.Put("crm-associations/v1/associations/delete", request)
	return err
}

// DeleteBulk - deletes multiple associations at once
func (api *Associations) DeleteBulk(data []Association) error {
	var request []map[string]interface{} = make([]map[string]interface{}, len(data))
	for index, ass := range data {
		request[index] = map[string]interface{}{
			"fromObjectId": ass.From,
			"toObjectId":   ass.To,
			"category":     hubspotDefinedType,
			"definitionId": int(ass.Type)}
	}

	_, err := api.rest.Put("crm-associations/v1/associations/delete-batch", request)
	return err
}
