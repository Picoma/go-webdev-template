package layouts_test

// This file documents the composition contract of the Base layout.
//
// The layout is responsible for:
//
//   - producing a complete HTML document;
//   - providing the application's common assets;
//   - rendering child content inside the main application container.
//
// Tests intentionally verify only stable structural guarantees. They avoid
// checking styling details, allowing the visual appearance to evolve without
// affecting the suite.

import (
	"bytes"
	"testing"

	"idp/internal/web/templates/layouts"

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

// TestBaseWithChild verifies that Base preserves caller-provided content.
//
// The layout is used as a composition primitive throughout the application.
// This test ensures that child components are not discarded and are rendered
// inside the expected application container.
func TestBaseWithChild(t *testing.T) {
	t.Parallel()

	doc := renderComponent(t, layouts.TestBaseWithChild())

	require.Equal(t, 1, doc.Find("html").Length())
	require.Equal(t, 1, doc.Find("head").Length())
	require.Equal(t, 1, doc.Find("body").Length())
	require.Equal(t, 1, doc.Find("main").Length())

	require.Equal(t, 1, doc.Find("main #test-child").Length())
}
