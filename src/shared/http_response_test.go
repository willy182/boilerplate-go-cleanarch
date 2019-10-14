package shared

import (
	"fmt"
	"math"
	"net/http"
	"reflect"
	"testing"
)

// writer implement http.ResponseWriter
type writer struct {
	http.ResponseWriter
}

func (w *writer) Header() http.Header {
	return http.Header{}
}

func (w *writer) WriteHeader(code int) {
}

func (w *writer) Write(b []byte) (int, error) {
	return len(b), nil
}

type ExampleModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestNewHTTPResponse(t *testing.T) {
	multiError := NewMultiError()
	multiError.Append("test", fmt.Errorf("error test"))
	type args struct {
		code    int
		message string
		params  []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "Testcase #1: Response data list (include meta)",
			args: args{
				code:    http.StatusOK,
				message: "Fetch all data",
				params: []interface{}{
					[]ExampleModel{{ID: 1, Name: "wahyu"}, {ID: 2, Name: "kandar"}},
					Meta{Page: 1, Limit: 10, TotalPages: 1, TotalRecords: 2},
				},
			},
			want: &Response{
				Success: true,
				Code:    200,
				Message: "Fetch all data",
				Meta:    Meta{Page: 1, Limit: 10, TotalPages: 1, TotalRecords: 2},
				Data:    []ExampleModel{{ID: 1, Name: "wahyu"}, {ID: 2, Name: "kandar"}},
			},
		},
		{
			name: "Testcase #2: Response data detail",
			args: args{
				code:    http.StatusOK,
				message: "Get detail data",
				params: []interface{}{
					ExampleModel{ID: 1, Name: "wahyu"},
				},
			},
			want: &Response{
				Success: true,
				Code:    200,
				Message: "Get detail data",
				Data:    ExampleModel{ID: 1, Name: "wahyu"},
			},
		},
		{
			name: "Testcase #3: Response only message (without data)",
			args: args{
				code:    http.StatusOK,
				message: "list data empty",
			},
			want: &Response{
				Success: true,
				Code:    200,
				Message: "list data empty",
			},
		},
		{
			name: "Testcase #4: Response failed (ex: Bad Request)",
			args: args{
				code:    http.StatusBadRequest,
				message: "id cannot be empty",
				params:  []interface{}{multiError},
			},
			want: &Response{
				Success: false,
				Code:    400,
				Message: "id cannot be empty",
				Errors:  map[string]string{"test": "error test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPResponse(tt.args.code, tt.args.message, tt.args.params...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestHTTPResponseJSON(t *testing.T) {
// 	resp := NewHTTPResponse(200, "success")
// 	w := new(writer)
// 	assert.NoError(t, resp.JSON(w))
// }

// func TestHTTPResponseXML(t *testing.T) {
// 	resp := NewHTTPResponse(200, "success")
// 	w := new(writer)
// 	assert.NoError(t, resp.XML(w))
// }

func TestCreateMeta(t *testing.T) {
	type args struct {
		total int
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want Meta
	}{
		{
			name: "test case 1",
			args: args{
				total: 1,
				page:  1,
				limit: 10,
			},
			want: Meta{
				Page:         1,
				Limit:        10,
				TotalRecords: int(math.Ceil(float64(1) / float64(10))),
				TotalPages:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateMeta(tt.args.total, tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}
