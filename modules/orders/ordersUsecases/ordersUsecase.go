package ordersUsecases

import (
	"github.com/Kamila3820/go-shop-tutorial/modules/orders"
	"github.com/Kamila3820/go-shop-tutorial/modules/orders/ordersRepositories"
	"github.com/Kamila3820/go-shop-tutorial/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
}

type ordersUsecase struct {
	orderRepository    ordersRepositories.IOrdersRepository
	productsRepository productsRepositories.IProductsRepository
}

func OrdersUsecase(orderRepository ordersRepositories.IOrdersRepository, productsRepository productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		orderRepository:    orderRepository,
		productsRepository: productsRepository,
	}
}

func (u *ordersUsecase) FindOneOrder(orderId string) (*orders.Order, error) {
	order, err := u.orderRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}
