package mocks

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/entity"
)

var ErrCRUD = errors.New("error crud")

type UserRepository struct {
	Items []entity.User
}

func (m *UserRepository) Get(_ context.Context, id uuid.UUID) (entity.User, error) {
	if reflect.DeepEqual(id, uuid.UUID{}) {
		return entity.User{}, ErrCRUD
	}

	for _, item := range m.Items {
		if item.ID == id {
			u := entity.User{
				ID:       item.ID,
				Username: item.Username,
			}

			return u, nil
		}
	}

	return entity.User{}, sql.ErrNoRows
}

func (m *UserRepository) Count(_ context.Context) (int64, error) {
	return int64(len(m.Items)), nil
}

func (m *UserRepository) Query(_ context.Context, page, size int) ([]entity.User, error) {
	if page <= 0 || size <= 0 {
		return []entity.User{}, ErrCRUD
	}

	users := []entity.User{}
	for _, v := range m.Items {
		users = append(users, entity.User{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	return users, nil
}

func (m *UserRepository) Create(_ context.Context, u entity.User) error {
	if u.Username == "error" {
		return ErrCRUD
	}

	id := uuid.New()

	m.Items = append(m.Items, entity.User{
		ID:        id,
		Username:  u.Username,
		Password:  u.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return nil
}

func (m *UserRepository) Update(_ context.Context, u entity.User) error {
	if u.Username == "error" {
		return ErrCRUD
	}

	isFound := false
	for i, item := range m.Items {
		if item.ID == u.ID {
			if u.Username != "" {
				m.Items[i].Username = u.Username
			}

			if u.Password != "" {
				m.Items[i].Password = u.Password
			}

			m.Items[i].UpdatedAt = time.Now()

			isFound = true
			break
		}
	}

	if !isFound {
		return sql.ErrNoRows
	}

	return nil
}

func (m *UserRepository) Delete(_ context.Context, id uuid.UUID) error {
	if reflect.DeepEqual(id, uuid.UUID{}) {
		return ErrCRUD
	}

	isFound := false
	for i, item := range m.Items {
		if item.ID == id {
			m.Items[i].DeletedAt = sql.NullTime{Time: time.Now(), Valid: true}

			isFound = true
			break
		}
	}

	if !isFound {
		return sql.ErrNoRows
	}

	return nil
}
