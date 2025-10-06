# Быстрый запуск Defects App

## Запуск бэкенда (Go API)

```bash
cd app/backend
go run cmd/api/main.go
```

API будет доступен на http://localhost:8080

## Запуск фронтенда (Vue 3)

```bash
cd frontend
npm install
npm run dev
```

Фронтенд будет доступен на http://localhost:5173

## Первый запуск

1. Откройте http://localhost:5173
2. Нажмите "Зарегистрироваться"
3. Создайте аккаунт (роль по умолчанию: engineer)
4. Создайте первый проект
5. Добавьте дефекты и управляйте ими

## Переменные окружения

Скопируйте `frontend/env.example` в `frontend/.env` и при необходимости измените:

```
VITE_API_BASE=http://localhost:8080
```

## Структура проекта

```
defects-app/
├── app/
│   ├── backend/          # Go API
│   └── frontend/         # Vue 3 Frontend
├── infra/
│   └── docker-compose.yml
└── README.md
```

## Функциональность

- ✅ Регистрация и авторизация
- ✅ Управление проектами
- ✅ Создание и просмотр дефектов
- ✅ Фильтрация дефектов
- ✅ Комментарии к дефектам
- ✅ Загрузка файлов
- ✅ Смена статусов
- ✅ Скачивание отчетов
