package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitEnv(t *testing.T) {
	err := LoadEnv()
	assert.Nil(t, err)

	assert.Equal(t, "development", Config.Env)
	assert.Equal(t, "user", Config.DBUser)
	assert.Equal(t, "localhost", Config.DBHost)
	assert.Equal(t, "term_keeper_db", Config.DBName)
	assert.Equal(t, "password", Config.DBPassword)
	assert.Equal(t, 3306, Config.DBPort)
	assert.Equal(t, []string{"http://localhost:3001"}, Config.APICorsAllowsOrigins)
}
