package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthorizationCode struct {
	Code      string
	UserID    string
	ClientID  string
	ExpiresAt time.Time
}
type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	TokenType    string
}

type TokenRequest struct {
	GrantType    string `json:"grant_type" db:"grant_type"`
	ClientID     string `json:"client_id" db:"client_id"`
	ClientSecret string `json:"client_secret" db:"client_secret"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	Code         string `json:"code" db:"code"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
	RedirectURL  string `json:"redirect_url" db:"redirect_url"`
}
type TokenClaims struct {
	Sub    	  string   `json:"sub"`
	ClientID  string   `json:"client_id"`
	Scopes    []string `json:"scopes"`
	Exp 	  int64    `json:"exp"`
	Iat		  int64    `json:"iat"`
	Nbf 	  int64    `json:"nbf"`
	Aud  	  string   `json:""`
	Iss    	  string   `json:"iss"`
	Type      string   `json:"type"`
	
}


func (c *TokenClaims) Valid() error {
    now := time.Now().Unix()
    if c.Exp != 0 && now > c.Exp {
        return fmt.Errorf("token expired")
    }
    if c.Nbf != 0 && now < c.Nbf {
        return fmt.Errorf("token not valid yet")
    }
    return nil
}

func (c *TokenClaims) GetIssuer() (string,error) {
    return c.Iss,nil
}

func (c *TokenClaims) GetSubject() (string,error) {
    return c.Sub,nil
}

func (c *TokenClaims) GetAudience() (jwt.ClaimStrings, error) {
    return jwt.ClaimStrings{c.Aud}, nil
}

func (c *TokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
    return jwt.NewNumericDate(time.Unix(c.Exp, 0)), nil
}

func (c *TokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
    return jwt.NewNumericDate(time.Unix(c.Iat, 0)), nil
}

func (c *TokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
    return jwt.NewNumericDate(time.Unix(c.Nbf, 0)), nil
}