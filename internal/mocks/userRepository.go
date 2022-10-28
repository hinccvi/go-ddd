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

func (m *UserRepository) Get(ctx context.Context, id uuid.UUID) (entity.GetUserRow, error) {
	if reflect.DeepEqual(id, uuid.UUID{}) {
		return entity.GetUserRow{}, ErrCRUD
	}

	for _, item := range m.Items {
		if item.ID == id {
			u := entity.GetUserRow{
				ID:       item.ID,
				Username: item.Username,
			}

			return u, nil
		}
	}

	return entity.GetUserRow{}, sql.ErrNoRows
}

func (m *UserRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.Items)), nil
}

func (m *UserRepository) Query(ctx context.Context, args entity.ListUserParams) ([]entity.ListUserRow, error) {
	if args.Offset < 0 {
		return []entity.ListUserRow{}, ErrCRUD
	}

	users := []entity.ListUserRow{}
	for _, v := range m.Items {
		users = append(users, entity.ListUserRow{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	return users, nil
}

func (m *UserRepository) Create(ctx context.Context, args entity.CreateUserParams) (entity.CreateUserRow, error) {
	if args.Username == "error" {
		return entity.CreateUserRow{}, ErrCRUD
	}

	id := uuid.New()

	row := entity.CreateUserRow{
		ID:       id,
		Username: args.Username,
	}

	m.Items = append(m.Items, entity.User{
		ID:        id,
		Username:  args.Username,
		Password:  args.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return row, nil
}

func (m *UserRepository) Update(ctx context.Context, args entity.UpdateUserParams) (entity.UpdateUserRow, error) {
	if args.Username == "error" {
		return entity.UpdateUserRow{}, ErrCRUD
	}

	var row entity.UpdateUserRow

	for i, item := range m.Items {
		if item.ID == args.ID {
			if args.Username != "" {
				m.Items[i].Username = args.Username
			}

			if args.Password != "" {
				m.Items[i].Password = args.Password
			}

			m.Items[i].UpdatedAt = time.Now()

			row.ID = m.Items[i].ID
			row.Username = m.Items[i].Username

			break
		}
	}

	if row.Username == "" {
		return entity.UpdateUserRow{}, sql.ErrNoRows
	}

	return row, nil
}

func (m *UserRepository) Delete(ctx context.Context, id uuid.UUID) (entity.SoftDeleteUserRow, error) {
	if reflect.DeepEqual(id, uuid.UUID{}) {
		return entity.SoftDeleteUserRow{}, ErrCRUD
	}

	var row entity.SoftDeleteUserRow

	for i, item := range m.Items {
		if item.ID == id {
			m.Items[i].DeletedAt = sql.NullTime{Time: time.Now(), Valid: true}

			row.ID = m.Items[i].ID
			row.Username = m.Items[i].Username

			break
		}
	}

	if row.Username == "" {
		return entity.SoftDeleteUserRow{}, sql.ErrNoRows
	}

	return row, nil
}
