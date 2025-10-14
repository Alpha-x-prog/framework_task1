# 📘 Defects App — система учёта дефектов

**Defects App** — полнофункциональное приложение для управления проектами и дефектами: создание задач/дефектов, комментарии, вложения, статусы, приоритеты и отчёты (JSON/CSV). Бэкенд на **Go (Gin) + PostgreSQL**, фронтенд на **Vue 3 (Vite)**.

---

## Возможности

- ✅ Регистрация и авторизация (JWT)
- ✅ Роли и разрешения
- ✅ Управление проектами
- ✅ Создание, просмотр и фильтрация дефектов
- ✅ Комментарии к дефектам
- ✅ Загрузка и хранение вложений
- ✅ Смена статусов и приоритетов
- ✅ Отчёты: сводка и тренды (JSON) и выгрузка CSV

---

## Технологии

- **Backend**: Go, Gin, PostgreSQL, pgx, golang-migrate
- **Frontend**: Vue 3, Vite 5, Pinia, Vue Router, Axios, Chart.js
- **Инфраструктура**: Docker Compose (опционально)

---

## Структура проекта

```text
defects-app/
├── app/
│   ├── backend/              # Go API (Gin)
│   │   ├── cmd/api/main.go   # вход в приложение
│   │   ├── internal/         # доменная логика, http, repo, миграции
│   │   └── uploads/          # файлы-вложения (локальное хранилище)
│   └── sql/                  # init/seed SQL (если применяете вручную)
├── frontend/                 # Vue 3 (Vite)
├── infra/
│   └── docker-compose.yml    # локальный запуск в контейнерах
└── QUICK_START.md            # короткая инструкция
```

---

## Требования

- Go 1.22+ (в `go.mod` указан 1.24 — подойдёт 1.22 и выше)
- Node.js 20 LTS или выше (Vite 5 поддерживает ^18.0.0 || >=20.0.0; рекомендуется 20+)
- npm (идёт вместе с Node.js)
- PostgreSQL (локально или в Docker)

> Если Vite предупредит о версии Node — обновите Node до актуальной LTS (20+).

---

## Быстрый старт (локально)

### 1) Бэкенд (API)

```bash
cd app/backend
go run cmd/api/main.go
```

По умолчанию API поднимется на `http://localhost:8080`.

### 2) Фронтенд (Vue 3)

```bash
cd frontend
npm install
npm run dev
```

Фронтенд будет доступен на `http://localhost:5173`.

### 3) Переменные окружения (фронтенд)

Скопируйте `frontend/env.example` в `frontend/.env` и при необходимости измените:

```dotenv
VITE_API_BASE=http://localhost:8080
```

### 4) Первый запуск в UI

1. Откройте `http://localhost:5173`
2. Зарегистрируйтесь (по умолчанию роль: engineer)
3. Создайте проект
4. Добавьте дефекты, комментируйте, меняйте статусы и приоритеты

---

## Запуск через Docker Compose (опционально)

Если удобнее работать в контейнерах, используйте `infra/docker-compose.yml`. Файл включает сервисы API, БД и (опционально) фронтенда.

Пример команд (из корня репозитория):

```bash
cd infra
# убедитесь, что переменные окружения для БД и API настроены
docker compose up -d --build
```

После запуска проверьте доступность:
- API: `http://localhost:8080`
- Frontend (если вынесен в compose): `http://localhost:5173`

---

## Миграции и данные

В репозитории предусмотрены два подхода:
- Миграции Go: `app/backend/internal/migrate/migrations/*` — применяются программно
- SQL-скрипты: `app/sql/init.sql`, `app/sql/seed.sql` — можно применить вручную

Выберите удобный способ для инициализации схемы и базовых данных.

---

## Отчёты

Основные эндпоинты:
- `GET /api/reports/summary` — сводка (фильтры: `project_id`, `from`, `to`)
- `GET /api/reports/trends` — тренды по дням/неделям/месяцам (`group=day|week|month` и те же фильтры)
- `GET /api/reports/summary.csv` — CSV-выгрузка (те же фильтры)

Даты передавайте в формате `YYYY-MM-DD`.

---

## Полезные пути в коде

- Хендлеры HTTP: `app/backend/internal/http/handlers/*`
- Роли/права: `app/backend/internal/http/mv/*`
- Доменная логика: `app/backend/internal/core/*`
- Репозитории (доступ к БД): `app/backend/internal/repo/*`
- JWT/аутентификация: `app/backend/internal/auth/*`

---

## Скрипты (Frontend)

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  }
}
```

Сборка production:

```bash
cd frontend
npm run build
npm run preview
```

---

## Вложения (uploads)

Загруженные файлы сохраняются в `app/backend/uploads/` (по умолчанию локально). Убедитесь, что у процесса есть права на запись.

---

## Лицензия

MIT (или укажите вашу актуальную лицензию).

---

## FAQ

- Почему в README не указан точный URL репозитория? — Укажите ваш реальный Git URL при публикации.
- TypeScript на фронтенде обязателен? — Сейчас фронтенд использует JavaScript с Vite и Vue 3. TypeScript можно добавить позже.
