DEPLOY.md
Что это за сервис

Небольшой Go API-сервис с несколькими endpoint’ами. Хранит данные в PostgreSQL.

Требования

Ubuntu Server (VM), без GUI
Docker + Docker Compose v2 (docker compose)
PostgreSQL установлен на хосте VM, не в контейнере
GitHub Actions + self-hosted runner на этой VM

1) Установка PostgreSQL на хосте VM
sudo apt update
sudo apt install -y postgresql postgresql-contrib


Запуск/проверка кластера:
sudo systemctl status postgresql@16-main --no-pager || true
pg_lsclusters

2) Создание БД и пользователя

Создаём пользователя и базу под сервис:

sudo -u postgres psql <<'SQL'
CREATE USER api WITH PASSWORD 'superpass';
CREATE DATABASE api OWNER api;
GRANT ALL PRIVILEGES ON DATABASE api TO api;
SQL

3) Настройка сети Postgres для доступа из контейнера
3.1 Разрешить слушать сеть

Файл: /etc/postgresql/16/main/postgresql.conf

Параметр:
listen_addresses = '*'

3.2 Разрешить подключения из docker-сетей

Файл: /etc/postgresql/16/main/pg_hba.conf

Добавить в конец:

host    api     api     172.17.0.0/16    scram-sha-256
host    api     api     172.20.0.0/16    scram-sha-256


Перезапуск:
sudo systemctl restart postgresql@16-main


Проверка, что Postgres слушает 5432:
ss -lntp | grep 5432

4) Миграция (создание таблицы users)

Применяем SQL (пример миграции):

cat > /tmp/001_create_users.sql <<'SQL'
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);
SQL

psql "host=127.0.0.1 port=5432 user=api password=superpass dbname=api sslmode=disable" -f /tmp/001_create_users.sql
psql "host=127.0.0.1 port=5432 user=api password=superpass dbname=api sslmode=disable" -c "\dt"

5) Dockerfile (multistage)

Сборка образа происходит multistage Dockerfile:
stage build: собираем Go бинарник
stage runtime: минимальный runtime образ и запуск приложения

6) Запуск сервиса на VM (Docker Compose)

Папка деплоя на VM: ~/ose-app
docker-compose.yml использует подключение контейнера к Postgres на хосте VM через host-gateway:
extra_hosts: host.docker.internal:host-gateway
DB_HOST=host.docker.internal

Запуск:

cd ~/ose-app
docker compose up -d



curl -s http://localhost:8080/health
curl -s http://localhost:8080/hello







7) CI/CD (GitHub Actions)

Используется GitHub Actions workflow .github/workflows/ci-cd.yml:

7.1 Lint job (cloud runner)

gofmt (проверка форматирования)

go vet

7.2 Build & Deploy job (self-hosted runner на VM)

сборка docker image на VM: docker build -t ose-api:latest ./Go

деплой через docker compose up -d --force-recreate

health-check через curl -f /health

Почему self-hosted runner

VM находится локально в VirtualBox и недоступна из интернета напрямую, поэтому деплой сделан через self-hosted runner, который выполняет шаги прямо на VM.

8) Сетевое решение “контейнер → Postgres на хосте”

Postgres работает на хосте VM и слушает порт 5432

Контейнер обращается к хосту по имени host.docker.internal

Для Linux добавлено: extra_hosts: host.docker.internal:host-gateway
Доступ разрешён в pg_hba.conf для docker-подсетей