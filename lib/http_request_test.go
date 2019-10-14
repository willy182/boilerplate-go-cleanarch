package lib

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

const (
	ContentType     = "Content-Type"
	ApplicationForm = "application/x-www-form-urlencoded"
	URLMock         = "http://mock.com"
	Message         = "message"
)

type argNewRequest struct {
	timeOut time.Duration
}

var testNewRequest = []struct {
	name string
	args argNewRequest
	want bool
}{
	{
		name: "Testcase #1: Positive NewRequest",
		args: argNewRequest{
			timeOut: 10 * time.Second,
		},
		want: true,
	},
}

type argNewReq struct {
	method   string
	fullPath string
	body     io.Reader
	headers  map[string]string
}

var testNewReq = []struct {
	name string
	args argNewReq
	want bool
}{
	{
		name: "Testcase #1: Positive NewReq",
		args: argNewReq{
			method:   "GET",
			fullPath: URLMock,
			headers:  map[string]string{ContentType: ApplicationForm},
		},
		want: true,
	}, {
		name: "Testcase #1: Negative NewReq",
		args: argNewReq{
			method:   "UNDEFINED",
			fullPath: URLMock,
			headers:  map[string]string{ContentType: ApplicationForm},
		},
		want: true,
	},
}

type resp struct {
	Message string `json:"message"`
}

type argExec struct {
	method     string
	url        string
	body       io.Reader
	target     interface{}
	headers    map[string]string
	httpResp   interface{}
	statusCode int
}

var testExec = []struct {
	name string
	args argExec
	want bool
}{
	{
		name: "Testcase #1: Positive Exec",
		args: argExec{
			method:     "GET",
			url:        URLMock,
			headers:    map[string]string{ContentType: ApplicationForm},
			target:     &resp{},
			httpResp:   map[string]string{Message: "success"},
			statusCode: 200,
		},
	},
	{
		name: "Testcase #2: Negative Exec",
		args: argExec{
			method:     "GET",
			url:        URLMock,
			headers:    map[string]string{ContentType: ApplicationForm},
			target:     &resp{},
			httpResp:   "<html></html>",
			statusCode: 500,
		},
		want: true,
	},
	{
		name: "Testcase #3: Negative Exec",
		args: argExec{
			method:     "UNDEFINED",
			url:        URLMock,
			headers:    map[string]string{ContentType: ApplicationForm},
			target:     &resp{},
			httpResp:   map[string]string{Message: "undefined"},
			statusCode: 400,
		},
		want: true,
	},
	{
		name: "Testcase #4: Negative Exec",
		args: argExec{
			method:     "GET",
			url:        "<:::>",
			headers:    map[string]string{ContentType: ApplicationForm},
			target:     &resp{},
			httpResp:   map[string]string{Message: "error"},
			statusCode: 400,
		},
		want: true,
	},
}

func testExecAndExecAsync(name string, t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	newReq := &httpRequest{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	// test Exec
	for _, x := range testExec {
		MockHTTP("GET", URLMock, x.args.statusCode, x.args.httpResp)
		t.Run(x.name, func(t *testing.T) {
			if name == "exec" {
				if err := newReq.Req(x.args.method, x.args.url, x.args.body, x.args.target, x.args.headers); (err != nil) != x.want {
					t.Errorf("newReq.Exec error = %v, wantErr %v", err, x.want)
				}
			}

			if err := <-newReq.ReqAsync(x.args.method, x.args.url, x.args.body, x.args.target, x.args.headers); (err != nil) != x.want {
				t.Errorf("newReq.Exec error = %v, wantErr %v", err, x.want)
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	// test NewRequest
	for _, tt := range testNewRequest {
		t.Run(tt.name, func(t *testing.T) {
			newReq := NewRequest(tt.args.timeOut)
			if (newReq != nil) != tt.want {
				t.Errorf("NewRequest = %v, want %v", newReq, tt.want)
			}

			// test NewReq
			for _, qq := range testNewReq {
				t.Run(qq.name, func(t *testing.T) {
					resReq, e := newReq.newReq(qq.args.method, qq.args.fullPath, qq.args.body, qq.args.headers)
					if e != nil {
						t.Errorf("newReq.NewReq = %v, want %v", resReq, qq.want)
					}
				})
			}
		})
	}
}

func TestExec(t *testing.T) {
	testExecAndExecAsync("exec", t)
}

func TestExecAsync(t *testing.T) {
	testExecAndExecAsync("execAsync", t)
}
