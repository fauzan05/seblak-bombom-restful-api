# API Testing Guide

This guide provides examples for testing the Seblak Bombom API using various tools like Postman, curl, or any HTTP client.

## Quick Start

### 1. Base URL
- **Production**: `https://api.fznh-dev.my.id/api`
- **Local**: `http://localhost:8000/api`

### 2. Demo Accounts
Use these accounts for testing:
- **Customer**: `cust1@email.com` / `Cust1Testing#`
- **Admin**: `admin1@email.com` / `Admin1Testing#`

## Common API Testing Scenarios

### Authentication Flow

#### 1. User Registration
```bash
curl -X POST https://api.fznh-dev.my.id/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com",
    "phone": "+1234567890",
    "password": "TestPassword123",
    "role": "customer"
  }'
```

#### 2. User Login
```bash
curl -X POST https://api.fznh-dev.my.id/api/users/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "cust1@email.com",
    "password": "Cust1Testing#"
  }'
```

#### 3. Get Current User (with cookies)
```bash
curl -X GET https://api.fznh-dev.my.id/api/users/current \
  -b cookies.txt
```

### Product Browsing

#### 1. Get All Categories
```bash
curl -X GET https://api.fznh-dev.my.id/api/categories
```

#### 2. Get All Products
```bash
curl -X GET "https://api.fznh-dev.my.id/api/products?page=1&limit=10"
```

#### 3. Get Product by ID
```bash
curl -X GET https://api.fznh-dev.my.id/api/products/1
```

### Shopping Cart

#### 1. Add to Cart
```bash
curl -X POST https://api.fznh-dev.my.id/api/carts \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "product_id": 1,
    "quantity": 2,
    "note": "Extra spicy"
  }'
```

#### 2. Get Current Cart
```bash
curl -X GET https://api.fznh-dev.my.id/api/carts \
  -b cookies.txt
```

### Order Management

#### 1. Create Order
```bash
curl -X POST https://api.fznh-dev.my.id/api/orders \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
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
  }'
```

#### 2. Get Order by ID
```bash
curl -X GET https://api.fznh-dev.my.id/api/orders/1 \
  -b cookies.txt
```

### Payment Processing

#### 1. Create QR Code Payment
```bash
curl -X POST https://api.fznh-dev.my.id/api/xendit/orders/qr-code/transaction \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "order_id": 1,
    "channel_code": "ID_DANA"
  }'
```

### Admin Operations

#### 1. Admin Login
```bash
curl -X POST https://api.fznh-dev.my.id/api/users/login \
  -H "Content-Type: application/json" \
  -c admin_cookies.txt \
  -d '{
    "email": "admin1@email.com",
    "password": "Admin1Testing#"
  }'
```

#### 2. Create Category (Admin)
```bash
curl -X POST https://api.fznh-dev.my.id/api/categories \
  -b admin_cookies.txt \
  -F "name=New Category" \
  -F "description=Category Description" \
  -F "image=@category-image.jpg"
```

#### 3. Create Product (Admin)
```bash
curl -X POST https://api.fznh-dev.my.id/api/products \
  -b admin_cookies.txt \
  -F "category_id=1" \
  -F "name=New Product" \
  -F "description=Product Description" \
  -F "price=25000" \
  -F "stock=100" \
  -F "images=@product-image1.jpg" \
  -F "images=@product-image2.jpg"
```

## Postman Collection

### Environment Variables
Set up these variables in Postman:
- `base_url`: `https://api.fznh-dev.my.id/api`
- `customer_email`: `cust1@email.com`
- `customer_password`: `Cust1Testing#`
- `admin_email`: `admin1@email.com`
- `admin_password`: `Admin1Testing#`

### Authentication Setup
1. Create a login request
2. In the Tests tab, add this script to automatically save the token:
```javascript
if (responseCode.code === 200) {
    // Save token from cookie if needed
    var cookies = pm.cookies.jar();
    // Tokens are automatically handled via cookies
}
```

### Common Test Scripts

#### Check Response Status
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has success status", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.code).to.eql(200);
});
```

#### Validate Response Structure
```javascript
pm.test("Response has required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('code');
    pm.expect(jsonData).to.have.property('status');
    pm.expect(jsonData).to.have.property('data');
});
```

## Error Handling Examples

### Common Error Responses

#### 401 Unauthorized
```json
{
  "code": 401,
  "status": "missing token",
  "data": null
}
```

#### 400 Bad Request
```json
{
  "code": 400,
  "status": "validation error: field is required",
  "data": null
}
```

#### 404 Not Found
```json
{
  "code": 404,
  "status": "resource not found",
  "data": null
}
```

## File Upload Testing

### Upload Product Image
Use `multipart/form-data` for file uploads:

```bash
curl -X POST https://api.fznh-dev.my.id/api/products \
  -b admin_cookies.txt \
  -H "Content-Type: multipart/form-data" \
  -F "category_id=1" \
  -F "name=Test Product" \
  -F "description=Test Description" \
  -F "price=15000" \
  -F "stock=50" \
  -F "images=@/path/to/image1.jpg" \
  -F "images=@/path/to/image2.jpg"
```

## Testing Tips

1. **Authentication**: Most endpoints require authentication via cookies
2. **File Uploads**: Use `multipart/form-data` for endpoints with file uploads
3. **Pagination**: Check pagination parameters for list endpoints
4. **Error Handling**: Always test with invalid data to verify error responses
5. **Admin Endpoints**: Use admin credentials for admin-only operations

## Automated Testing

For automated testing, consider:
1. Setting up test data fixtures
2. Using environment-specific configurations
3. Implementing cleanup procedures
4. Testing both success and failure scenarios
5. Validating response schemas

## Real-time Testing

The API includes Pusher integration for real-time features. Test the WebSocket connection:
```bash
curl -X GET "https://api.fznh-dev.my.id/api/test-pusher?message=Hello%20World"
```

## Demo Environment

Remember that the demo environment:
- Resets data periodically
- May have rate limiting
- Uses simulated payment processing
- Should not be used for production testing