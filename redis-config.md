# Redis Configuration for OTP Storage

## OTP Code Storage Strategy

### Key Naming Convention:
```
otp:{tenant_id}:{user_phone}
```

### Example Redis Commands:

```bash
# Store OTP code with TTL (expires automatically)
SET otp:1:+1234567890 "123456" EX 300

# Get OTP code
GET otp:1:+1234567890

# Check if OTP exists
EXISTS otp:1:+1234567890

# Delete OTP code manually (if needed)
DEL otp:1:+1234567890

# Set TTL for existing key
EXPIRE otp:1:+1234567890 300

# Get remaining TTL
TTL otp:1:+1234567890
```

### Retry Count Tracking:
```bash
# Store retry count with longer TTL
SET retry:1:+1234567890 "2" EX 3600

# Increment retry count
INCR retry:1:+1234567890

# Get retry count
GET retry:1:+1234567890
```

### Ban Status:
```bash
# Ban user for X minutes
SET ban:1:+1234567890 "1" EX 3600

# Check if banned
EXISTS ban:1:+1234567890
```

## Go Code Example:

```go
// Store OTP in Redis
func StoreOTP(tenantID int, userPhone, otpCode string, expirationSeconds int) error {
    key := fmt.Sprintf("otp:%d:%s", tenantID, userPhone)
    return redisClient.Set(ctx, key, otpCode, time.Duration(expirationSeconds)*time.Second).Err()
}

// Get OTP from Redis
func GetOTP(tenantID int, userPhone string) (string, error) {
    key := fmt.Sprintf("otp:%d:%s", tenantID, userPhone)
    return redisClient.Get(ctx, key).Result()
}

// Check if OTP exists
func OTPExists(tenantID int, userPhone string) (bool, error) {
    key := fmt.Sprintf("otp:%d:%s", tenantID, userPhone)
    result, err := redisClient.Exists(ctx, key).Result()
    return result > 0, err
}
```

## Benefits of Redis Storage:

1. **Automatic Expiration**: OTP codes are automatically deleted after TTL
2. **High Performance**: Fast read/write operations
3. **Memory Efficient**: No need to store expired codes
4. **Atomic Operations**: Built-in support for counters and flags
5. **Scalability**: Can handle high concurrent OTP requests 