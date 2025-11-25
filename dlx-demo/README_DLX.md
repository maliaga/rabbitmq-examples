# RabbitMQ Dead Letter Exchange (DLX) Demo

Demostración de mensajería confiable usando **Dead Letter Exchange (DLX)** en RabbitMQ. Este proyecto muestra cómo manejar mensajes que fallan en su procesamiento mediante el uso de colas de mensajes muertos (Dead Letter Queue - DLQ).

## ¿Qué es DLX?

**Dead Letter Exchange (DLX)** es una característica de RabbitMQ que permite redirigir mensajes que no pueden ser procesados exitosamente a una cola especial llamada **Dead Letter Queue (DLQ)**. Esto es útil para:

- 🔄 **Reintentos**: Procesar mensajes fallidos más tarde
- 🐛 **Debugging**: Analizar mensajes problemáticos
- 📊 **Monitoreo**: Detectar patrones de fallos
- 🛡️ **Resiliencia**: Evitar pérdida de mensajes importantes

## Arquitectura

```mermaid
graph LR
    A[Producer] -->|Publish| B[Main Queue<br/>messages]
    B -->|Consume OK| C[Consumer]
    B -->|Reject/Nack| D[DLX Exchange<br/>dlx.exchange]
    D -->|Route| E[Dead Letter Queue<br/>messages.dlq]
    E -->|Analyze/Retry| F[DLQ Consumer]
    
    style B fill:#90EE90
    style D fill:#FFB6C1
    style E fill:#FFB6C1
```

### Componentes

1. **Main Queue** (`messages`): Cola principal con DLX configurado
2. **DLX Exchange** (`dlx.exchange`): Exchange que recibe mensajes rechazados
3. **Dead Letter Queue** (`messages.dlq`): Cola donde se almacenan mensajes fallidos

## Requisitos

- Go 1.21 o superior
- Docker y Docker Compose (para RabbitMQ)
- RabbitMQ ejecutándose en `localhost:5672`

## Instalación

1. **Asegúrate de que RabbitMQ esté ejecutándose:**
   ```bash
   cd c:\wk\rabbitmq
   docker-compose up -d
   ```

2. **Navega al directorio del demo:**
   ```bash
   cd dlx-demo
   ```

3. **Descarga las dependencias:**
   ```bash
   go mod download
   ```

## Uso

### 1. Iniciar el servicio DLX

```bash
go run main.go
```

El servicio estará disponible en `http://localhost:8081`

**Salida esperada:**
```
========================================
RabbitMQ DLX Demo Service Started
========================================
HTTP Server: http://localhost:8081

Endpoints:
  POST http://localhost:8081/publish      - Publish a message
  GET  http://localhost:8081/consume      - Consume a message
  POST http://localhost:8081/reject       - Reject a message (simulate failure)
  GET  http://localhost:8081/dlq/consume  - Consume from Dead Letter Queue
  GET  http://localhost:8081/health       - Health check

RabbitMQ Management UI: http://localhost:15672 (guest/guest)
========================================
```

### 2. Ejecutar el test automático

**PowerShell (Windows):**
```powershell
.\test_dlx.ps1
```

**Bash (Linux/Mac/Git Bash):**
```bash
chmod +x test_dlx.sh
./test_dlx.sh
```

El script ejecutará automáticamente:
1. ✅ Publicar 5 mensajes
2. ✅ Consumir 2 mensajes exitosamente
3. ✅ Rechazar 2 mensajes (enviar a DLX)
4. ✅ Consumir 1 mensaje restante
5. ✅ Recuperar mensajes de la DLQ

## API Endpoints

### POST /publish
Publica un mensaje en la cola principal.

**Request:**
```bash
curl -X POST http://localhost:8081/publish \
  -H "Content-Type: application/json" \
  -d '{"message":"Test message"}'
```

**Response:**
```json
{
  "status": "success",
  "message": "Message published successfully"
}
```

---

### GET /consume
Consume un mensaje de la cola principal (procesamiento exitoso).

**Request:**
```bash
curl http://localhost:8081/consume
```

**Response:**
```json
{
  "status": "success",
  "message": "Test message"
}
```

---

### POST /reject
Consume y **rechaza** un mensaje, enviándolo al DLX (simula un fallo de procesamiento).

**Request:**
```bash
curl -X POST http://localhost:8081/reject
```

**Response:**
```json
{
  "status": "success",
  "message": "Message rejected and sent to DLX: Test message"
}
```

---

### GET /dlq/consume
Consume un mensaje de la Dead Letter Queue.

**Request:**
```bash
curl http://localhost:8081/dlq/consume
```

**Response:**
```json
{
  "status": "success",
  "message": "Message from DLQ: Test message"
}
```

---

### GET /health
Verifica el estado del servicio.

**Request:**
```bash
curl http://localhost:8081/health
```

**Response:**
```json
{
  "status": "healthy"
}
```

## Flujo de Trabajo Completo

### Escenario: Procesamiento con Fallos

1. **Publicar mensajes:**
   ```bash
   curl -X POST http://localhost:8081/publish -H "Content-Type: application/json" -d '{"message":"Order #1"}'
   curl -X POST http://localhost:8081/publish -H "Content-Type: application/json" -d '{"message":"Order #2"}'
   curl -X POST http://localhost:8081/publish -H "Content-Type: application/json" -d '{"message":"Order #3"}'
   ```

2. **Procesar exitosamente:**
   ```bash
   curl http://localhost:8081/consume
   # Response: {"status":"success","message":"Order #1"}
   ```

3. **Simular fallo (rechazar mensaje):**
   ```bash
   curl -X POST http://localhost:8081/reject
   # Response: {"status":"success","message":"Message rejected and sent to DLX: Order #2"}
   ```
   
   ⚠️ El mensaje "Order #2" ahora está en la DLQ

4. **Recuperar mensaje fallido de la DLQ:**
   ```bash
   curl http://localhost:8081/dlq/consume
   # Response: {"status":"success","message":"Message from DLQ: Order #2"}
   ```

## Verificación en RabbitMQ Management UI

1. Abre http://localhost:15672 (usuario: `guest`, contraseña: `guest`)

2. Ve a la pestaña **Queues**:
   - `messages`: Cola principal con argumentos DLX configurados
   - `messages.dlq`: Dead Letter Queue

3. Ve a la pestaña **Exchanges**:
   - `dlx.exchange`: Dead Letter Exchange de tipo `direct`

4. Haz clic en la cola `messages` para ver:
   - **Arguments**: `x-dead-letter-exchange` = `dlx.exchange`
   - **Arguments**: `x-dead-letter-routing-key` = `dlx.routing.key`

## Configuración

Puedes personalizar la configuración mediante variables de entorno (archivo `.env`):

```env
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE_NAME=messages
HTTP_PORT=8081
```

## Estructura del Proyecto

```
dlx-demo/
├── main.go                  # Aplicación principal
├── go.mod                   # Dependencias
├── .env.example             # Ejemplo de configuración
├── test_dlx.sh              # Script de prueba (Bash)
├── test_dlx.ps1             # Script de prueba (PowerShell)
├── handlers/
│   └── dlx_handlers.go      # Handlers HTTP
└── rabbitmq/
    ├── connection.go        # Conexión con DLX
    ├── dlx_setup.go         # Configuración DLX
    ├── publisher.go         # Publicación de mensajes
    └── consumer.go          # Consumo y rechazo de mensajes
```

## Casos de Uso Reales

### 1. Sistema de Procesamiento de Pedidos
- Pedidos que fallan en validación → DLQ
- Revisión manual de pedidos problemáticos
- Reintento después de corrección

### 2. Procesamiento de Pagos
- Pagos rechazados por el gateway → DLQ
- Análisis de patrones de fallos
- Notificación al equipo de soporte

### 3. Integración con APIs Externas
- Timeouts o errores de API → DLQ
- Reintentos con backoff exponencial
- Alertas para APIs caídas

## Troubleshooting

### Error: "Failed to connect to RabbitMQ"
- Verifica que RabbitMQ esté ejecutándose: `docker-compose ps`
- Verifica la URL de conexión en las variables de entorno

### Error: "Port 8081 already in use"
- Cambia el puerto en `.env`: `HTTP_PORT=8082`
- O detén el proceso que usa el puerto 8081

### Los mensajes no llegan a la DLQ
- Verifica que estés usando `/reject` (no `/consume`)
- Revisa los logs del servicio
- Verifica la configuración DLX en RabbitMQ Management UI

## Próximos Pasos

Posibles mejoras para este demo:

- ⏱️ Agregar TTL (Time To Live) a los mensajes
- 🔁 Implementar reintentos automáticos desde la DLQ
- 📧 Notificaciones cuando hay mensajes en la DLQ
- 📈 Métricas y monitoreo de mensajes fallidos
- 🔐 Autenticación en los endpoints HTTP

## Referencias

- [RabbitMQ Dead Letter Exchanges](https://www.rabbitmq.com/dlx.html)
- [RabbitMQ Go Client](https://github.com/rabbitmq/amqp091-go)
- [Reliable Messaging Patterns](https://www.rabbitmq.com/reliability.html)
