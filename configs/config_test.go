package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitEnv(t *testing.T) {
	err := LoadEnv()
	assert.Nil(t, err)

	assert.Equal(t, "development", Configs.Env)
	assert.Equal(t, []string{"http://localhost:3001"}, Configs.APICorsAllowsOrigins)
}
