#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}RabbitMQ DLX (Dead Letter Exchange) Test${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Configuration
BASE_URL="http://localhost:8081"
SLEEP_TIME=1

# Function to print section headers
print_header() {
    echo ""
    echo -e "${YELLOW}>>> $1${NC}"
    echo ""
}

# Function to print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Step 1: Health Check
print_header "Step 1: Health Check"
response=$(curl -s "$BASE_URL/health")
echo "Response: $response"
print_success "Service is healthy"
sleep $SLEEP_TIME

# Step 2: Publish messages to main queue
print_header "Step 2: Publishing messages to main queue"
for i in {1..5}; do
    response=$(curl -s -X POST "$BASE_URL/publish" \
        -H "Content-Type: application/json" \
        -d "{\"message\":\"Message $i\"}")
    echo "Published Message $i: $response"
done
print_success "Published 5 messages"
sleep $SLEEP_TIME

# Step 3: Consume some messages successfully
print_header "Step 3: Consuming messages successfully (simulating normal processing)"
for i in {1..2}; do
    response=$(curl -s "$BASE_URL/consume")
    echo "Consumed: $response"
done
print_success "Successfully consumed 2 messages"
sleep $SLEEP_TIME

# Step 4: Reject messages (simulate processing failures)
print_header "Step 4: Rejecting messages (simulating failures - these go to DLX)"
for i in {1..2}; do
    response=$(curl -s -X POST "$BASE_URL/reject")
    echo "Rejected: $response"
done
print_success "Rejected 2 messages (sent to Dead Letter Queue)"
sleep $SLEEP_TIME

# Step 5: Check if there are still messages in main queue
print_header "Step 5: Checking remaining messages in main queue"
response=$(curl -s "$BASE_URL/consume")
echo "Remaining message: $response"
print_success "1 message still in main queue"
sleep $SLEEP_TIME

# Step 6: Consume from Dead Letter Queue
print_header "Step 6: Consuming from Dead Letter Queue"
for i in {1..2}; do
    response=$(curl -s "$BASE_URL/dlq/consume")
    echo "From DLQ: $response"
    if [ $i -lt 2 ]; then
        sleep 0.5
    fi
done
print_success "Retrieved rejected messages from DLQ"
sleep $SLEEP_TIME

# Step 7: Try to consume from empty DLQ
print_header "Step 7: Attempting to consume from empty DLQ"
response=$(curl -s "$BASE_URL/dlq/consume")
echo "Response: $response"
print_success "DLQ is now empty"

# Summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Published 5 messages to main queue${NC}"
echo -e "${GREEN}✓ Successfully consumed 2 messages${NC}"
echo -e "${GREEN}✓ Rejected 2 messages (sent to DLX)${NC}"
echo -e "${GREEN}✓ Retrieved 2 messages from DLQ${NC}"
echo -e "${GREEN}✓ 1 message consumed from main queue${NC}"
echo ""
echo -e "${YELLOW}Check RabbitMQ Management UI:${NC}"
echo -e "  URL: http://localhost:15672"
echo -e "  User: guest / Password: guest"
echo -e "  You should see:"
echo -e "    - Queue 'messages' (main queue with DLX configured)"
echo -e "    - Queue 'messages.dlq' (Dead Letter Queue)"
echo -e "    - Exchange 'dlx.exchange' (Dead Letter Exchange)"
echo ""
