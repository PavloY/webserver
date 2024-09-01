package api

import (
	"WEB_SERVER/internal/app/middleware"
	"WEB_SERVER/storage"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	prefix string = "/api/v1"
)

func(a *API) configreLoggerField() error{
	log_level, err:= logrus.ParseLevel(a.config.LoggerLevel)
	if err!=nil{
		return err
	}
	a.logger.SetLevel(log_level)
	return nil
}

func (a *API) configreRouterField(){
	a.router.HandleFunc(prefix + "/products", a.GetAllProducts).Methods("GET")
	a.router.Handle(prefix + "/products/{id}", middleware.JwtMiddleware.Handler(http.HandlerFunc(a.GetProductById))).Methods("GET")
	a.router.HandleFunc(prefix + "/products/{id}", a.UpdateProductById).Methods("PUT")
	a.router.HandleFunc(prefix + "/products/{id}", a.DeleteProductById).Methods("DELETE")
	a.router.HandleFunc(prefix + "/products", a.CreateProduct).Methods("POST")
	a.router.HandleFunc(prefix + "/user/register", a.RegisterUser).Methods("POST")
	a.router.HandleFunc(prefix + "/user/auth", a.PostToAuth).Methods("POST")
}

func (a *API) configreStorageField() error {
	storage := storage.New(a.config.Storage)
	if err := storage.Open(); err !=nil {
		return err
	}
	a.storage = storage
	return nil
}