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
  '$2a$10$bOqXZ0l4YV5aWFcd03x/MuSrpuRpSa.O.0Hzq3dcnefhGlwxkRr4G',
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

8. Driver submits KYC
curl --request POST "http://localhost:8080/api/v1/drivers/kyc" \
  --header "Authorization: Bearer DRIVER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "national_id": "12345678",
    "license_number": "DL-90876",
    "truck_registration": "KDA123A"
  }'

9. Admin approves driver KYC
curl --request PATCH "http://localhost:8080/api/v1/drivers/DRIVER_USER_ID/kyc/review" \
  --header "Authorization: Bearer ADMIN_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "status": "approved"
  }'

10. Driver updates operational state for dispatch
curl --request PATCH "http://localhost:8080/api/v1/drivers/me/operations" \
  --header "Authorization: Bearer DRIVER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "years_experience": 8,
    "max_load_kg": 34000,
    "equipment_type": "lowbed",
    "is_online": true,
    "is_available": true,
    "latitude": -1.286389,
    "longitude": 36.817223
  }'

11. Customer posts a heavy-goods load
   The system will attempt automatic distance-based matching on create.
curl --request POST "http://localhost:8080/api/v1/loads" \
  --header "Authorization: Bearer CUSTOMER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "title": "Excavator Move",
    "description": "Lowbed movement for industrial equipment",
    "origin": "Mombasa",
    "destination": "Nairobi",
    "weight_kg": 12000,
    "pickup_latitude": -4.043477,
    "pickup_longitude": 39.668206,
    "dropoff_latitude": -1.286389,
    "dropoff_longitude": 36.817223,
    "pickup_at": "2026-04-11T08:00:00Z",
    "priority": 5,
    "cargo_type": "machinery",
    "equipment_type": "lowbed"
  }'

Save:
  LOAD_ID = response.id

If auto-matching succeeds, response will include:
  assigned_driver_id
  status = "matched"
  matching_score
  quoted_price_kes

12. Customer lists own loads
curl --request GET "http://localhost:8080/api/v1/loads" \
  --header "Authorization: Bearer CUSTOMER_TOKEN"

13. Driver lists available or assigned loads
curl --request GET "http://localhost:8080/api/v1/loads" \
  --header "Authorization: Bearer DRIVER_TOKEN"

14. If needed, admin runs batch dispatch optimization for tomorrow's loads
curl --request POST "http://localhost:8080/api/v1/dispatch/optimize" \
  --header "Authorization: Bearer ADMIN_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "pickup_date": "2026-04-11T00:00:00Z"
  }'

15. Driver picks matched load
curl --request POST "http://localhost:8080/api/v1/loads/LOAD_ID/pick" \
  --header "Authorization: Bearer DRIVER_TOKEN"

16. Driver sends tracking ping
curl --request POST "http://localhost:8080/api/v1/tracking/loads/LOAD_ID/events" \
  --header "Authorization: Bearer DRIVER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "latitude": -2.154007,
    "longitude": 37.990985,
    "speed_kph": 62,
    "heading_degrees": 135
  }'

17. WebSocket tracking stream
   Connect to:
   ws://localhost:8080/api/v1/tracking/ws?load_id=LOAD_ID

   Include Authorization header:
   Authorization: Bearer CUSTOMER_TOKEN

18. Driver marks load in transit
curl --request PATCH "http://localhost:8080/api/v1/loads/LOAD_ID/status" \
  --header "Authorization: Bearer DRIVER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "status": "in_transit"
  }'

19. Customer initiates M-Pesa checkout
curl --request POST "http://localhost:8080/api/v1/payments/loads/LOAD_ID/mpesa-checkout" \
  --header "Authorization: Bearer CUSTOMER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "phone_number": "254700000001"
  }'

20. Driver marks load delivered
curl --request PATCH "http://localhost:8080/api/v1/loads/LOAD_ID/status" \
  --header "Authorization: Bearer DRIVER_TOKEN" \
  --header "Content-Type: application/json" \
  --data-raw '{
    "status": "delivered"
  }'

21. Check current user profile
curl --request GET "http://localhost:8080/api/v1/auth/me" \
  --header "Authorization: Bearer CUSTOMER_TOKEN"

Notes:
- Apply migrations 000003, 000004, and 000005 before using dispatch, tracking, or payments.
- Drivers must be KYC-approved and operationally online before auto-matching can work.
- The M-Pesa integration is currently a sandbox-style provider abstraction that stores initiated payments.
- For future Swagger work, these routes and payloads are now stable enough to annotate.
*/
