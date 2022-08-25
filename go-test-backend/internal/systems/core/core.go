package core

import (
	"github.com/jmoiron/sqlx"
)

// CoreProvider provides *Core
type CoreProvider struct {
	Core *Core
}

// Core is capable of providing core access
type Core struct {
	dbh *sqlx.DB
}

// NewCoreProvider returns a new Core provider
func NewCoreProvider(dbh *sqlx.DB, urlPrefixBackend string) *CoreProvider {
	return &CoreProvider{
		&Core{
			dbh: dbh,
		},
	}
}

func (b *CoreProvider) NewCore() *Core {
	return b.Core
}
