#!/bin/bash
set -e

check_docker() {
    if ! command -v docker &> /dev/null; then
        echo "Error: Docker no esta instalado"
        exit 1
    fi
    if ! command -v docker-compose &> /dev/null; then
        echo "Error: Docker Compose no esta instalado"
        exit 1
    fi
}

check_files() {
    if [ ! -f "docker-compose.yml" ]; then
        echo "Error: docker-compose.yml no encontrado"
        exit 1
    fi
}

setup_env() {
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
        else
            cat > .env << 'ENVEOF'
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
GIN_MODE=release
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=catdb
DB_SSLMODE=disable
CATAAS_API_URL=https://cataas.com
TIMEOUT_SECONDS=30
ENVEOF
        fi
        echo "Archivo .env creado"
    fi
}

start_services() {
    echo "Iniciando servicios..."
    docker-compose up -d
    echo "Servicios iniciados"
    echo "Esperando..."
    sleep 10
    
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "Servicio funcionando"
        echo ""
        echo "Endpoints:"
        echo "  http://localhost:8080/api/cat"
        echo "  http://localhost:8080/api/count"
        echo "  http://localhost:8080/api/stats"
    else
        echo "Advertencia: Servicio no responde aun"
    fi
}

stop_services() {
    echo "Deteniendo servicios..."
    docker-compose down
    echo "Servicios detenidos"
}

restart_services() {
    echo "Reiniciando servicios..."
    docker-compose restart
    sleep 5
    echo "Servicios reiniciados"
}

test_endpoints() {
    echo "Probando endpoints..."
    
    if curl -s http://localhost:8080/health | grep -q "healthy"; then
        echo "Health: OK"
    else
        echo "Health: FAIL"
    fi
    
    if curl -s -o /tmp/test.jpg http://localhost:8080/api/cat; then
        echo "Get cat: OK"
        rm -f /tmp/test.jpg
    else
        echo "Get cat: FAIL"
    fi
    
    count=$(curl -s http://localhost:8080/api/count)
    echo "Count: $count"
    
    if curl -s http://localhost:8080/api/stats | grep -q "total_images"; then
        echo "Stats: OK"
    else
        echo "Stats: FAIL"
    fi
}

case "${1:-}" in
    start)
        check_docker
        check_files
        setup_env
        start_services
        ;;
    stop)
        check_docker
        stop_services
        ;;
    restart)
        check_docker
        restart_services
        ;;
    test)
        test_endpoints
        ;;
    *)
        echo "Uso: ./run.sh [start|stop|restart|test]"
        exit 1
        ;;
esac