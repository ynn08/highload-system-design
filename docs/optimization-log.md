# Журнал оптимизаций (Optimization Log)

Требования к среде: выделенная VM (2 vCPU, 8 GB RAM, 72 GB HDD).

## NFR (Целевые метрики)
- Read-операции: ≥ 100 RPS, p99 < 500ms
- Write-операции: ≥ 30 RPS, p99 < 1s
- Error rate: < 1% при устойчивой нагрузке

## Таблица прогресса

| Метрика | NFR (ДЗ1) | Iter 0 (Baseline) | Iter 1 | Iter 2 | Iter 3 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Latency p99** | < 800ms | 18.22ms | | | |
| **Max RPS** | > 30 | ~16 (20 VUs) | | | |
| **Error rate** | < 1% | 0.00% | | | |
| **CPU / RAM** | 70-90% | ~10% CPU / 1.5 GB | | | |
| **Bottleneck** | — | Пока нет | | | |
| **Что сделали** | — | — | | | |
| **NFR достигнут?** | — | Ожидание 100 VUs | | | |

---

## Итерации

### Iteration 0: Baseline (Initial PoC)
**Сценарий:** Смешанный (Search -> Menu -> Cart -> Order), см. `loadtest/k6-script.js`.
**Настройки среды:** Использованы лимиты ресурсов `deploy.resources.limits` (docker-compose).
**Запуск:** `docker run --rm --network highload-system-design_default -i grafana/k6 run -e BASE_URL=http://gateway:8080/api/v1 - < loadtest/k6-script-short.js`

**Метрики (RED/USE):**
- **Rate:** `15.86 RPS`
- **Errors:** `0.00 %`
- **Duration (p99):** `18.22 ms`
- **Utilization:** Контейнеры потребляют минимум CPU (в пределах 5-10%), RAM в рамках лимитов (Postgres ~100MB, Elastic ~1.2GB, сервисы на Go < 50MB).
- **Saturation:** Очередей на I/O или соединений к БД не наблюдается.
- **Bottleneck Analysis:** Тестовый прогон на 20 VU показал, что система справляется с нагрузкой легко. Узких мест пока нет, так как RPS ниже целевого. Для поиска bottleneck необходимо запустить полноценный скрипт `loadtest/k6-script.js` (до 100 VU).

**План для Iter 1:** 
- Увеличить нагрузку до 100 VUs. Если система начнет деградировать (например, из-за ElasticSearch CPU limits), настроим кэширование результатов поиска в Redis через Catalog Service или увеличим пул соединений в Postgres.

---

### Iteration 1: [Название патча/оптимизации]
**Действия:**
- `TODO`

**Метрики после применения:**
- `TODO`

---

### Iteration 2: [Название патча/оптимизации]
**Действия:**
- `TODO`

**Метрики после применения:**
- `TODO`
