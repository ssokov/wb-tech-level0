# Simple Order Service

Микросервис для управления заказами с кэшированием в Redis и асинхронной обработкой через Kafka.

## Технологии

- **Go 1.24.2** - основной язык
- **PostgreSQL** - основная БД
- **Redis** - кэш
- **Kafka** - message broker
- **Docker** - контейнеризация

## Архитектура

Сервис получает заказы через Kafka, сохраняет в PostgreSQL и кэширует в Redis. HTTP API предоставляет быстрый доступ к данным через кэш.

## Структура

```
wb/
├── cmd/main.go              # Точка входа
├── internal/
│   ├── domain/              # Модели данных
│   ├── http/                # HTTP обработчики
│   ├── repo/                # Репозитории
│   └── service/             # Бизнес-логика
├── pkg/db/                  # Подключение к БД
└── web/static/              # Веб-интерфейс
```

## Запуск

```bash
# Запуск всех сервисов
docker-compose up -d

# Проверка работы
curl http://localhost:8081/ping
```

## API

- `GET /ping` - проверка работоспособности
- `GET /order?id={id}` - получение заказа по ID

## Порты

- 8081 - HTTP API
- 5432 - PostgreSQL  
- 6379 - Redis
- 9092 - Kafka
- 9000 - Kafdrop UI