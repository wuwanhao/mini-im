package ctxdata

import "github.com/golang-jwt/jwt/v4"

const IdentityKey = "mini_im_identity_key"

// 生成 jwt token
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[IdentityKey] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
