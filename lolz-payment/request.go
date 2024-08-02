package lolzpayment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

type params map[string]interface{}

type request struct {
	method   string
	endpoint string
	query    url.Values
	header   http.Header
	fullURL  string
}

func (r *request) validate() (err error) {
	if r.query == nil {
		r.query = url.Values{}
	}

	return nil
}

// addParam add param with key/value to query string
func (r *request) addParam(key string, value interface{}) *request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Add(key, fmt.Sprintf("%v", value))
	return r
}

// setParam set param with key/value to query string
func (r *request) setParam(key string, value interface{}) *request {
	if r.query == nil {
		r.query = url.Values{}
	}

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		v, err := json.Marshal(value)
		if err == nil {
			value = string(v)
		}
	}

	r.query.Set(key, fmt.Sprintf("%v", value))
	return r
}

// setParams set params with key/values to query string
func (r *request) setParams(m params) *request {
	for k, v := range m {
		r.setParam(k, v)
	}
	return r
}
