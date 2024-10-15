package handler

import (
	"app/apps/im/ws/internal/svc"
	"app/pkg/ctxdata"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
)

// JWT 验证器

type JwtAuto struct {
	ctx    *svc.ServiceContext
	parser *token.TokenParser
}

func NewJwtAuto(srcCtx *svc.ServiceContext) *JwtAuto {
	return &JwtAuto{
		ctx:    srcCtx,
		parser: token.NewTokenParser(),
	}
}

// Auth 尝试从 HTTP 请求中解析 JWT。如果解析成功且令牌有效，它会从令牌中提取用户标识（claims），并将其添加到请求的上下文中。
// 如果解析失败或令牌无效，记录错误并返回 false。
func (j *JwtAuto) Auth(w http.ResponseWriter, r *http.Request) bool {
	// 解析token
	parseToken, err := j.parser.ParseToken(r, j.ctx.Config.JwtSecret.AccessSecret, "")
	if err != nil {
		return false
	}

	// token 有效
	if !parseToken.Valid {
		return false
	}

	// 拿到存储在 token 中的用户信息
	claims, ok := parseToken.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// add context into http request
	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.IdentityKey, claims[ctxdata.IdentityKey]))

	return true

}


// 从 http 请求中获取 uid
func (j *JwtAuto) UserId(r *http.Request) string {

	return ctxdata.GetUid(r.Context())

}
