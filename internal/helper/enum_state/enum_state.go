package enum_state

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
type WalletFlowType string
type WalletTransactionSource string
type WalletTransactionType string
type WalletPaymentMethod string
type WalletTransactionStatus string
type WalletWithdrawRequest string

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
	PAYMENT_METHOD_CASH    PaymentMethod = "CASH"

	PAYMENT_GATEWAY_XENDIT PaymentGateway = "XENDIT"
	PAYMENT_GATEWAY_SYSTEM PaymentGateway = "SYSTEM"

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

	WALLET_FLOW_TYPE_DEBIT  WalletFlowType = "debit"
	WALLET_FLOW_TYPE_CREDIT WalletFlowType = "credit"

	WALLET_TRANSACTION_TYPE_TOP_UP           WalletTransactionType = "top_up"
	WALLET_TRANSACTION_TYPE_ORDER_PAYMENT    WalletTransactionType = "order_payment"
	WALLET_TRANSACTION_TYPE_ORDER_REFUND     WalletTransactionType = "order_refund"
	WALLET_TRANSACTION_TYPE_WITHDRAW         WalletTransactionType = "withdraw"
	WALLET_TRANSACTION_TYPE_ADMIN_ADJUSTMENT WalletTransactionType = "admin_adjustment"
	WALLET_TRANSACTION_TYPE_CASHBACK         WalletTransactionType = "cashback"
	WALLET_TRANSACTION_TYPE_TRANSFER_IN      WalletTransactionType = "transfer_in"
	WALLET_TRANSACTION_TYPE_TRANSFER_OUT     WalletTransactionType = "transfer_out"

	WALLET_TRANSACTION_STATUS_PENDING    WalletTransactionStatus = "pending"
	WALLET_TRANSACTION_STATUS_PROCESSING WalletTransactionStatus = "processing"
	WALLET_TRANSACTION_STATUS_COMPLETED  WalletTransactionStatus = "completed"
	WALLET_TRANSACTION_STATUS_FAILED     WalletTransactionStatus = "failed"
	WALLET_TRANSACTION_STATUS_CANCELLED  WalletTransactionStatus = "cancelled"

	WALLET_WITHDRAW_REQUEST_METHOD_CASH          WalletWithdrawRequest = "cash"
	WALLET_WITHDRAW_REQUEST_METHOD_BANK_TRANSFER WalletWithdrawRequest = "bank_transfer"
	WALLET_WITHDRAW_REQUEST_STATUS_PENDING       WalletWithdrawRequest = "pending"
	WALLET_WITHDRAW_REQUEST_STATUS_APPROVED      WalletWithdrawRequest = "approved"
	WALLET_WITHDRAW_REQUEST_STATUS_REJECTED      WalletWithdrawRequest = "rejected"
)

func IsValidChannelCode(pc ChannelCode) bool {
	switch pc {
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

func IsValidPaymentGateway(pg PaymentGateway) bool {
	switch pg {
	case PAYMENT_GATEWAY_XENDIT, PAYMENT_GATEWAY_SYSTEM:
		return true
	default:
		return false
	}
}
