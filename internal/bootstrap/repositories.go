package bootstrap

import "database/sql"

type Repositories struct {
}

func GetRepositories(db *sql.DB) *Repositories {
	return &Repositories{}
}
