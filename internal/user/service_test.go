package user

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/uuid"
	"github.com/jackc/pgx/v4"
)

type mockRepository struct {
	items []entity.User
}

func (m *mockRepository) Get(ctx context.Context, id uuid.UUID) (entity.GetUserRow, error) {
	if reflect.DeepEqual(id, uuid.UUID{}) {
		return entity.GetUserRow{}, pgx.ErrNoRows
	}

	for _, item := range m.items {
		if item.ID == id {
			u := entity.GetUserRow{
				ID:       item.ID,
				Username: item.Username,
			}

			return u, nil
		}
	}

	return entity.GetUserRow{}, pgx.ErrNoRows
}

func (m *mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m *mockRepository) Query(ctx context.Context, args entity.ListUserParams) ([]entity.ListUserRow, error) {
	users := []entity.ListUserRow{}
	for _, v := range m.items {
		users = append(users, entity.ListUserRow{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	return users, nil
}

func (m *mockRepository) Create(ctx context.Context, args entity.CreateUserParams) (entity.CreateUserRow, error) {
	if args.Username == "error" {
		return entity.CreateUserRow{}, constants.ErrCRUD
	}

	id, err := tools.GenerateUUIDv4()
	for err != nil {
		id, err = tools.GenerateUUIDv4()
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	row := entity.CreateUserRow{
		ID:       id,
		Username: args.Username,
	}

	m.items = append(m.items, entity.User{
		ID:        id,
		Username:  args.Username,
		Password:  args.Password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})

	return row, nil
}

func (m *mockRepository) Update(ctx context.Context, args entity.UpdateUserParams) (entity.UpdateUserRow, error) {
	if args.Username == "error" {
		return entity.UpdateUserRow{}, constants.ErrCRUD
	}

	var row entity.UpdateUserRow

	for i, item := range m.items {
		if item.ID == args.ID {
			if args.Username != "" {
				m.items[i].Username = args.Username
			}

			if args.Password != "" {
				m.items[i].Password = args.Password
			}

			m.items[i].UpdatedAt = time.Now()

			row.ID = m.items[i].ID
			row.Username = m.items[i].Username

			break
		}
	}

	return row, nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) (entity.SoftDeleteUserRow, error) {
	if reflect.DeepEqual(id, uuid.UUID{}) {
		return entity.SoftDeleteUserRow{}, constants.ErrCRUD
	}

	var row entity.SoftDeleteUserRow

	for i, item := range m.items {
		if item.ID == id {
			m.items[i].DeletedAt = sql.NullTime{Time: time.Now(), Valid: true}

			row.ID = m.items[i].ID
			row.Username = m.items[i].Username

			break
		}
	}

	return row, nil
}
