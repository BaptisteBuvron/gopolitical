package gopolitical

import (
	"testing"
)

func TestSWFFactory(t *testing.T) {
	assert := NewAssert(t)
	_, err := LoadSimulation("../resources/data.json")
	assert.NoError(err)
}
