# RabbitMQ Go Service

Servicio en Go que integra RabbitMQ para publicar y consumir mensajes a través de endpoints HTTP.

Este repositorio demuestra diversas **estrategias de envío y publicación de mensajes**, exponiendo una API HTTP para facilitar la integración.
El proyecto abarca desde conceptos básicos hasta patrones avanzados de confiabilidad, incluyendo:
*   **Endpoints HTTP** para publicación y consumo.
*   **Dead Letter Exchange (DLX)** para gestión de errores y reintentos.
*   **Quorum Queues** para alta disponibilidad y tolerancia a fallos.
*   Configuración completa con **Docker Compose**.

## Características

- ✅ Publicar mensajes a RabbitMQ mediante endpoint HTTP POST
- ✅ Consumir mensajes de RabbitMQ mediante endpoint HTTP GET
- ✅ RabbitMQ ejecutándose en Docker Compose
- ✅ Interfaz de administración de RabbitMQ
- ✅ Manejo de errores y logging
- ✅ Configuración mediante variables de entorno

## Requisitos

- Go 1.21 o superior
- Docker y Docker Compose

## Instalación

1. **Clonar o navegar al directorio del proyecto:**
   ```bash
   cd c:\wk\rabbitmq
   ```

2. **Iniciar RabbitMQ con Docker Compose:**
   ```bash
   docker-compose up -d
   ```

3. **Verificar que RabbitMQ esté ejecutándose:**
   ```bash
   docker-compose ps
   ```

4. **Descargar dependencias de Go:**
   ```bash
   go mod download
   ```

## Configuración

Puedes configurar el servicio mediante variables de entorno. Crea un archivo `.env` basado en `.env.example`:

```env
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE_NAME=messages
HTTP_PORT=8080
```

## Uso

### 1. Iniciar el servicio

```bash
go run main.go
```

El servicio estará disponible en `http://localhost:8080`

### 2. Publicar un mensaje

Envía una petición POST al endpoint `/publish`:

```bash
curl -X POST http://localhost:8080/publish \
  -H "Content-Type: application/json" \
  -d '{"message":"Hola RabbitMQ desde Go!"}'
```

**Respuesta esperada:**
```json
{
  "status": "success",
  "message": "Message published successfully"
}
```

### 3. Consumir un mensaje

Envía una petición GET al endpoint `/consume`:

```bash
curl http://localhost:8080/consume
```

**Respuesta esperada (si hay mensajes):**
```json
{
  "status": "success",
  "message": "Hola RabbitMQ desde Go!"
}
```

**Respuesta si no hay mensajes:**
```json
{
  "status": "error",
  "error": "no messages available in queue"
}
```

### 4. Health Check

Verifica que el servicio esté funcionando:

```bash
curl http://localhost:8080/health
```

**Respuesta:**
```json
{
  "status": "healthy"
}
```

## Interfaz de Administración de RabbitMQ

Accede a la interfaz web de RabbitMQ en:
- **URL:** http://localhost:15672
- **Usuario:** guest
- **Contraseña:** guest

Desde aquí puedes:
- Ver las colas y sus mensajes
- Monitorear el estado del servidor
- Ver estadísticas de mensajes publicados/consumidos

## Estructura del Proyecto

```
c:\wk\rabbitmq\
├── docker-compose.yml      # Configuración de Docker para RabbitMQ
├── go.mod                  # Dependencias de Go
├── main.go                 # Punto de entrada de la aplicación
├── .env.example            # Ejemplo de variables de entorno
├── README.md               # Este archivo
├── handlers/
│   └── handlers.go         # Handlers HTTP para publish/consume
└── rabbitmq/
    ├── connection.go       # Gestión de conexión a RabbitMQ
    ├── publisher.go        # Lógica de publicación de mensajes
    └── consumer.go         # Lógica de consumo de mensajes
```

## API Endpoints

### POST /publish
Publica un mensaje en la cola de RabbitMQ.

**Request Body:**
```json
{
  "message": "Tu mensaje aquí"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Message published successfully"
}
```

### GET /consume
Consume un mensaje de la cola de RabbitMQ.

**Response (éxito):**
```json
{
  "status": "success",
  "message": "Contenido del mensaje"
}
```

**Response (sin mensajes):**
```json
{
  "status": "error",
  "error": "no messages available in queue"
}
```

### GET /health
Verifica el estado del servicio.

**Response:**
```json
{
  "status": "healthy"
}
```

## Detener el Servicio

1. **Detener la aplicación Go:** Presiona `Ctrl+C` en la terminal donde está corriendo

2. **Detener RabbitMQ:**
   ```bash
   docker-compose down
   ```

3. **Detener y eliminar volúmenes (borra todos los mensajes):**
   ```bash
   docker-compose down -v
   ```

## Troubleshooting

### Error: "Failed to connect to RabbitMQ"
- Verifica que RabbitMQ esté ejecutándose: `docker-compose ps`
- Verifica que el puerto 5672 esté disponible
- Revisa los logs de RabbitMQ: `docker-compose logs rabbitmq`

### Error: "Port already in use"
- Cambia el puerto HTTP en las variables de entorno
- Verifica que no haya otra aplicación usando el puerto 8080

### No se pueden consumir mensajes
- Verifica que hayas publicado mensajes primero
- Revisa la cola en la interfaz de administración de RabbitMQ

## Próximos Pasos

Posibles mejoras para el proyecto:
- Agregar autenticación a los endpoints
- Implementar reintentos automáticos en caso de fallo
- Agregar métricas y monitoreo
- Implementar diferentes tipos de exchanges (fanout, topic, headers)
- Agregar tests unitarios e integración

---

## 🔥 Demo: Mensajería Confiable con Dead Letter Exchange (DLX)

Este proyecto incluye una **demostración completa de Dead Letter Exchange (DLX)** que muestra cómo implementar mensajería confiable en RabbitMQ.

### ¿Qué es DLX?

**Dead Letter Exchange (DLX)** es una característica de RabbitMQ que permite manejar mensajes que fallan en su procesamiento, enviándolos a una cola especial (Dead Letter Queue) en lugar de perderlos.

### Características de la Demo

- ✅ Configuración automática de DLX y DLQ
- ✅ Endpoint para simular fallos de procesamiento
- ✅ Recuperación de mensajes desde la DLQ
- ✅ Scripts de prueba automatizados (PowerShell y Bash)
- ✅ Documentación completa con ejemplos

### Inicio Rápido

```bash
# 1. Navega al directorio de la demo
cd dlx-demo

# 2. Descarga dependencias
go mod download

# 3. Inicia el servicio (puerto 8081)
go run main.go

# 4. En otra terminal, ejecuta el test
.\test_dlx.ps1  # Windows PowerShell
# o
./test_dlx.sh   # Linux/Mac/Git Bash
```

### Documentación Completa

Para más información sobre la implementación de DLX, arquitectura, casos de uso y ejemplos detallados, consulta:

📖 **[dlx-demo/README_DLX.md](dlx-demo/README_DLX.md)**

### Arquitectura

```
Producer → Main Queue → Consumer (OK)
              ↓
           Reject/Nack
              ↓
        DLX Exchange → Dead Letter Queue → DLQ Consumer
```

La demo incluye endpoints para:
- `POST /publish` - Publicar mensajes
- `GET /consume` - Consumir mensajes exitosamente
- `POST /reject` - Rechazar mensajes (simular fallo) → envía a DLX
- `GET /dlq/consume` - Recuperar mensajes de la DLQ

---

## 🚀 Demo: Alta Disponibilidad con Quorum Queues

Este proyecto incluye una **demostración completa de Quorum Queues** que muestra cómo implementar alta disponibilidad y mensajería confiable mediante replicación en RabbitMQ.

### ¿Qué son las Quorum Queues?

**Quorum Queues** son colas modernas de RabbitMQ diseñadas para alta disponibilidad mediante replicación automática usando el algoritmo de consenso **Raft**.

### Características de la Demo

- 🔄 **Cluster de 3 nodos** RabbitMQ (replicación automática)
- ✅ **Publisher confirmations** (garantía de entrega al broker)
- 🎯 **Manual acknowledgments** (control fino de procesamiento)
- 🛡️ **Alta disponibilidad** (funciona si 1 nodo falla)
- 📊 **Algoritmo Raft** (consenso y elección de líder)
- 💾 **Durabilidad** (mensajes persistidos automáticamente)

### Diferencias entre Demos

| Demo | Propósito | Característica Principal |
|------|-----------|-------------------------|
| **Original** | Básico | Publisher/Consumer simple |
| **DLX** | Manejo de fallos | Dead Letter Queue para mensajes rechazados |
| **Quorum** | Alta disponibilidad | Replicación en cluster de 3 nodos |

### Inicio Rápido

```bash
# 1. Navega al directorio de la demo
cd quorum-demo

# 2. Inicia el cluster de RabbitMQ (3 nodos)
docker-compose up -d

# 3. Espera 30 segundos para que el cluster se forme

# 4. Descarga dependencias
go mod download

# 5. Inicia el servicio (puerto 8082)
go run main.go

# 6. En otra terminal, ejecuta el test
.\test_quorum.ps1  # Windows PowerShell
# o
./test_quorum.sh   # Linux/Mac/Git Bash
```

### Documentación Completa

Para más información sobre Quorum Queues, arquitectura del cluster, pruebas de failover y casos de uso, consulta:

📖 **[quorum-demo/README_QUORUM.md](quorum-demo/README_QUORUM.md)**

### Arquitectura del Cluster

```
Publisher → Node 1 (Leader) → Node 2 (Follower)
                ↓
            Node 3 (Follower)
                ↓
            Consumer (con ACK manual)
```

La demo incluye:
- **3 nodos RabbitMQ** en cluster (puertos 5672-5674)
- **Management UIs** para cada nodo (puertos 15672-15674)
- **Endpoints HTTP**:
  - `POST /publish` - Publica con confirmación del broker
  - `GET /consume` - Consume con ACK manual
  - `POST /consume/fail` - Consume con NACK (requeue)
  - `GET /stats` - Estadísticas de la cola

### Prueba de Alta Disponibilidad

```bash
# Detener un nodo
docker stop rabbitmq-node2

# El servicio sigue funcionando!
curl -X POST http://localhost:8082/publish \
  -H "Content-Type: application/json" \
  -d '{"message":"Still working!"}'

# Reiniciar el nodo
docker start rabbitmq-node2
```
