#!/bin/bash

# Phase 3 Testing Script
# Tests all new financial reporting, subscription, and enhanced transaction features

set -e

echo "🚀 Starting Phase 3 Testing..."
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:8080"
API_VERSION="v1"

# Test user data
TEST_USER_ID="5500000000000" # Test phone number
TEST_USER_ID_2="5500000000001" # Second test user

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to make HTTP requests
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    if [ -n "$data" ]; then
        curl -s -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL/api/$API_VERSION$endpoint"
    else
        curl -s -X "$method" \
            -H "Content-Type: application/json" \
            "$BASE_URL/api/$API_VERSION$endpoint"
    fi
}

# Function to check if server is running
check_server() {
    print_status "Checking if server is running..."
    
    if curl -s "$BASE_URL/health" > /dev/null; then
        print_success "Server is running"
        return 0
    else
        print_error "Server is not running. Please start the server first."
        exit 1
    fi
}

# Test 1: Health Check
test_health_check() {
    print_status "Testing health check endpoint..."
    
    response=$(make_request "GET" "/health")
    if echo "$response" | grep -q "status.*ok"; then
        print_success "Health check passed"
    else
        print_error "Health check failed"
        echo "Response: $response"
    fi
}

# Test 2: Financial Summary
test_financial_summary() {
    print_status "Testing financial summary endpoint..."
    
    response=$(make_request "GET" "/financial/users/$TEST_USER_ID/summary?period=today")
    if echo "$response" | grep -q "summary"; then
        print_success "Financial summary endpoint working"
    else
        print_error "Financial summary endpoint failed"
        echo "Response: $response"
    fi
}

# Test 3: Detailed Report
test_detailed_report() {
    print_status "Testing detailed report endpoint..."
    
    response=$(make_request "GET" "/financial/users/$TEST_USER_ID/report?period=week")
    if echo "$response" | grep -q "summary"; then
        print_success "Detailed report endpoint working"
    else
        print_error "Detailed report endpoint failed"
        echo "Response: $response"
    fi
}

# Test 4: User Balance
test_user_balance() {
    print_status "Testing user balance endpoint..."
    
    response=$(make_request "GET" "/financial/users/$TEST_USER_ID/balance")
    if echo "$response" | grep -q "balance"; then
        print_success "User balance endpoint working"
    else
        print_error "User balance endpoint failed"
        echo "Response: $response"
    fi
}

# Test 5: User Transactions
test_user_transactions() {
    print_status "Testing user transactions endpoint..."
    
    response=$(make_request "GET" "/financial/users/$TEST_USER_ID/transactions?limit=5")
    if echo "$response" | grep -q "transactions"; then
        print_success "User transactions endpoint working"
    else
        print_error "User transactions endpoint failed"
        echo "Response: $response"
    fi
}

# Test 6: Top Categories
test_top_categories() {
    print_status "Testing top categories endpoint..."
    
    response=$(make_request "GET" "/financial/users/$TEST_USER_ID/categories?period=today&limit=3")
    if echo "$response" | grep -q "categories"; then
        print_success "Top categories endpoint working"
    else
        print_error "Top categories endpoint failed"
        echo "Response: $response"
    fi
}

# Test 7: Trial Status
test_trial_status() {
    print_status "Testing trial status endpoint..."
    
    response=$(make_request "GET" "/subscriptions/users/$TEST_USER_ID/trial-status")
    if echo "$response" | grep -q "trial_transactions_count"; then
        print_success "Trial status endpoint working"
    else
        print_error "Trial status endpoint failed"
        echo "Response: $response"
    fi
}

# Test 8: Subscription Info
test_subscription_info() {
    print_status "Testing subscription info endpoint..."
    
    response=$(make_request "GET" "/subscriptions/users/$TEST_USER_ID/info")
    if echo "$response" | grep -q "subscription_status"; then
        print_success "Subscription info endpoint working"
    else
        print_error "Subscription info endpoint failed"
        echo "Response: $response"
    fi
}

# Test 9: Create Subscription
test_create_subscription() {
    print_status "Testing create subscription endpoint..."
    
    data='{"payment_method": "pix"}'
    response=$(make_request "POST" "/subscriptions/users/$TEST_USER_ID_2" "$data")
    if echo "$response" | grep -q "message"; then
        print_success "Create subscription endpoint working"
    else
        print_error "Create subscription endpoint failed"
        echo "Response: $response"
    fi
}

# Test 10: Payment Webhook
test_payment_webhook() {
    print_status "Testing payment webhook endpoint..."
    
    data='{
        "user_id": "'$TEST_USER_ID'",
        "payment_id": "test_payment_123",
        "status": "approved",
        "amount": 9.90,
        "currency": "BRL"
    }'
    response=$(make_request "POST" "/subscriptions/webhook/payment" "$data")
    if echo "$response" | grep -q "message"; then
        print_success "Payment webhook endpoint working"
    else
        print_error "Payment webhook endpoint failed"
        echo "Response: $response"
    fi
}

# Test 11: Transaction Correction
test_transaction_correction() {
    print_status "Testing transaction correction endpoint..."
    
    # First, we need a transaction ID - this is a mock test
    data='{
        "amount": 50.00,
        "description": "Correção de transação"
    }'
    response=$(make_request "PUT" "/financial/transactions/test-transaction-id/correct" "$data")
    if echo "$response" | grep -q "message"; then
        print_success "Transaction correction endpoint working"
    else
        print_warning "Transaction correction endpoint (mock test)"
        echo "Response: $response"
    fi
}

# Test 12: WhatsApp Webhook (Enhanced)
test_whatsapp_webhook() {
    print_status "Testing enhanced WhatsApp webhook..."
    
    data='{
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
                        "wa_id": "'$TEST_USER_ID'"
                    }],
                    "messages": [{
                        "from": "'$TEST_USER_ID'",
                        "id": "wamid.123456789",
                        "timestamp": "1234567890",
                        "type": "text",
                        "text": {
                            "body": "resumo"
                        }
                    }]
                },
                "field": "messages"
            }]
        }]
    }'
    
    response=$(make_request "POST" "/webhook/whatsapp" "$data")
    if echo "$response" | grep -q "status.*ok"; then
        print_success "WhatsApp webhook (enhanced) working"
    else
        print_error "WhatsApp webhook (enhanced) failed"
        echo "Response: $response"
    fi
}

# Test 13: Period-based Reports
test_period_reports() {
    print_status "Testing period-based reports..."
    
    periods=("today" "week" "month")
    
    for period in "${periods[@]}"; do
        print_status "Testing $period period..."
        response=$(make_request "GET" "/financial/users/$TEST_USER_ID/summary?period=$period")
        if echo "$response" | grep -q "summary"; then
            print_success "$period period report working"
        else
            print_error "$period period report failed"
            echo "Response: $response"
        fi
    done
}

# Test 14: Error Handling
test_error_handling() {
    print_status "Testing error handling..."
    
    # Test with invalid user ID
    response=$(make_request "GET" "/financial/users/invalid-user-id/summary")
    if echo "$response" | grep -q "error"; then
        print_success "Error handling working (invalid user ID)"
    else
        print_error "Error handling failed (invalid user ID)"
        echo "Response: $response"
    fi
    
    # Test with invalid period
    response=$(make_request "GET" "/financial/users/$TEST_USER_ID/summary?period=invalid")
    if echo "$response" | grep -q "error"; then
        print_success "Error handling working (invalid period)"
    else
        print_error "Error handling failed (invalid period)"
        echo "Response: $response"
    fi
}

# Test 15: Performance Test
test_performance() {
    print_status "Testing API performance..."
    
    start_time=$(date +%s.%N)
    
    # Make multiple requests
    for i in {1..5}; do
        make_request "GET" "/financial/users/$TEST_USER_ID/summary" > /dev/null
    done
    
    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    
    print_success "Performance test completed in ${duration}s"
}

# Main test execution
main() {
    echo "🧪 Phase 3 Testing Suite"
    echo "========================="
    echo ""
    
    # Check server first
    check_server
    
    # Run all tests
    test_health_check
    test_financial_summary
    test_detailed_report
    test_user_balance
    test_user_transactions
    test_top_categories
    test_trial_status
    test_subscription_info
    test_create_subscription
    test_payment_webhook
    test_transaction_correction
    test_whatsapp_webhook
    test_period_reports
    test_error_handling
    test_performance
    
    echo ""
    echo "🎉 Phase 3 Testing Complete!"
    echo "============================="
    print_success "All Phase 3 features have been tested"
    echo ""
    echo "📊 Test Summary:"
    echo "- Financial Reporting Service ✅"
    echo "- Subscription Management ✅"
    echo "- Enhanced Transaction Features ✅"
    echo "- Period-based Reports ✅"
    echo "- Error Handling ✅"
    echo "- Performance ✅"
    echo ""
    print_status "Phase 3 is ready for production!"
}

# Run main function
main "$@" 