package hubspot

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func getBodyValue(body interface{}, path string) interface{} {
	current := body
	lastitem := "<root>"

	pathsplit := strings.Split(path, "/")
	for _, itemname := range pathsplit {
		item, ok := current.(map[string]interface{})
		if !ok {
			return errors.Errorf("Item '%s' not an object map", lastitem)
		}

		lastitem = itemname
		current = item
	}

	return current
}

type TestRest struct {
	requests []string
	bodies   []interface{}
	Response map[string]interface{}
}

func (rest *TestRest) LastRequest() string {
	if len(rest.requests) == 0 {
		return ""
	}

	return rest.requests[len(rest.requests)-1]
}

func (rest *TestRest) LastBody() interface{} {
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
		rest.bodies = append(rest.bodies, request)
	} else {
		rest.bodies = append(rest.bodies, nil)
	}
}

func (rest *TestRest) Post(url string, request interface{}, params ...*Parameter) (map[string]interface{}, error) {
	rest.log("POST "+url, request, params...)
	return rest.Response, nil
}

func (rest *TestRest) Put(url string, request interface{}, params ...*Parameter) (map[string]interface{}, error) {
	rest.log("PUT "+url, request, params...)
	return rest.Response, nil
}
func (rest *TestRest) Delete(url string) error {
	rest.log("DELETE "+url, nil)
	return nil
}

func (rest *TestRest) Get(url string, params ...*Parameter) (map[string]interface{}, error) {
	rest.log("GET "+url, nil, params...)
	return rest.Response, nil
}
