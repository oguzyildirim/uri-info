package internal_test

import (
	"errors"
	"testing"

	"github.com/Oguzyildirim/url-info/internal"
)

func TestURL_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   internal.URL
		withErr bool
	}{
		{
			"OK",
			internal.URL{
				HTMLVersion:            "22",
				PageTitle:              "title",
				HeadingsCount:          "hc",
				LinksCount:             2,
				InaccessibleLinksCount: 5,
				HaveLoginForm:          true,
			},
			false,
		},
		{
			"ERR: HTMLVersion",
			internal.URL{
				PageTitle:              "title",
				HeadingsCount:          "hc",
				LinksCount:             2,
				InaccessibleLinksCount: 5,
				HaveLoginForm:          true,
			},
			true,
		},
		{
			"ERR: PageTitle",
			internal.URL{
				HTMLVersion:            "22",
				HeadingsCount:          "hc",
				LinksCount:             2,
				InaccessibleLinksCount: 5,
				HaveLoginForm:          true,
			},
			true,
		},
		{
			"ERR: HeadingsCount",
			internal.URL{
				HTMLVersion:            "22",
				PageTitle:              "title",
				LinksCount:             2,
				InaccessibleLinksCount: 5,
				HaveLoginForm:          true,
			},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualErr := tt.input.Validate()
			if (actualErr != nil) != tt.withErr {
				t.Fatalf("expected error %t, got %s", tt.withErr, actualErr)
			}

			var ierr *internal.Error
			if tt.withErr && !errors.As(actualErr, &ierr) {
				t.Fatalf("expected %T error, got %T", ierr, actualErr)
			}
		})
	}
}
