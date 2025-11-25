
# 📚 Plan de Acción y Estudio: RoadMap GGSS

[cite_start]Este documento resume los tópicos clave necesarios para comprender y trabajar en el ecosistema técnico del Gestor de Solicitudes (GGSS), basados en el plan de acción del RoadMap GGSS[cite: 1].

---

## 1. Servicios & Broker's: RabbitMQ (Mensajería Asíncrona)

[cite_start]**Objetivo:** Comprender el funcionamiento de RabbitMQ como broker de mensajería y sus patrones de fiabilidad[cite: 3].

### 1.1. [cite_start]Arquitectura Básica [cite: 5]
* **Exchange:** Es el punto de entrada de RabbitMQ. Recibe mensajes del productor y decide a dónde enrutarlos.
* **Queue (Cola):** Componente donde el mensaje espera hasta que es consumido. Es la unidad de almacenamiento.
* **Binding:** La regla o "puente" que conecta un Exchange con una Queue.
* **Routing Key:** La etiqueta que el productor adjunta al mensaje y que el Exchange usa para hacer la coincidencia con los Bindings.

### 1.2. [cite_start]Tipos de Exchange [cite: 6]
La forma en que el Exchange usa la Routing Key depende de su tipo:

| Tipo | Comportamiento Clave | Coincidencia de la Routing Key |
| :--- | :--- | :--- |
| **Direct** | Enrutamiento **exacto**. | Debe coincidir perfectamente con la clave del Binding. |
| **Fanout** | **Broadcast** (retransmisión). | **Ignora** la Routing Key; envía el mensaje a todas las colas conectadas. |
| **Topic** | Enrutamiento por **patrones**. | Permite comodines (`*` para una palabra, `#` para cero o más palabras) para filtrar por temas jerárquicos. |
| **Headers** | Enrutamiento por **metadatos**. | **Ignora** la Routing Key y filtra por los pares clave-valor que se encuentran dentro de los encabezados del mensaje. |

### 1.3. [cite_start]Mensajería Confiable y Alta Disponibilidad [cite: 7, 11]
* **ACK / NACK:** Mecanismo para confirmar (ACK) o rechazar (NACK) el procesamiento. Solo un ACK elimina el mensaje de la cola.
* [cite_start]**DLX (Dead Letter Exchange):** Actúa como un "cementerio de mensajes"[cite: 7]. [cite_start]Captura mensajes que fallaron repetidamente (NACKs) o expiraron, reenviándolos a una cola especial (DLQ) para **despejar la cola principal**[cite: 10].
* [cite_start]**Quorum Queues:** Soluciona el problema de **fallo de servidor** y asegura la **Alta Disponibilidad**[cite: 7, 11]. Utiliza un algoritmo de consenso para replicar la cola en varios nodos (cluster), garantizando que los mensajes no se pierdan.

---

## 2. gRPC y HTTP (Protocolos de Comunicación)

[cite_start]**Objetivo:** Comparar gRPC con HTTP/REST y entender cuándo elegir cada uno[cite: 13, 14, 16].

### 2.1. [cite_start]HTTP/REST [cite: 17, 18, 19, 21]
* [cite_start]**Arquitectura:** Basada en recursos[cite: 18].
* [cite_start]**Formato:** Principalmente **JSON** como formato común[cite: 19].
* [cite_start]**Ventajas:** Simplicidad y compatibilidad universal[cite: 21].

### [cite_start]2.2. gRPC [cite: 22, 23, 24, 25]
* [cite_start]**Tecnología:** Basado en **HTTP/2** y usa **Protocol Buffers** (Protobuf)[cite: 23].
* [cite_start]**Comunicación:** Formato **binario** eficiente y permite **streaming bidireccional**[cite: 24].
* [cite_start]**Contratos:** Usa contratos fuertes mediante IDL (Interface Definition Language)[cite: 25].
* **Ventajas Clave:** Mayor velocidad y menor tamaño de mensaje debido al formato binario.

---

## 3. Contenerización Local

[cite_start]**Objetivo:** Usar Docker y Docker Compose para levantar entornos locales integrados[cite: 30, 31].

* [cite_start]**Conceptos:** Entender qué es la contenerización y por qué usarla[cite: 28].
* [cite_start]**Diferencias:** Diferenciar contenedores de Máquinas Virtuales (VM)[cite: 29].
* [cite_start]**Docker Compose:** Usar para levantar entornos locales[cite: 30].
* [cite_start]**Demo Práctica (Docker Compose):** Creación de un `docker-compose.yml` que integre servicios como MySQL, Kong y una API REST[cite: 33, 34, 35, 36].
* **Nota de Arquitectura (Windows):** Docker Desktop requiere utilizar una máquina virtual (como Hyper-V) para ejecutar contenedores Linux, a diferencia de sistemas Linux que los ejecutan de forma nativa.

---

## 4. El Producto GGSS y Stack Técnico

[cite_start]**Objetivo:** Conocer el Gestor de Solicitudes (GGSS) y su arquitectura técnica[cite: 40].

### 4.1. [cite_start]Stack Técnico Clave [cite: 42, 47, 48]
* [cite_start]**Brokers:** Uso de **RabbitMQ y Kafka** en el GGSS[cite: 42].
* [cite_start]**Librerías Compartidas:** `gs-commons` [cite: 43] [cite_start]y `kafka-toolkit`[cite: 44].
* [cite_start]**API Gateway y Monitoreo:** **Kong** y **DataDog**[cite: 47].
* [cite_start]**Protocolos:** Implementación de **gRPC & HTTP**[cite: 48].

### 4.2. Herramientas y Despliegue Local
* [cite_start]**Herramientas:** Postman y ApiDog[cite: 50].
* [cite_start]**Montaje:** Capacidad para levantar el entorno de desarrollo del GGSS localmente usando `docker-local-dev-services`[cite: 51].