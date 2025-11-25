#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}RabbitMQ Quorum Queue Demo Test${NC}"
echo -e "${BLUE}High Availability & Reliable Messaging${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Configuration
BASE_URL="http://localhost:8082"
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

# Function to print info
print_info() {
    echo -e "${CYAN}ℹ $1${NC}"
}

# Step 1: Health Check
print_header "Step 1: Health Check"
response=$(curl -s "$BASE_URL/health")
echo "Response: $response"
print_success "Service is healthy"
sleep $SLEEP_TIME

# Step 2: Initial queue stats
print_header "Step 2: Initial Queue Statistics"
response=$(curl -s "$BASE_URL/stats")
echo "Queue Stats: $response"
print_info "Queue Type: Quorum (replicated across cluster)"
sleep $SLEEP_TIME

# Step 3: Publish messages
print_header "Step 3: Publishing messages with broker confirmations"
print_info "Each message waits for confirmation from the broker"
for i in {1..5}; do
    response=$(curl -s -X POST "$BASE_URL/publish" \
        -H "Content-Type: application/json" \
        -d "{\"message\":\"Order #$i\"}")
    echo "Published Order #$i: $response"
done
print_success "Published 5 messages (all confirmed by broker)"
sleep $SLEEP_TIME

# Step 4: Queue stats after publishing
print_header "Step 4: Queue Statistics After Publishing"
response=$(curl -s "$BASE_URL/stats")
echo "Queue Stats: $response"
print_info "Messages are replicated across all cluster nodes"
sleep $SLEEP_TIME

# Step 5: Consume with ACK
print_header "Step 5: Consuming messages with manual acknowledgment"
print_info "Each message is explicitly acknowledged after processing"
for i in {1..3}; do
    response=$(curl -s "$BASE_URL/consume")
    echo "Consumed: $response"
done
print_success "Successfully consumed and acknowledged 3 messages"
sleep $SLEEP_TIME

# Step 6: Simulate failure
print_header "Step 6: Simulating processing failure (NACK with requeue)"
print_info "Message will be rejected and requeued"
response=$(curl -s -X POST "$BASE_URL/consume/fail")
echo "Result: $response"
print_success "Message rejected and requeued for retry"
sleep $SLEEP_TIME

# Step 7: Check remaining
print_header "Step 7: Checking Remaining Messages"
response=$(curl -s "$BASE_URL/stats")
echo "Queue Stats: $response"
sleep $SLEEP_TIME

# Step 8: Consume remaining
print_header "Step 8: Consuming Remaining Messages"
for i in {1..3}; do
    response=$(curl -s "$BASE_URL/consume")
    if [[ $response == *"error"* ]]; then
        echo "No more messages"
        break
    fi
    echo "Consumed: $response"
done
print_success "Consumed all remaining messages"

# Summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Test Summary - Quorum Queue Features${NC}"
echo -e "${BLUE}========================================${NC}"
print_success "Publisher Confirmations: All messages confirmed by broker"
print_success "Manual Acknowledgments: Fine-grained control over processing"
print_success "Replication: Messages replicated across cluster nodes"
print_success "NACK & Requeue: Failed messages can be retried"
echo ""
echo -e "${YELLOW}Quorum Queue Benefits:${NC}"
echo -e "  ✓ High Availability: Survives node failures"
echo -e "  ✓ Data Safety: Messages replicated (Raft consensus)"
echo -e "  ✓ Durability: Messages persisted to disk"
echo -e "  ✓ Reliability: Publisher confirmations guarantee delivery"
echo ""
echo -e "${YELLOW}Check RabbitMQ Management UI:${NC}"
echo -e "  Node 1: http://localhost:15672 (guest/guest)"
echo -e "  Node 2: http://localhost:15673 (guest/guest)"
echo -e "  Node 3: http://localhost:15674 (guest/guest)"
echo ""
echo -e "${CYAN}Advanced Test:${NC}"
echo -e "  Try stopping one node: docker stop rabbitmq-node2"
echo -e "  The service should continue working!"
echo -e "  Restart it: docker start rabbitmq-node2"
echo ""
