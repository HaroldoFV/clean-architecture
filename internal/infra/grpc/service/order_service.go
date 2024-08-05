package service

import (
	"clean-architecture/internal/infra/grpc/pb"
	"clean-architecture/internal/usecase"
	"context"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewOrderService(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrdersUseCase *usecase.ListOrdersUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.Blank) (*pb.OrderList, error) {
	orders, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var ordersResponse []*pb.Order

	for _, order := range orders {
		response := &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		}
		ordersResponse = append(ordersResponse, response)
	}

	return &pb.OrderList{Orders: ordersResponse}, nil
}
