package auth

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	"github.com/jackc/pgx/v4"
)

type mockRepository struct {
	items []entity.User
}

func (m *mockRepository) GetUserByUsername(ctx context.Context, username string) (entity.GetByUsernameRow, error) {
	for _, item := range m.items {
		if item.Username == username {
			u := entity.GetByUsernameRow{
				ID:       item.ID,
				Username: item.Username,
				Password: item.Password,
			}

			return u, nil
		}
	}

	return entity.GetByUsernameRow{}, pgx.ErrNoRows
}
