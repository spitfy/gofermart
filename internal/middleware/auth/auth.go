package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

var (
	ErrUnAuth = errors.New("unauthorized")
)

type AuthManager struct {
	secretKey []byte
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func New(secret string) *AuthManager {
	return &AuthManager{secretKey: []byte(secret)}
}

func (a *AuthManager) GetTokenFromCookie(r *http.Request) (string, error) {
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

func (a *AuthManager) BuildJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})
	return token.SignedString(a.secretKey)
}

func (a *AuthManager) ParseUserID(tokenStr string) (int, error) {
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

func (a *AuthManager) CreateToken(w http.ResponseWriter, userID int) (string, error) {
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
