import http from "k6/http";
import { check, sleep } from "k6";
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";

// Configuration for the test
export const options = {
  scenarios: {
    api_load_test: {
      executor: "constant-arrival-rate",
      rate: 1000, //requests per second
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 100,
      maxVUs: 1000,
    },
  },
};

// Base URL for the API
const BASE_URL = "http://localhost:8080/postgres-crud/api/v1";

// Sample order data for creation
function createOrderData() {
  return {
    customerName: `Customer A`,
    status: "PENDING",
    items: [
      {
        productName: `Product A`,
        quantity: 5,
        price: 100.0,
      },
      {
        productName: `Product B`,
        quantity: 10,
        price: 200.0,
      },
    ],
  };
}

// Global variable to store created order IDs
let orderIds = [];

export default function () {
  // Create a new order
  const createPayload = JSON.stringify(createOrderData());
  const createRes = http.post(`${BASE_URL}/orders`, createPayload, {
    headers: { "Content-Type": "application/json" },
    tags: { name: "create" },
  });

  console.log("createRess", createRes);

  check(createRes, {
    "Create order status is 201": (r) => r.status === 201,
    "Create response has ID": (r) => JSON.parse(r.body).data.id !== undefined,
  });

  if (createRes.status === 201) {
    const orderId = JSON.parse(createRes.body).data.id;
    orderIds.push(orderId);
  }
}

// k6 run create-order.js
