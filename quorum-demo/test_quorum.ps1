# PowerShell Test Script for RabbitMQ Quorum Queue Demo

Write-Host "========================================" -ForegroundColor Blue
Write-Host "RabbitMQ Quorum Queue Demo Test" -ForegroundColor Blue
Write-Host "High Availability & Reliable Messaging" -ForegroundColor Blue
Write-Host "========================================" -ForegroundColor Blue
Write-Host ""

$BaseUrl = "http://localhost:8082"
$SleepTime = 1

function Print-Header {
    param($Message)
    Write-Host ""
    Write-Host ">>> $Message" -ForegroundColor Yellow
    Write-Host ""
}

function Print-Success {
    param($Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Print-Info {
    param($Message)
    Write-Host "ℹ $Message" -ForegroundColor Cyan
}

# Step 1: Health Check
Print-Header "Step 1: Health Check"
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/health" -Method Get
    Write-Host "Response: $($response | ConvertTo-Json -Compress)"
    Print-Success "Service is healthy"
}
catch {
    Write-Host "Error: Service not running. Please start with 'go run main.go'" -ForegroundColor Red
    exit 1
}
Start-Sleep -Seconds $SleepTime

# Step 2: Check Queue Statistics (initial state)
Print-Header "Step 2: Initial Queue Statistics"
$response = Invoke-RestMethod -Uri "$BaseUrl/stats" -Method Get
Write-Host "Queue Stats: $($response.data | ConvertTo-Json)"
Print-Info "Queue Type: Quorum (replicated across cluster)"
Start-Sleep -Seconds $SleepTime

# Step 3: Publish messages with confirmations
Print-Header "Step 3: Publishing messages with broker confirmations"
Print-Info "Each message waits for confirmation from the broker"
1..5 | ForEach-Object {
    $body = @{ message = "Order #$_" } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$BaseUrl/publish" -Method Post -Body $body -ContentType "application/json"
    Write-Host "Published Order #$_`: $($response.message)"
}
Print-Success "Published 5 messages (all confirmed by broker)"
Start-Sleep -Seconds $SleepTime

# Step 4: Check queue stats after publishing
Print-Header "Step 4: Queue Statistics After Publishing"
$response = Invoke-RestMethod -Uri "$BaseUrl/stats" -Method Get
Write-Host "Queue Stats: $($response.data | ConvertTo-Json)"
Print-Info "Messages are replicated across all cluster nodes"
Start-Sleep -Seconds $SleepTime

# Step 5: Consume messages with manual ACK
Print-Header "Step 5: Consuming messages with manual acknowledgment"
Print-Info "Each message is explicitly acknowledged after processing"
1..3 | ForEach-Object {
    $response = Invoke-RestMethod -Uri "$BaseUrl/consume" -Method Get
    Write-Host "Consumed: $($response.message)"
}
Print-Success "Successfully consumed and acknowledged 3 messages"
Start-Sleep -Seconds $SleepTime

# Step 6: Simulate processing failure (NACK with requeue)
Print-Header "Step 6: Simulating processing failure (NACK with requeue)"
Print-Info "Message will be rejected and requeued"
$response = Invoke-RestMethod -Uri "$BaseUrl/consume/fail" -Method Post
Write-Host "Result: $($response.message)"
Print-Success "Message rejected and requeued for retry"
Start-Sleep -Seconds $SleepTime

# Step 7: Check remaining messages
Print-Header "Step 7: Checking Remaining Messages"
$response = Invoke-RestMethod -Uri "$BaseUrl/stats" -Method Get
Write-Host "Queue Stats: $($response.data | ConvertTo-Json)"
Print-Info "Remaining messages: $($response.data.messages)"
Start-Sleep -Seconds $SleepTime

# Step 8: Consume remaining messages
Print-Header "Step 8: Consuming Remaining Messages"
$remainingCount = $response.data.messages
if ($remainingCount -gt 0) {
    1..$remainingCount | ForEach-Object {
        try {
            $response = Invoke-RestMethod -Uri "$BaseUrl/consume" -Method Get
            Write-Host "Consumed: $($response.message)"
        }
        catch {
            Write-Host "No more messages" -ForegroundColor Gray
        }
    }
    Print-Success "Consumed all remaining messages"
}
else {
    Print-Info "No messages remaining in queue"
}

# Summary
Write-Host ""
Write-Host "========================================" -ForegroundColor Blue
Write-Host "Test Summary - Quorum Queue Features" -ForegroundColor Blue
Write-Host "========================================" -ForegroundColor Blue
Print-Success "Publisher Confirmations: All messages confirmed by broker"
Print-Success "Manual Acknowledgments: Fine-grained control over processing"
Print-Success "Replication: Messages replicated across cluster nodes"
Print-Success "NACK & Requeue: Failed messages can be retried"
Write-Host ""
Write-Host "Quorum Queue Benefits:" -ForegroundColor Yellow
Write-Host "  ✓ High Availability: Survives node failures"
Write-Host "  ✓ Data Safety: Messages replicated (Raft consensus)"
Write-Host "  ✓ Durability: Messages persisted to disk"
Write-Host "  ✓ Reliability: Publisher confirmations guarantee delivery"
Write-Host ""
Write-Host "Check RabbitMQ Management UI:" -ForegroundColor Yellow
Write-Host "  Node 1: http://localhost:15672 (guest/guest)"
Write-Host "  Node 2: http://localhost:15673 (guest/guest)"
Write-Host "  Node 3: http://localhost:15674 (guest/guest)"
Write-Host ""
Write-Host "Advanced Test:" -ForegroundColor Cyan
Write-Host "  Try stopping one node: docker stop rabbitmq-node2"
Write-Host "  The service should continue working!"
Write-Host "  Restart it: docker start rabbitmq-node2"
Write-Host ""
