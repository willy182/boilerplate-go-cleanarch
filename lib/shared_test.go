package lib

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	length := 10
	rs := RandomString(length)

	// set invalid string that
	// should not be contained in random string
	invalidString := `!@#$%^&*()_+`

	if len(rs) != length {
		t.Errorf("length of random string is not equal %d", length)
	}

	if strings.Contains(rs, invalidString) {
		t.Fatal("random string contains symbols")
	}
}

func TestValidateEmail(t *testing.T) {
	var (
		email string
		err   error
	)

	// test valid email should be not error/valid
	email = "larry.page@google.com"
	if err = ValidateEmail(email); err != nil {
		t.Fatal("testing valid email is not valid")
	}

	// test invalid email should be error
	email = "larry.@page@google.com"
	if err = ValidateEmail(email); err == nil {
		t.Fatal("testing invalid email is not valid")
	}
}

func TestValidateURL(t *testing.T) {
	var (
		url string
		err error
	)

	url = "http://www.google.com"
	if err = ValidateURL(url); err != nil {
		t.Fatal("testing 1st valid URL is not valid")
	}

	url = "www.google.com"
	if err = ValidateURL(url); err != nil {
		t.Fatal("testing 2nd valid URL is not valid")
	}

	url = "ftp://www.google.com"
	if err = ValidateURL(url); err != nil {
		t.Fatal("testing 3rd valid URL is not valid")
	}

	url = "https:///www.google.com"
	if err = ValidateURL(url); err == nil {
		t.Fatal("testing invalid URL is not valid")
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	var (
		tel string
		err error
	)

	tel = "08119889788"
	if err = ValidatePhoneNumber(tel); err != nil {
		t.Fatal("testing valid phone number is not valid")
	}

	tel = "081-1988-9788"
	if err = ValidatePhoneNumber(tel); err == nil {
		t.Fatal("testing 1st invalid phone number is not valid")
	}

	tel = "0811"
	if err = ValidatePhoneNumber(tel); err == nil {
		t.Fatal("testing 2nd invalid phone number - not greater than 5 chars is not valid")
	}
}

func TestValidatePhoneAreaNumber(t *testing.T) {
	var (
		area string
		err  error
	)

	area = "+62"
	if err = ValidatePhoneAreaNumber(area); err != nil {
		t.Fatal("testing valid area number is not valid")
	}

	area = "+6 2"
	if err = ValidatePhoneAreaNumber(area); err == nil {
		t.Fatal("testing 1st invalid area number is not valid")
	}

	area = "+"
	if err = ValidatePhoneAreaNumber(area); err == nil {
		t.Fatal("testing 2nd invalid area number is not valid")
	}
}

func TestValidateNumeric(t *testing.T) {
	t.Run("Test Validate Numeric", func(t *testing.T) {
		boolFalse := ValidateNumeric("1.0.1")
		assert.False(t, boolFalse)

		boolTrue := ValidateNumeric("0123456789")
		assert.True(t, boolTrue)
	})
}

func TestValidateAlphabet(t *testing.T) {
	t.Run("Test Validate Alphabet", func(t *testing.T) {
		boolTrue := ValidateAlphabet("huFtBanGeT")
		assert.True(t, boolTrue)

		boolFalse := ValidateAlphabet("1FgH^*")
		assert.False(t, boolFalse)
	})
}

func TestValidateAlphabetWithSpace(t *testing.T) {
	t.Run("Test Validate Alphabet With Space", func(t *testing.T) {
		boolFalse := ValidateAlphabetWithSpace("huFtBanGeT*")
		assert.False(t, boolFalse)

		boolTrue := ValidateAlphabetWithSpace("huFt BanGeT")
		assert.True(t, boolTrue)
	})
}

func TestValidateAlphanumeric(t *testing.T) {
	t.Run("Test Validate Alphabet Numeric", func(t *testing.T) {
		boolTrue := ValidateAlphanumeric("okesip12", true)
		assert.True(t, boolTrue)

		boolTrue = ValidateAlphanumeric("okesip", false)
		assert.True(t, boolTrue)

		boolFalse := ValidateAlphanumeric("1FgH^*", false)
		assert.False(t, boolFalse)
	})
}

func TestValidateAlphanumericWithSpace(t *testing.T) {
	t.Run("Test Validate Alphabet Numeric With Space", func(t *testing.T) {
		boolTrue := ValidateAlphanumericWithSpace("oke sip1", false)
		assert.True(t, boolTrue)

		boolTrue = ValidateAlphanumericWithSpace("OKE sip1", false)
		assert.True(t, boolTrue)

		boolFalse := ValidateAlphanumericWithSpace("okesip1", true)
		assert.False(t, boolFalse)

		boolFalse = ValidateAlphanumericWithSpace("okesip1@", true)
		assert.False(t, boolFalse)
	})
}

func TestGenerateRandomID(t *testing.T) {
	t.Run("Test Generate Random ID", func(t *testing.T) {
		var res string
		randomID := GenerateRandomID(5)
		assert.IsType(t, res, randomID)

		randomID = GenerateRandomID(5, "00")
		assert.IsType(t, res, randomID)
	})
}

func TestRandomNumber(t *testing.T) {
	t.Run("Test Generate Random Number", func(t *testing.T) {
		var res string
		randomNumber := RandomNumber(5)

		assert.IsType(t, res, randomNumber)
	})
}

func TestStringInSlice(t *testing.T) {
	var positiveStr = "mantab"
	var positiveStrCheck = "mantul"
	type args struct {
		str           string
		list          []string
		caseSensitive []bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Testcase #1: Positive",
			args: args{str: positiveStr, list: []string{positiveStr, positiveStrCheck}},
			want: true,
		},
		{
			name: "Testcase #2: Positive",
			args: args{str: positiveStr, list: []string{positiveStr, positiveStrCheck}, caseSensitive: []bool{false}},
			want: true,
		},
		{
			name: "Testcase #3: Negative",
			args: args{str: positiveStr, list: []string{"mantap", positiveStrCheck}, caseSensitive: []bool{false}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.str, tt.args.list, tt.args.caseSensitive...); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProtocol(t *testing.T) {
	type args struct {
		isTLS bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test case http",
			args: args{isTLS: false},
			want: "http://",
		},
		{
			name: "test case https",
			args: args{isTLS: true},
			want: "https://",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProtocol(tt.args.isTLS); got != tt.want {
				t.Errorf("GetProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostURL(t *testing.T) {
	httpDummy := "http://google.com"
	urlDummy, _ := http.NewRequest("GET", httpDummy, nil)

	type args struct {
		req *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "always positive",
			args: args{req: urlDummy},
			want: httpDummy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHostURL(tt.args.req); got != tt.want {
				t.Errorf("GetHostURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSelfLink(t *testing.T) {
	httpDummy := "http://google.com"
	reqDummy, _ := http.NewRequest("GET", httpDummy, nil)

	type args struct {
		req *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "without secure protocol",
			args: args{req: reqDummy},
			want: httpDummy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSelfLink(tt.args.req); got != tt.want {
				t.Errorf("GetSelfLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockHTTP(*testing.T) {
	type modelMock struct {
		Members []struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"members"`
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	MockHTTP("GET", "http://google.mock", 200, modelMock{
		Members: []struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		}{
			{
				Email:    "foo@bar.com",
				Username: "foobar",
			},
		},
	})
}
