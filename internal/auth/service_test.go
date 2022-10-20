package auth

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/jackc/pgx/v4"
)

type mockRepository struct {
	items []models.User
}

func (m *mockRepository) GetUserByUsername(ctx context.Context, username string) (models.GetByUsernameRow, error) {
	for _, item := range m.items {
		if item.Username == username {
			u := models.GetByUsernameRow{
				ID:       item.ID,
				Username: item.Username,
				Password: item.Password,
			}

			return u, nil
		}
	}

	return models.GetByUsernameRow{}, pgx.ErrNoRows
}
