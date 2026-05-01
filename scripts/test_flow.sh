#!/bin/bash

# Food Delivery Platform E2E Flow Test
# Using Gateway on port 8080

BASE_URL="http://localhost:8080/api/v1"
CUSTOMER_ID="550e8400-e29b-41d4-a716-446655440000" 
RESTAURANT_ID="rest-1" # Wait, rest-1 is not a UUID.

# Let's use valid UUIDs everywhere.
RESTAURANT_ID="00000000-0000-0000-0000-000000000001"
ITEM_ID="00000000-0000-0000-0000-000000000101"

echo "--- 1. Searching for restaurants (lat=55.75, lon=37.62) ---"
curl -s -X GET "$BASE_URL/restaurants?lat=55.75&lon=37.62&radius=5000" | json_pp || echo "Search failed"
echo -e "\n"

echo "--- 2. Fetching menu for Burger King ---"
curl -s -X GET "$BASE_URL/restaurants/$RESTAURANT_ID/menu" | json_pp || echo "Menu failed"
echo -e "\n"

echo "--- 3. Adding item to cart ---"
curl -s -X POST "$BASE_URL/carts/$CUSTOMER_ID?restaurantId=$RESTAURANT_ID" \
     -H "Content-Type: application/json" \
     -d "{\"menu_item_id\": \"$ITEM_ID\", \"quantity\": 2}" | json_pp || echo "Add to cart failed"
echo -e "\n"

echo "--- 4. Placing an order (Triggering Saga) ---"
ORDER_RESP=$(curl -s -X POST "$BASE_URL/orders" \
     -H "Content-Type: application/json" \
     -d "{
       \"customer_id\": \"$CUSTOMER_ID\",
       \"restaurant_id\": \"$RESTAURANT_ID\",
       \"delivery_address\": \"Red Square, 1\",
       \"items\": [{\"menu_item_id\": \"$ITEM_ID\", \"quantity\": 2, \"price\": 300}]
     }")
echo $ORDER_RESP | json_pp || echo $ORDER_RESP

ORDER_ID=$(echo $ORDER_RESP | grep -o '"order_id":"[^"]*' | cut -d'"' -f4)

if [ -z "$ORDER_ID" ]; then
    echo "ERROR: Failed to get ORDER_ID"
    exit 1
fi

echo -e "\nORDER_ID: $ORDER_ID"
echo "Waiting 7 seconds for Payment -> Kafka Outbox -> Delivery Assignment..."
sleep 7

echo -e "\n--- 5. Checking order status ---"
curl -s -X GET "$BASE_URL/orders/$ORDER_ID/status" | json_pp || echo "Status check failed"
echo -e "\n"

echo "--- 6. (Optional) Batching test: Placing another order from rest-1 ---"
# This should ideally be assigned to the same courier if they have capacity
curl -s -X POST "$BASE_URL/orders" \
     -H "Content-Type: application/json" \
     -d "{
       \"customer_id\": \"other-user\",
       \"restaurant_id\": \"rest-1\",
       \"delivery_address\": \"Arbat St, 10\",
       \"items\": [{\"menu_item_id\": \"m1\", \"quantity\": 1, \"price\": 300}]
     }" > /dev/null

echo "Test Flow Execution Finished!"
echo "Check logs: docker-compose logs -f delivery-service"
