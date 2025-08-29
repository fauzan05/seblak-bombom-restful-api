# API Endpoints Reference

Quick reference guide for all available API endpoints.

## Base URL
- **Production**: `https://api.fznh-dev.my.id/api`
- **Development**: `http://localhost:8000/api`

## Authentication
- Method: Cookie-based JWT authentication
- Cookie Name: `access_token`
- Admin Creation Header: `X-Admin-Key`

## Guest Endpoints (No Authentication Required)

### User Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users/register` | Register new user |
| POST | `/users/login` | User login |
| POST | `/users/forgot-password` | Request password reset |
| POST | `/users/forgot-password/{passwordResetId}/validate` | Validate reset token |
| POST | `/users/forgot-password/{passwordResetId}/reset-password` | Reset password |
| GET | `/users/verify-email/{token}` | Verify email address |

### Public Content
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/categories` | Get all categories |
| GET | `/categories/{categoryId}` | Get category by ID |
| GET | `/products` | Get all products |
| GET | `/products/{productId}` | Get product by ID |
| GET | `/discount-coupons` | Get all discount coupons |
| GET | `/discount-coupons/{discountId}` | Get discount coupon by ID |
| GET | `/deliveries` | Get delivery options |
| GET | `/applications` | Get application info |

### Static Files
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/assets/*` | Static assets (CSS, JS, etc.) |
| GET | `/image/*` | Uploaded images |

### Testing
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/test-pusher` | Test real-time notifications |
| POST | `/applications-use-admin-key` | Create app config (requires X-Admin-Key) |

## Authenticated User Endpoints

### User Profile
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users/current` | Get current user profile |
| PATCH | `/users/current` | Update user profile |
| PATCH | `/users/current/password` | Change password |
| DELETE | `/users/current` | Delete account |
| DELETE | `/users/logout` | Logout |

### Address Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users/current/addresses` | Add new address |
| GET | `/users/current/addresses` | Get all user addresses |
| GET | `/users/current/addresses/{addressId}` | Get address by ID |
| PUT | `/users/current/addresses/{addressId}` | Update address |
| DELETE | `/users/current/addresses` | Delete addresses |

### Shopping Cart
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/carts` | Add item to cart |
| GET | `/carts` | Get current user's cart |
| PATCH | `/carts/cart-items/{cartItemId}` | Update cart item quantity |
| DELETE | `/carts/cart-items/{cartItemId}` | Remove item from cart |

### Orders
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/orders` | Create new order |
| GET | `/orders` | Get all orders |
| GET | `/orders/{orderId}` | Get order by ID |
| GET | `/orders/users/{userId}` | Get orders by user ID |
| PATCH | `/orders/{orderId}/status` | Update order status |
| GET | `/orders/{invoiceId}/invoice` | Get order invoice (HTML) |

### Product Reviews
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/reviews` | Create product review |

### Payments (Xendit)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/xendit/orders/qr-code/transaction` | Create QR code payment |
| GET | `/xendit/orders/{orderId}/qr-code/transaction` | Get QR transaction |
| POST | `/xendit/payouts/{userId}` | Create payout |
| POST | `/xendit/payout-request/{payoutId}/cancel` | Cancel payout |
| GET | `/xendit/payout-request/{payoutId}` | Get payout details |

### Wallet
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/wallets/withdraw-cust` | Request withdrawal |

### Standard Payouts
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/payouts/{userId}` | Create payout request |

## Admin Endpoints (Requires Admin Role)

### Category Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/categories` | Create category (multipart) |
| PUT | `/categories/{categoryId}` | Update category (multipart) |
| DELETE | `/categories` | Delete categories |

### Product Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/products` | Create product (multipart) |
| PUT | `/products/{productId}` | Update product (multipart) |
| DELETE | `/products` | Delete products |

### Discount Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/discount-coupons` | Create discount coupon |
| PUT | `/discount-coupons/{discountId}` | Update discount coupon |
| DELETE | `/discount-coupons` | Delete discount coupons |

### Delivery Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/deliveries` | Create delivery option |
| PUT | `/deliveries/{deliveryId}` | Update delivery option |
| DELETE | `/deliveries` | Delete delivery options |

### Application Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/applications` | Create/update app config |

### Financial Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/balance` | Get admin balance |
| PATCH | `/wallets/{withdrawRequestId}/withdraw-approval` | Approve/reject withdrawal |

## Xendit Webhook Endpoints

### Payment Callbacks
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/xendits/payment-request/notifications/callback` | Payment notification callback |
| POST | `/xendits/payout-request/notifications/callback` | Payout notification callback |

## Common Query Parameters

### Pagination
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)

### Product Filtering
- `category_id`: Filter by category
- `search`: Search by name
- `min_price`: Minimum price
- `max_price`: Maximum price

### User Registration
- `timezone`: User timezone (default: "UTC")
- `lang`: Language ("en" or "id", default: "en")

## Content Types

### JSON Endpoints
Most endpoints use `application/json`

### File Upload Endpoints
These endpoints use `multipart/form-data`:
- `POST /categories`
- `PUT /categories/{categoryId}`
- `POST /products`
- `PUT /products/{productId}`

## Response Format

### Success Response
```json
{
  "code": 200,
  "status": "success message",
  "data": <response_data>
}
```

### Paginated Response
```json
{
  "code": 200,
  "status": "success message",
  "data": <array_of_items>,
  "pagination": {
    "current_page": 1,
    "total_pages": 10,
    "total_items": 100,
    "data_per_pages": 10
  }
}
```

### Error Response
```json
{
  "code": 400,
  "status": "error message",
  "data": null
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Internal Server Error |

## Demo Credentials

### Customer Account
- Email: `cust1@email.com`
- Password: `Cust1Testing#`

### Admin Account
- Email: `admin1@email.com`
- Password: `Admin1Testing#`

## Payment Methods & Channels

### Payment Methods
- `qr_code`: QR Code payment
- `va`: Virtual Account
- `e_wallet`: E-Wallet

### Channel Codes (QR Code)
- `ID_DANA`: DANA
- `ID_LINKAJA`: LinkAja
- `ID_OVO`: OVO
- `ID_SHOPEEPAY`: ShopeePay

### Order Status Values
- `pending`: Order placed, awaiting payment
- `paid`: Payment confirmed
- `preparing`: Order being prepared
- `ready`: Order ready for pickup/delivery
- `completed`: Order completed
- `cancelled`: Order cancelled

### Payment Status Values
- `pending`: Payment pending
- `paid`: Payment successful
- `failed`: Payment failed
- `expired`: Payment expired