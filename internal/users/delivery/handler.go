package delivery

import (
	"log"
	"net/http"

	"backend/internal/users"
	"backend/internal/users/usecase"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.usecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser meng-handle permintaan untuk membuat pengguna baru
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req users.Pengguna

	// Binding JSON request ke struct Pengguna
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err) // Log error saat binding JSON
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash password sebelum menyimpannya
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Println("Error hashing password:", err) // Log error saat hashing password
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Ganti password yang ada dengan password yang di-hash
	req.Password = hashedPassword

	// Menyimpan user menggunakan usecase
	id, err := h.usecase.CreateUser(req)
	if err != nil {
		log.Println("Error creating user:", err) // Log error saat membuat user
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "id": id})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var req struct {
		ID int `json:"id"` // Ambil id dari JSON payload
	}

	// Binding JSON request ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Cek apakah ID kosong
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Cek apakah pengguna ada
	user, err := h.usecase.GetUserByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Panggil usecase untuk menghapus pengguna berdasarkan ID
	err = h.usecase.DeleteUser(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Masukkan token ke dalam blacklist setelah penghapusan pengguna
	tokenString := c.GetHeader("Authorization")
	// Menambahkan token ke dalam blacklist di database (atau memcached)
	err = h.usecase.BlacklistToken(tokenString)
	if err != nil {
		log.Printf("Failed to blacklist token: %v", err)
	}

	// Jika berhasil
	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
		"id":      req.ID,
		"user":    user,
	})
}

// Login meng-handle permintaan login untuk mendapatkan JWT
func (h *UserHandler) Login(c *gin.Context) {
	var req users.Pengguna
	// Binding JSON request ke struct Pengguna
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Mencari pengguna berdasarkan email
	user, err := h.usecase.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not foun"})
		return
	}

	// Verifikasi password (bandingkan dengan password yang di-hash)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Membuat token JWT setelah verifikasi berhasil
	token, err := utils.CreateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Kembalikan token ke pengguna
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func HashPassword(password string) (string, error) {
	// Debugging: Pastikan password yang diterima benar
	log.Println("Hashing password:", password)

	// Menggunakan bcrypt untuk hashing password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating hash:", err) // Log error
		return "", err
	}
	return string(bytes), nil
}
