package jwt

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/config"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/models"
)
var (
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
)

const (
	ttl_access = time.Hour
	ttl_refresh = time.Hour * 24 * 7
)

type JwtService interface {
	GenerateAccessToken(userID string, clientID string, scopes []string) (string, error)
	VerifyToken(ctx context.Context, token string) (*models.TokenClaims, error)
	GenerateRefreshToken(userID string) (string, error)
}

type jwtSvcImpl struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}


func NewJwtService(config *config.Config) (JwtService,error) {

	privBytes, err := ioutil.ReadFile(config.JWTPrivetKeyPath)
    if err != nil {
        return nil, fmt.Errorf("could not read private key: %w", err)
    }

    pubBytes, err := ioutil.ReadFile(config.JWTPublicKeyPath)
    if err != nil {
        return nil, fmt.Errorf("could not read public key: %w", err)
    }

    privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
    if err != nil {
        return nil, fmt.Errorf("could not parse private key: %w", err)
    }

    publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
    if err != nil {
        return nil, fmt.Errorf("could not parse public key: %w", err)
    }
	return &jwtSvcImpl{
        privateKey: privateKey,
        publicKey:  publicKey,
	},nil
}


func (j *jwtSvcImpl) GenerateAccessToken(userID string, clientID string, scopes []string) (string, error) {
    tokenClaims := models.TokenClaims{
        Sub: userID,
        ClientID: clientID,
        Scopes: scopes,
        Exp: time.Now().Add(ttl_access).Unix(),
        Iat: time.Now().Unix(),
        Nbf: time.Now().Unix(),
        Aud: clientID,
        Iss: "own",
        Type: "access_token",
    }
	claims := jwt.MapClaims{
        "sub":       tokenClaims.Sub,
        "client_id": tokenClaims.ClientID,
        "scopes":    tokenClaims.Scopes,
        "exp":       tokenClaims.Exp,
        "iat":       tokenClaims.Iat,
        "nbf":       tokenClaims.Nbf,
        "aud":       tokenClaims.Aud,
        "iss":       tokenClaims.Iss,
        "type":      tokenClaims.Type,
    }
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}
func (j *jwtSvcImpl) GenerateRefreshToken(userID string) (string, error) {
	tokenClaims := models.TokenClaims{
        Sub: userID,
        Exp: time.Now().Add(ttl_refresh).Unix(),
        Iat: time.Now().Unix(),
        Nbf: time.Now().Unix(),
        Iss: "own",
        Type: "refresh_token",
    }
	claims := jwt.MapClaims{
        "sub":       tokenClaims.Sub,
        "exp":       tokenClaims.Exp,
        "iat":       tokenClaims.Iat,
        "nbf":       tokenClaims.Nbf,
        "iss":       tokenClaims.Iss,
        "type":      tokenClaims.Type,
    }
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j *jwtSvcImpl) VerifyToken(ctx context.Context, tokenStr string) (*models.TokenClaims, error) {
    if j.publicKey == nil {
        return nil, errors.New("public key not loaded")
    }

    token, err := jwt.ParseWithClaims(tokenStr, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Faqat RSA imzolarini qabul qilamiz
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return j.publicKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
