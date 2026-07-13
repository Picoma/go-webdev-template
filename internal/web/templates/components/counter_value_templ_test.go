package components_test

// This file documents the semantic HTML contract of the CounterValue component.
//
// Component tests intentionally verify semantics rather than presentation.
// They assert the HTML elements and attributes that other parts of the
// application depend upon, while deliberately ignoring implementation details
// such as CSS classes, whitespace or formatting.

import (
	"bytes"
	"strings"
	"testing"

	"idp/internal/web/templates/components"

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

// TestCounterValue verifies the semantic contract of the reusable counter
// component.
//
// The component promises to expose the current counter value through a stable
// DOM identifier. Other templates and HTMX interactions rely on this contract.
func TestCounterValue(t *testing.T) {
	t.Parallel()

	doc := renderComponent(t, components.CounterValue(42))

	span := doc.Find("#counter-value")

	require.Equal(t, 1, span.Length())
	require.Equal(t, "42", strings.TrimSpace(span.Text()))
}
