package user

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/model"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/uuid"
	"github.com/jackc/pgx/v4"
)

type mockRepository struct {
	items []model.User
}

func (m *mockRepository) Get(ctx context.Context, id uuid.UUID) (model.GetUserRow, error) {
	if reflect.DeepEqual(id, uuid.UUID{}) {
		return model.GetUserRow{}, pgx.ErrNoRows
	}

	for _, item := range m.items {
		if item.ID == id {
			u := model.GetUserRow{
				ID:       item.ID,
				Username: item.Username,
			}

			return u, nil
		}
	}

	return model.GetUserRow{}, pgx.ErrNoRows
}

func (m *mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m *mockRepository) Query(ctx context.Context, args model.ListUserParams) ([]model.ListUserRow, error) {
	users := []model.ListUserRow{}
	for _, v := range m.items {
		users = append(users, model.ListUserRow{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	return users, nil
}

func (m *mockRepository) Create(ctx context.Context, args model.CreateUserParams) (model.CreateUserRow, error) {
	if args.Username == "error" {
		return model.CreateUserRow{}, constants.ErrCRUD
	}

	id, err := tools.GenerateUUIDv4()
	for err != nil {
		id, err = tools.GenerateUUIDv4()
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	row := model.CreateUserRow{
		ID:       id,
		Username: args.Username,
	}

	m.items = append(m.items, model.User{
		ID:        id,
		Username:  args.Username,
		Password:  args.Password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})

	return row, nil
}

func (m *mockRepository) Update(ctx context.Context, args model.UpdateUserParams) (model.UpdateUserRow, error) {
	if args.Username == "error" {
		return model.UpdateUserRow{}, constants.ErrCRUD
	}

	var row model.UpdateUserRow

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

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) (model.SoftDeleteUserRow, error) {
	if reflect.DeepEqual(id, uuid.UUID{}) {
		return model.SoftDeleteUserRow{}, constants.ErrCRUD
	}

	var row model.SoftDeleteUserRow

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
