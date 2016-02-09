package httphelper

import (
	"net"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errgo"
	"github.com/julienschmidt/httprouter"
)

type Handler func(http.ResponseWriter, *http.Request, httprouter.Params) *HandlerError
type HandlerError struct {
	Error error
	Code  int
}

func NewHandlerErrorDef(err error) *HandlerError {
	return NewHandlerError(err, http.StatusInternalServerError)
}

func NewHandlerError(err error, code int) *HandlerError {
	return &HandlerError{Error: err, Code: code}
}

func handlerLoggerHTTP(fn Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Stopping {
			handlerLogger(pageRouterStopping, w, r, nil)
		} else {
			handlerLogger(fn, w, r, nil)
		}
	},
	)
}

func handlerLoggerRouter(fn Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if Stopping {
			handlerLogger(pageRouterStopping, w, r, p)
		} else {
			handlerLogger(fn, w, r, p)
		}
	}
}

func handlerLogger(fn Handler, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	l := newHandlerLogEntry(r)
	l.Info("New request")

	starttime := time.Now()

	err := fn(w, r, p)
	if err != nil {
		if err.Error != nil {
			handlerError(w, r, err)
		}
	}
	duration := time.Since(starttime)

	l.Data["duration"] = duration
	l.Info("Finished")
}

func handlerError(w http.ResponseWriter, r *http.Request, e *HandlerError) {
	remote, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
		}).Warning(errgo.Notef(err, "can not get ip from remote addr"))
	}

	log.WithFields(log.Fields{
		"context":      "http",
		"responsecode": e.Code,
		"url":          r.URL,
		"ip":           remote,
	}).Error(e.Error)

	log.WithFields(log.Fields{
		"context":      "http",
		"responsecode": e.Code,
		"url":          r.URL,
		"ip":           remote,
	}).Debug(errgo.Details(e.Error))

	scode := strconv.Itoa(e.Code)
	http.Error(w, scode+" - "+e.Error.Error(), e.Code)
}

func newHandlerLogEntry(r *http.Request) *log.Entry {
	remote, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
		}).Warning(errgo.Notef(err, "can not get ip from remote addr"))
	}

	return log.WithFields(log.Fields{
		"context": "http",
		"url":     r.URL,
		"ip":      remote,
	})
}
