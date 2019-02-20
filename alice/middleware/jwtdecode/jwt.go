package jwtdecode

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/promoboxx/go-service/alice/middleware"
)

type jwtDecoder struct{}

// NewJWTDecoder returns a JWTDecoder that wraps the dgrijalva/jwt-go package.NewJWTDecoder
// this prevents the middleware package from forcing jwt-go imports
func NewJWTDecoder() middleware.JWTDecoder {
	return &jwtDecoder{}
}

func (j *jwtDecoder) DecodeSegment(seg string) ([]byte, error) {
	return jwt.DecodeSegment(seg)
}
