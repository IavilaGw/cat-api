# Cat Service API

API REST con Golang y Gin 

## Requisitos

- Docker
- Docker Compose

## Instalacion
```bash
git clone https://github.com/IavilaGw/cat-api.git
cd cat-api
```

## Ejecucion
```bash
# Iniciar servicios
./run.sh start

# Detener servicios
./run.sh stop
```

## Endpoints

- **GET** `/api/cat` - Obtener imagen aleatoria de gato
- **GET** `/api/count` - Obtener conteo de imagenes unicas
- **GET** `/api/stats` - Obtener estadisticas


## Docker Hub

La imagen esta disponible en Docker Hub:
```bash
docker pull iavilagw/cat-api:latest
```

**Enlace:** https://hub.docker.com/r/iavilagw/cat-api

## Tecnologias

- Go 1.21
- Gin Framework
- PostgreSQL 15
- Docker




