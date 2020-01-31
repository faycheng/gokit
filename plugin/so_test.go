package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSoPlugin(t *testing.T) {
	soPlugin := NewSoPlugin("../examples/plugin/so/echo.so")
	call, err := soPlugin.Lookup("Echo")
	assert.Nil(t, err)
	_, err = call(context.TODO(), nil)
	assert.Nil(t, err)
}
