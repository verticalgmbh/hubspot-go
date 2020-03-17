package hubspot

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strings"
)

type TestRest struct {
	requests []string
	bodies   []string
	Response map[string]interface{}
}

func (rest *TestRest) LastRequest() string {
	if len(rest.requests) == 0 {
		return ""
	}

	return rest.requests[len(rest.requests)-1]
}

func (rest *TestRest) LastBody() string {
	if len(rest.bodies) == 0 {
		return ""
	}

	return rest.bodies[len(rest.bodies)-1]
}

func (rest *TestRest) buildBaseURL(address string, params ...*Parameter) *strings.Builder {
	var builder strings.Builder
	builder.WriteString(address)
	builder.WriteString("?hapikey=xyz")

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

func (rest *TestRest) log(url string, request interface{}, params ...*Parameter) {
	builder := rest.buildBaseURL(url, params...)
	rest.requests = append(rest.requests, builder.String())

	if request != nil {
		buffer := new(bytes.Buffer)

		encoder := json.NewEncoder(buffer)
		encoder.Encode(request)
		rest.bodies = append(rest.bodies, buffer.String())
	} else {
		rest.bodies = append(rest.bodies, "")
	}
}

func (rest *TestRest) Post(url string, request interface{}, params ...*Parameter) (map[string]interface{}, error) {
	rest.log(url, request, params...)
	return rest.Response, nil
}

func (rest *TestRest) Put(url string, request interface{}, params ...*Parameter) (map[string]interface{}, error) {
	rest.log(url, request, params...)
	return rest.Response, nil
}
func (rest *TestRest) Delete(url string) error {
	rest.log(url, nil)
	return nil
}

func (rest *TestRest) Get(url string, params ...*Parameter) (map[string]interface{}, error) {
	rest.log(url, nil, params...)
	return rest.Response, nil
}
