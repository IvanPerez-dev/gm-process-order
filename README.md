# 🧪 Prueba Técnica - Grupo Mariposa

## 📌 Descripción

Esta prueba técnica consiste en la implementación de una arquitectura distribuida para el procesamiento de pedidos, utilizando:

* API en **Go** para gestión de recursos
* Worker en **Java (Spring WebFlux)** para procesamiento reactivo
* Integración con **Kafka**, **MongoDB** y **Redis**

El sistema permite crear pedidos, enriquecerlos con información externa y almacenarlos de forma resiliente.

---

## 🚀 Ejecución rápida

### 🔹 Levantar todo (infraestructura + aplicaciones)

```bash
docker compose -f docker-compose.yml -f docker-compose.apps.yml up -d --build
```

### 🔹 Parar todo (infraestructura + aplicaciones)

```bash
docker compose up -d
```

Esto levanta:

* Kafka
* MongoDB
* Redis
* Otros servicios necesarios

## 🏗️ Arquitectura

El sistema está compuesto por:

### API en Go

Expone 3 recursos principales:

## 📡 Endpoints

### Orders

#### Crear orden
POST http://localhost:8081/api/v1/orders

**Body:**
```json
{
  "customerId": "c1a2b3c4-0001-4000-a000-000000000001",
  "items": [
    {
      "productId": "p1a2b3c4-0001-4000-a000-000000000001",
      "quantity": 10
    },
    {
      "productId": "p1a2b3c4-0001-4000-a000-000000000002",
      "quantity": 5
    }
  ]
}

```
## Listar todas las órdenes

GET http://localhost:8081/api/v1/orders

## Obtener orden por ID

GET http://localhost:8081/api/v1/orders/{orderId}

# Products

- POST /api/v1/products
- GET /api/v1/products/{id}
- PUT /api/v1/products/{id}
- DELETE /api/v1/products/{id}

# Customers

- POST /api/v1/customers
- GET /api/v1/customers/{id}
- PUT /api/v1/customers/{id}
- DELETE /api/v1/customers/{id}

---

### 🔹 Flujo de procesamiento

1. Se crea una orden vía API (`orders`)

```json
{
 "customerId": "c1a2b3c4-0001-4000-a000-000000000001", 
   "items": [
     {
       "productId": "p1a2b3c4-0001-4000-a000-000000000001",
       "quantity": 10
     }, 
      {
       "productId": "p1a2b3c4-0001-4000-a000-000000000002",
       "quantity": 5
     }
   ]
  }
```
2. La orden se guarda en **MongoDB** con estado `PENDING`
3. Se emite un evento a **Kafka** (`order.created`)
4. Un Worker en Java consume el evento
5. El Worker:

   * Consulta APIs de **products** y **customers**
   * Valida datos
   * Enriquece la información
   * Guarda la orden procesada en MongoDB
6. validar la orden enriquecida

```url
http://localhost:8081/api/v1/orders/{orderId}

```

---

### 🔹 Worker en Java

Implementado con:

* **Java 21**
* **Spring Boot**
* **Spring WebFlux**

Responsabilidades:

* Consumir eventos desde Kafka
* Enriquecer datos del pedido
* Validar existencia de cliente y productos
* Manejar errores y reintentos
* Persistir en MongoDB

---

## 📦 Estructura del Pedido

```json
{
  "_id": "abe24763-6283-48a7-827a-aa034d113c19",
  "customer": {
    "_id": "c1a2b3c4-0001-4000-a000-000000000001",
    "name": "Juan Pérez",
    "email": "juan.perez@grupomariposa.com"
  },
  "products": [
    {
      "productId": "p1a2b3c4-0001-4000-a000-000000000001",
      "name": "Producto Alpha",
      "price": 99.99, 
      "quantity": 10,
      "total": 999.90
    }
  ],
  "status": "PROCESSED",
  "retryCount": 0,
  "createdAt": "2026-03-23T16:29:27.581Z",
  "processedAt": "2026-03-23T16:29:45.086Z"
}
```

---

## ⚙️ Manejo de Errores y Resiliencia

* 🔁 Reintentos exponenciales para llamadas HTTP
* 🧠 Uso de **Redis** para:

  * Almacenar mensajes fallidos
  * Contador de reintentos
* 🔒 Lock distribuido para evitar procesamiento duplicado

---

## 🧪 Pruebas

El proyecto incluye:

* ✅ Pruebas unitarias
* ✅ Cobertura de lógica principal

---

## 📡 Tecnologías utilizadas

* Go
* Java 21
* Spring Boot
* Spring WebFlux
* Apache Kafka
* MongoDB
* Redis
* Docker

---

## 💡 Consideraciones

* Arquitectura modular y desacoplada como parte de Clean Architecture
* Separación clara entre API y procesamiento
* Uso de programación reactiva para eficiencia
* Preparado para escalabilidad
* Separación clara entre API y procesamiento

---

## ✅ Notas finales

El sistema está diseñado para simular un flujo real de procesamiento de pedidos en entornos distribuidos, aplicando buenas prácticas de:

* Clean Architecture
* Manejo de eventos
* Resiliencia
* Integración entre servicios

