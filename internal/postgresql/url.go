package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Oguzyildirim/url-info/internal"
)

// URL represents the repository used for interacting with URL records
type URL struct {
	q *Queries
}

// NewURL instantiates the URL repository
func NewURL(db *sql.DB) *URL {
	return &URL{
		q: New(db),
	}
}

// Create inserts a new ıser record
func (u *URL) Create(ctx context.Context, HTMLVersion string, pageTitle string, headingsCount string, linksCount int,
	inaccessibleLinksCount int, haveLoginForm bool) (internal.URL, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("URLTracer").Start(ctx, "URL.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	id, err := u.q.InsertURL(ctx, InsertURLParams{
		Htmlversion:            HTMLVersion,
		Pagetitle:              pageTitle,
		Headingscount:          headingsCount,
		Linkscount:             int32(linksCount),
		Inaccessiblelinkscount: int32(inaccessibleLinksCount),
		Haveloginform:          haveLoginForm,
	})
	if err != nil {
		return internal.URL{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "insert URL")
	}

	return internal.URL{
		ID:                     id.String(),
		HTMLVersion:            HTMLVersion,
		PageTitle:              pageTitle,
		HeadingsCount:          headingsCount,
		LinksCount:             linksCount,
		InaccessibleLinksCount: inaccessibleLinksCount,
		HaveLoginForm:          haveLoginForm,
	}, nil
}

// Delete deletes the existing record matching the id
func (u *URL) Delete(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("URLTracer").Start(ctx, "URL.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	val, err := uuid.Parse(id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid uuid")
	}
	_, err = u.q.DeleteURL(ctx, val)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "URL not found")
		}

		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "delete URL")
	}
	return nil
}

// Find returns the requested URL by searching its id
func (u *URL) Find(ctx context.Context, id string) (internal.URL, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("URLTracer").Start(ctx, "URL.Find")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	val, err := uuid.Parse(id)
	if err != nil {
		return internal.URL{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid uuid")
	}
	res, err := u.q.SelectURL(ctx, val)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.URL{}, internal.WrapErrorf(err, internal.ErrorCodeNotFound, "URL not found")
		}

		return internal.URL{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "select URL")
	}
	return internal.URL{
		ID:                     res.ID.String(),
		HTMLVersion:            res.HtmlVersion,
		PageTitle:              res.PageTitle,
		HeadingsCount:          res.HeadingsCount,
		LinksCount:             int(res.LinksCount),
		InaccessibleLinksCount: int(res.InaccessibleLinksCount),
		HaveLoginForm:          res.HaveLoginForm,
	}, nil
}
