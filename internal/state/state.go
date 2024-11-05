package state

import (
	"gator/internal/config"
	"gator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}
