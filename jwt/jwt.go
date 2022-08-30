// Package jwt
// @Title  jwt.go
// @Description  jwt token
// @Author  Brandon     时间（2022/8/29）
// @Update  Brandon     时间（2022/8/29）
package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

type UserClaims struct {
	UId    int64  `json:"u_id"`    // 用户id
	UName  string `json:"u_name"`  // 用户名
	ULevel int    `json:"u_level"` // 用户等级
	Avater string `json:"avater"`  // 用户头像
	Mobile string `json:"mobile"`  // 用户手机号
}
type tokenStandardClaims struct {
	jwt.StandardClaims
	UserClaims
}

// TokenProvider jwt配置加载, token生成和读取
type TokenProvider interface {
	VerifyToken(tokenString, aud, iss string) (jti string, userClaims UserClaims, err error)
	CreateToken(aud, iss, sub, jti string, nbf int64, claims UserClaims) (signed string)
	Load(key string, exp int)
}

type jwtProvider struct {
	secret []byte // claim 密钥
	exp    int    // 过期时间
}

var JwtTokenProvider TokenProvider

func init() {
	JwtTokenProvider = new(jwtProvider)
}

func (a *jwtProvider) Load(secret string, exp int) {
	if len(secret) < 8 {
		panic("The key length of jwt must be greater than or equal to 8 bits")
	}
	a.secret = []byte(secret)
	a.exp = exp
}

// CreateToken
// @title CreateToken
// @description     创建token
// @auth            Brandon     时间（2022/8/29）
// @param           aud         JWT token 的收件人
// @param           iss         JWT token 的发件人
// @param           sub         JWT token 的主体
// @param           jti         JWT token 的唯一标识符
// @param           nbf         JWT token 的生效时间
// @return          signed      JWT token 签名
func (a jwtProvider) CreateToken(aud, iss, sub, jti string, nbf int64, claims UserClaims) (signed string) {
	now := time.Now()
	nowUnix := now.Unix()
	exp := now.Add(time.Second * time.Duration(a.exp)).Unix()

	tsc := &tokenStandardClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  aud,     // 认证 JWT 的收件人
			ExpiresAt: exp,     // 认证过期时间
			IssuedAt:  nowUnix, // 认证 JWT 的时间
			Issuer:    iss,     // 认证 JWT 的发件人
			NotBefore: nbf,     // 认证 JWT 的生效时间（时间戳）
			Subject:   sub,     // 认证 JWT 的主体
			Id:        jti,     // 认证 JWT 的唯一标识符（一般为用户id）这里为用户id
		},
		UserClaims: claims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tsc)
	signed, err := token.SignedString(a.secret)
	if err != nil {
		log.Err(err).Msg("签发token发生错误")
	}
	return
}

// VerifyToken
// @title VerifyToken
// @description     验证token
// @auth            Brandon     时间（2022/8/29）
// @param           tokenString JWT token
// @param           aud         JWT token 的收件人
// @param           iss         JWT token 的发件人
// @return          jti         JWT token 的唯一标识符
// @return          claims      用户信息
// @return          err         错误
func (a jwtProvider) VerifyToken(tokenString, aud, iss string) (jti string, userClaims UserClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenStandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secret, nil
	})
	if err != nil {
		return jti, userClaims, err
	}
	claims, ok := token.Claims.(*tokenStandardClaims)

	if ok && token.Valid {
		if !claims.VerifyAudience(aud, true) {
			log.Info().Msgf("JWT接收者不匹配[%s]", aud)
			return claims.Id, claims.UserClaims, jwt.NewValidationError("JWT接收者不匹配", jwt.ValidationErrorAudience)
		}
		if !claims.VerifyIssuer(iss, true) {
			log.Info().Msgf("JWT签发者不匹配[%s]", iss)
			return claims.Id, claims.UserClaims, jwt.NewValidationError("JWT签发者不匹配", jwt.ValidationErrorIssuer)
		}
		return claims.Id, claims.UserClaims, nil
	} else {
		return claims.Id, userClaims, jwt.NewValidationError("JWT验证错误", jwt.ValidationErrorClaimsInvalid)
	}
}
