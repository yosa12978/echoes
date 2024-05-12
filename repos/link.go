package repos

import (
	"context"
	"database/sql"

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
	q := "SELECT * FROM links ORDER BY created DESC;"
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return links, err
	}
	defer rows.Close()
	for rows.Next() {
		link := types.Link{}
		rows.Scan(&link.Id, &link.Name, &link.URL, &link.Created, &link.Icon)
		links = append(links, link)
	}
	return links, nil
}

func (repo *linkPostgres) FindById(ctx context.Context, id string) (*types.Link, error) {
	var link types.Link
	q := "SELECT * FROM links WHERE id=$1;"
	err := repo.db.QueryRowContext(ctx, q, id).Scan(&link.Id, &link.Name, &link.URL, &link.Created, &link.Icon)
	return &link, err
}

func (repo *linkPostgres) Create(ctx context.Context, link types.Link) (*types.Link, error) {
	q := "INSERT INTO links (id, name, url, created, icon) VALUES ($1, $2, $3, $4, $5);"
	_, err := repo.db.ExecContext(ctx, q, link.Id, link.Name, link.URL, link.Created, link.Icon)
	return &link, err
}

func (repo *linkPostgres) Update(ctx context.Context, id string, link types.Link) (*types.Link, error) {
	q := "UPDATE links SET name=$1, url=$2, created=$3 icon=$4 WHERE id=$5;"
	_, err := repo.db.ExecContext(ctx, q, link.Name, link.URL, link.Created, link.Icon, id)
	return &link, err
}

func (repo *linkPostgres) Delete(ctx context.Context, id string) (*types.Link, error) {
	link, err := repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	q := "DELETE FROM links WHERE id=$1;"
	_, err = repo.db.ExecContext(ctx, q, id)
	return link, err
}
