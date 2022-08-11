package jwt

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestToken(t *testing.T) {
	Convey("Test Jwt Token", t, func() {
		jp := jwtProvider{}
		jp.Load("fobuzkd&yshKl#Si", 8600)
		//设置当前时间为：2021-06-28 11:31:07 +0800 CST
		now := time.Unix(1624851067, 0)
		jwt.TimeFunc = func() time.Time {
			return now
		}
		nowUnix := now.Unix()
		exp := now.Add(time.Second * time.Duration(jp.exp)).Unix()
		audience := "passport"
		issuer := "passport"
		subject := "o7hLH01PWOsQSja3_Nmmrm3UnKnQ"
		jti := "64"
		claims := PassportClaims{
			Vid:    16,
			Openid: "oKKg76YzwrmwG_CGWYT3VEC-_iYA",
			App:    "风变编程",
			AppId:  82,
			Scene:  "test-jwt",
			Puid:   "3ba5ee62-0602-44cf-b727-da0b92b2297a",
			Bs:     "gxq322ToElEdE",
		}
		expected := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwYXNzcG9ydCIsImV4cCI6MTYyNDg1OTY2NywianRpIjoiNjQiLCJpYXQiOjE2MjQ4NTEwNjcsImlzcyI6InBhc3Nwb3J0IiwibmJmIjoxNjI0ODUxMDY3LCJzdWIiOiJvN2hMSDAxUFdPc1FTamEzX05tbXJtM1VuS25RIiwiaWQiOjE2LCJvaWQiOiJvS0tnNzZZendybXdHX0NHV1lUM1ZFQy1faVlBIiwiYXBwaWQiOjgyLCJhcHAiOiLpo47lj5jnvJbnqIsiLCJzY2VuZSI6InRlc3Qtand0IiwicHVpZCI6IjNiYTVlZTYyLTA2MDItNDRjZi1iNzI3LWRhMGI5MmIyMjk3YSIsImJzIjoiZ3hxMzIyVG9FbEVkRSJ9.4tKqrvbCgFX1OF_2FNdyleV_8Sfdfy4FhINhh5lsXqA`
		Convey("Test Mark Token", func() {
			actual := jp.maskToken(audience, issuer, subject, jti, nowUnix, nowUnix, exp, claims)
			So(actual, ShouldEqual, expected)
		})
		Convey("Test Read Token", func() {
			_jti, _claims, err := jp.ReadToken(expected, audience, issuer)
			So(err, ShouldBeNil)
			So(_jti, ShouldEqual, jti)
			So(_claims.Bs, ShouldEqual, claims.Bs)
		})
	})
}
