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
		"vid": 61574,
		"properties": map[string]interface{}{
			"name":     map[string]interface{}{"value": "Peter"},
			"email":    map[string]interface{}{"value": "peter@lack.de"},
			"humanage": map[string]interface{}{"value": 28}}}

	contacts := NewContacts(rest, NewModel(reflect.TypeOf(Person{})))

	person := &Person{
		Name:  "Peter",
		EMail: "peter@lack.de",
		Age:   28}

	created, err := contacts.CreateOrUpdate("peter@lack.de", person)
	createdperson := created.(*Person)
	require.NoError(t, err)
	require.Equal(t, "POST contacts/v1/contact/createOrUpdate/email/peter@lack.de?hapikey=xyz", rest.LastRequest())

	request := extractRequestProperties(rest.LastBody())
	require.Equal(t, "Peter", request["name"])
	require.Equal(t, "peter@lack.de", request["email"])
	require.Equal(t, 28, request["humanage"])

	require.Equal(t, int64(61574), createdperson.ID)
	require.Equal(t, "Peter", createdperson.Name)
	require.Equal(t, "peter@lack.de", createdperson.EMail)
	require.Equal(t, 28, createdperson.Age)
}

func TestUpdate(t *testing.T) {
	rest := &TestRest{}

	contacts := NewContacts(rest, NewModel(reflect.TypeOf(Person{})))

	person := &Person{
		ID:    61574,
		Name:  "Peter",
		EMail: "peter@lack.de",
		Age:   28}

	err := contacts.Update(61574, person)
	require.NoError(t, err)
	require.Equal(t, "POST contacts/v1/contact/vid/61574/profile?hapikey=xyz", rest.LastRequest())

	request := extractRequestProperties(rest.LastBody())
	require.Equal(t, "Peter", request["name"])
	require.Equal(t, "peter@lack.de", request["email"])
	require.Equal(t, 28, request["humanage"])
}

func TestDelete(t *testing.T) {
	rest := &TestRest{}

	contacts := NewContacts(rest, NewModel(reflect.TypeOf(Person{})))

	err := contacts.Delete(61574)

	require.NoError(t, err)
	require.Equal(t, "DELETE contacts/v1/contact/vid/61574?hapikey=xyz", rest.LastRequest())
}

func TestGetByID(t *testing.T) {
	rest := &TestRest{}
	rest.Response = map[string]interface{}{
		"vid": 61574,
		"properties": map[string]interface{}{
			"name":     map[string]interface{}{"value": "Peter"},
			"email":    map[string]interface{}{"value": "peter@lack.de"},
			"humanage": map[string]interface{}{"value": 28}}}

	contacts := NewContacts(rest, NewModel(reflect.TypeOf(Person{})))

	data, err := contacts.GetByID(61574)
	person := data.(*Person)
	require.NoError(t, err)
	require.Equal(t, "GET contacts/v1/contact/vid/61574/profile?hapikey=xyz", rest.LastRequest())

	require.Equal(t, int64(61574), person.ID)
	require.Equal(t, "Peter", person.Name)
	require.Equal(t, "peter@lack.de", person.EMail)
	require.Equal(t, 28, person.Age)
}

func TestGetByEmail(t *testing.T) {
	rest := &TestRest{}
	rest.Response = map[string]interface{}{
		"vid": 61574,
		"properties": map[string]interface{}{
			"name":     map[string]interface{}{"value": "Peter"},
			"email":    map[string]interface{}{"value": "peter@lack.de"},
			"humanage": map[string]interface{}{"value": 28}}}

	contacts := NewContacts(rest, NewModel(reflect.TypeOf(Person{})))

	data, err := contacts.GetByEmail("peter@lack.de")
	person := data.(*Person)
	require.NoError(t, err)
	require.Equal(t, "GET contacts/v1/contact/email/peter@lack.de/profile?hapikey=xyz", rest.LastRequest())

	require.Equal(t, int64(61574), person.ID)
	require.Equal(t, "Peter", person.Name)
	require.Equal(t, "peter@lack.de", person.EMail)
	require.Equal(t, 28, person.Age)
}

func TestListDefaultPage(t *testing.T) {
	rest := &TestRest{}
	rest.Response = map[string]interface{}{
		"has-more":   true,
		"vid-offset": 12345,
		"contacts": []map[string]interface{}{
			map[string]interface{}{
				"vid": 61574,
				"properties": map[string]interface{}{
					"name":     map[string]interface{}{"value": "Peter"},
					"email":    map[string]interface{}{"value": "peter@lack.de"},
					"humanage": map[string]interface{}{"value": 28}}},
			map[string]interface{}{
				"vid": 51157,
				"properties": map[string]interface{}{
					"name":     map[string]interface{}{"value": "Monika"},
					"email":    map[string]interface{}{"value": "monika@left.de"},
					"humanage": map[string]interface{}{"value": 24}}}}}

	contacts := NewContacts(rest, NewModel(reflect.TypeOf(Person{})))

	page, err := contacts.ListPage(nil, "name", "email", "humanage")

	require.NoError(t, err)
	require.Equal(t, "GET contacts/v1/lists/all/contacts/all?hapikey=xyz&property=name&property=email&property=humanage", rest.LastRequest())

	require.True(t, page.HasMore)
	require.Equal(t, int64(12345), page.Offset)
	require.Equal(t, 2, len(page.Data))

	person := page.Data[0].(*Person)
	require.Equal(t, int64(61574), person.ID)
	require.Equal(t, "Peter", person.Name)
	require.Equal(t, "peter@lack.de", person.EMail)
	require.Equal(t, 28, person.Age)

	person = page.Data[1].(*Person)
	require.Equal(t, int64(51157), person.ID)
	require.Equal(t, "Monika", person.Name)
	require.Equal(t, "monika@left.de", person.EMail)
	require.Equal(t, 24, person.Age)
}

func TestListCustomPage(t *testing.T) {
	rest := &TestRest{}
	rest.Response = map[string]interface{}{
		"has-more": false,
		"contacts": []map[string]interface{}{
			map[string]interface{}{
				"vid": 61574,
				"properties": map[string]interface{}{
					"name":     map[string]interface{}{"value": "Peter"},
					"email":    map[string]interface{}{"value": "peter@lack.de"},
					"humanage": map[string]interface{}{"value": 28}}},
			map[string]interface{}{
				"vid": 51157,
				"properties": map[string]interface{}{
					"name":     map[string]interface{}{"value": "Monika"},
					"email":    map[string]interface{}{"value": "monika@left.de"},
					"humanage": map[string]interface{}{"value": 24}}}}}

	contacts := NewContacts(rest, NewModel(reflect.TypeOf(Person{})))

	page, err := contacts.ListPage(NewPage(int64(544), 54), "name", "email", "humanage")

	require.NoError(t, err)
	require.Equal(t, "GET contacts/v1/lists/all/contacts/all?hapikey=xyz&count=54&vidOffset=544&property=name&property=email&property=humanage", rest.LastRequest())

	require.False(t, page.HasMore)
	require.Equal(t, int64(0), page.Offset)
	require.Equal(t, 2, len(page.Data))

	person := page.Data[0].(*Person)
	require.Equal(t, int64(61574), person.ID)
	require.Equal(t, "Peter", person.Name)
	require.Equal(t, "peter@lack.de", person.EMail)
	require.Equal(t, 28, person.Age)

	person = page.Data[1].(*Person)
	require.Equal(t, int64(51157), person.ID)
	require.Equal(t, "Monika", person.Name)
	require.Equal(t, "monika@left.de", person.EMail)
	require.Equal(t, 24, person.Age)
}
