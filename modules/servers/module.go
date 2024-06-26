package servers

import (
	"github.com/Kamila3820/go-shop-tutorial/modules/appinfo/appinfoHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/appinfo/appinfoRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/appinfo/appinfoUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/files/filesHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/files/filesUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/middlewares/middlewaresUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/monitor/monitorHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/orders/ordersHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/orders/ordersRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/orders/ordersUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/products/productsHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/products/productsRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/products/productsUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/users/usersHandlers"
	"github.com/Kamila3820/go-shop-tutorial/modules/users/usersRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
	OrdersModule()
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

	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)
	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) FilesModule() {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	router := m.r.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)
}

func (m *moduleFactory) ProductsModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)

	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)
	productsUsecase := productsUsecases.ProductsUsecase(productsRepository)
	productsHandler := productsHandlers.ProductsHandler(m.s.cfg, productsUsecase, filesUsecase)

	router := m.r.Group("/products")

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.AddProduct)

	router.Patch("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.UpdateProduct)

	router.Get("/", m.mid.ApiKeyAuth(), productsHandler.FindProduct)
	router.Get("/:product_id", m.mid.ApiKeyAuth(), productsHandler.FindOneProduct)

	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.DeleteProduct)
}

func (m *moduleFactory) OrdersModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)

	ordersRepository := ordersRepositories.OrdersRepository(m.s.db)
	ordersUsecase := ordersUsecases.OrdersUsecase(ordersRepository, productsRepository)
	ordersHandler := ordersHandlers.OrdersHandler(m.s.cfg, ordersUsecase)

	router := m.r.Group("/orders")

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:order_id", m.mid.JwtAuth(), ordersHandler.FindOneOrder)

}
