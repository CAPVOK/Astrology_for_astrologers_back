package middleware

import (
	"errors"
	"net/http"
	"space/internal/app/ds"
	"space/internal/app/repository"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	rr repository.Redis
}

func New(redisRepo repository.Redis) *Middleware {
	return &Middleware{
		rr: redisRepo,
	}
}

func (m *Middleware) IsAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, err := ctx.Request.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
				return
			}

			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if c.Expires.Before(time.Now()) {
			// Check if the session has expired
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sessionToken := c.Value
		sc, err := m.rr.SessionExists(sessionToken)
		if err != nil {
			ctx.AbortWithStatusJSON(ds.GetHttpStatusCode(err), err.Error())
			return
		}

		if sc.UserID == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userID", sc.UserID)
		ctx.Set("role", sc.Role)
		ctx.Set("sessionContext", sc)

		ctx.Next()
	}
}
