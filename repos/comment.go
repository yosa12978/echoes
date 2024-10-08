package repos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/types"
)

type Comment interface {
	FindAll(ctx context.Context) ([]types.Comment, error)
	GetPage(ctx context.Context, postId string, page, size int) (*types.Page[types.Comment], error)
	GetPageTime(ctx context.Context, time, postId string, page, size int) (*types.Page[types.Comment], error)
	FindById(ctx context.Context, id string) (*types.Comment, error)
	FindByPostId(ctx context.Context, postId string) ([]types.Comment, error)
	Create(ctx context.Context, comment types.Comment) (*types.Comment, error)
	Update(ctx context.Context, id string, comment types.Comment) (*types.Comment, error)
	Delete(ctx context.Context, id string) (*types.Comment, error)
	GetCommentsCount(ctx context.Context, postId string) (int, error)
}

type commentPostgres struct {
	db *sql.DB
}

func NewCommentPostgres() Comment {
	repo := new(commentPostgres)
	repo.db = data.Postgres()
	return repo
}

/*

type Comment struct {
	Id      string
	Email   string
	Name    string
	Content string
	Created string
}

*/

func (repo *commentPostgres) FindAll(ctx context.Context) ([]types.Comment, error) {
	comments := []types.Comment{}
	q := "SELECT * FROM comments ORDER BY created DESC;"
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return comments, nil
		}
		return nil, types.NewErrInternalFailure(err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment types.Comment
		rows.Scan(
			&comment.Id,
			&comment.Email,
			&comment.Name,
			&comment.Content,
			&comment.Created,
			&comment.PostId,
		)
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repo *commentPostgres) FindById(ctx context.Context, id string) (*types.Comment, error) {
	var comment types.Comment
	q := "SELECT * FROM comments WHERE id=$1;"
	err := repo.db.QueryRowContext(ctx, q, id).
		Scan(&comment.Id,
			&comment.Email,
			&comment.Name,
			&comment.Content,
			&comment.Created,
			&comment.PostId,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrNotFound
		}
		return nil, types.NewErrInternalFailure(err)
	}
	return &comment, nil
}

func (repo *commentPostgres) FindByPostId(ctx context.Context, postId string) ([]types.Comment, error) {
	comments := []types.Comment{}
	q := "SELECT * FROM comments WHERE postid=$1;"
	rows, err := repo.db.QueryContext(ctx, q, postId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return comments, nil
		}
		return nil, types.NewErrInternalFailure(err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment types.Comment
		rows.Scan(
			&comment.Id,
			&comment.Email,
			&comment.Name,
			&comment.Content,
			&comment.Created,
			&comment.PostId,
		)
		comments = append(comments, comment)
	}
	return comments, err
}

func (repo *commentPostgres) Create(ctx context.Context, comment types.Comment) (*types.Comment, error) {
	q := "INSERT INTO comments (id, email, name, content, created, postid) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := repo.db.ExecContext(ctx, q,
		comment.Id,
		comment.Email,
		comment.Name,
		comment.Content,
		comment.Created,
		comment.PostId,
	)
	if err != nil {
		return nil, types.NewErrInternalFailure(err)
	}
	return &comment, err
}

func (repo *commentPostgres) Update(ctx context.Context, id string, comment types.Comment) (*types.Comment, error) {
	q := "UPDATE comments SET email=$1, name=$2, content=$3 WHERE id=$4;"
	_, err := repo.db.ExecContext(ctx, q, comment.Email, comment.Name, comment.Content, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrNotFound
		}
		return nil, types.NewErrInternalFailure(err)
	}
	return &comment, err
}

func (repo *commentPostgres) Delete(ctx context.Context, id string) (*types.Comment, error) {
	// delete this and similar checks
	comment, err := repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrNotFound
		}
		return nil, types.NewErrInternalFailure(err)
	}
	q := "DELETE FROM comments WHERE id=$1;"
	_, err = repo.db.ExecContext(ctx, q, id)
	if err != nil {
		return nil, types.NewErrInternalFailure(err)
	}
	return comment, nil
}

func (repo *commentPostgres) GetPage(ctx context.Context, postId string, page, size int) (*types.Page[types.Comment], error) {
	comments := []types.Comment{}
	qcount := "SELECT COUNT(*) FROM comments WHERE postId=$1;"
	var count int
	repo.db.QueryRowContext(ctx, qcount, postId).Scan(&count)
	hasNext := true
	if (page-1)*size+size >= count {
		hasNext = false
	}
	q := "SELECT * FROM comments WHERE postId=$1 ORDER BY created DESC LIMIT $2 OFFSET $3;"
	rows, err := repo.db.QueryContext(ctx, q, postId, size, (page-1)*size)
	if err != nil {
		return &types.Page[types.Comment]{
			Content:  comments,
			HasNext:  false,
			Size:     size,
			NextPage: 1,
			Total:    0,
		}, err
	}
	defer rows.Close()
	for rows.Next() {
		comment := types.Comment{}
		rows.Scan(
			&comment.Id,
			&comment.Email,
			&comment.Name,
			&comment.Content,
			&comment.Created,
			&comment.PostId,
		)
		comments = append(comments, comment)
	}
	return &types.Page[types.Comment]{
		Content:  comments,
		HasNext:  hasNext,
		Size:     size,
		NextPage: page + 1,
		Total:    count,
	}, nil
}

func (repo *commentPostgres) GetCommentsCount(ctx context.Context, postId string) (int, error) {
	q := "SELECT COUNT(*) FROM comments WHERE postId=$1;"
	var count int
	err := repo.db.QueryRowContext(ctx, q, postId).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, types.ErrNotFound
		}
		return 0, types.NewErrInternalFailure(err)
	}
	return count, nil
}

func (repo *commentPostgres) GetPageTime(
	ctx context.Context, time, postId string, page, size int) (*types.Page[types.Comment], error) {
	comments := []types.Comment{}
	qcount := "SELECT COUNT(*) FROM comments WHERE postId=$1 AND created <= $2;"
	var count int
	repo.db.QueryRowContext(ctx, qcount, postId, time).Scan(&count)
	hasNext := true
	if (page-1)*size+size >= count {
		hasNext = false
	}
	q := "SELECT * FROM comments WHERE postId=$1 AND created <= $4 ORDER BY created DESC LIMIT $2 OFFSET $3;"
	rows, err := repo.db.QueryContext(ctx, q, postId, size, (page-1)*size, time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &types.Page[types.Comment]{
				Content:  comments,
				HasNext:  false,
				Size:     size,
				NextPage: 1,
				Total:    0,
			}, nil
		}
		return &types.Page[types.Comment]{
			Content:  comments,
			HasNext:  false,
			Size:     size,
			NextPage: 1,
			Total:    0,
		}, types.NewErrInternalFailure(err)
	}
	defer rows.Close()
	for rows.Next() {
		comment := types.Comment{}
		rows.Scan(
			&comment.Id,
			&comment.Email,
			&comment.Name,
			&comment.Content,
			&comment.Created,
			&comment.PostId,
		)
		comments = append(comments, comment)
	}
	return &types.Page[types.Comment]{
		Content:  comments,
		HasNext:  hasNext,
		Size:     size,
		NextPage: page + 1,
		Total:    count,
	}, nil
}
