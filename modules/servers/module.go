package servers

import (
	"github.com/Kamila3820/go-shop-tutorial/modules/appinfo/appinfoHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/appinfo/appinfoRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/appinfo/appinfoUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/monitor/monitorHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/users/usersHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/users/usersRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
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

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignUpAdmin)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

	// Initial admin ขึ้นมา 1 คน ใน Db (Insert ใน SQL)
	// Generate Admin Key
	// ทุกครั้งที่ทำการสมัคร Admin เพิ่ม ให้ส่ง Admin Token มาด้วยทุกครั้ง ผ่าน Middleware
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

}
