package shortener

import (
	"database/sql"
	"errors"
	"github.com/VoRaX00/shortener/internal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Add(link string) (string, error) {
	id := xid.New().String()
	query := `INSERT INTO links (id, link) VALUES ($1, $2) 
				 ON CONFLICT (link) DO UPDATE 
				 SET id=links.id
				 RETURNING id;`

	var existsId string
	err := r.db.QueryRow(query, id, link).Scan(&existsId)
	if err != nil {
		return "", err
	}
	return existsId, err
}

func (r *Repository) Get(id string) (string, error) {
	query := `SELECT link FROM links WHERE id = $1;`
	var link string
	err := r.db.QueryRow(query, id).Scan(&link)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return link, err
}
