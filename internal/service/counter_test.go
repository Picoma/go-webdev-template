package service_test

// This file documents the contract of the service layer.
//
// The CountService sits between transport and persistence. Its responsibility
// is intentionally small:
//
//   - delegate business operations to the persistence layer;
//   - preserve successful results unchanged;
//   - translate persistence failures into service-level errors;
//   - enrich the request context with database telemetry.
//
// The tests intentionally do not verify SQL execution. SQL semantics belong to
// the repository layer and are already covered by the database package tests.

import (
	"context"
	"errors"
	"testing"

	"idp/internal/service"

	"github.com/stretchr/testify/require"
)

// fakeStore is a minimal handwritten test double.
//
// The service depends on a very small interface, making a dedicated fake easier
// to understand than a generic mocking framework. The fake records which method
// was invoked while allowing each test to control returned values and errors.
//
// fakeStore implements [service.CounterStore].
type fakeStore struct {
	getValue int64
	getErr   error

	incrementValue int64
	incrementErr   error

	getCalled       bool
	incrementCalled bool
}

func (f *fakeStore) GetCounter(context.Context) (int64, error) {
	f.getCalled = true
	return f.getValue, f.getErr
}

func (f *fakeStore) IncrementAndGetCounter(context.Context) (int64, error) {
	f.incrementCalled = true
	return f.incrementValue, f.incrementErr
}

// TestGet verifies the behavior of the read operation.
//
// The service must faithfully propagate successful results while translating
// persistence errors into a service-specific error. This allows callers to
// distinguish the failing operation while preserving the original cause via
// [errors.Is].
func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns the counter value", func(t *testing.T) {
		t.Parallel()

		store := &fakeStore{
			getValue: 42,
		}

		svc := service.NewCountService(store)

		value, err := svc.Get(t.Context())

		require.NoError(t, err)
		require.EqualValues(t, 42, value)
		require.True(t, store.getCalled)
		require.False(t, store.incrementCalled)
	})

	t.Run("wraps store errors", func(t *testing.T) {
		t.Parallel()

		expected := errors.New("database unavailable")

		store := &fakeStore{
			getErr: expected,
		}

		svc := service.NewCountService(store)

		value, err := svc.Get(t.Context())

		require.Zero(t, value)
		require.Error(t, err)
		require.ErrorContains(t, err, "error getting counter")
		require.ErrorIs(t, err, expected)

		require.True(t, store.getCalled)
		require.False(t, store.incrementCalled)
	})
}

// TestIncrement verifies the behavior of the write operation.
//
// Besides delegating to the persistence layer, the service is responsible for
// preserving the returned value and translating persistence failures into
// contextual service errors.
func TestIncrement(t *testing.T) {
	t.Parallel()

	t.Run("returns the incremented counter value", func(t *testing.T) {
		t.Parallel()

		store := &fakeStore{
			incrementValue: 43,
		}

		svc := service.NewCountService(store)

		value, err := svc.Increment(t.Context())

		require.NoError(t, err)
		require.EqualValues(t, 43, value)
		require.False(t, store.getCalled)
		require.True(t, store.incrementCalled)
	})

	t.Run("wraps store errors", func(t *testing.T) {
		t.Parallel()

		expected := errors.New("write failed")

		store := &fakeStore{
			incrementErr: expected,
		}

		svc := service.NewCountService(store)

		value, err := svc.Increment(t.Context())

		require.Zero(t, value)
		require.Error(t, err)
		require.ErrorContains(t, err, "error incrementing counter")
		require.ErrorIs(t, err, expected)

		require.False(t, store.getCalled)
		require.True(t, store.incrementCalled)
	})
}

/*

Future test: request telemetry
==============================

The service currently enriches the request context with database telemetry via
httplog.SetAttrs(). This is observable behavior because downstream middleware
uses these attributes when emitting structured logs.

Unfortunately, httplog currently exposes no public API allowing tests to inspect
the accumulated attributes attached to a context. As a result, the behavior
cannot be verified without relying on httplog internals.

This test is intentionally left as a placeholder because the application plans
to replace httplog.SetAttrs() with its own event abstraction. Once that
abstraction exposes an inspectable API, this test should be implemented.

func TestGet_AttachesDatabaseTelemetry(t *testing.T) {
	t.Parallel()

	store := &fakeStore{
		getValue: 42,
	}

	svc := NewCountService(store)

	ctx := event.NewContext(t.Context())

	_, err := svc.Get(ctx)
	require.NoError(t, err)

	ev := event.FromContext(ctx)

	require.Equal(t, queries.GetCounter, ev.String("db.query.text"))
	require.Equal(t, "counter", ev.String("db.namespace"))
	require.Equal(t, "SELECT", ev.String("db.operation.name"))
	require.Equal(t, "GetCounter", ev.String("db.stored_procedure.name"))
	require.Equal(t, "1", ev.String("db.response.returned_rows"))
}

*/
