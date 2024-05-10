package repos

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/types"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type Post interface {
	FindAll(ctx context.Context) ([]types.Post, error)
	GetPage(ctx context.Context, page, size int) (*types.Page[types.Post], error)
	FindById(ctx context.Context, id string) (*types.Post, error)
	Create(ctx context.Context, post types.Post) (*types.Post, error)
	Update(ctx context.Context, id string, post types.Post) (*types.Post, error)
	Delete(ctx context.Context, id string) (*types.Post, error)
	GetPageTime(ctx context.Context, time string, page, size int) (*types.Page[types.Post], error)
	Search(ctx context.Context, query string, page, size int) (*types.Page[types.Post], error)
}

type postMock struct {
	posts []types.Post
}

func NewPostMock() Post {
	return new(postMock)
}

func (repo *postMock) GetPage(ctx context.Context, page, size int) (*types.Page[types.Post], error) {
	return nil, nil
}

func (repo *postMock) FindAll(ctx context.Context) ([]types.Post, error) {
	return repo.posts, nil
}

func (repo *postMock) FindById(ctx context.Context, id string) (*types.Post, error) {
	for i := 0; i < len(repo.posts); i++ {
		if repo.posts[i].Id == id {
			return &repo.posts[i], nil
		}
	}
	return nil, ErrPostNotFound
}

func (repo *postMock) GetPageTime(ctx context.Context, time string, page, size int) (*types.Page[types.Post], error) {
	return nil, nil
}

func (repo *postMock) Create(ctx context.Context, post types.Post) (*types.Post, error) {
	repo.posts = append(repo.posts, post)
	return &post, nil
}

func (repo *postMock) Update(ctx context.Context, id string, post types.Post) (*types.Post, error) {
	return nil, nil
}

func (repo *postMock) Delete(ctx context.Context, id string) (*types.Post, error) {
	for i := 0; i < len(repo.posts); i++ {
		if repo.posts[i].Id == id {
			temp := repo.posts[i]
			repo.posts = append(repo.posts[:i], repo.posts[i+1:]...)
			return &temp, nil
		}
	}
	return nil, ErrPostNotFound
}
func (repo *postMock) Search(ctx context.Context, query string, page, size int) (*types.Page[types.Post], error) {
	return nil, nil
}

type postPostgres struct {
	db *sql.DB
}

func NewPostPostgres() Post {
	repo := new(postPostgres)
	repo.db = data.Postgres()
	return repo
}

func (repo *postPostgres) FindAll(ctx context.Context) ([]types.Post, error) {
	posts := []types.Post{}
	q := "SELECT id, title, content, created FROM posts ORDER BY pinned, created DESC;"
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id      string
			title   string
			content string
			created string
		)
		rows.Scan(&id, &title, &content, &created)
		post := types.Post{
			Id:      id,
			Title:   title,
			Content: content,
			Created: created,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *postPostgres) FindById(ctx context.Context, id string) (*types.Post, error) {
	var post types.Post
	q := "SELECT * FROM posts WHERE id=$1;"
	err := repo.db.QueryRowContext(ctx, q, id).Scan(
		&post.Id,
		&post.Title,
		&post.Content,
		&post.Created,
		&post.Pinned,
	)
	return &post, err
}

func (repo *postPostgres) Create(ctx context.Context, post types.Post) (*types.Post, error) {
	q := "INSERT INTO posts (id, title, content, created) VALUES ($1, $2, $3, $4);"
	_, err := repo.db.ExecContext(ctx, q, post.Id, post.Title, post.Content, post.Created)
	return &post, err
}

func (repo *postPostgres) Update(ctx context.Context, id string, post types.Post) (*types.Post, error) {
	q := "UPDATE posts SET title=$1, content=$2, pinned=$3 WHERE id=$4;"
	_, err := repo.db.ExecContext(ctx, q, post.Title, post.Content, post.Pinned, id)
	return &post, err
}

func (repo *postPostgres) Delete(ctx context.Context, id string) (*types.Post, error) {
	post, err := repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	q := "DELETE FROM posts WHERE id=$1;"
	_, err = repo.db.ExecContext(ctx, q, id)
	return post, err
}

func (repo *postPostgres) GetPage(ctx context.Context, page, size int) (*types.Page[types.Post], error) {
	posts := []types.Post{}
	qcount := "SELECT COUNT(*) FROM posts;"
	var count int
	repo.db.QueryRowContext(ctx, qcount).Scan(&count)
	hasNext := true
	if (page-1)*size+size >= count {
		hasNext = false
	}
	q := "SELECT * FROM posts ORDER BY pinned DESC, created DESC LIMIT $1 OFFSET $2;"
	rows, err := repo.db.QueryContext(ctx, q, size, (page-1)*size)
	if err != nil {
		return &types.Page[types.Post]{
			Content:  posts,
			HasNext:  false,
			Size:     size,
			NextPage: 1,
			Total:    0,
		}, err
	}
	defer rows.Close()
	for rows.Next() {
		post := types.Post{}
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.Created, &post.Pinned)
		posts = append(posts, post)
	}
	return &types.Page[types.Post]{
		Content:  posts,
		HasNext:  hasNext,
		Size:     size,
		NextPage: page + 1,
		Total:    count,
	}, nil
}

func (repo *postPostgres) GetPageTime(
	ctx context.Context,
	time string,
	page, size int) (*types.Page[types.Post], error) {
	posts := []types.Post{}
	qcount := "SELECT COUNT(*) FROM posts WHERE created <= $1;"
	var count int
	repo.db.QueryRowContext(ctx, qcount, time).Scan(&count)
	hasNext := true
	if (page-1)*size+size >= count {
		hasNext = false
	}
	q := "SELECT * FROM posts WHERE created <= $3 ORDER BY pinned DESC, created DESC LIMIT $1 OFFSET $2;"
	rows, err := repo.db.QueryContext(ctx, q, size, (page-1)*size, time)
	if err != nil {
		return &types.Page[types.Post]{
			Content:  posts,
			HasNext:  false,
			Size:     size,
			NextPage: 1,
			Total:    0,
		}, err
	}
	defer rows.Close()
	for rows.Next() {
		post := types.Post{}
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.Created, &post.Pinned)
		posts = append(posts, post)
	}
	return &types.Page[types.Post]{
		Content:  posts,
		HasNext:  hasNext,
		Size:     size,
		NextPage: page + 1,
		Total:    count,
	}, nil
}

func (repo *postPostgres) Search(ctx context.Context, query string, page, size int) (*types.Page[types.Post], error) {
	return nil, nil
}

type PostSearcher interface {
	Search(ctx context.Context, q string, page, size int) (*types.Page[types.Post], error)
	Append(ctx context.Context, p types.Post) error
	Delete(ctx context.Context, id string) error
	Bulk(ctx context.Context, p ...types.Post) error
}

type postES struct {
	es *elasticsearch.Client
}

func NewPostSearcherES(es *elasticsearch.Client) PostSearcher {
	return &postES{
		es: es,
	}
}

func (repo *postES) Search(ctx context.Context, q string, page, size int) (*types.Page[types.Post], error) {
	skip := (page - 1) * size
	req := esapi.SearchRequest{
		Index: []string{"posts"},
		Query: q,
		Size:  &size,
		From:  &skip,
	}
	resp, err := req.Do(ctx, repo.es)
	if err != nil {
		return nil, err
	}
	// i don't think this will work (it must be sone elasticsearch response type instead of []types.Post)
	var posts []types.Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}
	pageRes := types.Page[types.Post]{
		HasNext:  true,
		Size:     size,
		NextPage: page + 1,
		Content:  posts,
		Total:    0, // todo complete
	}

	return &pageRes, err
}

func (repo *postES) Append(ctx context.Context, p types.Post) error {
	data, _ := json.Marshal(p)
	req := esapi.IndexRequest{
		Index:      "posts",
		DocumentID: p.Id,
		Body:       bytes.NewReader(data),
	}
	_, err := req.Do(ctx, repo.es)
	return err
}

func (repo *postES) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{Index: "posts", DocumentID: id}
	_, err := req.Do(ctx, repo.es)
	return err
}

func (repo *postES) Bulk(ctx context.Context, p ...types.Post) error {
	data, _ := json.Marshal(p)
	req := esapi.BulkRequest{
		Index: "posts",
		Body:  bytes.NewReader(data),
	}
	_, err := req.Do(ctx, repo.es)
	return err
}
