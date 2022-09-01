package api

import (
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)



type createUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"fullname" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type createUserResponse struct{ 
	Username          string    `json:"username"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var request createUserRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	hashedPassword, err := util.HashPassword(request.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return;
	}

	arg := db.CreateUserParams{
		Username: request.Username,
		FullName: request.Fullname,
		HashedPassword: hashedPassword,
		Email: request.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil{
		if pqError, ok := err.(*pq.Error); ok{
			switch pqError.Code.Name(){
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return;
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return;
	}

	res := createUserResponse{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, res)
}
