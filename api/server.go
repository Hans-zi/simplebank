package api

import (
	"fmt"

	db "github.com/Hans-zi/simple_bank/db/sqlc"
	"github.com/Hans-zi/simple_bank/token"
	"github.com/Hans-zi/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: maker,
	}

	registerValidation("currency", validCurrency)

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoute := router.Group("/").Use(authMiddleWare(server.tokenMaker))

	authRoute.POST("/accounts", server.createAccount)
	authRoute.GET("/accounts/:id", server.getAccount)
	authRoute.GET("/accounts", server.listAccounts)
	authRoute.DELETE("/accounts/:id", server.deleteAccount)

	authRoute.POST("/transfer", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func registerValidation(name string, fn validator.Func) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation(name, fn)
		if err != nil {
			return
		}
	}
}
