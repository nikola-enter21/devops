package server

import "github.com/gofiber/fiber/v3"

func wrapRoutes(router fiber.Router, middlewares []fiber.Handler, fn func(r fiber.Router)) {
	r := &routeWrapper{Router: router, middlewares: middlewares}
	fn(r)
}

type routeWrapper struct {
	fiber.Router
	middlewares []fiber.Handler
}

func (rw *routeWrapper) chain(handlers ...fiber.Handler) []fiber.Handler {
	return append(rw.middlewares, handlers...)
}

func (rw *routeWrapper) Get(path string, handler fiber.Handler, handlers ...fiber.Handler) fiber.Router {
	all := rw.chain(append([]fiber.Handler{handler}, handlers...)...)
	return rw.Router.Get(path, all[0], all[1:]...)
}

func (rw *routeWrapper) Post(path string, handler fiber.Handler, handlers ...fiber.Handler) fiber.Router {
	all := rw.chain(append([]fiber.Handler{handler}, handlers...)...)
	return rw.Router.Post(path, all[0], all[1:]...)
}

func (rw *routeWrapper) Put(path string, handler fiber.Handler, handlers ...fiber.Handler) fiber.Router {
	all := rw.chain(append([]fiber.Handler{handler}, handlers...)...)
	return rw.Router.Put(path, all[0], all[1:]...)
}

func (rw *routeWrapper) Patch(path string, handler fiber.Handler, handlers ...fiber.Handler) fiber.Router {
	all := rw.chain(append([]fiber.Handler{handler}, handlers...)...)
	return rw.Router.Patch(path, all[0], all[1:]...)
}

func (rw *routeWrapper) Delete(path string, handler fiber.Handler, handlers ...fiber.Handler) fiber.Router {
	all := rw.chain(append([]fiber.Handler{handler}, handlers...)...)
	return rw.Router.Delete(path, all[0], all[1:]...)
}

func (rw *routeWrapper) Options(path string, handler fiber.Handler, handlers ...fiber.Handler) fiber.Router {
	all := rw.chain(append([]fiber.Handler{handler}, handlers...)...)
	return rw.Router.Options(path, all[0], all[1:]...)
}

func (rw *routeWrapper) Head(path string, handler fiber.Handler, handlers ...fiber.Handler) fiber.Router {
	all := rw.chain(append([]fiber.Handler{handler}, handlers...)...)
	return rw.Router.Head(path, all[0], all[1:]...)
}
