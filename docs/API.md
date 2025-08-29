# Seblak Bombom Restful API Documentation

## Table of Contents
- [Overview](#overview)
- [Base URL](#base-url)
- [Authentication](#authentication)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [API Endpoints](#api-endpoints)
  - [Authentication & User Management](#authentication--user-management)
  - [Products & Categories](#products--categories)
  - [Orders & Shopping Cart](#orders--shopping-cart)
  - [Payments (Xendit Integration)](#payments-xendit-integration)
  - [Admin Operations](#admin-operations)
  - [File Management](#file-management)

## Overview

The Seblak Bombom Restful API is a comprehensive backend service for a food delivery application. It provides functionality for user management, product catalog, order processing, payment integration with Xendit, and administrative operations.

### Key Features
- User registration and authentication
- Product catalog with categories
- Shopping cart management
- Order processing and tracking
- Payment integration with Xendit
- File upload and management
- Real-time notifications via Pusher
- Admin panel operations

## Base URL

### Production
```
https://api.fznh-dev.my.id/api
```

### Local Development
```
http://localhost:8000/api
```

## Authentication

The API uses cookie-based authentication with JWT tokens. After successful login, the access token is stored in an HTTP-only cookie named `access_token`.

### Authentication Headers
- **Cookie**: `access_token=<jwt_token>`

### User Roles
- **Customer**: Regular users who can place orders
- **Admin**: Administrators with additional privileges

### Admin Creation
To create an admin account, include the following header in the registration request:
- **X-Admin-Key**: `<admin_creation_key>` (set in environment variable `ADMIN_CREATION_KEY`)

## Response Format

All API responses follow a consistent JSON structure:

```json
{
  "code": 200,
  "status": "success message",
  "data": <response_data>
}
```

### Paginated Responses
For paginated endpoints:

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

## Error Handling

### HTTP Status Codes
- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

### Error Response Format
```json
{
  "code": 400,
  "status": "error message",
  "data": null
}
```

## API Endpoints

## Authentication & User Management

### Register User
Create a new user account.

**Endpoint:** `POST /users/register`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "password": "SecurePassword123",
  "role": "customer"
}
```

**Query Parameters:**
- `timezone` (optional): User timezone (default: "UTC")
- `lang` (optional): Language preference ("en" or "id", default: "en")

**Headers (for admin creation):**
- `X-Admin-Key`: Required for creating admin accounts

**Response:**
```json
{
  "code": 201,
  "status": "success to register an user",
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "role": "customer",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Login
Authenticate user and receive access token.

**Endpoint:** `POST /users/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Response:**
```json
{
  "code": 200,
  "status": "success to login",
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "role": "customer",
    "wallet": {
      "id": 1,
      "balance": 0.0
    },
    "cart": {
      "id": 1,
      "cart_items": []
    }
  }
}
```

### Get Current User
Get current authenticated user's profile.

**Endpoint:** `GET /users/current`

**Authentication:** Required

**Response:**
```json
{
  "code": 200,
  "status": "success to get current user",
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "role": "customer",
    "addresses": [],
    "wallet": {
      "id": 1,
      "balance": 0.0
    },
    "cart": {
      "id": 1,
      "cart_items": []
    }
  }
}
```

### Update User Profile
Update current user's profile information.

**Endpoint:** `PATCH /users/current`

**Authentication:** Required

**Request Body:**
```json
{
  "first_name": "John Updated",
  "last_name": "Doe Updated",
  "email": "john.updated@example.com",
  "phone": "+1234567891"
}
```

### Update Password
Change user's password.

**Endpoint:** `PATCH /users/current/password`

**Authentication:** Required

**Request Body:**
```json
{
  "old_password": "CurrentPassword123",
  "new_password": "NewPassword123"
}
```

### Logout
Logout current user.

**Endpoint:** `DELETE /users/logout`

**Authentication:** Required

### Delete Account
Delete current user's account.

**Endpoint:** `DELETE /users/current`

**Authentication:** Required

### Forgot Password
Request password reset.

**Endpoint:** `POST /users/forgot-password`

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

### Validate Password Reset
Validate password reset token.

**Endpoint:** `POST /users/forgot-password/{passwordResetId}/validate`

**Request Body:**
```json
{
  "token": "reset_token"
}
```

### Reset Password
Reset password with valid token.

**Endpoint:** `POST /users/forgot-password/{passwordResetId}/reset-password`

**Request Body:**
```json
{
  "token": "reset_token",
  "new_password": "NewPassword123"
}
```

### Verify Email
Verify email address with token.

**Endpoint:** `GET /users/verify-email/{token}`

## Address Management

### Add Address
Add new address for current user.

**Endpoint:** `POST /users/current/addresses`

**Authentication:** Required

**Request Body:**
```json
{
  "label": "Home",
  "recipient_name": "John Doe",
  "phone": "+1234567890",
  "street": "123 Main St",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "is_primary": true
}
```

### Get All Addresses
Get all addresses for current user.

**Endpoint:** `GET /users/current/addresses`

**Authentication:** Required

### Get Address by ID
Get specific address by ID.

**Endpoint:** `GET /users/current/addresses/{addressId}`

**Authentication:** Required

### Update Address
Update specific address.

**Endpoint:** `PUT /users/current/addresses/{addressId}`

**Authentication:** Required

**Request Body:** Same as Add Address

### Delete Addresses
Delete addresses by IDs.

**Endpoint:** `DELETE /users/current/addresses`

**Authentication:** Required

**Request Body:**
```json
{
  "ids": [1, 2, 3]
}
```

## Products & Categories

### Get All Categories
Get list of all product categories.

**Endpoint:** `GET /categories`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Response:**
```json
{
  "code": 200,
  "status": "success to get all categories",
  "data": [
    {
      "id": 1,
      "name": "Seblak",
      "description": "Spicy Indonesian crackers",
      "image": "/api/image/categories/seblak.jpg",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 1,
    "data_per_pages": 10
  }
}
```

### Get Category by ID
Get specific category by ID.

**Endpoint:** `GET /categories/{categoryId}`

### Get All Products
Get list of all products with optional filtering.

**Endpoint:** `GET /products`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `category_id` (optional): Filter by category ID
- `search` (optional): Search by product name
- `min_price` (optional): Minimum price filter
- `max_price` (optional): Maximum price filter

**Response:**
```json
{
  "code": 200,
  "status": "success to get all products",
  "data": [
    {
      "id": 1,
      "category": {
        "id": 1,
        "name": "Seblak"
      },
      "name": "Seblak Original",
      "description": "Traditional spicy seblak",
      "price": 15000,
      "stock": 100,
      "images": [
        {
          "id": 1,
          "url": "/api/image/products/seblak-original.jpg"
        }
      ],
      "reviews": [],
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Get Product by ID
Get specific product by ID.

**Endpoint:** `GET /products/{productId}`

## Orders & Shopping Cart

### Create Order
Create a new order.

**Endpoint:** `POST /orders`

**Authentication:** Required

**Request Body:**
```json
{
  "discount_id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "payment_method": "qr_code",
  "channel_code": "ID_DANA",
  "payment_gateway": "xendit",
  "is_delivery": true,
  "delivery_id": 1,
  "complete_address": "123 Main St, Jakarta",
  "note": "Extra spicy please",
  "order_products": [
    {
      "product_id": 1,
      "quantity": 2,
      "note": "No onions"
    }
  ]
}
```

### Get Order by ID
Get specific order by ID.

**Endpoint:** `GET /orders/{orderId}`

**Authentication:** Required

### Get Orders by User ID
Get all orders for specific user.

**Endpoint:** `GET /orders/users/{userId}`

**Authentication:** Required

### Update Order Status
Update order status (Admin only).

**Endpoint:** `PATCH /orders/{orderId}/status`

**Authentication:** Required (Admin)

**Request Body:**
```json
{
  "order_status": "preparing"
}
```

### Get All Orders
Get all orders with pagination.

**Endpoint:** `GET /orders`

**Authentication:** Required

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page
- `status` (optional): Filter by order status

### Show Invoice
Get order invoice as HTML.

**Endpoint:** `GET /orders/{invoiceId}/invoice`

**Authentication:** Required

### Add to Cart
Add item to shopping cart.

**Endpoint:** `POST /carts`

**Authentication:** Required

**Request Body:**
```json
{
  "product_id": 1,
  "quantity": 2,
  "note": "Extra spicy"
}
```

### Get Current Cart
Get current user's cart.

**Endpoint:** `GET /carts`

**Authentication:** Required

### Update Cart Item
Update quantity of cart item.

**Endpoint:** `PATCH /carts/cart-items/{cartItemId}`

**Authentication:** Required

**Request Body:**
```json
{
  "quantity": 3
}
```

### Delete Cart Item
Remove item from cart.

**Endpoint:** `DELETE /carts/cart-items/{cartItemId}`

**Authentication:** Required

## Product Reviews

### Create Review
Add review for a product.

**Endpoint:** `POST /reviews`

**Authentication:** Required

**Request Body:**
```json
{
  "product_id": 1,
  "rating": 5,
  "comment": "Excellent product!"
}
```

## Payments (Xendit Integration)

### Create QR Code Transaction
Create QR code payment transaction.

**Endpoint:** `POST /xendit/orders/qr-code/transaction`

**Authentication:** Required

**Request Body:**
```json
{
  "order_id": 1,
  "channel_code": "ID_DANA"
}
```

### Get QR Transaction
Get QR code transaction details.

**Endpoint:** `GET /xendit/orders/{orderId}/qr-code/transaction`

**Authentication:** Required

### Create Payout
Create payout request.

**Endpoint:** `POST /xendit/payouts/{userId}`

**Authentication:** Required

### Cancel Payout
Cancel payout request.

**Endpoint:** `POST /xendit/payout-request/{payoutId}/cancel`

**Authentication:** Required

### Get Payout by ID
Get payout details.

**Endpoint:** `GET /xendit/payout-request/{payoutId}`

**Authentication:** Required

### Payment Callback (Webhook)
Xendit payment notification callback.

**Endpoint:** `POST /xendits/payment-request/notifications/callback`

**Authentication:** Xendit Token Required

### Payout Callback (Webhook)
Xendit payout notification callback.

**Endpoint:** `POST /xendits/payout-request/notifications/callback`

**Authentication:** Xendit Token Required

## Wallet Management

### Withdraw Request (Customer)
Create withdrawal request.

**Endpoint:** `POST /wallets/withdraw-cust`

**Authentication:** Required

**Request Body:**
```json
{
  "amount": 100000,
  "bank_code": "BCA",
  "account_number": "1234567890",
  "account_holder_name": "John Doe"
}
```

### Withdraw Approval (Admin)
Approve or reject withdrawal request.

**Endpoint:** `PATCH /wallets/{withdrawRequestId}/withdraw-approval`

**Authentication:** Required (Admin)

**Request Body:**
```json
{
  "status": "approved"
}
```

## Discount Coupons

### Get All Discount Coupons
Get list of available discount coupons.

**Endpoint:** `GET /discount-coupons`

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page

### Get Discount Coupon by ID
Get specific discount coupon.

**Endpoint:** `GET /discount-coupons/{discountId}`

## Delivery Options

### Get All Delivery Options
Get available delivery options.

**Endpoint:** `GET /deliveries`

## Admin Operations

### Create Category
Create new product category.

**Endpoint:** `POST /categories`

**Authentication:** Required (Admin)

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `name`: Category name
- `description`: Category description
- `image`: Category image file

### Update Category
Update existing category.

**Endpoint:** `PUT /categories/{categoryId}`

**Authentication:** Required (Admin)

### Delete Categories
Delete categories by IDs.

**Endpoint:** `DELETE /categories`

**Authentication:** Required (Admin)

**Request Body:**
```json
{
  "ids": [1, 2, 3]
}
```

### Create Product
Create new product.

**Endpoint:** `POST /products`

**Authentication:** Required (Admin)

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `category_id`: Category ID
- `name`: Product name
- `description`: Product description
- `price`: Product price
- `stock`: Stock quantity
- `images`: Product image files (multiple)

### Update Product
Update existing product.

**Endpoint:** `PUT /products/{productId}`

**Authentication:** Required (Admin)

### Delete Products
Delete products by IDs.

**Endpoint:** `DELETE /products`

**Authentication:** Required (Admin)

### Create Discount Coupon
Create new discount coupon.

**Endpoint:** `POST /discount-coupons`

**Authentication:** Required (Admin)

### Update Discount Coupon
Update existing discount coupon.

**Endpoint:** `PUT /discount-coupons/{discountId}`

**Authentication:** Required (Admin)

### Delete Discount Coupons
Delete discount coupons.

**Endpoint:** `DELETE /discount-coupons`

**Authentication:** Required (Admin)

### Create Delivery Option
Create new delivery option.

**Endpoint:** `POST /deliveries`

**Authentication:** Required (Admin)

### Update Delivery Option
Update existing delivery option.

**Endpoint:** `PUT /deliveries/{deliveryId}`

**Authentication:** Required (Admin)

### Delete Delivery Options
Delete delivery options.

**Endpoint:** `DELETE /deliveries`

**Authentication:** Required (Admin)

### Get Admin Balance
Get admin balance from Xendit.

**Endpoint:** `GET /balance`

**Authentication:** Required (Admin)

### Application Configuration
Create or update application configuration.

**Endpoint:** `POST /applications`

**Authentication:** Required (Admin)

### Application Configuration (with Admin Key)
Create or update application configuration using admin key.

**Endpoint:** `POST /applications-use-admin-key`

**Headers:**
- `X-Admin-Key`: Admin creation key

## File Management

### Static Assets
Serve static assets (CSS, JS, images).

**Endpoint:** `GET /assets/*`

### Image Files
Serve uploaded images.

**Endpoint:** `GET /image/*`

**Example:** `GET /image/products/seblak-original.jpg`

## Application Info

### Get Application Info
Get general application information.

**Endpoint:** `GET /applications`

### Test Pusher
Test real-time notifications.

**Endpoint:** `GET /test-pusher`

**Query Parameters:**
- `message`: Test message to broadcast

## Status Codes Reference

### Order Status
- `pending` - Order placed, awaiting payment
- `paid` - Payment confirmed
- `preparing` - Order being prepared
- `ready` - Order ready for pickup/delivery
- `completed` - Order completed
- `cancelled` - Order cancelled

### Payment Status
- `pending` - Payment pending
- `paid` - Payment successful
- `failed` - Payment failed
- `expired` - Payment expired

### Payment Methods
- `qr_code` - QR Code payment
- `va` - Virtual Account
- `e_wallet` - E-Wallet

### Channel Codes (for QR Code)
- `ID_DANA` - DANA
- `ID_LINKAJA` - LinkAja
- `ID_OVO` - OVO
- `ID_SHOPEEPAY` - ShopeePay

## Demo Account

For testing purposes, you can use these demo accounts:

### Customer Account
- **Email:** `cust1@email.com`
- **Password:** `Cust1Testing#`

### Admin Account
- **Email:** `admin1@email.com`
- **Password:** `Admin1Testing#`

## Live Demo

- **Frontend:** [https://seblak.fznh-dev.my.id](https://seblak.fznh-dev.my.id)
- **Backend API:** [https://api.fznh-dev.my.id/api](https://api.fznh-dev.my.id/api)

**Note:** Demo data is reset periodically.