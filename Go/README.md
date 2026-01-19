# Simple Go API Example

Минимальное API на Go, содержащее:

- `/health` — проверка состояния приложения и подключения к базе данных  
- `/hello` — вывод `"Hello, world!"` в JSON  
- `GET /user/{id}` — получение пользователя по ID  
- `POST /user` — создание пользователя  
- Миграция для создания таблицы `users`

## 📦 Требования

- Go 1.20+
- PostgreSQL
- Git (опционально)

## 🛠 Переменные окружения

Перед запуском необходимо задать параметры подключения к базе данных:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=testdb
export DB_SSLMODE=disable
```

## 🗄 Миграции

В проекте есть миграция:

```
migrations/001_create_users.sql
```

Содержимое:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL
);
```

Миграции можно применить вручную:

```bash
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/001_create_users.sql
```

## 🚀 Запуск приложения

```bash
go run .
```

После запуска сервер будет доступен на:

```
http://localhost:8080
```

## 📚 Примеры запросов

### Healthcheck
```
GET /health
```

### Hello
```
GET /hello
```

### Получить пользователя
```
GET /user/1
```

### Создать пользователя
```
POST /user
Content-Type: application/json

{
  "name": "Alex",
  "email": "alex@example.com"
}
```

## 🧩 Зависимости

Используются:

- `github.com/gorilla/mux`
- `github.com/lib/pq`

Установка:

```bash
go mod tidy
```