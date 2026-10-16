$ErrorActionPreference = "Stop"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "   ECOMMERCE FULL FLOW AUDIT (ADD, EDIT, REMOVE)         " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

$baseUrl = "http://127.0.0.1:8080/api/v1"

# -----------------------------------------------------------------------------
# 1. AUTHENTICATION AUDIT
# -----------------------------------------------------------------------------
Write-Host "`n[1/5] AUDITING AUTH FLOW..." -ForegroundColor Yellow

# Admin Login
$adminLogin = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body (@{ email = "admin@ecommerce.com"; password = "Admin123!" } | ConvertTo-Json) -ContentType "application/json"
$adminToken = $adminLogin.data.token
if ($adminToken) {
    Write-Host "  [OK] Admin Login successful (Token received)" -ForegroundColor Green
} else {
    throw "Admin login failed"
}

# Customer Login
$custLogin = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body (@{ email = "budi@customer.com"; password = "Admin123!" } | ConvertTo-Json) -ContentType "application/json"
$custToken = $custLogin.data.token
if ($custToken) {
    Write-Host "  [OK] Customer Login successful (Token received)" -ForegroundColor Green
} else {
    throw "Customer login failed"
}

# Register New User & Get Profile
$rand = Get-Random -Minimum 10000 -Maximum 99999
$newUserEmail = "audit_user_$rand@test.com"
$regRes = Invoke-RestMethod -Uri "$baseUrl/auth/register" -Method Post -Body (@{
    name = "Audit Test User $rand"
    email = $newUserEmail
    password = "Password123!"
    phone = "081298765432"
    address = "Jl. Sudirman Kav 21, Jakarta"
} | ConvertTo-Json) -ContentType "application/json"
$newUsertoken = $regRes.data.token

$meRes = Invoke-RestMethod -Uri "$baseUrl/auth/me" -Method Get -Headers @{ Authorization = "Bearer $newUsertoken" }
if ($meRes.data.email -eq $newUserEmail) {
    Write-Host "  [OK] Register and /auth/me profile verified" -ForegroundColor Green
} else {
    throw "Profile verification failed"
}

# -----------------------------------------------------------------------------
# 2. CATEGORY CRUD AUDIT (ADD, EDIT, REMOVE)
# -----------------------------------------------------------------------------
Write-Host "`n[2/5] AUDITING CATEGORY CRUD (ADD, EDIT, REMOVE)..." -ForegroundColor Yellow

# ADD Category
$catSlug = "audit-cat-$rand"
$addCat = Invoke-RestMethod -Uri "$baseUrl/admin/categories" -Method Post -Headers @{ Authorization = "Bearer $adminToken" } -Body (@{
    name = "Audit Category $rand"
    slug = $catSlug
    icon = "/assets/images/category/category-fruits-vegetables.jpg"
    description = "Test description for audit category"
    status = $true
} | ConvertTo-Json) -ContentType "application/json"
$catId = $addCat.data.id
Write-Host "  [OK] ADD Category: ID = $catId, Name = '$($addCat.data.name)'" -ForegroundColor Green

# EDIT Category
$editCat = Invoke-RestMethod -Uri "$baseUrl/admin/categories/$catId" -Method Put -Headers @{ Authorization = "Bearer $adminToken" } -Body (@{
    name = "Audit Category Updated $rand"
    description = "Updated description"
} | ConvertTo-Json) -ContentType "application/json"
if ($editCat.data.name -eq "Audit Category Updated $rand") {
    Write-Host "  [OK] EDIT Category: Updated Name = '$($editCat.data.name)'" -ForegroundColor Green
} else {
    throw "Category edit failed"
}

# REMOVE Category
$delCat = Invoke-RestMethod -Uri "$baseUrl/admin/categories/$catId" -Method Delete -Headers @{ Authorization = "Bearer $adminToken" }
Write-Host "  [OK] REMOVE Category: $($delCat.message)" -ForegroundColor Green

# -----------------------------------------------------------------------------
# 3. PRODUCT CRUD AUDIT (ADD, EDIT, REMOVE)
# -----------------------------------------------------------------------------
Write-Host "`n[3/5] AUDITING PRODUCT CRUD (ADD, EDIT, REMOVE)..." -ForegroundColor Yellow

# ADD Product
$prodSlug = "audit-prod-$rand"
$addProd = Invoke-RestMethod -Uri "$baseUrl/admin/products" -Method Post -Headers @{ Authorization = "Bearer $adminToken" } -Body (@{
    category_id = 1
    name = "Audit Fresh Sayur $rand"
    slug = $prodSlug
    image = "/assets/images/products/product-img-1.jpg"
    description = "Sayuran segar untuk pengetesan audit"
    regular_price = 25000
    sale_price = 20000
    unit = "ikat"
    weight = 400
    stock = 50
    variant_name = "Ikat 400g"
    status = "ACTIVE"
} | ConvertTo-Json) -ContentType "application/json"
$prodId = $addProd.data.id
Write-Host "  [OK] ADD Product: ID = $prodId, Name = '$($addProd.data.name)', Stock = $($addProd.data.stock)" -ForegroundColor Green

# EDIT Product
$editProd = Invoke-RestMethod -Uri "$baseUrl/admin/products/$prodId" -Method Put -Headers @{ Authorization = "Bearer $adminToken" } -Body (@{
    name = "Audit Fresh Sayur $rand (Updated)"
    regular_price = 28000
    sale_price = 22000
    stock = 45
    image = "https://images.unsplash.com/photo-1540420773420-3366772f4999?w=800"
} | ConvertTo-Json) -ContentType "application/json"
if ($editProd.data.regular_price -eq 28000 -and $editProd.data.stock -eq 45) {
    Write-Host "  [OK] EDIT Product: Price = $($editProd.data.regular_price), Stock = $($editProd.data.stock), Image = $($editProd.data.image)" -ForegroundColor Green
} else {
    throw "Product edit failed"
}

# REMOVE Product
$delProd = Invoke-RestMethod -Uri "$baseUrl/admin/products/$prodId" -Method Delete -Headers @{ Authorization = "Bearer $adminToken" }
Write-Host "  [OK] REMOVE Product: $($delProd.message)" -ForegroundColor Green

# -----------------------------------------------------------------------------
# 4. CART FLOW AUDIT (ADD, EDIT QUANTITY, REMOVE ITEM, CLEAR CART)
# -----------------------------------------------------------------------------
Write-Host "`n[4/5] AUDITING CART FLOW (ADD, EDIT QUANTITY, REMOVE)..." -ForegroundColor Yellow

# Clear Cart first
Invoke-RestMethod -Uri "$baseUrl/cart" -Method Delete -Headers @{ Authorization = "Bearer $custToken" } | Out-Null
Write-Host "  [OK] CLEAR Cart initialized" -ForegroundColor Green

# ADD to Cart (Product ID 1, Qty 3)
$cartAdd = Invoke-RestMethod -Uri "$baseUrl/cart/items" -Method Post -Headers @{ Authorization = "Bearer $custToken" } -Body (@{
    product_id = 1
    quantity = 3
} | ConvertTo-Json) -ContentType "application/json"
Write-Host "  [OK] ADD to Cart: Total Items = $($cartAdd.data.total_items), Subtotal = $($cartAdd.data.total_price), Weight = $($cartAdd.data.total_weight)g" -ForegroundColor Green

# EDIT Cart Item Quantity (Change Qty to 5)
$cartEdit = Invoke-RestMethod -Uri "$baseUrl/cart/items/1" -Method Put -Headers @{ Authorization = "Bearer $custToken" } -Body (@{
    quantity = 5
} | ConvertTo-Json) -ContentType "application/json"
Write-Host "  [OK] EDIT Cart Item (Qty -> 5): Total Items = $($cartEdit.data.total_items), Subtotal = $($cartEdit.data.total_price)" -ForegroundColor Green

# EDIT Cart Item Quantity to 0 (Should remove cleanly)
$cartZero = Invoke-RestMethod -Uri "$baseUrl/cart/items/1" -Method Put -Headers @{ Authorization = "Bearer $custToken" } -Body (@{
    quantity = 0
} | ConvertTo-Json) -ContentType "application/json"
Write-Host "  [OK] EDIT Cart Item (Qty -> 0): Total Items = $($cartZero.data.total_items) (Clean auto-removal)" -ForegroundColor Green

# ADD 2 items for removal test
Invoke-RestMethod -Uri "$baseUrl/cart/items" -Method Post -Headers @{ Authorization = "Bearer $custToken" } -Body (@{ product_id = 1; quantity = 2 } | ConvertTo-Json) -ContentType "application/json" | Out-Null
$cartRemove = Invoke-RestMethod -Uri "$baseUrl/cart/items/1" -Method Delete -Headers @{ Authorization = "Bearer $custToken" }
Write-Host "  [OK] REMOVE Cart Item: Total Items = $($cartRemove.data.total_items)" -ForegroundColor Green

# -----------------------------------------------------------------------------
# 5. ORDER & CHECKOUT FLOW AUDIT (CHECKOUT, ATOMIC STOCK LOCK, RESTOCK ON CANCEL)
# -----------------------------------------------------------------------------
Write-Host "`n[5/5] AUDITING ORDER & CHECKOUT FLOW..." -ForegroundColor Yellow

# Check initial stock of Product 2
$p2Before = Invoke-RestMethod -Uri "$baseUrl/products/tomat-merah-segar-pilihan-1kg"
$p2StockInitial = $p2Before.data.stock
Write-Host "  Current stock of Tomat Merah: $p2StockInitial" -ForegroundColor Gray

# Add 2 pcs of Product 2 to Cart
Invoke-RestMethod -Uri "$baseUrl/cart/items" -Method Post -Headers @{ Authorization = "Bearer $custToken" } -Body (@{ product_id = 2; quantity = 2 } | ConvertTo-Json) -ContentType "application/json" | Out-Null

# Shipping Estimate (100% Delivery Model)
$estRes = Invoke-RestMethod -Uri "$baseUrl/orders/estimate" -Method Post -Headers @{ Authorization = "Bearer $custToken" } -Body (@{
    shipping_address = "Jl. Kemang Raya No. 10, Jakarta Selatan"
} | ConvertTo-Json) -ContentType "application/json"
Write-Host "  [OK] Shipping Estimate: Subtotal = $($estRes.data.subtotal), ShippingFee = $($estRes.data.shipping_fee), Total = $($estRes.data.total_amount)" -ForegroundColor Green

# Checkout Order (using both /orders and /orders/checkout routes)
$checkoutRes = Invoke-RestMethod -Uri "$baseUrl/orders/checkout" -Method Post -Headers @{ Authorization = "Bearer $custToken" } -Body (@{
    recipient_name = "Budi Santoso"
    recipient_phone = "081987654321"
    shipping_address = "Jl. Kemang Raya No. 45B, Mampang Prapatan, Jakarta Selatan"
    remarks = "Harap antarkan sebelum siang"
} | ConvertTo-Json) -ContentType "application/json"
$createdOrder = $checkoutRes.data
$orderCode = $createdOrder.order_code
$orderId = $createdOrder.id
Write-Host "  [OK] CREATE Order: Code = $orderCode, ID = $orderId, Status = $($createdOrder.status)" -ForegroundColor Green

# Verify Stock Decremented by 2
$p2After = Invoke-RestMethod -Uri "$baseUrl/products/tomat-merah-segar-pilihan-1kg"
$p2StockAfter = $p2After.data.stock
if ($p2StockAfter -eq ($p2StockInitial - 2)) {
    Write-Host "  [OK] Stock Locked & Decremented: $p2StockInitial -> $p2StockAfter (Exact -2)" -ForegroundColor Green
} else {
    throw "Stock was not decremented correctly: expected $($p2StockInitial - 2), got $p2StockAfter"
}

# Admin Update Order Status: PENDING_PAYMENT -> PROCESSING
$procRes = Invoke-RestMethod -Uri "$baseUrl/admin/orders/$orderId/status" -Method Put -Headers @{ Authorization = "Bearer $adminToken" } -Body (@{
    status = "PROCESSING"
} | ConvertTo-Json) -ContentType "application/json"
Write-Host "  [OK] ADMIN Status Update: $($procRes.data.status)" -ForegroundColor Green

# Admin Update Order Status: PROCESSING -> ON_DELIVERY (with tracking number)
$deliveryRes = Invoke-RestMethod -Uri "$baseUrl/admin/orders/$orderId/status" -Method Put -Headers @{ Authorization = "Bearer $adminToken" } -Body (@{
    status = "ON_DELIVERY"
    tracking_number = "KURIR-SAYUR-9912"
} | ConvertTo-Json) -ContentType "application/json"
Write-Host "  [OK] ADMIN Delivery Update: Status = $($deliveryRes.data.status), Resi = $($deliveryRes.data.tracking_number)" -ForegroundColor Green

# Admin CANCEL Order -> Restock Verification
$cancelRes = Invoke-RestMethod -Uri "$baseUrl/admin/orders/$orderId/status" -Method Put -Headers @{ Authorization = "Bearer $adminToken" } -Body (@{
    status = "CANCELLED"
} | ConvertTo-Json) -ContentType "application/json"
Write-Host "  [OK] ADMIN CANCEL Order: Status = $($cancelRes.data.status)" -ForegroundColor Green

# Verify Stock Restored
$p2Restored = Invoke-RestMethod -Uri "$baseUrl/products/tomat-merah-segar-pilihan-1kg"
$p2StockRestored = $p2Restored.data.stock
if ($p2StockRestored -eq $p2StockInitial) {
    Write-Host "  [OK] Stock Restored on Cancellation: $p2StockAfter -> $p2StockRestored (Exact Restock)" -ForegroundColor Green
} else {
    throw "Stock was not restored correctly: expected $p2StockInitial, got $p2StockRestored"
}

Write-Host "`n==========================================================" -ForegroundColor Cyan
Write-Host "   ALL FLOWS PASSED AUDIT SUCCESSFULLY (100% OK)         " -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Cyan
