# Guía de Inicio Rápido

Este documento te ayudará a poner en marcha el proyecto rápidamente.

## Pasos para Ejecutar

### 1. Asegúrate de que Docker Desktop esté ejecutándose

Verifica que Docker Desktop esté corriendo en tu sistema.

### 2. Inicia RabbitMQ

```bash
cd c:\wk\rabbitmq
docker-compose up -d
```

Espera unos segundos para que RabbitMQ esté completamente iniciado.

### 3. Inicia la aplicación Go

```bash
go run main.go
```

Verás un mensaje indicando que el servidor está corriendo en el puerto 8080.

### 4. Prueba los endpoints

**Publicar un mensaje:**
```bash
curl -X POST http://localhost:8080/publish -H "Content-Type: application/json" -d "{\"message\":\"Hola RabbitMQ!\"}"
```

**Consumir un mensaje:**
```bash
curl http://localhost:8080/consume
```

## Interfaz de Administración

Abre tu navegador y ve a: http://localhost:15672
- Usuario: `guest`
- Contraseña: `guest`

¡Listo! Tu proyecto está funcionando.
