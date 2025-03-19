package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"time"
	"todo-api/internal/config"
	"todo-api/internal/models"
	"todo-api/internal/repositories"
	"todo-api/internal/utils"

	"github.com/dgrijalva/jwt-go"
)

func MiddlewareAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Missing token")
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := uint(claims["user_id"].(float64))
				ctx := context.WithValue(r.Context(), "user_id", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			}
		})
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = hashedPassword

	repo := repositories.NewUserRepository(r.Context().Value("db").(*gorm.DB))
	if err := repo.CreateUser(&user); err != nil {
		slog.Error("Failed to create user", "error", err)
		if err.Error() == fmt.Sprintf("user with email %s already exists", user.Email) {
			utils.RespondWithError(w, http.StatusConflict, "User with this email already exists")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	repo := repositories.NewUserRepository(r.Context().Value("db").(*gorm.DB))
	user, err := repo.GetUserByEmail(loginData.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !checkPasswordHash(loginData.Password, user.Password) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token := generateToken(user.ID, r.Context().Value("config").(*config.Config).JWTSecret)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func generateToken(userID uint, jwtSecret string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
