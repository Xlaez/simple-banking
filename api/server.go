package api

import (
	"fmt"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// serves all http request for the app
type Server struct {
	config util.Config
	store *db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

// initialize server
func NewServer(config util.Config, store *db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasteoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// custom parameter for validating currency
		v.RegisterValidation("currency", validCurrency) 
	}
	server.Router()
	return server, nil
}

// start a server on address
func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

// custom error handler
func errorResponse(err error) gin.H {
	return gin.H{"An error occured": err.Error() }
}