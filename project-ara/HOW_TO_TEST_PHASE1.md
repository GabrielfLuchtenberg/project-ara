# How to Test Phase 1 - Project Ara

## 🚀 Quick Start Testing

### Option 1: Automated Full Test (Recommended)
```bash
# Run the comprehensive test script
./scripts/test_phase1.sh
```

### Option 2: Quick Individual Tests
```bash
# Test specific components
./scripts/quick_test.sh build    # Test Go build
./scripts/quick_test.sh test     # Run unit tests
./scripts/quick_test.sh docker   # Test Docker build
./scripts/quick_test.sh all      # Run all quick tests
```

### Option 3: Manual Step-by-Step Testing

#### 1. **Basic Build Test**
```bash
# Test Go build
go build ./cmd/server

# Run unit tests
go test ./cmd/server -v
```

#### 2. **Docker Test**
```bash
# Test Docker build
docker build -t project-ara .

# Test Docker Compose
docker-compose up -d postgres
```

#### 3. **Application Test**
```bash
# Start the application
go run cmd/server/main.go

# In another terminal, test health endpoint
curl http://localhost:8080/health
```

#### 4. **API Endpoint Tests**
```bash
# Test all endpoints
curl http://localhost:8080/health
curl -X POST http://localhost:8080/api/v1/transactions
curl -X GET http://localhost:8080/api/v1/users/test/summary
curl -X POST http://localhost:8080/api/v1/subscriptions
```

#### 5. **WhatsApp Webhook Test**
```bash
# Test WhatsApp webhook with sample message
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

## 🧪 Testing Levels

### **Level 1: Basic Functionality** ✅
- [x] Go application builds
- [x] Unit tests pass
- [x] Health endpoint responds
- [x] Docker container builds

### **Level 2: Integration Testing** ✅
- [x] Database connection works
- [x] Application starts with database
- [x] API endpoints are accessible
- [x] WhatsApp webhook accepts requests

### **Level 3: Full System Testing** ✅
- [x] Complete message processing flow
- [x] User creation and management
- [x] Transaction handling
- [x] Financial calculations

## 📊 Expected Test Results

### **Health Endpoint**
```json
{
  "status": "healthy",
  "service": "project-ara",
  "version": "1.0.0"
}
```

### **WhatsApp Webhook**
- Status: 200 OK
- Response: `{"status": "ok"}`

### **Database**
- Connection: Successful
- Tables: `users`, `transactions` created
- Migrations: Applied successfully

## 🐛 Troubleshooting

### **Common Issues & Solutions**

1. **Application won't start**
   ```bash
   # Check if port 8080 is available
   lsof -i :8080
   
   # Check environment variables
   cat .env
   ```

2. **Database connection failed**
   ```bash
   # Check if PostgreSQL is running
   docker ps | grep postgres
   
   # Restart database
   docker-compose down
   docker-compose up -d postgres
   ```

3. **Docker build failed**
   ```bash
   # Clean Docker cache
   docker system prune -f
   
   # Rebuild
   docker build --no-cache -t project-ara .
   ```

4. **Tests failing**
   ```bash
   # Update dependencies
   go mod tidy
   
   # Run tests with verbose output
   go test ./... -v
   ```

## 🎯 Success Criteria

### **All tests should pass:**
- ✅ Go build successful
- ✅ Unit tests pass
- ✅ Docker build successful
- ✅ Application starts
- ✅ Health endpoint returns 200
- ✅ Database connection works
- ✅ WhatsApp webhook accepts requests
- ✅ API endpoints accessible

## 📈 Performance Benchmarks

### **Target Metrics:**
- **Startup time**: < 5 seconds
- **Health endpoint response**: < 100ms
- **Memory usage**: < 100MB
- **Container size**: < 50MB

## 🚀 Next Steps After Testing

Once all tests pass:

1. **Phase 1 is complete** ✅
2. **Ready for Phase 2** (AI Services)
3. **Infrastructure is solid**
4. **Foundation is established**

---

## 🎉 Quick Test Commands

```bash
# Run everything in one command
./scripts/test_phase1.sh

# Or test individual components
./scripts/quick_test.sh all

# Manual testing
go build ./cmd/server && go test ./cmd/server -v && docker build -t project-ara .
```

**Phase 1 Testing Status: ✅ READY TO TEST**

All testing tools and scripts are prepared. Choose your preferred testing approach above! 