package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// serves all http request for the app
type Server struct {
	store *db.Store
	router *gin.Engine
}

// initialize server
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// custom parameter for validating currency
		v.RegisterValidation("currency", validCurrency) 
	}

	// add routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.getAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.PATCH("/accounts", server.updateAccount)
	router.DELETE("/account/:id", server.deleteAccount)
	router.POST("/entries", server.createEntry)
	router.GET("/entry/:id", server.getEntry)
	router.GET("/entries", server.getEntries)
	router.PATCH("/entries", server.updateEntry)
	router.DELETE("/entry/:id", server.deleteEntry)
	router.POST("/transactions", server.createTransaction)

	server.router = router
	return server
}

// start a server on address
func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

// custom error handler
func errorResponse(err error) gin.H {
	return gin.H{"An error occured": err.Error() }
}