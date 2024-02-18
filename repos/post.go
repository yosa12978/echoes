package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/types"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type Post interface {
	FindAll(ctx context.Context) []types.Post
	GetPage(ctx context.Context, page, size int) *types.Page[types.Post]
	FindById(ctx context.Context, id string) (*types.Post, error)
	Create(ctx context.Context, post types.Post) (*types.Post, error)
	Update(ctx context.Context, id string, post types.Post) (*types.Post, error)
	Delete(ctx context.Context, id string) (*types.Post, error)
	Seed(ctx context.Context) error
}

type postMock struct {
	posts []types.Post
}

func NewPostMock() Post {
	return new(postMock)
}

func (repo *postMock) GetPage(ctx context.Context, page, size int) *types.Page[types.Post] {
	return nil
}

func (repo *postMock) FindAll(ctx context.Context) []types.Post {
	return repo.posts
}

func (repo *postMock) FindById(ctx context.Context, id string) (*types.Post, error) {
	for i := 0; i < len(repo.posts); i++ {
		if repo.posts[i].Id == id {
			return &repo.posts[i], nil
		}
	}
	return nil, ErrPostNotFound
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

func (repo *postMock) Seed(ctx context.Context) error {
	posts := []types.Post{
		types.NewPost("first post", "first post content"),
		types.NewPost("second post", "second post content"),
		types.NewPost("third post", "third post content"),
		types.NewPost("fourth post", "fourth post content"),
	}
	repo.posts = append(repo.posts, posts...)
	return nil
}

type postPostgres struct {
	db *sql.DB
}

func NewPostPostgres() Post {
	repo := new(postPostgres)
	repo.db = data.Postgres()
	return repo
}

func (repo *postPostgres) FindAll(ctx context.Context) []types.Post {
	posts := []types.Post{}
	q := "SELECT id, title, content, created FROM posts ORDER BY pinned, created DESC;"
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		log.Println(err.Error())
		return posts
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
	return posts
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

	// move this block to service layer
	id := uuid.NewString()
	created := time.Now().Format(time.RFC3339)
	post.Id = id
	post.Created = created

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

	// move this block to service layer
	post, err := repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	q := "DELETE FROM posts WHERE id=$1;"
	_, err = repo.db.ExecContext(ctx, q, id)
	return post, err
}

func (repo *postPostgres) Seed(ctx context.Context) error {
	for i := 0; i < 50; i++ {
		time.Sleep(50 * time.Millisecond)
		_, err := repo.Create(ctx, types.NewPost(fmt.Sprintf("post #%d", i), fmt.Sprintf("post content #%d", i)))
		if err != nil {
			return err
		}
	}
	return nil
}

// test this function (i'm not sure if it works)
func (repo *postPostgres) GetPage(ctx context.Context, page, size int) *types.Page[types.Post] {
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
		log.Println(err.Error())
		return &types.Page[types.Post]{
			Content:  posts,
			HasNext:  false,
			Size:     size,
			NextPage: 1,
			Total:    0,
		}
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
	}
}
