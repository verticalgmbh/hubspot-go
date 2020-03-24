package hubspot

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var responseDealCreate string = `{
	"portalId": 62515,
	"dealId": 151088,
	"isDeleted": false,
	"associations": {
	  "associatedVids": [
		27136
	  ],
	  "associatedCompanyIds": [
		8954037
	  ],
	  "associatedDealIds": [
		
	  ]
	},
	"properties": {
	  "amount": {
		"value": "60000",
		"timestamp": 1410381338943,
		"source": "API",
		"sourceId": null,
		"versions": [
		  {
			"name": "amount",
			"value": "60000",
			"timestamp": 1410381338943,
			"source": "API",
			"sourceVid": [
			  
			]
		  }
		]
	  },
	  "dealstage": {
		"value": "appointmentscheduled",
		"timestamp": 1410381338943,
		"source": "API",
		"sourceId": null,
		"versions": [
		  {
			"name": "dealstage",
			"value": "appointmentscheduled",
			"timestamp": 1410381338943,
			"source": "API",
			"sourceVid": [
			  
			]
		  }
		]
	  },
	  "pipeline": {
		"value": "default",
		"timestamp": 1410381338943,
		"source": "API",
		"sourceId": null,
		"versions": [
		  {
			"name": "pipeline",
			"value": "default",
			"timestamp": 1410381338943,
			"source": "API",
			"sourceVid": [
			  
			]
		  }
		]
	  },
	  "closedate": {
		"value": "1409443200000",
		"timestamp": 1410381338943,
		"source": "API",
		"sourceId": null,
		"versions": [
		  {
			"name": "closedate",
			"value": "1409443200000",
			"timestamp": 1410381338943,
			"source": "API",
			"sourceVid": [
			  
			]
		  }
		]
	  },
	  "createdate": {
		"value": "1410381339020",
		"timestamp": 1410381339020,
		"source": null,
		"sourceId": null,
		"versions": [
		  {
			"name": "createdate",
			"value": "1410381339020",
			"timestamp": 1410381339020,
			"sourceVid": [
			  
			]
		  }
		]
	  },
	  "hubspot_owner_id": {
		"value": "24",
		"timestamp": 1410381338943,
		"source": "API",
		"sourceId": null,
		"versions": [
		  {
			"name": "hubspot_owner_id",
			"value": "24",
			"timestamp": 1410381338943,
			"source": "API",
			"sourceVid": [
			  
			]
		  }
		]
	  },
	  "hs_createdate": {
		"value": "1410381339020",
		"timestamp": 1410381339020,
		"source": null,
		"sourceId": null,
		"versions": [
		  {
			"name": "hs_createdate",
			"value": "1410381339020",
			"timestamp": 1410381339020,
			"sourceVid": [
			  
			]
		  }
		]
	  },
	  "dealtype": {
		"value": "newbusiness",
		"timestamp": 1410381338943,
		"source": "API",
		"sourceId": null,
		"versions": [
		  {
			"name": "dealtype",
			"value": "newbusiness",
			"timestamp": 1410381338943,
			"source": "API",
			"sourceVid": [
			  
			]
		  }
		]
	  },
	  "dealname": {
		"value": "A new Deal",
		"timestamp": 1410381338943,
		"source": "API",
		"sourceId": null,
		"versions": [
		  {
			"name": "dealname",
			"value": "A new Deal",
			"timestamp": 1410381338943,
			"source": "API",
			"sourceVid": [
			  
			]
		  }
		]
	  }
	}
  }`

type Deal struct {
	ID        int64     `hubspot:"id"`
	IsDeleted bool      `hubspot:"deleted"`
	Contacts  []int64   `hubspot:"contacts"`
	Companies []int64   `hubspot:"companies"`
	Name      string    `hubspot:"name=dealname"`
	Stage     string    `hubspot:"name=dealstage"`
	CloseDate time.Time `hubspot:"name=closedate"`
}

func TestDealCreate(t *testing.T) {
	rest := &TestRest{Response: readTestResponse(responseDealCreate)}
	api := NewDeals(rest, NewModel(reflect.TypeOf(Deal{})))
	dealresponse, err := api.Create(&Deal{
		Name:      "TestDeal",
		Stage:     "closedwon",
		CloseDate: time.Date(2020, 3, 24, 00, 00, 00, 00, time.Local),
	})
	require.NoError(t, err)
	require.Equal(t, "POST deals/v1/deal?hapikey=xyz", rest.LastRequest())

	deal, ok := dealresponse.(*Deal)
	require.True(t, ok)
	require.Equal(t, int64(151088), deal.ID)
	require.False(t, deal.IsDeleted)
	require.Equal(t, "A new Deal", deal.Name)
	require.Equal(t, "appointmentscheduled", deal.Stage)
	require.Equal(t, 2014, deal.CloseDate.Year())
	require.Equal(t, time.Month(8), deal.CloseDate.Month())
	require.Equal(t, 31, deal.CloseDate.Day())
	require.Equal(t, 0, deal.CloseDate.Hour())
	require.Equal(t, 0, deal.CloseDate.Minute())
	require.Equal(t, 0, deal.CloseDate.Second())

	require.Equal(t, 1, len(deal.Contacts))
}
