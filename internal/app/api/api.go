package api

import (
	"WEB_SERVER/storage"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type API struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	storage *storage.Storage
}

func New(config *Config) *API {
	return &API{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (api *API) Start() error {
	if err:= api.configreLoggerField(); err !=nil {
		return err
	}
	api.logger.Info("starting api server at port:", api.config.BindAddr)
	api.configreRouterField()
	if err := api.configreStorageField(); err !=nil{
		return err
	}
	return http.ListenAndServe(api.config.BindAddr, api.router)
}
