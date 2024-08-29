package repos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/types"
)

type Link interface {
	FindAll(ctx context.Context) ([]types.Link, error)
	FindById(ctx context.Context, id string) (*types.Link, error)
	Create(ctx context.Context, link types.Link) (*types.Link, error)
	Update(ctx context.Context, id string, link types.Link) (*types.Link, error)
	Delete(ctx context.Context, id string) (*types.Link, error)
}

type linkPostgres struct {
	db *sql.DB
}

/*

type Link struct {
	Id      string
	Name    string
	URL     string
	Created string
}

*/

func NewLinkPostgres() Link {
	repo := new(linkPostgres)
	repo.db = data.Postgres()
	return repo
}

func (repo *linkPostgres) FindAll(ctx context.Context) ([]types.Link, error) {
	links := []types.Link{}
	q := "SELECT * FROM links ORDER BY place ASC"
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return links, nil
		}
		return links, errors.Join(err, ErrInternalFailure)
	}
	defer rows.Close()
	for rows.Next() {
		link := types.Link{}
		rows.Scan(&link.Id, &link.Name, &link.URL, &link.Created, &link.Icon, &link.Place)
		links = append(links, link)
	}
	return links, nil
}

func (repo *linkPostgres) FindById(ctx context.Context, id string) (*types.Link, error) {
	var link types.Link
	q := "SELECT * FROM links WHERE id=$1;"
	err := repo.db.QueryRowContext(ctx, q, id).
		Scan(
			&link.Id,
			&link.Name,
			&link.URL,
			&link.Created,
			&link.Icon,
			&link.Place,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &link, nil
}

func (repo *linkPostgres) Create(ctx context.Context, link types.Link) (*types.Link, error) {
	q := "INSERT INTO links (id, name, url, created, icon, place) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := repo.db.ExecContext(ctx, q, link.Id, link.Name, link.URL, link.Created, link.Icon, link.Place)
	if err != nil {
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &link, nil
}

func (repo *linkPostgres) Update(ctx context.Context, id string, link types.Link) (*types.Link, error) {
	q := "UPDATE links SET name=$1, url=$2, created=$3, icon=$4, place=$5 WHERE id=$6;"
	_, err := repo.db.ExecContext(ctx, q, link.Name, link.URL, link.Created, link.Icon, link.Place, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &link, nil
}

func (repo *linkPostgres) Delete(ctx context.Context, id string) (*types.Link, error) {
	link, err := repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	q := "DELETE FROM links WHERE id=$1;"
	_, err = repo.db.ExecContext(ctx, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return link, nil
}
