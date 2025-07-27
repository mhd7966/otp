# otp


## Load Testing with k6

This project supports load testing with [k6](https://k6.io/) for various OTP service models. You can benchmark the following endpoints:

1. **Hello World**  
2. **Build OTP Code** (`/OTP`)
3. **OTPWithRedis** (`/OTPWithRedis`)
4. **OTPWithRedisAndPostgres** (`/OTPWithRedisAndPostgres`)
5. **OTPWithRedisAndPostgresAndSMS** (`/OTPWithRedisAndPostgresAndSMS`)
6. **OTPWithRedisAndPostgresMemoized** (`/OTPWithRedisAndPostgresMemoized`)
7. **OTPWithRedisAndPostgresAndSMSMemoized** (`/OTPWithRedisAndPostgresAndSMSMemoized`)

Each function is exposed as an HTTP endpoint in your Go `main` package.

---

### Benchmark Endpoints and Metrics
| Function/Endpoint                        |   avg    |   min    |   med    |   max    |  p(90)   |  p(95)   | fail rate |
|------------------------------------------|----------|----------|----------|----------|----------|----------|-----------|
| `/hi` (Hello World)                      | 218.59ms |   300µs  | 160.56ms |   1.26s  | 540.41ms | 627.79ms |   0.00%   |
| `/OTP` (Build OTP Code)                  | 608.24ms |   298µs  | 362.69ms |   4.43s  |   1.47s  |   1.95s  |   0.00%   |
| `/OTPWithRedis`                          | 511.11ms |   332µs  | 467.80ms |   2.76s  |   1.03s  |   1.25s  |   0.00%   |
| `/OTPWithRedisAndPostgres`               |   6.50s  |   394µs  |   6.06ms |  36.79s  |  31.76s  |  32.72s  |   0.00%   |
| `/OTPWithRedisAndPostgresAndSMS`         |   7.94s  |   424µs  | 194.01ms |  54.66s  |  32.48s  |  33.29s  |   0.00%   |

---

### How to Run k6 Load Tests

1. **Install k6**  
   See [k6 installation guide](https://k6.io/docs/getting-started/installation/).

2. **Start your Go server**  
   Make sure your server is running locally (default: `http://localhost:8081`).

3. **Edit the k6 script**  
   Use or modify the provided `load-test.js` file.  
   Example for testing the "hello world" endpoint:
   ```js
   import http from 'k6/http';
   import { sleep } from 'k6';

   export let options = {
     stages: [
       { duration: '2s', target: 500 },
       { duration: '5s', target: 3000 },
       { duration: '5s', target: 6000 },
       { duration: '5s', target: 9000 },
       { duration: '10s', target: 9000 },
       { duration: '3s', target: 500 },
     ],
   };

   export default function () {
     http.get('http://localhost:8081/hi'); // Replace with your endpoint
     sleep(1);
   }
   ```
   To test other models, change the URL in `http.get()` to:
   - `/OTP`
   - `/OTPWithRedis`
   - `/OTPWithRedisAndPostgres`
   - `/OTPWithRedisAndPostgresAndSMS`
   - `/OTPWithRedisAndPostgresMemoized`
   - `/OTPWithRedisAndPostgresAndSMSMemoized`

4. **Run the test**
   ```bash
   k6 run load-test.js
   ```

### Metrics

k6 will output metrics such as:
- **avg**: Average response time
- **min**: Minimum response time
- **med**: Median response time
- **max**: Maximum response time
- **p(90)**: 90th percentile response time
- **p(95)**: 95th percentile response time
- **fail rate**: Error rate (non-2xx/3xx responses)

Example output:



