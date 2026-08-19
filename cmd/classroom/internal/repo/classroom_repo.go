package repo

import "github.com/abozorov/school_online/pkg/postgres"

type Repo struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		pg: pg,
	}
}
