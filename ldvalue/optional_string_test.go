package ldvalue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyOptionalString(t *testing.T) {
	o := OptionalString{}
	assert.False(t, o.IsDefined())
	assert.Equal(t, "", o.StringValue())
	assert.Nil(t, o.AsPointer())
	assert.Equal(t, Null(), o.AsValue())
	assert.True(t, o == o)
}

func TestOptionalStringWithValue(t *testing.T) {
	o := NewOptionalStringWithValue("value")
	assert.True(t, o.IsDefined())
	assert.Equal(t, "value", o.StringValue())
	assert.NotNil(t, o.AsPointer())
	assert.Equal(t, "value", *o.AsPointer())
	assert.Equal(t, String("value"), o.AsValue())
	assert.True(t, o == o)
	assert.False(t, o == OptionalString{})
}

func TestOptionalStringFromNilPointer(t *testing.T) {
	o := NewOptionalStringFromPointer(nil)
	assert.False(t, o.IsDefined())
	assert.Equal(t, "", o.StringValue())
	assert.Nil(t, o.AsPointer())
	assert.Equal(t, Null(), o.AsValue())
	assert.True(t, o == o)
	assert.True(t, o == OptionalString{})
}

func TestOptionalStringFromNonNilPointer(t *testing.T) {
	v := "value"
	p := &v
	o := NewOptionalStringFromPointer(p)
	assert.True(t, o.IsDefined())
	assert.Equal(t, "value", o.StringValue())
	assert.NotNil(t, o.AsPointer())
	assert.Equal(t, "value", *o.AsPointer())
	assert.False(t, p == o.AsPointer()) // should not be the same pointer, just the same underlying string
	assert.Equal(t, String("value"), o.AsValue())
	assert.True(t, o == o)
	assert.True(t, o == NewOptionalStringWithValue("value"))
}

func TestOptionalStringAsStringer(t *testing.T) {
	assert.Equal(t, "[none]", OptionalString{}.String())
	assert.Equal(t, "[empty]", NewOptionalStringWithValue("").String())
	assert.Equal(t, "x", NewOptionalStringWithValue("x").String())
}
