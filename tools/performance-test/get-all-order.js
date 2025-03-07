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

export default function () {
  const page = randomIntBetween(1, 1000);
  const perPage = randomIntBetween(100, 1000);
  const getAllRes = http.get(
    `${BASE_URL}/orders?page=${page}&per_page=${perPage}`,
    {
      tags: { name: "getAll" },
    }
  );

  check(getAllRes, {
    "Get all orders status is 200": (r) => r.status === 200,
    "Get all response has data": (r) => JSON.parse(r.body).data !== undefined,
  });
}

// k6 run get-all-order.js
