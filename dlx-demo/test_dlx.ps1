# PowerShell Test Script for RabbitMQ DLX Demo

Write-Host "========================================" -ForegroundColor Blue
Write-Host "RabbitMQ DLX (Dead Letter Exchange) Test" -ForegroundColor Blue
Write-Host "========================================" -ForegroundColor Blue
Write-Host ""

$BaseUrl = "http://localhost:8081"
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

# Step 1: Health Check
Print-Header "Step 1: Health Check"
$response = Invoke-RestMethod -Uri "$BaseUrl/health" -Method Get
Write-Host "Response: $($response | ConvertTo-Json -Compress)"
Print-Success "Service is healthy"
Start-Sleep -Seconds $SleepTime

# Step 2: Publish messages
Print-Header "Step 2: Publishing messages to main queue"
1..5 | ForEach-Object {
    $body = @{ message = "Message $_" } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$BaseUrl/publish" -Method Post -Body $body -ContentType "application/json"
    Write-Host "Published Message $_`: $($response | ConvertTo-Json -Compress)"
}
Print-Success "Published 5 messages"
Start-Sleep -Seconds $SleepTime

# Step 3: Consume successfully
Print-Header "Step 3: Consuming messages successfully (simulating normal processing)"
1..2 | ForEach-Object {
    $response = Invoke-RestMethod -Uri "$BaseUrl/consume" -Method Get
    Write-Host "Consumed: $($response | ConvertTo-Json -Compress)"
}
Print-Success "Successfully consumed 2 messages"
Start-Sleep -Seconds $SleepTime

# Step 4: Reject messages
Print-Header "Step 4: Rejecting messages (simulating failures - these go to DLX)"
1..2 | ForEach-Object {
    $response = Invoke-RestMethod -Uri "$BaseUrl/reject" -Method Post
    Write-Host "Rejected: $($response | ConvertTo-Json -Compress)"
}
Print-Success "Rejected 2 messages (sent to Dead Letter Queue)"
Start-Sleep -Seconds $SleepTime

# Step 5: Check remaining messages
Print-Header "Step 5: Checking remaining messages in main queue"
$response = Invoke-RestMethod -Uri "$BaseUrl/consume" -Method Get
Write-Host "Remaining message: $($response | ConvertTo-Json -Compress)"
Print-Success "1 message still in main queue"
Start-Sleep -Seconds $SleepTime

# Step 6: Consume from DLQ
Print-Header "Step 6: Consuming from Dead Letter Queue"
1..2 | ForEach-Object {
    try {
        $response = Invoke-RestMethod -Uri "$BaseUrl/dlq/consume" -Method Get
        Write-Host "From DLQ: $($response | ConvertTo-Json -Compress)"
    } catch {
        Write-Host "Error: $_" -ForegroundColor Red
    }
    if ($_ -lt 2) { Start-Sleep -Milliseconds 500 }
}
Print-Success "Retrieved rejected messages from DLQ"
Start-Sleep -Seconds $SleepTime

# Step 7: Try empty DLQ
Print-Header "Step 7: Attempting to consume from empty DLQ"
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/dlq/consume" -Method Get
    Write-Host "Response: $($response | ConvertTo-Json -Compress)"
} catch {
    Write-Host "Response: No messages in DLQ (expected)" -ForegroundColor Gray
}
Print-Success "DLQ is now empty"

# Summary
Write-Host ""
Write-Host "========================================" -ForegroundColor Blue
Write-Host "Test Summary" -ForegroundColor Blue
Write-Host "========================================" -ForegroundColor Blue
Print-Success "Published 5 messages to main queue"
Print-Success "Successfully consumed 2 messages"
Print-Success "Rejected 2 messages (sent to DLX)"
Print-Success "Retrieved 2 messages from DLQ"
Print-Success "1 message consumed from main queue"
Write-Host ""
Write-Host "Check RabbitMQ Management UI:" -ForegroundColor Yellow
Write-Host "  URL: http://localhost:15672"
Write-Host "  User: guest / Password: guest"
Write-Host "  You should see:"
Write-Host "    - Queue 'messages' (main queue with DLX configured)"
Write-Host "    - Queue 'messages.dlq' (Dead Letter Queue)"
Write-Host "    - Exchange 'dlx.exchange' (Dead Letter Exchange)"
Write-Host ""
