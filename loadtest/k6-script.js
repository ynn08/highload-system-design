import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
    stages: [
        { duration: '30s', target: 50 },  // Ramp up to 50 users
        { duration: '2m', target: 50 },   // Stay at 50 users for 2 mins (steady state)
        { duration: '30s', target: 100 }, // Spike to 100 users
        { duration: '1m', target: 100 },  // Hold spike
        { duration: '30s', target: 0 },   // Ramp down
    ],
    thresholds: {
        http_req_duration: ['p(99)<800'], // 99% of requests must complete below 800ms
        http_req_failed: ['rate<0.01'],   // Error rate should be less than 1%
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080/api/v1';
const RESTAURANT_ID = '00000000-0000-0000-0000-000000000001';
const MENU_ITEM_ID = '00000000-0000-0000-0000-000000000101';

export default function () {
    const customerId = uuidv4();

    // 1. Search Restaurants (Read-Heavy)
    let res = http.get(`${BASE_URL}/restaurants?lat=55.75&lon=37.62&radius=5000`);
    check(res, { 'search status was 200': (r) => r.status === 200 });

    sleep(1);

    // 2. Get Menu (Read-Heavy)
    res = http.get(`${BASE_URL}/restaurants/${RESTAURANT_ID}/menu`);
    check(res, { 'menu status was 200': (r) => r.status === 200 });

    sleep(1);

    // 3. Add to Cart (Write, Fast Redis)
    const cartPayload = JSON.stringify({
        menu_item_id: MENU_ITEM_ID,
        quantity: 2,
    });
    const headers = { 'Content-Type': 'application/json' };
    res = http.post(`${BASE_URL}/carts/${customerId}?restaurantId=${RESTAURANT_ID}`, cartPayload, { headers });
    check(res, { 'cart status was 200': (r) => r.status === 200 });

    sleep(1);

    // 4. Place Order (Write, Triggers Saga & Outbox)
    const orderPayload = JSON.stringify({
        customer_id: customerId,
        restaurant_id: RESTAURANT_ID,
        delivery_address: 'Load Test Street, 42',
        items: [{ menu_item_id: MENU_ITEM_ID, quantity: 2, price: 300 }]
    });
    res = http.post(`${BASE_URL}/orders`, orderPayload, { headers });
    
    check(res, { 
        'order status was 201': (r) => r.status === 201,
        'order created successfully': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body && body.order_id !== undefined;
            } catch (e) {
                return false;
            }
        }
    });

    sleep(1);
}
