package servers

import (
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	r   fiber.Router //router
	s   *server      //server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	handler := middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
	return handler
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}
