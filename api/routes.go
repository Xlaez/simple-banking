package api

import "github.com/gin-gonic/gin"


func (server *Server) Router(){
	// add routes
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts", server.getAccounts)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.PATCH("/accounts", server.updateAccount)
	authRoutes.DELETE("/account/:id", server.deleteAccount)

	router.POST("/entries", server.createEntry)
	router.GET("/entry/:id", server.getEntry)
	router.GET("/entries", server.getEntries)
	router.PATCH("/entries", server.updateEntry)
	router.DELETE("/entry/:id", server.deleteEntry)

	authRoutes.POST("/transactions", server.createTransaction)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}
