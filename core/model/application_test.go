package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version(0x00010203)
	assert.Equal(t, "1.2.3", v.String())
	v2 := VersionFromString("1.2.3")
	assert.Equal(t, v, v2)
	v3 := VersionFromString("1.2.4")
	assert.True(t, v3 > v2)
	assert.False(t, v3 < v2)
	assert.Equal(t, Version(0), VersionFromString(""))
}
