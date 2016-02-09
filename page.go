package httphelper

import (
	"net/http"

	"github.com/juju/errgo"
	"github.com/julienschmidt/httprouter"
)

func pageRouterNotFound(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HandlerError {
	return NewHandlerError(errgo.New("page not found"), http.StatusNotFound)
}

func pageRouterMethodNotAllowed(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HandlerError {
	return NewHandlerError(errgo.New("method not allowed"), http.StatusMethodNotAllowed)
}

func pageRouterStopping(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HandlerError {
	return NewHandlerError(errgo.New("service is stopping"), http.StatusInternalServerError)
}

func pageFiles(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HandlerError {
	return NewHandlerError(errgo.New("not implemented"), http.StatusInternalServerError)
}
