package cards

import (
	"database/sql"
)

type CardsDBRepository struct {
	dtb *sql.DB
}

func NewDBRepo(sdb *sql.DB) *CardsDBRepository {
	return &CardsDBRepository{dtb: sdb}
}
