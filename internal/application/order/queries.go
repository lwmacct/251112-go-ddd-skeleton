package order

import (
	"context"
	"github.com/yourusername/go-ddd-skeleton/internal/shared/pagination"
)

// GetOrder 获取订单（查询）
func (s *Service) GetOrder(ctx context.Context, orderID string) (*OrderDTO, error) {
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return domainOrderToDTO(o), nil
}

// GetOrderByNumber 根据订单号获取订单（查询）
func (s *Service) GetOrderByNumber(ctx context.Context, orderNumber string) (*OrderDTO, error) {
	o, err := s.orderRepo.FindByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, err
	}

	return domainOrderToDTO(o), nil
}

// ListOrders 列出用户订单（查询）
func (s *Service) ListOrders(ctx context.Context, userID string, req ListOrdersRequest) (*ListOrdersResponse, error) {
	// 解析分页参数
	offset, limit := pagination.ParsePaginationParams(req.Page, req.PageSize)

	// 查询订单列表
	orders, total, err := s.orderRepo.ListByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	orderDTOs := make([]*OrderDTO, len(orders))
	for i, o := range orders {
		orderDTOs[i] = domainOrderToDTO(o)
	}

	// 创建分页对象
	pg := pagination.NewPagination(req.Page, req.PageSize, total)

	return &ListOrdersResponse{
		Orders:     orderDTOs,
		Total:      total,
		Page:       pg.Page,
		PageSize:   pg.PageSize,
		TotalPages: pg.TotalPages,
	}, nil
}

// GetPayment 获取支付（查询）
func (s *Service) GetPayment(ctx context.Context, orderID string) (*PaymentDTO, error) {
	payment, err := s.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return domainPaymentToDTO(payment), nil
}

// GetShipment 获取发货（查询）
func (s *Service) GetShipment(ctx context.Context, orderID string) (*ShipmentDTO, error) {
	shipment, err := s.shipmentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return domainShipmentToDTO(shipment), nil
}

// GetShipmentByTrackingNumber 根据追踪号获取发货（查询）
func (s *Service) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*ShipmentDTO, error) {
	shipment, err := s.shipmentRepo.FindByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		return nil, err
	}

	return domainShipmentToDTO(shipment), nil
}

