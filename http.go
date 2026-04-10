package main

/*
HTTP flow reference for localhost testing.

Base URL:
  http://localhost:8080

1. Health check
curl --request GET "http://localhost:8080/health"

2. Register a customer
curl --request POST "http://localhost:8080/api/v1/auth/register" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "name": "Customer One",
    "email": "customer@example.com",
    "password": "Password123",
    "role": "customer"
  }'

3. Register a driver
curl --request POST "http://localhost:8080/api/v1/auth/register" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "name": "Driver One",
    "email": "driver@example.com",
    "password": "Password123",
    "role": "driver"
  }'

4. Seed or update an admin directly in Postgres
   Password: Password123
INSERT INTO users (id, name, email, password_hash, role, created_at)
VALUES (
  gen_random_uuid(),
  'Admin One',
  'admin@example.com',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
  'admin',
  NOW()
)
ON CONFLICT (email)
DO UPDATE SET
  name = EXCLUDED.name,
  password_hash = EXCLUDED.password_hash,
  role = EXCLUDED.role;

5. Log in as customer
curl --request POST "http://localhost:8080/api/v1/auth/login" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "email": "customer@example.com",
    "password": "Password123"
  }'

6. Log in as driver
curl --request POST "http://localhost:8080/api/v1/auth/login" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "email": "driver@example.com",
    "password": "Password123"
  }'

7. Log in as admin
curl --request POST "http://localhost:8080/api/v1/auth/login" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "email": "admin@example.com",
    "password": "Password123"
  }'

Save these values from the login responses:
  CUSTOMER_TOKEN = customer.token
  DRIVER_TOKEN   = driver.token
  ADMIN_TOKEN    = admin.token
  DRIVER_USER_ID = driver.user.id

8. Customer posts a load
curl --request POST "http://localhost:8080/api/v1/loads" \
  --header "Authorization: Bearer CUSTOMER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "title": "Fresh Produce",
    "description": "Tomatoes for Nairobi market",
    "origin": "Eldoret",
    "destination": "Nairobi",
    "weight_kg": 1200
  }'

Save:
  LOAD_ID = response.id

9. Driver submits KYC
curl --request POST "http://localhost:8080/api/v1/drivers/kyc" \
  --header "Authorization: Bearer DRIVER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "national_id": "12345678",
    "license_number": "DL-90876",
    "truck_registration": "KDA123A"
  }'

10. Admin approves driver KYC
curl --request PATCH "http://localhost:8080/api/v1/drivers/DRIVER_USER_ID/kyc/review" \
  --header "Authorization: Bearer ADMIN_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "status": "approved"
  }'

11. Driver lists loads
curl --request GET "http://localhost:8080/api/v1/loads" \
  --header "Authorization: Bearer DRIVER_TOKEN"

12. Driver picks the posted load
curl --request POST "http://localhost:8080/api/v1/loads/LOAD_ID/pick" \
  --header "Authorization: Bearer DRIVER_TOKEN"

13. Driver marks the load in transit
curl --request PATCH "http://localhost:8080/api/v1/loads/LOAD_ID/status" \
  --header "Authorization: Bearer DRIVER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "status": "in_transit"
  }'

14. Driver marks the load delivered
curl --request PATCH "http://localhost:8080/api/v1/loads/LOAD_ID/status" \
  --header "Authorization: Bearer DRIVER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "status": "delivered"
  }'

15. Check current user profile
curl --request GET "http://localhost:8080/api/v1/auth/me" \
  --header "Authorization: Bearer CUSTOMER_TOKEN"
*/
