package api

import (
	"github.com/gin-gonic/gin"
	db "simplebank/db/sqlc"
)

//Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

//NewServer builds a Server struct
func NewServer(store db.Store) *Server {
	router := gin.Default()
	accHandler := newAccountHandler(store)

	router.POST("/accounts", accHandler.post)
	router.GET("/accounts/:id", accHandler.get)
	router.GET("/accounts", accHandler.list)

	return &Server{
		store:  store,
		router: router,
	}
}

//Start runs the HTTP server on specific address.
func (s Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
