package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/mux"

	"github.com/Oguzyildirim/url-info/internal"
	"github.com/Oguzyildirim/url-info/internal/rest"
	"github.com/Oguzyildirim/url-info/internal/rest/resttesting"
)

func TestURLs_Delete(t *testing.T) {
	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		setup  func(*resttesting.FakeURLService)
		output output
	}{
		{
			"OK: 200",
			func(s *resttesting.FakeURLService) {},
			output{
				http.StatusOK,
				&struct{}{},
				&struct{}{},
			},
		},
		{
			"ERR: 404",
			func(s *resttesting.FakeURLService) {
				s.DeleteReturns(internal.NewErrorf(internal.ErrorCodeNotFound, "not found"))
			},
			output{
				http.StatusNotFound,
				&struct{}{},
				&struct{}{},
			},
		},
		{
			"ERR: 500",
			func(s *resttesting.FakeURLService) {
				s.DeleteReturns(errors.New("service failed"))
			},
			output{
				http.StatusInternalServerError,
				&struct{}{},
				&struct{}{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := mux.NewRouter()
			svc := &resttesting.FakeURLService{}
			tt.setup(svc)

			rest.NewURLHandler(svc).Register(router)

			res := doRequest(router,
				httptest.NewRequest(http.MethodDelete, "/URLs/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", nil))

			assertResponse(t, res, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestURLs_Search(t *testing.T) {
	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		setup  func(*resttesting.FakeURLService)
		input  []byte
		output output
	}{
		{
			"OK: 201",
			func(s *resttesting.FakeURLService) {
				s.SearchReturns(
					internal.URL{
						ID:                     "1-2-3",
						HTMLVersion:            "22",
						PageTitle:              "url.PageTitle",
						HeadingsCount:          "url.HeadingsCount",
						LinksCount:             2,
						InaccessibleLinksCount: 3,
						HaveLoginForm:          true,
					},
					nil)
			},
			func() []byte {
				b, _ := json.Marshal(&rest.CreateURLsRequest{
					URL: "exampleurl",
				})

				return b
			}(),
			output{
				http.StatusCreated,
				&rest.CreateURLsResponse{
					URL: rest.URL{
						ID:                     "1-2-3",
						HTMLVersion:            "22",
						PageTitle:              "url.PageTitle",
						HeadingsCount:          "url.HeadingsCount",
						LinksCount:             2,
						InaccessibleLinksCount: 3,
						HaveLoginForm:          true,
					},
				},
				&rest.CreateURLsResponse{},
			},
		},
		{
			"ERR: 400",
			func(*resttesting.FakeURLService) {},
			[]byte(`{"invalid":"json`),
			output{
				http.StatusBadRequest,
				&rest.ErrorResponse{
					Error: "invalid request",
				},
				&rest.ErrorResponse{},
			},
		},
		{
			"ERR: 500",
			func(s *resttesting.FakeURLService) {
				s.SearchReturns(internal.URL{},
					errors.New("service error"))
			},
			[]byte(`{}`),
			output{
				http.StatusInternalServerError,
				&rest.ErrorResponse{
					Error: "internal error",
				},
				&rest.ErrorResponse{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := mux.NewRouter()
			svc := &resttesting.FakeURLService{}
			tt.setup(svc)

			rest.NewURLHandler(svc).Register(router)

			res := doRequest(router,
				httptest.NewRequest(http.MethodPost, "/URLs", bytes.NewReader(tt.input)))

			assertResponse(t, res, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestURLs_Find(t *testing.T) {
	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		setup  func(*resttesting.FakeURLService)
		output output
	}{
		{
			"OK: 200",
			func(s *resttesting.FakeURLService) {
				s.FindReturns(
					internal.URL{
						ID:                     "a-b-c",
						HTMLVersion:            "22",
						PageTitle:              "url.PageTitle",
						HeadingsCount:          "url.HeadingsCount",
						LinksCount:             2,
						InaccessibleLinksCount: 3,
						HaveLoginForm:          true,
					},
					nil)
			},
			output{
				http.StatusOK,
				&rest.ReadURLResponse{
					URL: rest.URL{
						ID:                     "a-b-c",
						HTMLVersion:            "22",
						PageTitle:              "url.PageTitle",
						HeadingsCount:          "url.HeadingsCount",
						LinksCount:             2,
						InaccessibleLinksCount: 3,
						HaveLoginForm:          true,
					},
				},
				&rest.ReadURLResponse{},
			},
		},
		{
			"OK: 200",
			func(s *resttesting.FakeURLService) {
				s.FindReturns(internal.URL{},
					internal.NewErrorf(internal.ErrorCodeNotFound, "not found"))
			},
			output{
				http.StatusNotFound,
				&rest.ErrorResponse{
					Error: "find failed",
				},
				&rest.ErrorResponse{},
			},
		},
		{
			"ERR: 500",
			func(s *resttesting.FakeURLService) {
				s.FindReturns(internal.URL{},
					errors.New("service error"))
			},
			output{
				http.StatusInternalServerError,
				&rest.ErrorResponse{
					Error: "internal error",
				},
				&rest.ErrorResponse{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := mux.NewRouter()
			svc := &resttesting.FakeURLService{}
			tt.setup(svc)

			rest.NewURLHandler(svc).Register(router)

			res := doRequest(router,
				httptest.NewRequest(http.MethodGet, "/URLs/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", nil))

			assertResponse(t, res, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

type test struct {
	expected interface{}
	target   interface{}
}

func doRequest(router *mux.Router, req *http.Request) *http.Response {
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr.Result()
}

func assertResponse(t *testing.T, res *http.Response, test test) {
	t.Helper()

	if err := json.NewDecoder(res.Body).Decode(test.target); err != nil {
		t.Fatalf("couldn't decode %s", err)
	}
	defer res.Body.Close()

	if !cmp.Equal(test.expected, test.target, cmpopts.IgnoreUnexported()) {
		t.Fatalf("expected results don't match: %s", cmp.Diff(test.expected, test.target, cmpopts.IgnoreUnexported()))
	}
}
