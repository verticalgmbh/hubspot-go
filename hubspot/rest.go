package hubspot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

var httpclient http.Client

// Parameter - parameter for a rest client
type Parameter struct {
	Key   string
	Value string
}

// IRestClient - interface for a client sending rest requests to hubspot
type IRestClient interface {
	Post(url string, request interface{}, params ...*Parameter) (map[string]interface{}, error)
	Put(url string, request interface{}, params ...*Parameter) (map[string]interface{}, error)
	Delete(url string) error
	Get(url string, params ...*Parameter) (map[string]interface{}, error)
}

// RestClient - client used to send rest requests to hubspot
type RestClient struct {
	apikey  string
	address string
}

// NewParameter - creates a new parameter
func NewParameter(key string, value string) *Parameter {
	return &Parameter{
		Key:   key,
		Value: value}
}

func (client *RestClient) buildBaseURL(address string, params ...*Parameter) *strings.Builder {
	var builder strings.Builder
	builder.WriteString(client.address)
	builder.WriteString(address)
	builder.WriteString("?hapikey=")
	builder.WriteString(client.apikey)

	if len(params) > 0 {
		for _, param := range params {
			builder.WriteRune('&')
			builder.WriteString(param.Key)
			builder.WriteRune('=')
			builder.WriteString(url.QueryEscape(param.Value))
		}
	}
	return &builder
}

func (client *RestClient) checkError(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if response.ContentLength == 0 {
			return errors.Errorf("%d: %s", response.StatusCode, response.Status)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Request.Body)
		return errors.Errorf("%s", buf.String())
	}

	return nil
}

func (client *RestClient) readResponse(response *http.Response) (map[string]interface{}, error) {
	if response.ContentLength > 0 {
		jresponse := make(map[string]interface{})
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&jresponse)
		if err != nil {
			return nil, err
		}

		return jresponse, nil
	}

	return nil, errors.New("No response body")
}

// Get - send a GET request to hubspot
func (client *RestClient) Get(address string, params ...*Parameter) (map[string]interface{}, error) {
	builder := client.buildBaseURL(address, params...)

	response, err := http.Get(builder.String())
	if err != nil {
		return nil, err
	}

	err = client.checkError(response)
	if err != nil {
		return nil, err
	}

	return client.readResponse(response)
}

// Post - send a POST request to hubspot
func (client *RestClient) Post(address string, request interface{}, params ...*Parameter) (map[string]interface{}, error) {
	builder := client.buildBaseURL(address, params...)

	buffer := new(bytes.Buffer)

	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(request)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(builder.String(), "application/json", buffer)
	if err != nil {
		return nil, err
	}

	err = client.checkError(response)
	if err != nil {
		return nil, err
	}

	return client.readResponse(response)
}

// Put - send a PUT request to hubspot
func (client *RestClient) Put(address string, request interface{}, params ...*Parameter) (map[string]interface{}, error) {
	builder := client.buildBaseURL(address, params...)

	buffer := new(bytes.Buffer)

	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(request)
	if err != nil {
		return nil, err
	}

	putrequest, err := http.NewRequest("PUT", builder.String(), buffer)
	if err != nil {
		return nil, err
	}

	response, err := httpclient.Do(putrequest)
	if err != nil {
		return nil, err
	}

	err = client.checkError(response)
	if err != nil {
		return nil, err
	}

	return client.readResponse(response)
}

// Delete - send a DELETE request to hubspot
func (client *RestClient) Delete(address string) error {
	builder := client.buildBaseURL(address)

	request, err := http.NewRequest("DELETE", builder.String(), nil)
	if err != nil {
		return err
	}

	response, err := httpclient.Do(request)
	if err != nil {
		return err
	}

	return client.checkError(response)
}
