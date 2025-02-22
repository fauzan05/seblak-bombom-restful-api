package helper

type Role int
type PaymentMethod string
type PaymentStatus int
type OrderStatus int
type DiscountType int
type NotificationType int
type WalletStatus int
type TransactionStatus string
type RequestMethod string
type XenditTransactionStatus string
type ChannelCode string
type PaymentGateway string
type ItemType string

const (
	// role
	ADMIN    Role = 1
	CUSTOMER Role = 2
	// payment status
	PAID_PAYMENT    PaymentStatus = 1
	PENDING_PAYMENT PaymentStatus = 0
	CANCEL_PAYMENT  PaymentStatus = -1
	EXPIRED_PAYMENT PaymentStatus = -2
	FAILED_PAYMENT  PaymentStatus = -3
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
	PAYMENT_METHOD_QR_CODE PaymentMethod = "QR_CODE"
	PAYMENT_METHOD_EWALLET PaymentMethod = "EWALLET"
	PAYMENT_METHOD_WALLET  PaymentMethod = "WALLET"

	PAYMENT_GATEWAY_MIDTRANS PaymentGateway = "MIDTRANS"
	PAYMENT_GATEWAY_XENDIT   PaymentGateway = "XENDIT"
	PAYMENT_GATEWAY_SYSTEM   PaymentGateway = "SYSTEM"

	CAPTURE        TransactionStatus = "capture"
	SETTLEMENT     TransactionStatus = "settlement"
	PENDING        TransactionStatus = "pending"
	DENY           TransactionStatus = "deny"
	CANCEL         TransactionStatus = "cancel"
	EXPIRE         TransactionStatus = "expire"
	REFUND         TransactionStatus = "refund"
	PARTIAL_REFUND TransactionStatus = "partial_refund"
	AUTHORIZE      TransactionStatus = "authorize"

	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	PUT    RequestMethod = "PUT"
	PATCH  RequestMethod = "PATCH"
	DELETE RequestMethod = "DELETE"

	X_SUCCESS  XenditTransactionStatus = "SUCCESS"
	X_PENDING  XenditTransactionStatus = "PENDING"
	X_FAILED   XenditTransactionStatus = "FAILED"
	X_VOIDED   XenditTransactionStatus = "VOIDED"
	X_REVERSED XenditTransactionStatus = "REVERSED"

	XENDIT_QR_DANA_CHANNEL_CODE           ChannelCode = "QR_DANA"
	XENDIT_QR_LINKAJA_CHANNEL_CODE        ChannelCode = "QR_LINKAJA"
	XENDIT_EWALLET_LINKAJA_CHANNEL_CODE   ChannelCode = "EWALLET_LINKAJA"
	XENDIT_EWALLET_DANA_CHANNEL_CODE      ChannelCode = "EWALLET_DANA"
	XENDIT_EWALLET_OVO_CHANNEL_CODE       ChannelCode = "EWALLET_OVO"
	XENDIT_EWALLET_SHOPEEPAY_CHANNEL_CODE ChannelCode = "EWALLET_SHOPEEPAY"
	WALLET_CHANNEL_CODE                   ChannelCode = "WALLET"

	ITEM_TYPE_DIGITAL_PRODUCT  ItemType = "DIGITAL_PRODUCT"
	ITEM_TYPE_PHYSICAL_PRODUCT ItemType = "PHYSICAL_PRODUCT"
	ITEM_TYPE_DIGITAL_SERVICE  ItemType = "DIGITAL_SERVICE"
	ITEM_TYPE_PHYSICAL_SERVICE ItemType = "PHYSICAL_SERVICE"
	ITEM_TYPE_FEE              ItemType = "FEE"
	ITEM_TYPE_DELIVERY_FEE     ItemType = "DELIVERY_FEE"
	ITEM_TYPE_DISCOUNT         ItemType = "DISCOUNT"
)

func IsValidChannelCode(pm ChannelCode) bool {
	switch pm {
	case XENDIT_QR_DANA_CHANNEL_CODE, XENDIT_QR_LINKAJA_CHANNEL_CODE, XENDIT_EWALLET_LINKAJA_CHANNEL_CODE, XENDIT_EWALLET_DANA_CHANNEL_CODE, XENDIT_EWALLET_OVO_CHANNEL_CODE, XENDIT_EWALLET_SHOPEEPAY_CHANNEL_CODE, WALLET_CHANNEL_CODE:
		return true
	default:
		return false
	}
}

func IsValidPaymentMethod(pm PaymentMethod) bool {
	switch pm {
	case PAYMENT_METHOD_QR_CODE, PAYMENT_METHOD_EWALLET, PAYMENT_METHOD_WALLET:
		return true
	default:
		return false
	}
}

func IsValidPaymentGateway(pm PaymentGateway) bool {
	switch pm {
	case PAYMENT_GATEWAY_MIDTRANS, PAYMENT_GATEWAY_XENDIT, PAYMENT_GATEWAY_SYSTEM:
		return true
	default:
		return false
	}
}
