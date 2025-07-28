package images

import (
	"database/sql"
)

type ImagesDBRepository struct {
	dtb *sql.DB
}

func NewDBRepo(sdb *sql.DB) *ImagesDBRepository {
	return &ImagesDBRepository{dtb: sdb}
}
