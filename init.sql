-- Initialize database schema
-- This file will be executed when the PostgreSQL container starts for the first time

-- Create extensions if needed
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create tenant table for OTP-as-a-Service platform
CREATE TABLE IF NOT EXISTS tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    api_key VARCHAR(255) UNIQUE,
    webhook_url VARCHAR(500),
    sms_text_template TEXT DEFAULT 'Your OTP code is: {code}',
    max_otp_length INTEGER DEFAULT 6,
    otp_expiration_seconds INTEGER DEFAULT 300, -- 5 minutes in seconds
    max_retry_count INTEGER DEFAULT 3,
    daily_otp_limit INTEGER DEFAULT 1000,
    monthly_otp_limit INTEGER DEFAULT 30000,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    active BOOLEAN DEFAULT TRUE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_tenant_phone ON tenant(phone);
CREATE INDEX IF NOT EXISTS idx_tenant_email ON tenant(email);
CREATE INDEX IF NOT EXISTS idx_tenant_api_key ON tenant(api_key);
CREATE INDEX IF NOT EXISTS idx_tenant_active ON tenant(active);
CREATE INDEX IF NOT EXISTS idx_tenant_deleted_at ON tenant(deleted_at);


DO $$
BEGIN
    FOR i IN 1..1000 LOOP
        INSERT INTO tenant (
            name,
            phone,
            email,
            api_key,
            webhook_url,
            sms_text_template,
            max_otp_length,
            otp_expiration_seconds,
            max_retry_count,
            daily_otp_limit,
            monthly_otp_limit,
            created_at,
            updated_at,
            deleted_at,
            active
        ) VALUES (
            'Tenant ' || i,
            '+1000000' || LPAD(i::text, 4, '0'),
            'tenant' || i || '@example.com',
            md5(random()::text || clock_timestamp()::text || i::text),
            'https://webhook.example.com/tenant' || i,
            'Your OTP code is: {code}',
            6,
            300,
            3,
            1000,
            30000,
            NOW(),
            NOW(),
            NULL,
            TRUE
        );
    END LOOP;
END
$$;
