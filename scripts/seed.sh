#!/bin/bash

# Food Delivery Platform Seeder
# Populates Elasticsearch and InMemory repositories via API Gateway

GATEWAY_URL="http://localhost:8080/api/v1"
ES_URL="http://localhost:9200"

echo "Creating Elasticsearch Index with Mapping..."
curl -s -X DELETE "$ES_URL/restaurants" > /dev/null
curl -s -X PUT "$ES_URL/restaurants" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "location": { "type": "geo_point" },
      "cuisines": { "type": "keyword" },
      "rating": { "type": "float" }
    }
  }
}'
echo " Index Created."

echo "Seeding Restaurants..."

# Burger King
curl -s -X POST "$GATEWAY_URL/restaurants" \
     -H "Content-Type: application/json" \
     -d '{
       "restaurant_id": "00000000-0000-0000-0000-000000000001",
       "name": "Burger King",
       "location": { "lat": 55.75, "lon": 37.62 },
       "cuisines": ["burgers", "fastfood"],
       "rating": 4.5,
       "menu": [
         { "item_id": "00000000-0000-0000-0000-000000000101", "name": "Whopper", "price": 300, "available": true },
         { "item_id": "00000000-0000-0000-0000-000000000102", "name": "French Fries", "price": 150, "available": true }
       ]
     }'
echo " rest-1 OK"

# Pizza Hut
curl -s -X POST "$GATEWAY_URL/restaurants" \
     -H "Content-Type: application/json" \
     -d '{
       "restaurant_id": "00000000-0000-0000-0000-000000000002",
       "name": "Pizza Hut",
       "location": { "lat": 55.76, "lon": 37.63 },
       "cuisines": ["pizza", "italian"],
       "rating": 4.2,
       "menu": [
         { "item_id": "00000000-0000-0000-0000-000000000201", "name": "Pepperoni", "price": 600, "available": true },
         { "item_id": "00000000-0000-0000-0000-000000000202", "name": "Margherita", "price": 500, "available": true }
       ]
     }'
echo " rest-2 OK"

echo -e "\nSeeding Couriers..."

# Ivan
curl -s -X POST "$GATEWAY_URL/couriers" \
     -H "Content-Type: application/json" \
     -d '{ "id": "00000000-0000-0000-0000-000000000901", "name": "Ivan", "location": { "lat": 55.75, "lon": 37.62 }, "is_available": true }'
echo " courier-1 OK"

# Petr
curl -s -X POST "$GATEWAY_URL/couriers" \
     -H "Content-Type: application/json" \
     -d '{ "id": "00000000-0000-0000-0000-000000000902", "name": "Petr", "location": { "lat": 55.80, "lon": 37.70 }, "is_available": true }'
echo " courier-2 OK"

echo -e "\nSeeding Completed!"
