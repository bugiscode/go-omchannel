package usecase

import (
	"backend/internal/users"
	"backend/internal/users/repository"
	"errors"
)

type UserUsecase interface {
	GetAllUsers() ([]users.Pengguna, error) // Perhatikan penggunaan singular "User"
	GetUserByID(id int) (*users.Pengguna, error)
	CreateUser(u users.Pengguna) (int, error)
	GetUserByEmail(email string) (*users.Pengguna, error)
	UpdateUser(u users.Pengguna) error
	DeleteUser(id int) error
	BlacklistToken(token string) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetAllUsers() ([]users.Pengguna, error) {
	return u.repo.FetchAll()
}

// func (u *userUsecase) GetUserByID(id int) (*users.Pengguna, error) {
// 	return u.repo.GetByID(id)
// }

func (u *userUsecase) CreateUser(uData users.Pengguna) (int, error) {
	if uData.Username == "" || uData.Email == "" {
		return 0, nil
	}
	return u.repo.Create(uData)
}

func (u *userUsecase) UpdateUser(uData users.Pengguna) error {
	return u.repo.Update(uData)
}

// GetUserByID mencari pengguna berdasarkan ID
func (u *userUsecase) GetUserByID(id int) (*users.Pengguna, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *userUsecase) BlacklistToken(token string) error {
	err := u.repo.SaveBlacklistedToken(token)
	if err != nil {
		return errors.New("failed to blacklist token")
	}
	return nil
}

func (u *userUsecase) DeleteUser(id int) error {
	// Memanggil repository untuk menghapus pengguna berdasarkan ID
	err := u.repo.Delete(id)
	if err != nil {
		return errors.New("failed to delete user")
	}
	return nil
}

// GetUserByEmail mencari pengguna berdasarkan email
func (u *userUsecase) GetUserByEmail(email string) (*users.Pengguna, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
