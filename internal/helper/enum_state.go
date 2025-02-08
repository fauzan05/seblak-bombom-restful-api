package helper

type Role int
type PaymentMethod int
type PaymentStatus int
type OrderStatus int
type DiscountType int
type NotificationType int

const (
	// role
	ADMIN    Role = 1
	CUSTOMER Role = 2
	// payment method
	ONLINE PaymentMethod = 1
	ONSITE PaymentMethod = 2
	// payment status
	PENDING_PAYMENT PaymentStatus = 1
	PAID_PAYMENT    PaymentStatus = 2
	FAILED_PAYMENT  PaymentStatus = 0
	// order status
	ORDER_PENDING         OrderStatus = 1
	ORDER_RECEIVED        OrderStatus = 2
	ORDER_BEING_DELIVERED OrderStatus = 3
	ORDER_DELIVERED       OrderStatus = 4
	READY_FOR_PICKUP      OrderStatus = 5
	ORDER_REJECTED        OrderStatus = 0
	// discount type
	NOMINAL DiscountType = 1
	PERCENT DiscountType = 2
	// notification type
	TRANSACTION NotificationType = 1
	PROMOTION   NotificationType = 2
)
