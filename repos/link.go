package repos

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/types"
)

type Link interface {
	FindAll(ctx context.Context) []types.Link
	FindById(ctx context.Context, id string) (*types.Link, error)
	Create(ctx context.Context, link types.Link) (*types.Link, error)
	Update(ctx context.Context, id string, link types.Link) (*types.Link, error)
	Delete(ctx context.Context, id string) (*types.Link, error)
	Seed(ctx context.Context) error
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

func (repo *linkPostgres) FindAll(ctx context.Context) []types.Link {
	links := []types.Link{}
	q := "SELECT * FROM links ORDER BY created DESC;"
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return links
	}
	defer rows.Close()
	for rows.Next() {
		link := types.Link{}
		rows.Scan(&link.Id, &link.Name, &link.URL, &link.Created)
		links = append(links, link)
	}
	return links
}

func (repo *linkPostgres) FindById(ctx context.Context, id string) (*types.Link, error) {
	var link types.Link
	q := "SELECT * FROM links WHERE id=$1;"
	err := repo.db.QueryRowContext(ctx, q, id).Scan(&link.Id, &link.Name, &link.URL, &link.Created)
	return &link, err
}

func (repo *linkPostgres) Create(ctx context.Context, link types.Link) (*types.Link, error) {
	id := uuid.NewString()
	created := time.Now().Format(time.RFC3339)
	link.Id = id
	link.Created = created
	q := "INSERT INTO links (id, name, url, created) VALUES ($1, $2, $3, $4);"
	_, err := repo.db.ExecContext(ctx, q, link.Id, link.Name, link.URL, link.Created)
	return &link, err
}

func (repo *linkPostgres) Update(ctx context.Context, id string, link types.Link) (*types.Link, error) {
	q := "UPDATE links SET name=$1, url=$2, created=$3 WHERE id=$4;"
	_, err := repo.db.ExecContext(ctx, q, link.Name, link.URL, link.Created, id)
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

func (repo *linkPostgres) Seed(ctx context.Context) error {
	_, err := repo.Create(ctx, types.Link{
		Id:      "09741221-7ea7-4106-ac19-8d2c2c90afbc",
		Name:    "reddit",
		URL:     "https://reddit.com",
		Created: time.Now().Format(time.RFC3339),
	})
	_, err = repo.Create(ctx, types.Link{
		Id:      "c46428bd-a807-4042-812b-f3b56f047732",
		Name:    "my github",
		URL:     "https://github.com/yosa12978",
		Created: time.Now().Format(time.RFC3339),
	})
	_, err = repo.Create(ctx, types.Link{
		Id:      "60a9f6e8-8fda-480a-832a-3e3a07ae8890",
		Name:    "wow forum (icy veins)",
		URL:     "https://www.icy-veins.com/",
		Created: time.Now().Format(time.RFC3339),
	})
	return err
}
