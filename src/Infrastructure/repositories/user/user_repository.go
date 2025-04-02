package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/crypto"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/repositories"
)

const (
	// userTable = "users"
)
type userRepoImpl struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repositories.UserRepository {
    return &userRepoImpl{db: db}
}

func (r *userRepoImpl) ValidateUser(username, password string) (*models.User, error) {
	ctx := context.Background()

	var user models.User
	var hashedPassword string

	err := r.db.QueryRow(ctx,
		`SELECT id, username, password FROM users WHERE username=$1`,
		username,
	).Scan(&user.ID, &user.Username, &hashedPassword)

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Parolni tekshiramiz
	if err := crypto.PasswordMatch(password, hashedPassword); !err {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}

func (r *userRepoImpl) GetByUsername(username string) (*models.User, error) {
    ctx := context.Background()
    var user models.User

    row := r.db.QueryRow(ctx, `
        SELECT id, username, password FROM users WHERE username = $1
    `, username)

    err := row.Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        return nil, err
    }

    return &user, nil
}
func (r *userRepoImpl) GetByID(id string) (*models.User, error) {
    ctx := context.Background()
	var user models.User
	err := r.db.QueryRow(ctx, `SELECT id, username FROM users WHERE id=$1`, id).
		Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepoImpl) Create(user *models.User) error {
    ctx := context.Background()

    _, err := r.db.Exec(ctx, `
        INSERT INTO users (id, username, password) VALUES ($1, $2, $3)
    `, user.ID, user.Username, user.Password)

    return err
}
