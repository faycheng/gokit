package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginLoader_Load(t *testing.T) {
	loader := NewLoader("../examples/plugin/echo")
	_, err := loader.Load()
	assert.Equal(t, nil, err)
}
