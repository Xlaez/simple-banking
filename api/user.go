package api

import (
	"database/sql"
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

type UserResponse struct{ 
	Username          string    `json:"username"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func newUserResponse (user db.User) UserResponse{
	return UserResponse{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
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

	res := newUserResponse(user)

	ctx.JSON(http.StatusCreated, res)
}



type loginUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password  	string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken 	string `json:"access_token"`
	User 			UserResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return;
		}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return;
	}

	if err = util.ComparePassword(req.Password, user.HashedPassword); err!= nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return;
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return;
	}

	res := loginUserResponse{
		AccessToken: accessToken,
		User: newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, res)
}