package handler

import (
	"app/apps/im/ws/internal/svc"
	"app/pkg/ctxdata"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
)

// JWT 验证器
type JwtAuth struct {
	ctx    *svc.ServiceContext
	parser *token.TokenParser   // 解析 token 用的解析器
	logx.Logger
}

func NewJwtAuto(srcCtx *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		ctx:    srcCtx,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

// Auth 尝试从 HTTP 请求中解析 JWT。如果解析成功且令牌有效，它会从令牌中提取用户标识（claims），并将其添加到请求的上下文中。
// 如果解析失败或令牌无效，记录错误并返回 false。
func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	tok, err := j.parser.ParseToken(r, j.ctx.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("parse token err %v ", err)
		return false
	}

	if !tok.Valid {
		return false
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.IdentityKey, claims[ctxdata.IdentityKey]))

	return true

}


// 从 http 请求中获取 uid
func (j *JwtAuth) UserId(r *http.Request) string {

	return ctxdata.GetUid(r.Context())

}
