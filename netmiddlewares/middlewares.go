package middlewares

import (
	"net/http"
	"net/url"
)

// Middleware is a net/http middleware.
type Middleware = func(http.Handler) http.Handler

// RouteFinder finds Route by given URL.
type RouteFinder func(method string, u *url.URL) (Route, bool)

// Server is a generic ogen server type.
type Server[R Route] interface {
	FindPath(method string, u *url.URL) (r R, _ bool)
}

// Route is a generic ogen route type.
type Route interface {
	Name() string
	OperationID() string
	PathPattern() string
}

// MakeRouteFinder creates RouteFinder from given server.
func MakeRouteFinder[R Route, S Server[R]](server S) RouteFinder {
	return func(method string, u *url.URL) (Route, bool) {
		return server.FindPath(method, u)
	}
}

// Wrap handler using given middlewares.
func Wrap(h http.Handler, middlewares ...Middleware) http.Handler {
	switch len(middlewares) {
	case 0:
		return h
	case 1:
		return middlewares[0](h)
	default:
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}
