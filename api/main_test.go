package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/zohaibAsif/simple_bank_management_system/db/sqlc"
	"github.com/zohaibAsif/simple_bank_management_system/util"
)

func makeNewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomName(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
