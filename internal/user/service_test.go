package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/uuid"
)

type mockRepository struct {
	items []models.User
}

func (m *mockRepository) Get(ctx context.Context, id uuid.UUID) (models.GetUserRow, error) {
	for _, item := range m.items {
		if item.ID == id {
			u := models.GetUserRow{
				ID:       item.ID,
				Username: item.Username,
			}

			return u, nil
		}
	}

	return models.GetUserRow{}, errors.New("testing")
}

func (m *mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m *mockRepository) Query(ctx context.Context, args models.ListUserParams) ([]models.ListUserRow, error) {
	users := []models.ListUserRow{}
	for _, v := range m.items {
		users = append(users, models.ListUserRow{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	return users, nil
}

func (m *mockRepository) Create(ctx context.Context, args models.CreateUserParams) (models.CreateUserRow, error) {
	if args.Username == "error" {
		return models.CreateUserRow{}, constants.ErrCRUD
	}

	id, err := tools.GenerateUUIDv4()
	for err != nil {
		id, err = tools.GenerateUUIDv4()
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	row := models.CreateUserRow{
		ID:       id,
		Username: args.Username,
	}

	m.items = append(m.items, models.User{
		ID:        id,
		Username:  args.Username,
		Password:  args.Password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})

	return row, nil
}

func (m *mockRepository) Update(ctx context.Context, args models.UpdateUserParams) (models.UpdateUserRow, error) {
	if args.Username == "error" {
		return models.UpdateUserRow{}, constants.ErrCRUD
	}

	var row models.UpdateUserRow

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

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) (models.SoftDeleteUserRow, error) {
	var row models.SoftDeleteUserRow

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
