package crypto

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imrany/smart_spore_hub/server/database/processes/profile"
	"github.com/spf13/viper"
)

var (
	JWT_SECRET     = viper.GetString("JWT_SECRET")
	JWT_EXPIRATION = viper.GetString("JWT_EXPIRATION")
)

func GenerateToken(userID string) (string, error) {
	expirationTimeStr := JWT_EXPIRATION
	if expirationTimeStr == "" {
		expirationTimeStr = "1h" // Default to 1 hour if empty
	}
	expirationTime, err := time.ParseDuration(expirationTimeStr)
	if err != nil {
		return "", fmt.Errorf("could not parse JWT_EXPIRATION: %w", err)
	}

	claims := jwt.MapClaims{
		"exp":     time.Now().Add(expirationTime).Unix(),
		"user_id": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (bool, error) {
	ctx := context.Background()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if expirationTime.Before(time.Now()) {
			return false, fmt.Errorf("token has expired")
		}

		userID := claims["user_id"].(string)
		if userID, err := profile.GetByID(ctx, userID); err != nil || userID != userID {
			return false, fmt.Errorf("invalid user ID")
		}

		return true, nil
	} else {
		return false, fmt.Errorf("invalid token")
	}
}
