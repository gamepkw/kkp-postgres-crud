import http from "k6/http";
import { check, sleep } from "k6";
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";

// Configuration for the test
export const options = {
  scenarios: {
    api_load_test: {
      executor: "constant-arrival-rate",
      rate: 10000,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 1000,
      maxVUs: 1000,
    },
  },
};

// Base URL for the API
const BASE_URL = "http://localhost:8080/postgres-crud/api/v1";

// Global variable to store created order IDs

export default function () {
  // Generate a random order ID between 1 and 100000
  const orderId = Math.floor(Math.random() * 100000) + 1;

  const getByIdRes = http.get(`${BASE_URL}/orders/${orderId}`, {
    tags: { name: "getById" },
  });

  check(getByIdRes, {
    "Get order by ID status is 200": (r) => r.status === 200,
    "Get by ID response has correct ID": (r) => {
      const data = JSON.parse(r.body).data;
      return data && data.id === orderId;
    },
  });
}

// k6 run get-by-id.js
