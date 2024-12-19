package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key untuk signing dan verifying token
var secretKey = []byte("mysecretkey")

// Claims adalah struct custom untuk payload JWT
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// CreateToken untuk membuat JWT dari user ID dan username
func CreateToken(userID int, username string) (string, error) {
	// Atur klaim (payload)
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "myapp",                                            // Pengeluarnya (issuer)
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid selama 24 jam
		},
	}

	// Membuat token JWT dengan algoritma HMAC dan klaim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secretKey
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}
	return tokenString, nil
}

// VerifikasiToken memverifikasi token JWT dan mengembalikan klaim (claims)
func VerifyToken(tokenString string) (*Claims, error) {
	// Memparsing dan memverifikasi token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Memverifikasi bahwa metode signing adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		// Mengembalikan secret key untuk verifikasi
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Mengecek apakah klaim token valid
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
