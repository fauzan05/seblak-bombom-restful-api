package helper

type Role string
type PaymentMethod string
type PaymentStatus string
type DeliveryStatus string
type DiscountType string
type NotificationType string

const (
	// role
	ADMIN    Role = "admin"
	CUSTOMER Role = "customer"
	// payment method
	ONLINE PaymentMethod = "online"
	ONSITE PaymentMethod = "onsite"
	// payment status
	PENDING_PAYMENT PaymentStatus = "pending"
	PAID_PAYMENT    PaymentStatus = "paid"
	FAILED_PAYMENT  PaymentStatus = "failed"
	// delivery status
	PREPARE_DELIVERY DeliveryStatus = "prepare"
	ON_THE_WAY       DeliveryStatus = "on_the_way"
	SENT             DeliveryStatus = "sent"
	TAKE_AWAY        DeliveryStatus = "take_away"
	// discount type
	PERCENT DiscountType = "percent"
	NOMINAL DiscountType = "nominal"
	// notification type
	TRANSACTION NotificationType = "transaction"
	PROMOTION   NotificationType = "promotion"
)
