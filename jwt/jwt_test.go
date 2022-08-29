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
		jp.Load("buzkd&yshKl#Si", 3600)
		now := time.Now()
		jwt.TimeFunc = func() time.Time {
			return now
		}
		aud := "gateway"
		iss := "web"
		subject := "o7hLH01PWOsQSja3_Nmmrm3UnKnQ"
		jti := "107"
		nbf := now.Unix()

		signed := jp.CreateToken(aud, iss, subject, jti, nbf)
		_jti, err := jp.VerifyToken(signed, aud, iss)
		So(err, ShouldBeNil)
		So(_jti, ShouldEqual, jti)
	})
}
