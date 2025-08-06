# Phase 1 Testing Guide
## Project Ara - WhatsApp Financial Bot

This guide covers all the testing approaches for Phase 1 of the Project Ara WhatsApp financial bot.

---

## 🧪 Testing Approaches

### 1. **Unit Tests** ✅
```bash
# Run all unit tests
go test ./... -v

# Run specific test
go test ./cmd/server -v

# Run with coverage
go test ./... -cover
```

### 2. **Integration Tests** ✅
```bash
# Test with database
docker-compose up -d postgres
go run cmd/server/main.go

# Test health endpoint
curl http://localhost:8080/health
```

### 3. **API Endpoint Tests** ✅
```bash
# Test all endpoints
curl -X GET http://localhost:8080/health
curl -X POST http://localhost:8080/api/v1/transactions
curl -X GET http://localhost:8080/api/v1/users/test/summary
curl -X POST http://localhost:8080/api/v1/subscriptions
```

### 4. **WhatsApp Webhook Tests** ✅
```bash
# Test WhatsApp webhook (simulated)
curl -X POST http://localhost:8080/api/v1/webhook/whatsapp \
  -H "Content-Type: application/json" \
  -d '{
    "object": "whatsapp_business_account",
    "entry": [{
      "id": "123456789",
      "changes": [{
        "value": {
          "messaging_product": "whatsapp",
          "metadata": {
            "display_phone_number": "5511999999999",
            "phone_number_id": "123456789"
          },
          "contacts": [{
            "profile": {
              "name": "Test User"
            },
            "wa_id": "5511999999999"
          }],
          "messages": [{
            "from": "5511999999999",
            "id": "wamid.123456789",
            "timestamp": "1234567890",
            "type": "text",
            "text": {
              "body": "vendi 3 cachorros quentes por R$ 45"
            }
          }]
        },
        "field": "messages"
      }]
    }]
  }'
```

---

## 🐳 Docker Testing

### 1. **Build Test**
```bash
# Test Docker build
docker build -t project-ara .

# Test Docker Compose
docker-compose up --build
```

### 2. **Container Health Test**
```bash
# Test container startup
docker run --rm -p 8080:8080 project-ara

# Test with environment variables
docker run --rm -p 8080:8080 \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=project_ara \
  project-ara
```

---

## 🗄️ Database Testing

### 1. **Database Connection Test**
```bash
# Start PostgreSQL
docker-compose up -d postgres

# Test connection
docker exec project-ara-postgres-1 psql -U postgres -d project_ara -c "SELECT version();"
```

### 2. **Schema Migration Test**
```bash
# Run application (will auto-migrate)
go run cmd/server/main.go

# Check tables created
docker exec project-ara-postgres-1 psql -U postgres -d project_ara -c "\dt"
```

---

## 🔧 Manual Testing Scenarios

### 1. **Health Check Test**
```bash
# Expected response:
{
  "status": "healthy",
  "service": "project-ara",
  "version": "1.0.0"
}
```

### 2. **WhatsApp Message Processing Test**
```bash
# Test text message processing
curl -X POST http://localhost:8080/api/v1/webhook/whatsapp \
  -H "Content-Type: application/json" \
  -d '{
    "object": "whatsapp_business_account",
    "entry": [{
      "id": "123456789",
      "changes": [{
        "value": {
          "messaging_product": "whatsapp",
          "metadata": {
            "display_phone_number": "5511999999999",
            "phone_number_id": "123456789"
          },
          "contacts": [{
            "profile": {
              "name": "Test User"
            },
            "wa_id": "5511999999999"
          }],
          "messages": [{
            "from": "5511999999999",
            "id": "wamid.123456789",
            "timestamp": "1234567890",
            "type": "text",
            "text": {
              "body": "vendi 3 cachorros quentes por R$ 45"
            }
          }]
        },
        "field": "messages"
      }]
    }]
  }'
```

### 3. **Audio Message Test**
```bash
# Test audio message processing
curl -X POST http://localhost:8080/api/v1/webhook/whatsapp \
  -H "Content-Type: application/json" \
  -d '{
    "object": "whatsapp_business_account",
    "entry": [{
      "id": "123456789",
      "changes": [{
        "value": {
          "messaging_product": "whatsapp",
          "metadata": {
            "display_phone_number": "5511999999999",
            "phone_number_id": "123456789"
          },
          "contacts": [{
            "profile": {
              "name": "Test User"
            },
            "wa_id": "5511999999999"
          }],
          "messages": [{
            "from": "5511999999999",
            "id": "wamid.123456789",
            "timestamp": "1234567890",
            "type": "audio",
            "audio": {
              "id": "audio_id_123",
              "mime_type": "audio/ogg; codecs=opus"
            }
          }]
        },
        "field": "messages"
      }]
    }]
  }'
```

### 4. **Image Message Test**
```bash
# Test image message processing
curl -X POST http://localhost:8080/api/v1/webhook/whatsapp \
  -H "Content-Type: application/json" \
  -d '{
    "object": "whatsapp_business_account",
    "entry": [{
      "id": "123456789",
      "changes": [{
        "value": {
          "messaging_product": "whatsapp",
          "metadata": {
            "display_phone_number": "5511999999999",
            "phone_number_id": "123456789"
          },
          "contacts": [{
            "profile": {
              "name": "Test User"
            },
            "wa_id": "5511999999999"
          }],
          "messages": [{
            "from": "5511999999999",
            "id": "wamid.123456789",
            "timestamp": "1234567890",
            "type": "image",
            "image": {
              "id": "image_id_123",
              "mime_type": "image/jpeg",
              "sha256": "image_sha256_hash"
            }
          }]
        },
        "field": "messages"
      }]
    }]
  }'
```

---

## 🚀 Performance Testing

### 1. **Load Test**
```bash
# Install hey (load testing tool)
go install github.com/rakyll/hey@latest

# Test health endpoint
hey -n 1000 -c 10 http://localhost:8080/health

# Test webhook endpoint
hey -n 100 -c 5 -m POST -H "Content-Type: application/json" \
  -d '{"object":"whatsapp_business_account","entry":[]}' \
  http://localhost:8080/api/v1/webhook/whatsapp
```

### 2. **Memory Usage Test**
```bash
# Monitor memory usage
docker stats project-ara-app-1

# Or with local process
ps aux | grep "go run"
```

---

## 🔍 Debugging Tests

### 1. **Log Level Testing**
```bash
# Set debug mode
export GIN_MODE=debug
go run cmd/server/main.go

# Check logs
docker logs project-ara-app-1
```

### 2. **Database Debugging**
```bash
# Connect to database
docker exec -it project-ara-postgres-1 psql -U postgres -d project_ara

# Check tables
\dt

# Check data
SELECT * FROM users;
SELECT * FROM transactions;
```

---

## ✅ Test Checklist

### Phase 1 Core Features
- [x] **Go application builds successfully**
- [x] **Docker container builds and runs**
- [x] **Health endpoint responds correctly**
- [x] **Database connection works**
- [x] **WhatsApp webhook endpoint accepts requests**
- [x] **Message parsing works for text messages**
- [x] **Message parsing works for audio messages**
- [x] **Message parsing works for image messages**
- [x] **User creation and retrieval works**
- [x] **Transaction creation works**
- [x] **Financial summary calculation works**

### Infrastructure Tests
- [x] **PostgreSQL database starts correctly**
- [x] **Docker Compose setup works**
- [x] **Environment variables are loaded**
- [x] **API endpoints are accessible**
- [x] **Error handling works correctly**
- [x] **Logging is functional**

---

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check if PostgreSQL is running
   docker ps | grep postgres
   
   # Check database logs
   docker logs project-ara-postgres-1
   ```

2. **Application Won't Start**
   ```bash
   # Check environment variables
   cat .env
   
   # Check application logs
   docker logs project-ara-app-1
   ```

3. **WhatsApp Webhook Not Working**
   ```bash
   # Test webhook endpoint directly
   curl -X POST http://localhost:8080/api/v1/webhook/whatsapp \
     -H "Content-Type: application/json" \
     -d '{"test": "data"}'
   ```

---

## 📊 Test Results Summary

### Success Criteria
- ✅ Application starts without errors
- ✅ Health endpoint returns 200 OK
- ✅ Database migrations run successfully
- ✅ WhatsApp webhook accepts valid payloads
- ✅ All API endpoints are accessible
- ✅ Docker containerization works
- ✅ Unit tests pass

### Performance Benchmarks
- **Startup time**: < 5 seconds
- **Health endpoint response**: < 100ms
- **Memory usage**: < 100MB
- **Container size**: < 50MB

---

**Phase 1 Testing Status: ✅ COMPLETED**

All core functionality is working correctly and ready for Phase 2 development! 