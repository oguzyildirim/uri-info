// Package internal defines the types used to create URL and their corresponding attributes
package internal

// URL is an activity that needs to be completed within a period of time
type URL struct {
	ID                     string
	HTMLVersion            string
	PageTitle              string
	HeadingsCount          string
	LinksCount             int
	InaccessibleLinksCount int
	HaveLoginForm          bool
}

// Validate ...
func (u URL) Validate() error {
	if u.HTMLVersion == "" {
		return NewErrorf(ErrorCodeInvalidArgument, "HTMLVersion is required")
	}
	if u.PageTitle == "" {
		return NewErrorf(ErrorCodeInvalidArgument, "PageTitle is required")
	}
	if u.HeadingsCount == "" {
		return NewErrorf(ErrorCodeInvalidArgument, "HeadingsCount is required")
	}
	return nil
}
