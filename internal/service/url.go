package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.opentelemetry.io/otel/trace"

	"github.com/Oguzyildirim/url-info/internal"
)

var doctypes = make(map[string]string)

func init() {
	doctypes["HTML 4.01 Strict"] = `"-//W3C//DTD HTML 4.01//EN"`
	doctypes["HTML 4.01 Transitional"] = `"-//W3C//DTD HTML 4.01 Transitional//EN"`
	doctypes["HTML 4.01 Frameset"] = `"-//W3C//DTD HTML 4.01 Frameset//EN"`
	doctypes["XHTML 1.0 Strict"] = `"-//W3C//DTD XHTML 1.0 Strict//EN"`
	doctypes["XHTML 1.0 Transitional"] = `"-//W3C//DTD XHTML 1.0 Transitional//EN"`
	doctypes["XHTML 1.0 Frameset"] = `"-//W3C//DTD XHTML 1.0 Frameset//EN"`
	doctypes["XHTML 1.1"] = `"-//W3C//DTD XHTML 1.1//EN"`
	doctypes["HTML 5"] = `<!DOCTYPE html>`
}

// URLRepository defines the datastore handling persisting URL records
type URLRepository interface {
	Create(ctx context.Context, HTMLVersion string, pageTitle string, headingsCount string, linksCount int,
		inaccessibleLinksCount int, haveLoginForm bool) (internal.URL, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (internal.URL, error)
}

// URL defines the application service in charge of interacting with URLs
type URL struct {
	repo URLRepository
}

// NewURL
func NewURL(repo URLRepository) *URL {
	return &URL{
		repo: repo,
	}
}

// Create stores a new record
func (u *URL) Search(ctx context.Context, URL string) (internal.URL, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("URLTracer").Start(ctx, "URL.Create")
	defer span.End()

	// do crawling here
	resp, err := http.Get(URL)

	if err != nil {
		return internal.URL{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return internal.URL{}, errors.New("Error Retrieving Document")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return internal.URL{}, fmt.Errorf("NewDocumentFromReader: %w", err)
	}

	html, err := doc.Html()

	if err != nil {
		return internal.URL{}, fmt.Errorf("doc.Html: %w", err)
	}

	// get html version
	htmlVersion := detectHTMLVersion(html)

	// get page title
	pageTitle := detectPageTitle(doc)

	// get headings count by level
	headings := detectHeadingsCountByLevel(doc)

	// get internal links
	linksCount := detectLinks(doc)

	// get internal links
	inaccessibleLinksCount := detectInaccessibleLinksCount(doc)

	// get internal links
	haveLoginForm := detectHaveLoginForm(doc)

	info, err := u.repo.Create(ctx, htmlVersion, pageTitle, headings, linksCount, inaccessibleLinksCount, haveLoginForm)
	if err != nil {
		return internal.URL{}, fmt.Errorf("repo create: %w", err)
	}
	return info, nil
}

// Delete removes an existing URL from the datastore
func (u *URL) Delete(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("URLTracer").Start(ctx, "URL.Delete")
	defer span.End()

	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo delete: %w", err)
	}

	return nil
}

// Find gets an existing URL from the datastore
func (u *URL) Find(ctx context.Context, id string) (internal.URL, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("URLTracer").Start(ctx, "URL.Find")
	defer span.End()

	URL, err := u.repo.Find(ctx, id)
	if err != nil {
		return internal.URL{}, fmt.Errorf("repo find: %w", err)
	}

	return URL, nil
}

func detectHTMLVersion(html string) string {
	version := "UNKNOWN"

	for doctype, matcher := range doctypes {
		match := strings.Contains(html, matcher)

		if match == true {
			version = doctype
		}
	}

	return version
}

func detectPageTitle(query *goquery.Document) string {
	title := query.Find("title").Contents()
	return title.Text()
}

func detectHeadingsCountByLevel(query *goquery.Document) string {
	h1 := query.Find("h1")
	h1Count := h1.Length()
	header1 := fmt.Sprintf("%s%d", "h1: ", h1Count)

	h2 := query.Find("h2")
	h2Count := h2.Length()
	header2 := fmt.Sprintf("%s%d", "h2: ", h2Count)

	h3 := query.Find("h3")
	h3Count := h3.Length()
	header3 := fmt.Sprintf("%s%d", "h3: ", h3Count)

	h4 := query.Find("h4")
	h4Count := h4.Length()
	header4 := fmt.Sprintf("%s%d", "h4: ", h4Count)

	h5 := query.Find("h5")
	h5Count := h5.Length()
	header5 := fmt.Sprintf("%s%d", "h5: ", h5Count)

	h6 := query.Find("h6")
	h6Count := h6.Length()
	header6 := fmt.Sprintf("%s %d", "h6:  ", h6Count)

	return header1 + "  " + header2 + "  " + header3 + "  " + header4 + "  " + header5 + "  " + header6
}

func detectLinks(query *goquery.Document) int {
	return query.Find("body a").Length()
}

func detectInaccessibleLinksCount(query *goquery.Document) int {
	var links []string
	query.Find("body a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		links = append(links, link)
	})

	var inaccessibleLinksCount int
	accessChan := make(chan bool)

	for _, value := range links {
		go isAccessible(value, accessChan)
	}

	for {
		select {
		case isAccessible := <-accessChan:
			if !isAccessible {
				inaccessibleLinksCount++
			}
		case <-time.After(time.Millisecond * 10):
			return inaccessibleLinksCount
		}
	}
}

func isAccessible(url string, accessChan chan<- bool) {
	_, err := http.Get(url)
	if err != nil {
		accessChan <- false
	}
	accessChan <- true
}

func detectHaveLoginForm(query *goquery.Document) bool {
	html, _ := query.Html()
	var loginRelatedWords = []string{"login", "password", "signup", "signin", "logout"}
	for _, word := range loginRelatedWords {
		if strings.Contains(html, word) {
			return true
		}
	}
	return false
}
