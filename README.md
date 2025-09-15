# 📘 Система учёта дефектов

Проект на **Go (Gin)** и **Vue 3 (Vite, TypeScript)** для управления проектами, дефектами, комментариями и вложениями.

---

## Используемые технологии

### Backend (Go + Gin)
- Go (Golang)
- Gin Web Framework
- PostgreSQL (хранение данных)

### Frontend (Vue.js)
- Vue 3
- Vite
- TypeScript
- Vue Router

### База данных
- PostgreSQL  
- Диаграмма БД: https://dbdiagram.io/d/68c886e21ff9c616bdcdd20d

---

## Требования

- **Go** 1.22+  
- **Node.js**: либо **20.19+**, либо **22.12+ и новее** (Vite этого требует)  
- **npm** (идёт вместе с Node.js)
- (Опционально) установленная локально PostgreSQL — если вы уже подключаете БД

> Если при запуске Vite увидите предупреждение про версию Node — обновите Node до 22.12+.

---

## Запуск локально

### 1. Клонирование репозитория
```bash
git clone https://github.com/Alpha-x-prog/framework_task1.git
cd framework_task1
```

### 2. Запуск Backend (Go + Gin)

Перейдите в папку **backend**, установите зависимости и запустите сервер.

```bash
cd backend
go mod tidy
go run ./cmd/api
```

### 3. Запуск Frontend (Vue)

Откройте новый терминал, перейдите в папку **frontend**, установите зависимости и запустите дев-сервер.

```bash
cd frontend
npm install
npm run dev
```
