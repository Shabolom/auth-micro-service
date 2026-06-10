# Swagger documentation

Документация Swagger/OpenAPI для `notification_service`.

В комплект входят:

- `openapi.yaml` — OpenAPI 3.0.3 спецификация;
- `swagger.html` — HTML-страница для просмотра Swagger UI в браузере;
- `README.md` — инструкция по запуску.

## Назначение

Так как `notification_service` является event-driven сервисом, основная логика работает не через HTTP API, а через:

```text
RabbitMQ -> notification_service -> Resend
RabbitMQ -> notification_service -> Kafka
```

Swagger здесь используется как документация для:

- структуры auth events;
- формата сообщений;
- health/readiness endpoints;
- описания payload, который сервис читает из RabbitMQ и пишет в Kafka.

## Структура файлов

```text
swagger/
├── openapi.yaml
├── swagger.html
└── README.md
```

## Как открыть Swagger в браузере

Перейди в папку со Swagger-файлами:

```bash
cd swagger
```

Запусти простой локальный HTTP-сервер:

```bash
python -m http.server 8088
```

Открой в браузере:

```text
http://localhost:8088/swagger.html
```

## Почему лучше запускать через HTTP-сервер

Не стоит просто открывать `swagger.html` двойным кликом как файл:

```text
file:///...
```

Некоторые браузеры блокируют загрузку локального `openapi.yaml` из-за CORS/security restrictions.

Поэтому надежнее запускать через:

```bash
python -m http.server 8088
```

## Как положить в проект

Рекомендуемый вариант:

```text
notification_service/
├── cmd/
├── internal/
├── docs/
│   └── swagger/
│       ├── openapi.yaml
│       ├── swagger.html
│       └── README.md
├── build/
└── README.md
```

То есть файлы можно положить сюда:

```text
docs/swagger/
```

## Как открыть после добавления в проект

```bash
cd docs/swagger
python -m http.server 8088
```

Затем:

```text
http://localhost:8088/swagger.html
```

## Что описано в OpenAPI

### Health endpoints

```http
GET /health
```

Проверяет, что процесс сервиса жив.

```http
GET /ready
```

Проверяет, что сервис готов обрабатывать сообщения.

### Event schema

```http
GET /docs/events/auth
```

Документационный endpoint с примером auth event.

## Основная схема события

```go
type Event struct {
    Time  time.Time `json:"time"`
    Email string    `json:"email"`
    Type  string    `json:"type"`
}
```

Пример JSON:

```json
{
  "time": "2026-05-16T12:00:00Z",
  "email": "user@example.com",
  "type": "register"
}
```

## Типы событий

Поддерживаются два типа:

```text
register
login
```

## RabbitMQ

Документация описывает входящие сообщения RabbitMQ.

Ожидаемый content type:

```text
text/plain
```

Тело сообщения:

```text
user@example.com
```

Очереди:

```text
auth.register
auth.login.logs
```

Exchange:

```text
auth.events
```

## Kafka

После обработки события сервис публикует сообщение в Kafka.

Пример Kafka payload:

```json
{
  "time": "2026-05-16T12:00:00Z",
  "email": "user@example.com",
  "type": "login"
}
```

## Важное замечание

Эта Swagger-документация описывает контракт и формат данных.

Она не означает, что все endpoints уже физически реализованы в сервисе.  
Если в коде нет HTTP-сервера, то `/health`, `/ready` и `/docs/events/auth` нужно добавить отдельно.

## Возможное подключение Swagger через Go HTTP server

Если в сервис будет добавлен HTTP server, можно раздавать Swagger как static files.

Пример для Echo:

```go
e.Static("/swagger", "docs/swagger")
```

После этого Swagger будет доступен по адресу:

```text
http://localhost:8080/swagger/swagger.html
```

## Возможные улучшения

- добавить реальные `/health` и `/ready` endpoints;
- добавить проверку RabbitMQ connection в `/ready`;
- добавить проверку Kafka producer в `/ready`;
- добавить проверку доступности Resend;
- добавить Swagger UI в основной HTTP-сервер сервиса;
- добавить JSON Schema файл отдельно для Kafka Schema Registry;
- добавить примеры ошибок и retry/DLQ сценариев.
