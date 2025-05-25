package helper

type Role string
type PaymentMethod string
type PaymentStatus string
type OrderStatus string
type DiscountType string
type NotificationType string
type WalletStatus string
type TransactionStatus string
type RequestMethod string
type XenditTransactionStatus string
type ChannelCode string
type PaymentGateway string
type ItemType string
type PayoutStatus string
type PayoutMethod string
type PDFPageSize string
type PDFOrientation string
type Languange string

const (
	// role
	ADMIN    Role = "admin"
	CUSTOMER Role = "customer"
	// Payment Status
	PAID_PAYMENT      PaymentStatus = "paid"      // Pembayaran sukses
	PENDING_PAYMENT   PaymentStatus = "pending"   // Menunggu konfirmasi
	CANCELLED_PAYMENT PaymentStatus = "cancelled" // Dibatalkan oleh pengguna
	EXPIRED_PAYMENT   PaymentStatus = "expired"   // Waktu pembayaran habis
	FAILED_PAYMENT    PaymentStatus = "failed"    // Pembayaran gagal

	// order status
	ORDER_PENDING                OrderStatus = "pending_order"
	ORDER_RECEIVED               OrderStatus = "order_received"
	ORDER_BEING_DELIVERED        OrderStatus = "order_being_delivered"
	ORDER_DELIVERED              OrderStatus = "order_delivered"
	READY_FOR_PICKUP             OrderStatus = "ready_for_pickup"
	ORDER_REJECTED               OrderStatus = "order_rejected"
	ORDER_CANCELLED              OrderStatus = "order_cancelled"
	ORDER_CANCELLATION_REQUESTED OrderStatus = "order_cancellation_requested"
	DELIVERY_FAILED              OrderStatus = "delivery_failed"

	// discount type
	NOMINAL DiscountType = "nominal"
	PERCENT DiscountType = "percent"
	// notification type
	AUTHENTICATION NotificationType = "authentication"
	TRANSACTION    NotificationType = "transaction"
	PROMOTION      NotificationType = "promotion"
	// wallet status
	ACTIVE_WALLET  WalletStatus = "active"
	INACIVE_WALLET WalletStatus = "inactive"

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

	PAYOUT_PENDING   PayoutStatus = "pending"
	PAYOUT_ACCEPTED  PayoutStatus = "accepted"
	PAYOUT_CANCELLED PayoutStatus = "cancelled"
	PAYOUT_FAILED    PayoutStatus = "failed"
	PAYOUT_SUCCEEDED PayoutStatus = "succeeded"
	PAYOUT_EXPIRED   PayoutStatus = "expired"
	PAYOUT_REFUNDED  PayoutStatus = "refunded"

	PAYOUT_METHOD_ONLINE  PayoutMethod = "online"
	PAYOUT_METHOD_OFFLINE PayoutMethod = "offline"

	A4     PDFPageSize = "A4"
	LETTER PDFPageSize = "Letter"

	PORTRAIT  PDFOrientation = "Portrait"
	LANDSCAPE PDFOrientation = "Landscape"

	INDONESIA Languange = "id"
	ENGLISH   Languange = "en"
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
