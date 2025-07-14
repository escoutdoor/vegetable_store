package token

import (
	"errors"
	"time"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/user_service/internal/entity"
	apperrors "github.com/escoutdoor/vegetable_store/user_service/internal/errors"
	"github.com/golang-jwt/jwt/v5"
)

type Provider interface {
	ValidateAccessToken(accessToken string) (string, error)
	ValidateRefreshToken(refreshToken string) (string, error)
	GenerateTokens(userID string) (entity.Tokens, error)
}

type tokenProvider struct {
	accessTokenSecretKey  string
	refreshTokenSecretKey string
	accessTokenTTL        time.Duration
	refreshTokenTTL       time.Duration
}

func NewTokenProvider(
	accessTokenSecretKey string,
	refreshTokenSecretKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) Provider {
	return &tokenProvider{
		accessTokenSecretKey:  accessTokenSecretKey,
		refreshTokenSecretKey: refreshTokenSecretKey,
		accessTokenTTL:        accessTokenTTL,
		refreshTokenTTL:       refreshTokenTTL,
	}
}

type accessTokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func (p *tokenProvider) GenerateTokens(userID string) (entity.Tokens, error) {
	accessToken, err := p.generateAccessToken(userID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate access token", err)
	}

	refreshToken, err := p.generateRefreshToken(userID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate refresh token", err)
	}

	return entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (p *tokenProvider) generateAccessToken(userID string) (string, error) {
	claims := accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.accessTokenTTL)),
		},
		UserID: userID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(p.accessTokenSecretKey))
	if err != nil {
		return "", errwrap.Wrap("new jwt token with claims", err)
	}

	return token, nil
}

func (p *tokenProvider) generateRefreshToken(userID string) (string, error) {
	claims := refreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshTokenTTL)),
		},
		UserID: userID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(p.refreshTokenSecretKey))
	if err != nil {
		return "", errwrap.Wrap("new jwt token with claims", err)
	}

	return token, nil
}

func (p *tokenProvider) ValidateAccessToken(accessToken string) (string, error) {
	jwtToken, err := jwt.ParseWithClaims(accessToken, &accessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(p.accessTokenSecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", apperrors.ErrJwtTokenExpired
		}

		return "", errwrap.Wrap("parse with claims", err)
	}

	if !jwtToken.Valid {
		return "", apperrors.ErrInvalidJwtToken
	}

	claims, ok := jwtToken.Claims.(*accessTokenClaims)
	if !ok {
		return "", errwrap.Wrap("get claims", err)
	}

	return claims.UserID, nil
}

func (p *tokenProvider) ValidateRefreshToken(refreshToken string) (string, error) {
	jwtToken, err := jwt.ParseWithClaims(refreshToken, &refreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(p.refreshTokenSecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", apperrors.ErrJwtTokenExpired
		}

		return "", errwrap.Wrap("parse with claims", err)
	}

	if !jwtToken.Valid {
		return "", apperrors.ErrInvalidJwtToken
	}

	claims, ok := jwtToken.Claims.(*refreshTokenClaims)
	if !ok {
		return "", errwrap.Wrap("get claims", err)
	}

	return claims.UserID, nil
}
