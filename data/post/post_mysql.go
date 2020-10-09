package post

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ptflp/go-chi-rest/data"

	models "github.com/ptflp/go-chi-rest"
)

// NewPostService retunrs implement of post data interface
func NewPostService(Conn *sqlx.DB) data.PostRepo {
	return &postService{
		db: Conn,
	}
}

type postService struct {
	db *sqlx.DB
}

func (m *postService) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Post, error) {
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*models.Post, 0)
	for rows.Next() {
		post := new(models.Post)

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, post)
	}
	return payload, nil
}

func (m *postService) Fetch(ctx context.Context, num int64) ([]*models.Post, error) {
	query := "Select id, title, content From posts limit ?"

	return m.fetch(ctx, query, num)
}

func (m *postService) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	query := "Select id, title, content From posts where id=?"

	rows, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &models.Post{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, models.ErrNotFound
	}

	return payload, nil
}

func (m *postService) Create(ctx context.Context, p *models.Post) (int64, error) {
	query := "Insert posts SET title=?, content=?"

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, p.Title, p.Content)
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (m *postService) Update(ctx context.Context, p *models.Post) (*models.Post, error) {
	query := "Update posts set title=?, content=? where id=?"

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		p.Title,
		p.Content,
		p.ID,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return p, nil
}

func (m *postService) Delete(ctx context.Context, id int64) (bool, error) {
	query := "Delete From posts Where id=?"

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
