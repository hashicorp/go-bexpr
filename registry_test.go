package bexpr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func benchRegistry(b *testing.B, registry Registry) {
	rtype := derefType(reflect.TypeOf((*testNestedTypes)(nil)))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		registry.GetFieldConfigurations(rtype)
	}
}

// Used as a bit of a baseline of what using a registry buys us
func BenchmarkNilRegistry(b *testing.B) {
	benchRegistry(b, NilRegistry)
}

// Used to show the speedup for repeated usage of the same type
// Additionally this ensures that we are not regenerating fields
// because if we were we wouldn't see the speedup over the nil
// registry
func BenchmarkRegistry(b *testing.B) {
	benchRegistry(b, NewSyncRegistry())
}

// This test is mostly to ensure that the SyncRegistry will return
// the same FieldConfigurations as calling GenerateFieldConfigurations
// or the internal generateFieldConfigurations. Then it ensures
// that the fields are stored in the internal map
func TestSyncRegistry(t *testing.T) {
	var ttype *testNestedTypes

	direct, err := GenerateFieldConfigurations(ttype)
	require.NoError(t, err)

	registry := NewSyncRegistry()
	rtype := derefType(reflect.TypeOf(ttype))
	fromRegistry, err := registry.GetFieldConfigurations(rtype)
	require.NoError(t, err)

	validateFieldConfigurations(t, direct, fromRegistry)

	fields, ok := registry.configurations[rtype]
	require.True(t, ok)
	validateFieldConfigurations(t, fromRegistry, fields)
}

// This just tests that the NilRegistry is a true pass through
// to the GenerateFieldConfigurations function
func TestNilRegistry(t *testing.T) {
	var ttype *testNestedTypes

	direct, err := GenerateFieldConfigurations(ttype)
	require.NoError(t, err)

	rtype := derefType(reflect.TypeOf(ttype))
	fromRegistry, err := NilRegistry.GetFieldConfigurations(rtype)
	require.NoError(t, err)

	validateFieldConfigurations(t, direct, fromRegistry)
}
