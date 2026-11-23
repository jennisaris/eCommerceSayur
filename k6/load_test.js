import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Trend, Rate } from 'k6/metrics';

// Custom performance metrics
export const errorRate = new Rate('errors');
export const productListLatency = new Trend('product_list_duration');
export const checkoutFlowLatency = new Trend('checkout_flow_duration');

export const options = {
  stages: [
    { duration: '20s', target: 200 },   // Warm-up ramp to 200 VUs
    { duration: '40s', target: 1000 },  // Ramp-up to 1,000 VUs peak load
    { duration: '1m',  target: 1000 },  // Sustained peak test at 1,000 VUs
    { duration: '20s', target: 0 },     // Ramp-down recovery
  ],
  thresholds: {
    // 95% of requests must complete under 200ms (leveraging Redis cache & connection pooling)
    http_req_duration: ['p(95)<200'],
    // Error rate must remain under 1%
    http_req_failed: ['rate<0.01'],
    'product_list_duration': ['p(95)<100'], // Redis cached product query under 100ms
  },
};

const BASE_URL = __ENV.API_URL || 'http://localhost:8080/api/v1';
const CUSTOMER_EMAIL = __ENV.E2E_CUSTOMER_EMAIL;
const CUSTOMER_PASSWORD = __ENV.E2E_CUSTOMER_PASSWORD;

if (!CUSTOMER_EMAIL || !CUSTOMER_PASSWORD) {
  throw new Error('Set E2E_CUSTOMER_EMAIL and E2E_CUSTOMER_PASSWORD before running the load test.');
}

export default function () {
  let authenticated = false;
  let csrfToken = null;

  // 1. Authenticate user to obtain JWT
  group('01_Authentication', function () {
    const loginPayload = JSON.stringify({
      email: CUSTOMER_EMAIL,
      password: CUSTOMER_PASSWORD,
    });

    const csrfRes = http.get(`${BASE_URL}/auth/csrf`);
    csrfToken = JSON.parse(csrfRes.body).data.token;
    const loginRes = http.post(`${BASE_URL}/auth/login`, loginPayload, {
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': csrfToken,
      },
    });

    const isLoginOk = check(loginRes, {
      'login status 200': (r) => r.status === 200,
      'has authenticated user': (r) => JSON.parse(r.body).data && JSON.parse(r.body).data.user !== undefined,
    });
    errorRate.add(!isLoginOk);

    if (isLoginOk) {
      authenticated = true;
    }
  });

  sleep(0.1);

  // 2. High-throughput Catalog Browsing (Exercises Redis Cache-Aside)
  group('02_Catalog_Redis_Cache', function () {
    const start = new Date();
    const prodRes = http.get(`${BASE_URL}/products?page=1&limit=20&status=ACTIVE`);
    productListLatency.add(new Date() - start);

    const isCatalogOk = check(prodRes, {
      'get products status 200': (r) => r.status === 200,
      'products list returned': (r) => JSON.parse(r.body).data && JSON.parse(r.body).data.length > 0,
    });
    errorRate.add(!isCatalogOk);

    // Detail product by slug
    const detailRes = http.get(`${BASE_URL}/products/sony-wh-1000xm5-wireless-headphones`);
    check(detailRes, {
      'get product detail status 200': (r) => r.status === 200,
    });
  });

  sleep(0.2);

  // 3. Authenticated Cart & Checkout Estimation Flow
  if (authenticated) {
    group('03_Cart_and_Checkout_Flow', function () {
      const authHeaders = {
        'Content-Type': 'application/json',
        'X-CSRF-Token': csrfToken,
      };

      const start = new Date();

      // Add product #1 to cart
      const cartPayload = JSON.stringify({
        product_id: 1,
        quantity: 1,
      });

      const addCartRes = http.post(`${BASE_URL}/cart/items`, cartPayload, { headers: authHeaders });
      check(addCartRes, {
        'add to cart status 200': (r) => r.status === 200,
      });

      // Checkout Estimation (100% Delivery fee calculation)
      const estimatePayload = JSON.stringify({
        shipping_address: 'Jl. Sudirman Kav 52-53, Senayan, Jakarta Selatan',
      });

      const estRes = http.post(`${BASE_URL}/orders/estimate`, estimatePayload, { headers: authHeaders });
      const isEstOk = check(estRes, {
        'estimate checkout status 200': (r) => r.status === 200,
        'estimate total calculated': (r) => JSON.parse(r.body).data.total_amount > 0,
      });
      errorRate.add(!isEstOk);

      checkoutFlowLatency.add(new Date() - start);
    });
  }

  // Random realistic user think time between actions (100ms - 300ms)
  sleep(0.2);
}
