package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

// PassportClaims
// {
//  "app": "pyclass",
//  "appid": 82,
//  "exp": 1624995458,
//  "id": 49,
//  "oid": "oRb8OwokrjVzMAXGOvvlEzpJwaf4",
//  "puid": "",
//  "sub": "o7hLH06V6t1JNjukXclK2nQq8EvE", //unioind
//}
type PassportClaims struct {
	Vid    int64  `json:"id"`    //for_os.users.id
	Openid string `json:"oid"`   // for_os.users.openid
	AppId  int64  `json:"appid"` //for_apis.wx_apps.id
	App    string `json:"app"`   //for_apis.wx_apps.app_name
	Scene  string `json:"scene"` //场景
	Puid   string `json:"puid"`  //v2系统唯一性标识，只有绑定手机号才会存在
	Bs     string `json:"bs"`    //v3系统存在于缓存中的唯一标识
}
type tokenStandardClaims struct {
	jwt.StandardClaims
	PassportClaims
}

// TokenProvider jwt配置加载,token生成和读取
type TokenProvider interface {
	ReadToken(tokenString, audience, issuer string) (jti string, bean PassportClaims, err error)
	MaskToken(audience, issuer, subject, jti string, nbf int64, claims PassportClaims) (signed string)
	Load(key string, exp int)
}

type jwtProvider struct {
	claimKey []byte // claim密钥
	exp      int    // 过期时间
}

func NewTokenProvider() TokenProvider {
	return &jwtProvider{}
}
func (a *jwtProvider) Load(key string, exp int) {
	if len(key) < 8 {
		panic("jwt的密钥长度必须大于等于8位")
	}
	a.claimKey = []byte(key)
	a.exp = exp
}

var (
	VerifyAudienceError = errors.New("接收JWT的一方不匹配")
	VerifyIssuerError   = errors.New("JWT签发者不匹配")
	CalaimsAssertError  = errors.New("claims类型不匹配")
)

func (a jwtProvider) ReadToken(tokenString, audience, issuer string) (jti string, bean PassportClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenStandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.claimKey, nil
	})
	if err != nil {
		return jti, bean, err
	}
	if claims, ok := token.Claims.(*tokenStandardClaims); ok {
		if token.Valid { //VerifyExpiresAt,VerifyIssuedAt,VerifyNotBefore,
			if !claims.VerifyAudience(audience, true) {
				log.Info().Msgf("接收JWT的一方不匹配[%s]", audience)
				return claims.Id, claims.PassportClaims, VerifyAudienceError
			}
			if !claims.VerifyIssuer(issuer, true) {
				log.Info().Msgf("JWT签发者不匹配[%s]", audience)
				return claims.Id, claims.PassportClaims, VerifyIssuerError
			}
			return claims.Id, claims.PassportClaims, nil
		} else {
			return claims.Id, claims.PassportClaims, err
		}
	} else {
		return jti, bean, CalaimsAssertError
	}
}

// MaskToken 签发token
//aud：接收 JWT 的一方
//exp：JWT 的过期时间，这个过期时间必须要大于签发时间
//jti：JWT 的唯一身份标识，主要用来作为一次性 token, 从而回避重放攻击。
//iat：JWT 的签发时间
//iss：JWT 签发者
//nbf：定义在什么时间之前，该 JWT 都是不可用的
//sub：JWT 所面向的用户
func (a jwtProvider) MaskToken(audience, issuer, subject, jti string, nbf int64, claims PassportClaims) (signed string) {
	now := time.Now()
	nowUnix := now.Unix()
	exp := now.Add(time.Second * time.Duration(a.exp)).Unix()
	return a.maskToken(audience, issuer, subject, jti, nowUnix, nbf, exp, claims)
}
func (a jwtProvider) maskToken(audience, issuer, subject, jti string, now, nbf, exp int64, claims PassportClaims) (signed string) {
	tsc := &tokenStandardClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  audience,
			ExpiresAt: exp,
			IssuedAt:  now,
			Issuer:    issuer,
			NotBefore: nbf,
			Subject:   subject,
			Id:        jti,
		},
		PassportClaims: claims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tsc)
	if ss, err := token.SignedString(a.claimKey); err == nil {
		return ss
	} else {
		log.Error("签发token发生错误:%+v", err)
		return ""
	}
}
