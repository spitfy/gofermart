package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrUnAuth = errors.New("unauthorized")
)

type Service struct {
	secretKey []byte
}

//go:generate mockgen -destination=servicer_mock.go -package=auth github.com/spitfy/gofermart/internal/middleware/auth Servicer
type Servicer interface {
	GetTokenFromCookie(r *http.Request) (string, error)
	BuildJWT(userID int) (string, error)
	ParseUserID(tokenStr string) (int, error)
	CreateToken(w http.ResponseWriter, userID int) (string, error)
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func New(secret string) *Service {
	return &Service{secretKey: []byte(secret)}
}

func (a *Service) GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("ID")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("%w: cookie not found", ErrUnAuth)
		}
		return "", fmt.Errorf("%w: error reading cookie", ErrUnAuth)
	}
	if cookie.Value == "" {
		return "", fmt.Errorf("%w: empty cookie value", ErrUnAuth)
	}
	return cookie.Value, nil
}

func (a *Service) BuildJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})
	return token.SignedString(a.secretKey)
}

func (a *Service) ParseUserID(tokenStr string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return a.secretKey, nil
	})
	if err != nil || !token.Valid {
		return -1, fmt.Errorf("invalid token: %w", err)
	}
	return claims.UserID, nil
}

func (a *Service) CreateToken(w http.ResponseWriter, userID int) (string, error) {
	tokenString, err := a.BuildJWT(userID)
	if err != nil {
		return "", err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "ID",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
	})
	return tokenString, nil
}
