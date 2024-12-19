package repository

import (
	"backend/internal/users"
	"database/sql"
	"fmt"
)

// UserRepository adalah interface untuk repository User
type UserRepository interface {
	FetchAll() ([]users.Pengguna, error) // Menggunakan slice
	GetByID(id int) (*users.Pengguna, error)
	GetByEmail(email string) (*users.Pengguna, error)
	Create(u users.Pengguna) (int, error)
	Update(u users.Pengguna) error
	Delete(id int) error
	SaveBlacklistedToken(token string) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

// SaveBlacklistedToken menyimpan token yang diblacklist
func (r *userRepo) SaveBlacklistedToken(token string) error {
	_, err := r.db.Exec("INSERT INTO blacklisted_tokens (token) VALUES ($1)", token)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	return nil
}

// FetchAll mengambil semua pengguna dari database
func (r *userRepo) FetchAll() ([]users.Pengguna, error) {
	// Menjalankan query untuk mengambil data pengguna
	rows, err := r.db.Query("SELECT user_id, username, email, role_id, client_id, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []users.Pengguna // Menggunakan slice untuk menampung pengguna
	for rows.Next() {
		var u users.Pengguna
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.RoleID, &u.ClientID, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		userList = append(userList, u) // Menambahkan pengguna ke slice
	}
	return userList, nil
}

// Method lainnya tetap sama
func (r *userRepo) GetByID(id int) (*users.Pengguna, error) {
	var u users.Pengguna
	err := r.db.QueryRow("SELECT user_id, username, email, role_id, client_id, created_at FROM users WHERE user_id = $1", id).
		Scan(&u.ID, &u.Username, &u.Email, &u.RoleID, &u.ClientID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Create(u users.Pengguna) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO users (username, email, password_hash, role_id, client_id) VALUES ($1, $2, $3, $4, $5) RETURNING user_id",
		u.Username, u.Email, u.Password, u.RoleID, u.ClientID,
	).Scan(&id)
	return id, err
}

func (r *userRepo) Update(u users.Pengguna) error {
	_, err := r.db.Exec(
		"UPDATE users SET username = $1, email = $2, role_id = $3, client_id = $4 WHERE user_id = $5",
		u.Username, u.Email, u.RoleID, u.ClientID, u.ID,
	)
	return err
}

func (r *userRepo) Delete(id int) error {
	// Query untuk menghapus pengguna berdasarkan ID
	_, err := r.db.Exec("DELETE FROM users WHERE user_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user with id %d: %w", id, err)
	}
	return nil
}

// GetByEmail mencari pengguna berdasarkan email
func (r *userRepo) GetByEmail(email string) (*users.Pengguna, error) {
	var u users.Pengguna
	err := r.db.QueryRow("SELECT user_id, username, email, password_hash, role_id, client_id, created_at FROM users WHERE email = $1", email).
		Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.RoleID, &u.ClientID, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}
