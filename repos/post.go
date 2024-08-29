package repos

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/types"
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
	return nil, ErrNotFound
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
	return nil, ErrNotFound
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
	q := `
		SELECT p.id, p.title, p.content, p.created, p.pinned, p.tweet, COUNT(c) comment_count 
		FROM posts p LEFT JOIN comments c ON c.postid = p.id GROUP BY p.id ORDER BY p.pinned, p.created DESC;
	`
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return posts, ErrNotFound
		}
		return posts, errors.Join(err, ErrInternalFailure)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id            string
			title         string
			content       string
			created       string
			pinned        bool
			tweet         bool
			comment_count int
		)
		rows.Scan(&id, &title, &content, &created, &pinned, &tweet, &comment_count)
		post := types.Post{
			Id:       id,
			Title:    title,
			Content:  content,
			Created:  created,
			Pinned:   pinned,
			Tweet:    tweet,
			Comments: comment_count,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *postPostgres) FindById(ctx context.Context, id string) (*types.Post, error) {
	var post types.Post
	q := `
		SELECT p.id, p.title, p.content, p.created, p.pinned, p.tweet, COUNT(c) comment_count 
		FROM posts p LEFT JOIN comments c ON c.postid = p.id GROUP BY p.id HAVING p.id = $1;
	`
	err := repo.db.QueryRowContext(ctx, q, id).Scan(
		&post.Id,
		&post.Title,
		&post.Content,
		&post.Created,
		&post.Pinned,
		&post.Tweet,
		&post.Comments,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &post, nil
}

func (repo *postPostgres) Create(ctx context.Context, post types.Post) (*types.Post, error) {
	q := "INSERT INTO posts (id, title, content, created, tweet) VALUES ($1, $2, $3, $4, $5);"
	_, err := repo.db.ExecContext(ctx, q, post.Id, post.Title, post.Content, post.Created, post.Tweet)
	if err != nil {
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &post, nil
}

func (repo *postPostgres) Update(ctx context.Context, id string, post types.Post) (*types.Post, error) {
	q := "UPDATE posts SET title=$1, content=$2, pinned=$3 WHERE id=$4;"
	_, err := repo.db.ExecContext(ctx, q, post.Title, post.Content, post.Pinned, id)
	if err != nil {
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &post, nil
}

func (repo *postPostgres) Delete(ctx context.Context, id string) (*types.Post, error) {
	post, err := repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}

	q := "DELETE FROM posts WHERE id=$1;"
	_, err = repo.db.ExecContext(ctx, q, id)
	if err != nil {
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return post, nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return &types.Page[types.Post]{
				Content:  posts,
				HasNext:  false,
				Size:     size,
				NextPage: 1,
				Total:    0,
			}, nil
		}
		return &types.Page[types.Post]{
			Content:  posts,
			HasNext:  false,
			Size:     size,
			NextPage: 1,
			Total:    0,
		}, errors.Join(err, ErrInternalFailure)
	}
	defer rows.Close()
	for rows.Next() {
		post := types.Post{}
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.Created, &post.Pinned, &post.Tweet)
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
	page, size int,
) (*types.Page[types.Post], error) {
	posts := []types.Post{}
	qcount := "SELECT COUNT(*) FROM posts WHERE created <= $1;"
	var count int
	repo.db.QueryRowContext(ctx, qcount, time).Scan(&count)
	hasNext := true
	if (page-1)*size+size >= count {
		hasNext = false
	}
	q := `
		SELECT p.id, p.title, p.content, p.created, p.pinned, p.tweet, COUNT(c) comment_count 
		FROM posts p LEFT JOIN comments c ON c.postid = p.id GROUP BY p.id HAVING p.created <= $3 
		ORDER BY pinned DESC, created DESC LIMIT $1 OFFSET $2;
	`
	//q := "SELECT * FROM posts WHERE created <= $3 ORDER BY pinned DESC, created DESC LIMIT $1 OFFSET $2;"
	rows, err := repo.db.QueryContext(ctx, q, size, (page-1)*size, time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &types.Page[types.Post]{
				Content:  posts,
				HasNext:  false,
				Size:     size,
				NextPage: 1,
				Total:    0,
			}, nil
		}
		return &types.Page[types.Post]{
			Content:  posts,
			HasNext:  false,
			Size:     size,
			NextPage: 1,
			Total:    0,
		}, errors.Join(err, ErrInternalFailure)
	}
	defer rows.Close()
	for rows.Next() {
		post := types.Post{}
		rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.Created,
			&post.Pinned,
			&post.Tweet,
			&post.Comments,
		)
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

type postSearcherPostgres struct {
	db *sql.DB
}

func NewPostSearcherPostgres() PostSearcher {
	repo := new(postSearcherPostgres)
	repo.db = data.Postgres()
	return repo
}

func (repo *postSearcherPostgres) Search(ctx context.Context, q string, page, size int) (*types.Page[types.Post], error) {
	q = strings.ToLower(q)
	qcount := "SELECT COUNT(*) FROM posts WHERE LOWER(title) LIKE '%' || $1 || '%';"
	var count int
	repo.db.QueryRowContext(ctx, qcount, q).Scan(&count)
	hasNext := true
	if (page-1)*size+size >= count {
		hasNext = false
	}
	posts := []types.Post{}
	sqlq := `
		SELECT p.id, p.title, p.content, p.created, p.pinned, p.tweet, COUNT(c) comment_count 
		FROM posts p LEFT JOIN comments c ON c.postid = p.id GROUP BY p.id 
		HAVING LOWER(p.title) LIKE '%' || $1 || '%' ORDER BY p.pinned DESC, p.created DESC OFFSET $2 LIMIT $3;
	`
	//sqlq := "SELECT * FROM posts WHERE LOWER(title) LIKE '%' || $1 || '%' ORDER BY pinned DESC, created DESC OFFSET $2 LIMIT $3;"
	rows, err := repo.db.QueryContext(ctx, sqlq, q, (page-1)*size, size)
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
		rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.Created,
			&post.Pinned,
			&post.Tweet,
			&post.Comments,
		)
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
func (repo *postSearcherPostgres) Append(ctx context.Context, p types.Post) error {
	return nil
}
func (repo *postSearcherPostgres) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *postSearcherPostgres) Bulk(ctx context.Context, p ...types.Post) error {
	return nil
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
