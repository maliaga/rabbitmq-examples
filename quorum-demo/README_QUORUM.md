# RabbitMQ Quorum Queues Demo

Demostración completa de **Quorum Queues** en RabbitMQ para alta disponibilidad y mensajería confiable mediante replicación basada en el algoritmo de consenso **Raft**.

## ¿Qué son las Quorum Queues?

**Quorum Queues** son un tipo moderno de cola en RabbitMQ diseñadas para proporcionar:

- 🔄 **Alta Disponibilidad**: Replicación automática en múltiples nodos
- 🛡️ **Durabilidad**: Mensajes persistidos en disco automáticamente
- ✅ **Confiabilidad**: Algoritmo Raft garantiza consenso entre nodos
- 📊 **Data Safety**: No se pierden mensajes incluso si falla un nodo
- ⚡ **Publisher Confirmations**: Garantía de que el broker aceptó el mensaje

## Arquitectura del Cluster

```mermaid
graph TB
    subgraph "RabbitMQ Cluster"
        N1[Node 1<br/>Leader<br/>:5672]
        N2[Node 2<br/>Follower<br/>:5673]
        N3[Node 3<br/>Follower<br/>:5674]
    end
    
    P[Publisher<br/>with Confirmations] -->|Publish| N1
    N1 -->|Replicate| N2
    N1 -->|Replicate| N3
    N1 -->|Consume| C[Consumer<br/>with Manual ACK]
    
    style N1 fill:#90EE90
    style N2 fill:#87CEEB
    style N3 fill:#87CEEB
    style P fill:#FFD700
    style C fill:#FFA500
```

### Componentes del Cluster

1. **3 Nodos RabbitMQ**: Cluster con replicación automática
2. **Quorum Queue** (`orders-quorum`): Cola replicada en los 3 nodos
3. **Raft Consensus**: Algoritmo que garantiza consistencia
4. **Leader-Follower**: Un nodo líder, dos followers

## Requisitos

- Go 1.21 o superior
- Docker y Docker Compose
- Puertos disponibles: 5672-5674, 15672-15674, 8082

## Instalación y Configuración

### 1. Iniciar el Cluster de RabbitMQ

```bash
cd quorum-demo
docker-compose up -d
```

Esto iniciará:
- **rabbitmq-node1** (puerto 5672, UI 15672)
- **rabbitmq-node2** (puerto 5673, UI 15673)
- **rabbitmq-node3** (puerto 5674, UI 15674)

### 2. Verificar el Cluster

```bash
docker-compose ps
```

Deberías ver 3 contenedores ejecutándose. Espera unos 30 segundos para que el cluster se forme completamente.

### 3. Verificar el Estado del Cluster

```bash
docker exec rabbitmq-node1 rabbitmqctl cluster_status
```

Deberías ver los 3 nodos en el cluster.

### 4. Descargar Dependencias de Go

```bash
go mod download
```

## Uso

### 1. Iniciar el Servicio

```bash
go run main.go
```

El servicio estará disponible en `http://localhost:8082`

**Salida esperada:**
```
========================================
RabbitMQ Quorum Queue Demo Service
========================================
HTTP Server: http://localhost:8082

Endpoints:
  POST http://localhost:8082/publish       - Publish with confirmation
  GET  http://localhost:8082/consume       - Consume with ACK
  POST http://localhost:8082/consume/fail  - Consume with NACK (requeue)
  GET  http://localhost:8082/stats         - Queue statistics
  GET  http://localhost:8082/health        - Health check

RabbitMQ Cluster Management UIs:
  Node 1: http://localhost:15672 (guest/guest)
  Node 2: http://localhost:15673 (guest/guest)
  Node 3: http://localhost:15674 (guest/guest)
========================================
```

### 2. Ejecutar el Test Automático

**PowerShell (Windows):**
```powershell
.\test_quorum.ps1
```

**Bash (Linux/Mac/Git Bash):**
```bash
chmod +x test_quorum.sh
./test_quorum.sh
```

El script demostrará:
1. ✅ Publicación con confirmaciones del broker
2. ✅ Consumo con acknowledgments manuales
3. ✅ Simulación de fallo (NACK con requeue)
4. ✅ Estadísticas de la cola

## API Endpoints

### POST /publish
Publica un mensaje con confirmación del broker.

**Request:**
```bash
curl -X POST http://localhost:8082/publish \
  -H "Content-Type: application/json" \
  -d '{"message":"Order #12345"}'
```

**Response:**
```json
{
  "status": "success",
  "message": "Message published and confirmed by broker"
}
```

**Características:**
- ✅ Espera confirmación del broker antes de retornar
- ✅ Mensaje replicado en los 3 nodos
- ✅ Persistido en disco automáticamente

---

### GET /consume
Consume un mensaje con acknowledgment manual.

**Request:**
```bash
curl http://localhost:8082/consume
```

**Response:**
```json
{
  "status": "success",
  "message": "Order #12345"
}
```

**Características:**
- ✅ ACK manual después de consumir
- ✅ Mensaje eliminado solo después del ACK
- ✅ Si el consumer falla antes del ACK, el mensaje se mantiene

---

### POST /consume/fail
Consume un mensaje y lo rechaza (simula fallo de procesamiento).

**Request:**
```bash
curl -X POST http://localhost:8082/consume/fail
```

**Response:**
```json
{
  "status": "success",
  "message": "Message consumed and rejected (requeued): Order #12345"
}
```

**Características:**
- ✅ NACK con requeue=true
- ✅ Mensaje vuelve a la cola para reintento
- ✅ Útil para errores transitorios

---

### GET /stats
Muestra estadísticas de la cola.

**Request:**
```bash
curl http://localhost:8082/stats
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "queue_name": "orders-quorum",
    "queue_type": "quorum",
    "messages": 5,
    "consumers": 0
  }
}
```

---

### GET /health
Verifica el estado del servicio.

**Request:**
```bash
curl http://localhost:8082/health
```

**Response:**
```json
{
  "status": "healthy",
  "queue_type": "quorum"
}
```

## Pruebas de Alta Disponibilidad

### Prueba de Failover

Demuestra que el servicio continúa funcionando incluso si un nodo falla.

**1. Publicar mensajes:**
```bash
curl -X POST http://localhost:8082/publish -H "Content-Type: application/json" -d '{"message":"Test HA 1"}'
curl -X POST http://localhost:8082/publish -H "Content-Type: application/json" -d '{"message":"Test HA 2"}'
curl -X POST http://localhost:8082/publish -H "Content-Type: application/json" -d '{"message":"Test HA 3"}'
```

**2. Detener un nodo:**
```bash
docker stop rabbitmq-node2
```

**3. Verificar que el servicio sigue funcionando:**
```bash
# Publicar más mensajes
curl -X POST http://localhost:8082/publish -H "Content-Type: application/json" -d '{"message":"After node2 stopped"}'

# Consumir mensajes
curl http://localhost:8082/consume
```

✅ **El servicio sigue funcionando!** Los mensajes están replicados en node1 y node3.

**4. Reiniciar el nodo:**
```bash
docker start rabbitmq-node2
```

El nodo se re-sincronizará automáticamente con el cluster.

## Verificación en RabbitMQ Management UI

### Ver la Cola Quorum

1. Abre http://localhost:15672 (o 15673, 15674)
2. Login: `guest` / `guest`
3. Ve a la pestaña **Queues**
4. Busca `orders-quorum`

**Verás:**
- **Type**: `quorum`
- **Features**: `D` (Durable)
- **Replicas**: Información sobre los 3 nodos
- **Leader**: Nodo que actúa como líder

### Ver el Cluster

1. Ve a **Admin** → **Cluster**
2. Verás los 3 nodos: `rabbit@rabbitmq-node1`, `rabbit@rabbitmq-node2`, `rabbit@rabbitmq-node3`
3. Estado: `running`

## Comparación: Quorum vs Classic Queues

| Característica | Classic Queue | Quorum Queue |
|----------------|---------------|--------------|
| **Replicación** | Opcional (mirroring) | Automática (Raft) |
| **Durabilidad** | Configurable | Siempre durable |
| **Confirmaciones** | Opcional | Recomendado |
| **Alta Disponibilidad** | Mirroring manual | Built-in |
| **Consistencia** | Eventual | Fuerte (Raft) |
| **Performance** | Más rápido (sin replicación) | Optimizado para durabilidad |
| **Uso Recomendado** | Mensajes efímeros | Mensajes críticos |

## Casos de Uso Reales

### 1. Sistema de Procesamiento de Pedidos
- **Problema**: No se pueden perder pedidos
- **Solución**: Quorum Queue con confirmaciones
- **Beneficio**: Pedidos replicados, sobreviven a fallos de nodos

### 2. Sistema de Pagos
- **Problema**: Transacciones críticas
- **Solución**: Quorum Queue + Manual ACK
- **Beneficio**: Garantía de procesamiento, no se pierden transacciones

### 3. Sistema de Notificaciones Importantes
- **Problema**: Emails/SMS críticos deben enviarse
- **Solución**: Quorum Queue con reintento (NACK)
- **Beneficio**: Alta disponibilidad, reintentos automáticos

## Configuración Avanzada

### Variables de Entorno

Crea un archivo `.env`:

```env
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE_NAME=orders-quorum
HTTP_PORT=8082
```

### Ajustar Tamaño del Quorum

En `quorum_setup.go`, puedes especificar el tamaño inicial:

```go
args := amqp.Table{
    "x-queue-type": "quorum",
    "x-quorum-initial-group-size": 3, // Número de réplicas
}
```

## Estructura del Proyecto

```
quorum-demo/
├── docker-compose.yml         # Cluster de 3 nodos
├── main.go                    # Aplicación principal
├── go.mod                     # Dependencias
├── .env.example               # Configuración ejemplo
├── test_quorum.ps1           # Test PowerShell
├── test_quorum.sh            # Test Bash
├── handlers/
│   └── quorum_handlers.go    # Handlers HTTP
└── rabbitmq/
    ├── connection.go         # Conexión con confirmaciones
    ├── quorum_setup.go       # Setup de Quorum Queue
    ├── publisher.go          # Publisher con confirmaciones
    └── consumer.go           # Consumer con ACK manual
```

## Troubleshooting

### Error: "Failed to connect to RabbitMQ"
- Verifica que el cluster esté ejecutándose: `docker-compose ps`
- Espera 30 segundos después de `docker-compose up -d`
- Verifica logs: `docker-compose logs rabbitmq-node1`

### Error: "Port already in use"
- Cambia el puerto en `.env`: `HTTP_PORT=8083`
- O detén el proceso que usa el puerto 8082

### Los nodos no forman cluster
- Verifica que `RABBITMQ_ERLANG_COOKIE` sea igual en todos los nodos
- Reinicia el cluster: `docker-compose down && docker-compose up -d`
- Verifica logs: `docker-compose logs cluster-setup`

### Mensajes no se replican
- Verifica que la cola sea tipo `quorum` en Management UI
- Verifica que los 3 nodos estén en el cluster
- Revisa logs del servicio

## Detener el Cluster

```bash
# Detener servicios
docker-compose down

# Detener y eliminar volúmenes (borra todos los datos)
docker-compose down -v
```

## Próximos Pasos

Mejoras posibles:

- 🔐 Agregar autenticación TLS
- 📊 Implementar métricas de replicación
- 🔁 Configurar políticas de retry automático
- 📈 Monitoreo de salud del cluster
- 🚀 Optimización de performance

## Referencias

- [RabbitMQ Quorum Queues](https://www.rabbitmq.com/quorum-queues.html)
- [Raft Consensus Algorithm](https://raft.github.io/)
- [RabbitMQ Clustering](https://www.rabbitmq.com/clustering.html)
- [Publisher Confirms](https://www.rabbitmq.com/confirms.html)
