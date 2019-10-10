package shared

import (
	"encoding/json"
	"encoding/xml"
	"math"
	"net/http"
	"reflect"
)

// HTTPResponse abstract interface
type HTTPResponse interface {
	JSON(w http.ResponseWriter)
	XML(w http.ResponseWriter)
}

type (
	// Response model
	Response struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Meta    interface{} `json:"meta,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Errors  interface{} `json:"errors,omitempty"`
	}

	// Meta model
	Meta struct {
		Page         int `json:"page"`
		Limit        int `json:"limit"`
		TotalRecords int `json:"totalRecords"`
		TotalPages   int `json:"totalPages"`
	}
)

// NewHTTPResponse for create common response, data must in first params and meta in second params
func NewHTTPResponse(code int, message string, params ...interface{}) HTTPResponse {
	commonResponse := new(Response)

	for _, param := range params {
		// get value param if type is pointer
		refValue := reflect.ValueOf(param)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		param = refValue.Interface()

		switch val := param.(type) {
		case Meta:
			commonResponse.Meta = val
		case MultiError:
			commonResponse.Errors = val.ToMap()
		default:
			commonResponse.Data = param
		}
	}

	if code < http.StatusBadRequest && message != ErrorDataNotFound {
		commonResponse.Success = true
	} else {
		commonResponse.Success = false
	}
	commonResponse.Code = code
	commonResponse.Message = message
	return commonResponse
}

// JSON for set http JSON response (Content-Type: application/json) with parameter is http response writer
func (resp *Response) JSON(w http.ResponseWriter) {
	if resp.Data == nil {
		resp.Data = struct{}{}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)

	result, err := json.Marshal(&resp)
	if err != nil {
		panic(err)
	}

	w.Write(result)
}

// XML for set http XML response (Content-Type: application/xml)
func (resp *Response) XML(w http.ResponseWriter) {
	if resp.Data == nil {
		resp.Data = struct{}{}
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(resp.Code)

	result, err := xml.Marshal(&resp)
	if err != nil {
		panic(err)
	}

	w.Write(result)
}

// CreateMeta method to create meta response
func CreateMeta(total int, page int, limit int) Meta {
	var totalPages float64

	totalPages = 1

	if total > 0 && limit > 0 {
		totalPages = math.Ceil(float64(total) / float64(limit))
	}

	meta := Meta{
		Limit:        limit,
		Page:         page,
		TotalPages:   int(totalPages),
		TotalRecords: total,
	}

	return meta
}
