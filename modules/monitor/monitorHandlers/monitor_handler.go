package monitorHandlers

import (
	"github.com/Kamila3820/go-shop-tutorial/config"
	"github.com/Kamila3820/go-shop-tutorial/modules/entities"
	"github.com/Kamila3820/go-shop-tutorial/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

//Handler: จะรับ parameter แค่ contextของfiber แล้วresponseออกมาเป็น error

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}
