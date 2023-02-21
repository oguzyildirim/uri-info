package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Oguzyildirim/url-info/internal"
)

const uuidRegEx string = `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`

//go:generate counterfeiter -o resttesting/URL_service.gen.go . URLService

// URLService
type URLService interface {
	Search(ctx context.Context, URL string) (internal.URL, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (internal.URL, error)
}

// URLHandler
type URLHandler struct {
	svc URLService
}

// NewURLHandler
func NewURLHandler(svc URLService) *URLHandler {
	return &URLHandler{
		svc: svc,
	}
}

// Register connects the handlers to the router.
func (u *URLHandler) Register(r *mux.Router) {
	r.HandleFunc("/URLs", u.search).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/URLs/{id:%s}", uuidRegEx), u.find).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/URLs/{id:%s}", uuidRegEx), u.delete).Methods(http.MethodDelete)
}

// URL is one of the key concepts of the Web. It is the mechanism used by browsers to retrieve any published resource on the web
type URL struct {
	ID                     string `json:"id"`
	HTMLVersion            string `json:"HTMLVersion"`
	PageTitle              string `json:"pageTitle"`
	HeadingsCount          string `json:"headingsCount"`
	LinksCount             int    `json:"linksCount"`
	InaccessibleLinksCount int    `json:"inaccessibleLinksCount"`
	HaveLoginForm          bool   `json:"haveLoginForm"`
}

// CreateURLsRequest defines the request used for creating URLs.
type CreateURLsRequest struct {
	URL string `json:"url"`
}

// CreateURLsResponse defines the response returned back after creating URLs.
type CreateURLsResponse struct {
	URL URL `json:"URL"`
}

func (u *URLHandler) search(w http.ResponseWriter, r *http.Request) {
	var req CreateURLsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer r.Body.Close()

	url, err := u.svc.Search(r.Context(), req.URL)
	fmt.Println(err)
	if err != nil {
		renderErrorResponse(r.Context(), w, "search failed", err)
		return
	}

	renderResponse(w,
		&CreateURLsResponse{
			URL: URL{
				ID:                     url.ID,
				HTMLVersion:            url.HTMLVersion,
				PageTitle:              url.PageTitle,
				HeadingsCount:          url.HeadingsCount,
				LinksCount:             url.LinksCount,
				InaccessibleLinksCount: url.InaccessibleLinksCount,
				HaveLoginForm:          url.HaveLoginForm,
			},
		},
		http.StatusCreated)
}

func (u *URLHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]

	if err := u.svc.Delete(r.Context(), id); err != nil {
		renderErrorResponse(r.Context(), w, "delete failed", err)
		return
	}

	renderResponse(w, struct{}{}, http.StatusOK)
}

// ReadURLsResponse defines the response returned back after searching one URL.
type ReadURLResponse struct {
	URL URL `json:"URL"`
}

func (u *URLHandler) find(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]

	url, err := u.svc.Find(r.Context(), id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "find failed", err)
		return
	}

	renderResponse(w,
		&ReadURLResponse{
			URL: URL{
				ID:                     url.ID,
				HTMLVersion:            url.HTMLVersion,
				PageTitle:              url.PageTitle,
				HeadingsCount:          url.HeadingsCount,
				LinksCount:             url.LinksCount,
				InaccessibleLinksCount: url.InaccessibleLinksCount,
				HaveLoginForm:          url.HaveLoginForm,
			},
		},
		http.StatusOK)
}
