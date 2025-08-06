#!/bin/bash

# Quick Test Script for Project Ara
# Run individual tests quickly

echo "🚀 Quick Test Script for Project Ara"
echo "===================================="

case "$1" in
    "build")
        echo "🔨 Testing Go build..."
        go build ./cmd/server
        echo "✅ Build successful"
        ;;
    "test")
        echo "🧪 Running unit tests..."
        go test ./cmd/server -v
        ;;
    "docker")
        echo "🐳 Testing Docker build..."
        docker build -t project-ara .
        echo "✅ Docker build successful"
        ;;
    "health")
        echo "🏥 Testing health endpoint..."
        curl -s http://localhost:8080/health | jq .
        ;;
    "webhook")
        echo "📱 Testing WhatsApp webhook..."
        curl -X POST http://localhost:8080/api/v1/webhook/whatsapp \
          -H "Content-Type: application/json" \
          -d '{"object":"whatsapp_business_account","entry":[]}' \
          -w "\nStatus: %{http_code}\n"
        ;;
    "db")
        echo "🗄️ Testing database..."
        docker-compose up -d postgres
        sleep 3
        docker exec project-ara-postgres-1 psql -U postgres -d project_ara -c "SELECT version();"
        ;;
    "start")
        echo "🚀 Starting application..."
        docker-compose up -d postgres
        sleep 3
        go run cmd/server/main.go &
        echo "✅ Application started on http://localhost:8080"
        ;;
    "stop")
        echo "🛑 Stopping application..."
        pkill -f "go run cmd/server/main.go" || true
        docker-compose down
        echo "✅ Application stopped"
        ;;
    "all")
        echo "🎯 Running all quick tests..."
        ./scripts/quick_test.sh build
        ./scripts/quick_test.sh test
        ./scripts/quick_test.sh docker
        ./scripts/quick_test.sh db
        ./scripts/quick_test.sh start
        sleep 3
        ./scripts/quick_test.sh health
        ./scripts/quick_test.sh webhook
        ./scripts/quick_test.sh stop
        echo "✅ All tests completed"
        ;;
    *)
        echo "Usage: $0 {build|test|docker|health|webhook|db|start|stop|all}"
        echo ""
        echo "Commands:"
        echo "  build   - Test Go build"
        echo "  test    - Run unit tests"
        echo "  docker  - Test Docker build"
        echo "  health  - Test health endpoint"
        echo "  webhook - Test WhatsApp webhook"
        echo "  db      - Test database connection"
        echo "  start   - Start application"
        echo "  stop    - Stop application"
        echo "  all     - Run all tests"
        exit 1
        ;;
esac 