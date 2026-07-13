package pages_test

// This file documents the composition contract of the Counter page.
//
// Unlike component tests, page tests verify how reusable pieces are assembled.
// They intentionally avoid re-testing individual components. For example,
// CounterValue is tested independently and therefore this suite only verifies
// that the page contains the component rather than asserting how the component
// itself renders.

import (
	"bytes"
	"testing"

	"idp/internal/web/templates/pages"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/stretchr/testify/require"
)

func renderComponent(t *testing.T, c templ.Component) *goquery.Document {
	t.Helper()

	var buf bytes.Buffer

	require.NoError(t, c.Render(t.Context(), &buf))

	doc, err := goquery.NewDocumentFromReader(&buf)
	require.NoError(t, err)

	return doc
}

// TestCounterPage verifies the page composition contract.
//
// The page assembles reusable UI components and HTMX wiring into a complete
// user interface. The exact HTML structure is intentionally not asserted so
// that layout refactorings remain inexpensive.
func TestCounterPage(t *testing.T) {
	t.Parallel()

	doc := renderComponent(t, pages.CounterPage(42))

	require.Equal(t, "Counter", doc.Find("h1").First().Text())
	require.Contains(t, doc.Text(), "Current value")

	// Verify the page embeds the reusable counter component.
	require.Equal(t, 1, doc.Find("#counter-value").Length())

	button := doc.Find(`button[hx-post="/counter"]`)

	require.Equal(t, 1, button.Length())
	require.Equal(t, "#counter-value", button.AttrOr("hx-target", ""))
	require.Equal(t, "outerHTML", button.AttrOr("hx-swap", ""))
}
