package hubspot

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var responseTicketCreate string = `{
	"objectType": "TICKET",
	"portalId": 62515,
	"objectId": 177769,
	"properties": {
	  "hs_lastmodifieddate": {
		"versions": [
		  {
			"name": "hs_lastmodifieddate",
			"value": "0",
			"timestamp": 0,
			"source": "CALCULATED",
			"sourceVid": []
		  }
		],
		"value": "0",
		"timestamp": 0,
		"source": "CALCULATED",
		"sourceId": null
	  },
	  "subject": {
		"versions": [
		  {
			"name": "subject",
			"value": "Problem hier",
			"timestamp": 1522862571921,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "Problem hier",
		"timestamp": 1522862571921,
		"source": "API",
		"sourceId": null
	  },
	  "hs_pipeline": {
		"versions": [
		  {
			"name": "hs_pipeline",
			"value": "0",
			"timestamp": 1522862571921,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "0",
		"timestamp": 1535746946950,
		"source": "API",
		"sourceId": null
	  },
	  "createdate": {
		"versions": [
		  {
			"name": "createdate",
			"value": "0",
			"timestamp": 0,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "0",
		"timestamp": 0,
		"source": "API",
		"sourceId": null
	  },
	  "content": {
		"versions": [
		  {
			"name": "content",
			"value": "Will bestellen, geht ni",
			"timestamp": 1522862571921,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "Will bestellen, geht ni",
		"timestamp": 1522862571921,
		"source": "API",
		"sourceId": null
	  },
	  "hs_pipeline_stage": {
		"versions": [
		  {
			"name": "hs_pipeline_stage",
			"value": "1",
			"timestamp": 1522862571921,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "1",
		"timestamp": 1522862571921,
		"source": "API",
		"sourceId": ""
	  },
	  "time_to_close": {
		"versions": [
		  {
			"name": "time_to_close",
			"value": "",
			"timestamp": 1535746947086,
			"source": "CALCULATED",
			"sourceVid": [],
			"sourceMetadata": ""
		  }
		],
		"value": "",
		"timestamp": 1535746947086,
		"source": "CALCULATED",
		"sourceId": null
	  }
	},
	"isDeleted": false
  }`

var responseTicketGet string = `{
	"objectType": "TICKET",
	"portalId": 62515,
	"objectId": 176602,
	"properties": {
	  "subject": {
		"versions": [
		  {
			"name": "subject",
			"value": "This is an example ticket",
			"timestamp": 1522870759430,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "This is an example ticket",
		"timestamp": 1522870759430,
		"source": "API",
		"sourceId": null
	  },
	  "created_by": {
		"versions": [
		  {
			"name": "created_by",
			"value": "496346",
			"timestamp": 1522870759430,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "496346",
		"timestamp": 1522870759430,
		"source": "API",
		"sourceId": null
	  },
	  "content": {
		"versions": [
		  {
			"name": "content",
			"value": "These are the details of the ticket.",
			"timestamp": 1522872562016,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "These are the details of the ticket.",
		"timestamp": 1522872562016,
		"source": "API",
		"sourceId": null
	  },
	  "hs_pipeline": {
		"versions": [
		  {
			"name": "hs_pipeline",
			"value": "0",
			"timestamp": 1522870759430,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "0",
		"timestamp": 1522870759430,
		"source": "API",
		"sourceId": null
	  },
	  "hs_pipeline_stage": {
		"versions": [
		  {
			"name": "hs_pipeline_stage",
			"value": "4",
			"timestamp": 1522870759430,
			"source": "API",
			"sourceVid": []
		  }
		],
		"value": "4",
		"timestamp": 1522870759430,
		"source": "API",
		"sourceId": null
	  }
	},
	"isDeleted": false
  }`

type TestTicket struct {
	ID       int64  `hubspot:"id"`
	Subject  string `hubspot:"name=subject"`
	Text     string `hubspot:"name=content"`
	Pipeline int    `hubspot:"name=hs_pipeline"`
	Stage    int    `hubspot:"name=hs_pipeline_stage"`
}

func TestTicketsInterfaceImpl(t *testing.T) {
	var tickets ITickets = &Tickets{}

	if tickets != nil {
		return
	}
}

func TestTicketCreate(t *testing.T) {
	rest := &TestRest{Response: readTestResponse(responseTicketCreate)}
	api := NewTickets(rest, NewModel(reflect.TypeOf(TestTicket{})))

	ticket := &TestTicket{
		Subject: "Problem hier",
		Text:    "Will bestellen, geht ni"}
	created, err := api.Create(ticket)

	require.NoError(t, err)
	require.Equal(t, "POST crm-objects/v1/objects/tickets?hapikey=xyz", rest.LastRequest())

	request := extractRequestProperties("name", rest.LastBody())
	require.Equal(t, "Problem hier", request["subject"])
	require.Equal(t, "Will bestellen, geht ni", request["content"])

	createdticket := created.(*TestTicket)
	require.Equal(t, ticket.Subject, createdticket.Subject)
	require.Equal(t, ticket.Text, createdticket.Text)
	require.NotEqual(t, int64(0), createdticket.ID)
}

func TestTicketGet(t *testing.T) {
	rest := &TestRest{Response: readTestResponse(responseTicketGet)}
	api := NewTickets(rest, NewModel(reflect.TypeOf(TestTicket{})))

	created, err := api.Get(176602)

	require.NoError(t, err)
	require.Equal(t, "GET crm-objects/v1/objects/tickets/176602?hapikey=xyz", rest.LastRequest())

	createdticket := created.(*TestTicket)
	require.Equal(t, "This is an example ticket", createdticket.Subject)
	require.Equal(t, "These are the details of the ticket.", createdticket.Text)
	require.Equal(t, 0, createdticket.Pipeline)
	require.Equal(t, 4, createdticket.Stage)
	require.NotEqual(t, int64(0), createdticket.ID)
}
