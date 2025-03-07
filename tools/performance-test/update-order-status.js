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
      maxVUs: 3000,
    },
  },
};

// Base URL for the API
const BASE_URL = "http://localhost:8080/postgres-crud/api/v1";

// Sample order statuses for updates
const statuses = [
  "PENDING",
  "PROCESSING",
  "PAID",
  "SHIPPED",
  "COMPLETED",
  "REFUNDED",
  "CANCELED",
];

export default function () {
  // Generate a random order ID between 1 and 100000
  const orderId = Math.floor(Math.random() * 100000) + 1;

  const randomStatus = statuses[Math.floor(Math.random() * statuses.length)];
  const updatePayload = JSON.stringify({
    status: randomStatus,
  });

  const updateRes = http.put(
    `${BASE_URL}/orders/${orderId}/status`,
    updatePayload,
    {
      headers: { "Content-Type": "application/json" },
      tags: { name: "updateStatus" },
    }
  );

  check(updateRes, {
    "Update order status is 200": (r) => r.status === 200,
  });
}

// k6 run update-order-status.js
