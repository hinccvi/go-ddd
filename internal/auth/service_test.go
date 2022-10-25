package auth

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/model"
	"github.com/jackc/pgx/v4"
)

type mockRepository struct {
	items []model.User
}

func (m *mockRepository) GetUserByUsername(ctx context.Context, username string) (model.GetByUsernameRow, error) {
	for _, item := range m.items {
		if item.Username == username {
			u := model.GetByUsernameRow{
				ID:       item.ID,
				Username: item.Username,
				Password: item.Password,
			}

			return u, nil
		}
	}

	return model.GetByUsernameRow{}, pgx.ErrNoRows
}
