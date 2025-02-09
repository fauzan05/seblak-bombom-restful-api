package helper

type Role int
type PaymentMethod string
type PaymentStatus int
type OrderStatus int
type DiscountType int
type NotificationType int
type WalletStatus int

const (
	// role
	ADMIN    Role = 1
	CUSTOMER Role = 2
	// payment status
	PAID_PAYMENT    PaymentStatus = 2
	PENDING_PAYMENT PaymentStatus = 1
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
	// wallet status
	ACTIVE  WalletStatus = 1
	INACIVE WalletStatus = 2
	SUSPEND WalletStatus = 3
	// payment method
	GOPAY     PaymentMethod = "gopay"
	SHOPEEPAY PaymentMethod = "shopeepay"
	DANA      PaymentMethod = "dana"
	QRIS      PaymentMethod = "qris"
	WALLET    PaymentMethod = "wallet"
)
