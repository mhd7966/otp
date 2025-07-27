import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2s', target: 500 },   // Ramp-up to 500 VUs
    { duration: '5s', target: 2000 },  // Ramp-up to 2000 VUs
    { duration: '5s', target: 5000 },  // Ramp-up to 5000 VUs
    { duration: '5s', target: 7000 },  // Ramp-up to 7000 VUs
    { duration: '10s', target: 7000 }, // Hold at 7000 VUs
    { duration: '3s', target: 500 },   // Ramp-down to 500 VUs
  ],
};

export default function () {
  http.get('http://localhost:8081/hi'); // Replace with your actual URL
  sleep(1); // Wait between iterations
}
