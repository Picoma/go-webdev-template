package handler_test

// This file documents the HTTP contract of CounterHandler.
//
// CounterHandler is a thin transport adapter responsible for:
//
//   - translating HTTP requests into service calls;
//   - selecting the appropriate template;
//   - translating service failures into HTTP 500 responses.
//
// The tests intentionally verify observable HTTP behavior rather than HTML
// implementation details. They assert the presence of stable, user-visible
// elements instead of the complete rendered document, allowing the page layout
// and styling to evolve without breaking the suite.

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"idp/internal/handler"

	"github.com/stretchr/testify/require"
)

// fakeCounterService is a minimal handwritten test double.
//
// The handler depends on a very small interface, making a dedicated fake easier
// to understand than a generic mocking framework. It records which operation
// was invoked while allowing each test to control the returned value or error.
//
// fakeCounterService implements [CounterService].
type fakeCounterService struct {
	getValue int64
	getErr   error

	incrementValue int64
	incrementErr   error

	getCalled       bool
	incrementCalled bool
}

func (f *fakeCounterService) Get(context.Context) (int64, error) {
	f.getCalled = true
	return f.getValue, f.getErr
}

func (f *fakeCounterService) Increment(context.Context) (int64, error) {
	f.incrementCalled = true
	return f.incrementValue, f.incrementErr
}

// TestDisplayCounter verifies the HTTP contract of the page endpoint.
//
// The endpoint is expected to:
//
//   - retrieve the current value from the service;
//   - render the complete page;
//   - expose the counter value within the rendered document;
//   - translate service failures into HTTP 500 responses.
func TestDisplayCounter(t *testing.T) {
	t.Parallel()

	t.Run("renders the counter page", func(t *testing.T) {
		t.Parallel()

		service := &fakeCounterService{
			getValue: 42,
		}

		handler := handler.NewCounterHandler(service)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.DisplayCounter(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		body := rec.Body.String()

		require.Contains(t, body, "Counter")
		require.Contains(t, body, "Current value")
		require.Contains(t, body, `id="counter-value"`)
		require.Contains(t, body, ">42<")

		require.True(t, service.getCalled)
		require.False(t, service.incrementCalled)
	})

	t.Run("returns internal server error when the service fails", func(t *testing.T) {
		t.Parallel()

		expected := errors.New("database unavailable")

		service := &fakeCounterService{
			getErr: expected,
		}

		handler := handler.NewCounterHandler(service)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.DisplayCounter(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "database unavailable")

		require.True(t, service.getCalled)
		require.False(t, service.incrementCalled)
	})
}

// TestIncreaseCounter verifies the HTMX endpoint contract.
//
// Unlike DisplayCounter, this endpoint intentionally renders only the HTML
// fragment required to update the current counter value.
func TestIncreaseCounter(t *testing.T) {
	t.Parallel()

	t.Run("renders the counter fragment", func(t *testing.T) {
		t.Parallel()

		service := &fakeCounterService{
			incrementValue: 43,
		}

		handler := handler.NewCounterHandler(service)

		req := httptest.NewRequest(http.MethodPost, "/counter", nil)
		rec := httptest.NewRecorder()

		handler.IncreaseCounter(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		body := rec.Body.String()

		require.Contains(t, body, `id="counter-value"`)
		require.Contains(t, body, ">43<")

		// The HTMX endpoint must only return the replacement fragment.
		require.NotContains(t, body, "Current value")
		require.NotContains(t, body, "<html")

		require.False(t, service.getCalled)
		require.True(t, service.incrementCalled)
	})

	t.Run("returns internal server error when the service fails", func(t *testing.T) {
		t.Parallel()

		expected := errors.New("write failed")

		service := &fakeCounterService{
			incrementErr: expected,
		}

		handler := handler.NewCounterHandler(service)

		req := httptest.NewRequest(http.MethodPost, "/counter", nil)
		rec := httptest.NewRecorder()

		handler.IncreaseCounter(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "write failed")

		require.False(t, service.getCalled)
		require.True(t, service.incrementCalled)
	})
}

/*

Future test: request telemetry
==============================

The handler currently enriches the request context using httplog.SetAttrs().
This is observable behavior because downstream middleware consumes these
attributes when producing structured logs.

httplog currently exposes no public API allowing tests to inspect the
attributes attached to a request context, making this behavior impossible to
verify without relying on implementation details.

Once telemetry is extracted into the application's event abstraction, the
following behavior should be verified:

- DisplayCounter attaches:
      counter.value

- IncreaseCounter attaches:
      counter.value

- Service failures do *not* attach counter.value because no value was produced.

*/

/*

Future test: template rendering failures
========================================

The handler explicitly handles errors returned by templ.Render(). Generated
templ components, however, cannot realistically be forced to fail without
modifying generated code or introducing an artificial seam.

If rendering is later abstracted behind an interface (or templates become
runtime-provided), add tests ensuring that rendering failures are translated
into HTTP 500 responses.

*/
