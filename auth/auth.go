package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	"github.com/satyarth42/chatter/storage"
)

const (
	ACCESS_TOKEN_SECRET  = "ACCESS_TOKEN_SECRET"
	REFRESH_TOKEN_SECRET = "REFRESH_TOKEN_SECRET"
	accessTokenDuration  = 7 * 24 * time.Hour
	refreshTokenDuration = 30 * 24 * time.Hour
	BEARER               = "Bearer"
	REFRESH              = "Refresh"
)

type Claims struct {
	jwt.RegisteredClaims
	TokenType string
}

func getAccessTokenSecretString() []byte {
	secret := os.Getenv(ACCESS_TOKEN_SECRET)
	if secret == "" {
		return []byte("test-secret")
	}
	return []byte(secret)
}
func getRefreshTokenSecretString() []byte {
	secret := os.Getenv(REFRESH_TOKEN_SECRET)
	if secret == "" {
		return []byte("test-secret")
	}
	return []byte(secret)
}

func GetToken(userID string) (string, string, int64, error) {
	currentTime := time.Now()
	claims := &Claims{}
	claims.Issuer = "chatter-be"
	claims.Subject = userID
	claims.Audience = jwt.ClaimStrings{"chatter-app"}
	claims.NotBefore = jwt.NewNumericDate(currentTime)
	claims.IssuedAt = jwt.NewNumericDate(currentTime)

	accessTokenClaims := *claims
	accessTokenClaims.ExpiresAt = jwt.NewNumericDate(currentTime.Add(getDuration(BEARER)))
	accessTokenClaims.ID = uuid.NewString()
	accessTokenClaims.TokenType = BEARER

	refreshTokenClaims := *claims
	refreshTokenClaims.ExpiresAt = jwt.NewNumericDate(currentTime.Add(getDuration(REFRESH)))
	refreshTokenClaims.ID = uuid.NewString()
	refreshTokenClaims.TokenType = REFRESH

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessTokenString, err := accessToken.SignedString(getAccessTokenSecretString())
	if err != nil {
		return "", "", 0, err
	}

	refreshTokenString, err := refreshToken.SignedString(getRefreshTokenSecretString())
	if err != nil {
		return "", "", 0, err
	}

	err = storeTokens(accessTokenString, accessTokenClaims.ExpiresAt.Time, refreshTokenString, refreshTokenClaims.ExpiresAt.Time)
	if err != nil {
		return "", "", 0, err
	}

	return accessTokenString, refreshTokenString, int64(getDuration(BEARER)), nil
}

func getDuration(tokenType string) time.Duration {
	switch tokenType {
	case REFRESH:
		return refreshTokenDuration
	case BEARER:
		return accessTokenDuration
	default:
		return 0
	}
}

func VerifyToken(tokenString, tokenType string) (string, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		switch tokenType {
		case BEARER:
			return getAccessTokenSecretString(), nil
		case REFRESH:
			return getRefreshTokenSecretString(), nil
		default:
			return "", fmt.Errorf("invalid token type")
		}
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return "", false
	}

	claimsData, ok := token.Claims.(*Claims)
	if !ok {
		return "", false
	}

	client := storage.GetRedis()
	res, _ := client.HGet(context.Background(), tokenString, "VALID").Result()
	if res != "1" {
		return "", false
	}

	return claimsData.Subject, true
}

func storeTokens(accessToken string, accessTokenExpiresAt time.Time, refreshToken string, refreshTokenExpiresAt time.Time) error {
	ctx := context.Background()
	client := storage.GetRedis()

	_, err := client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		var res = make([]redis.Cmder, 4)
		res[0] = pipe.HSet(ctx, accessToken, map[string]interface{}{
			"IS_BEARER": true,
			"TOKEN":     refreshToken,
			"VALID":     true,
		})

		res[1] = pipe.ExpireAt(ctx, accessToken, accessTokenExpiresAt)

		res[2] = pipe.HSet(context.Background(), refreshToken, map[string]interface{}{
			"IS_BEARER": false,
			"TOKEN":     accessToken,
			"VALID":     true,
		})

		res[3] = pipe.ExpireAt(ctx, refreshToken, refreshTokenExpiresAt)
		for _, cmd := range res {
			if cmd.Err() != nil {
				return cmd.Err()
			}
		}
		return nil
	})
	return err
}

func InvalidateToken(tokenString string) error {
	ctx := context.Background()
	client := storage.GetRedis()

	_, err := client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		result, err := pipe.HGetAll(ctx, tokenString).Result()
		if err != nil {
			return err
		}

		err = pipe.HSet(ctx, tokenString, "VALID", false).Err()
		if err != nil {
			return err
		}

		err = pipe.HSet(ctx, result["TOKEN"], "VALID", false).Err()
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
