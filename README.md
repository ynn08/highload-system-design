# Food Delivery Platform

Микросервисная платформа доставки еды. 
Стек: Go 1.24, PostgreSQL, Redis, Elasticsearch, Kafka.
Архитектура: Clean Architecture, Event-Driven (Saga, Outbox).

Проект адаптирован под Задание 3 (PoC, оптимизация под лимиты 2 vCPU / 8 GB RAM).

## Структура проекта

*   `api-gateway` — точка входа, rate limiting, auth middleware, reverse proxy.
*   `catalog-service` — геопоиск ресторанов и отдача меню (Elasticsearch).
*   `order-service` — процессинг корзины (Redis) и заказов (PostgreSQL). Реализует Outbox для отправки событий в Kafka.
*   `payment-service` — эмуляция оплаты, проверка идемпотентности по ключу (PostgreSQL).
*   `delivery-service` — трекинг доставки, алгоритм батчинга курьеров (In-Memory/Kafka).

## Запуск

1. Поднять инфраструктуру и приложения (лимиты ресурсов прописаны в `docker-compose.yml`):
   ```bash
   docker-compose up -d --build
   ```
2. Дождаться готовности (Postgres healthcheck). Проверить доступность шлюза:
   ```bash
   curl http://localhost:8080/health
   ```
3. Инициализировать БД и поисковый индекс тестовыми данными:
   ```bash
   chmod +x scripts/*.sh
   ./scripts/seed.sh
   ```

## Проверка API (Happy Path)

Примеры curl-запросов к API Gateway:

Поиск ресторанов (Geo-query):
```bash
curl -s "http://localhost:8080/api/v1/restaurants?lat=55.75&lon=37.62&radius=5000"
```

Меню ресторана:
```bash
curl -s "http://localhost:8080/api/v1/restaurants/00000000-0000-0000-0000-000000000001/menu"
```

Добавление в корзину:
```bash
curl -X POST "http://localhost:8080/api/v1/carts/550e8400-e29b-41d4-a716-446655440000?restaurantId=00000000-0000-0000-0000-000000000001" \
     -H "Content-Type: application/json" \
     -d '{"menu_item_id": "00000000-0000-0000-0000-000000000101", "quantity": 2}'
```

Создание заказа (запускает Saga, возвращает order_id):
```bash
curl -X POST "http://localhost:8080/api/v1/orders" \
     -H "Content-Type: application/json" \
     -d '{
       "customer_id": "550e8400-e29b-41d4-a716-446655440000",
       "restaurant_id": "00000000-0000-0000-0000-000000000001",
       "delivery_address": "Red Square, 1",
       "items": [{"menu_item_id": "00000000-0000-0000-0000-000000000101", "quantity": 2, "price": 300}]
     }'
```

Статус заказа (проверка назначения курьера после прохождения Saga):
```bash
curl -s "http://localhost:8080/api/v1/orders/<ORDER_ID>/status"
```

Автоматизированный прогон флоу:
```bash
./scripts/test_flow.sh
```

## Нагрузочное тестирование

Нагрузочный тест написан на k6 и эмулирует смешанный профиль трафика (Read/Write).
Запускать с отдельной машины, указав IP сервера в переменной `BASE_URL` внутри `loadtest/k6-script.js`.

Запуск:
```bash
k6 run loadtest/k6-script.js
```

Метрики:
- RED-метрики (latency, error rate, RPS) выводятся в консоль k6.
- USE-метрики сервера (CPU/RAM/IO) контролируются через `docker stats`, `htop`, `iostat -dx 1`.

## Паттерны

### Проектирование
1. **API Gateway:** Единая точка входа, скрывающая топологию сети. Включает rate limiting.
   Код: [`api-gateway/main.go`](./api-gateway/main.go)
2. **Transactional Outbox:** Исключает проблему Dual-Write при сохранении заказа и отправке события. Транзакция БД фиксирует заказ и event, после чего воркер доставляет event в Kafka.
   Код: [`order-service/internal/usecase/create_order.go`](./order-service/internal/usecase/create_order.go) (запись) и [`order-service/internal/infrastructure/messaging/outbox_processor.go`](./order-service/internal/infrastructure/messaging/outbox_processor.go) (воркер).
3. **Saga (Choreography/Orchestration):** Асинхронное выполнение транзакции заказа.
   Код: [`delivery-service/main.go`](./delivery-service/main.go) (обработка события оплаты и запуск доставки).

### Устойчивость (Resilience)
1. **Rate Limiting:** Защита от спайков трафика (алгоритм Token Bucket).
   Код: [`api-gateway/main.go`](./api-gateway/main.go) (Middleware `rateLimitMiddleware`).
2. **Idempotency:** Защита от дублей доставки сообщений (At-Least-Once) со стороны Kafka, предотвращение повторного списания средств.
   Код: [`payment-service/internal/usecase/process_payment.go`](./payment-service/internal/usecase/process_payment.go).
3. **Health Check Dependency:** Контейнеры ожидают перехода базы данных в статус healthy (утилита `pg_isready`) перед стартом.
   Код: [`docker-compose.yml`](./docker-compose.yml) (блок `healthcheck` и `depends_on`).

## Журнал оптимизаций

Логирование итераций тестирования и выявления bottlenecks: 
**[docs/optimization-log.md](./docs/optimization-log.md)**
