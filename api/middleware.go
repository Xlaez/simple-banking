package api

import (
	"errors"
	"fmt"
	"net/http"
	"simple-bank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorixationHeaderKey = "Authorization"
	authorizationPayloadKey = "Authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorixationHeader := ctx.GetHeader(authorixationHeaderKey)
		if len(authorixationHeader) == 0{
			err:= errors.New("provide an authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return;
		}

		fields := strings.Fields(authorixationHeader)
		if len(fields) < 2 {
			err:= errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return;
		}

		authorizationType := strings.ToLower(fields[0]) 
		if authorizationType != "bearer" {
			err:= fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return;
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil{

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return;
		}

		// store payload to context
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}