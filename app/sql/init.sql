-- =========================================
--  Postgres schema for Defects System
--  Tables: roles, statuses, users, projects,
--          defects, comments, attachments
-- =========================================

-- Роли пользователей
CREATE TABLE IF NOT EXISTS roles (
  id   SERIAL PRIMARY KEY,
  name VARCHAR(64) UNIQUE NOT NULL,   -- manager | engineer | viewer | lead
  description VARCHAR(255)
);

-- Статусы дефектов
CREATE TABLE IF NOT EXISTS statuses (
  id   SERIAL PRIMARY KEY,
  name VARCHAR(64) UNIQUE NOT NULL,   -- new | in_work | review | closed | canceled
  description VARCHAR(255)
);

-- Пользователи
CREATE TABLE IF NOT EXISTS users (
  id            SERIAL PRIMARY KEY,
  email         TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role_id       INT  NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Проекты
CREATE TABLE IF NOT EXISTS projects (
  id         SERIAL PRIMARY KEY,
  name       TEXT NOT NULL,
  customer   TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Дефекты
CREATE TABLE IF NOT EXISTS defects (
  id          SERIAL PRIMARY KEY,
  project_id  INT NOT NULL REFERENCES projects(id) ON DELETE RESTRICT,
  title       TEXT NOT NULL,
  description TEXT,
  priority    INT NOT NULL DEFAULT 3,   -- 1..5 (1=high)
  assignee_id INT REFERENCES users(id) ON DELETE SET NULL,
  status_id   INT NOT NULL REFERENCES statuses(id) ON DELETE RESTRICT,
  due_date    DATE,
  created_by  INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Комментарии
CREATE TABLE IF NOT EXISTS comments (
  id         SERIAL PRIMARY KEY,
  defect_id  INT NOT NULL REFERENCES defects(id) ON DELETE CASCADE,
  author_id  INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  body       TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Вложения
CREATE TABLE IF NOT EXISTS attachments (
  id          SERIAL PRIMARY KEY,
  defect_id   INT NOT NULL REFERENCES defects(id) ON DELETE CASCADE,
  file_path   TEXT NOT NULL,   -- путь/ключ файла
  mime        TEXT,
  uploaded_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- -----------------
-- Индексы под фильтры/связи
-- -----------------
CREATE INDEX IF NOT EXISTS idx_defects_project   ON defects(project_id);
CREATE INDEX IF NOT EXISTS idx_defects_status    ON defects(status_id);
CREATE INDEX IF NOT EXISTS idx_defects_assignee  ON defects(assignee_id);
CREATE INDEX IF NOT EXISTS idx_defects_priority  ON defects(priority);
CREATE INDEX IF NOT EXISTS idx_defects_due_date  ON defects(due_date);
CREATE INDEX IF NOT EXISTS idx_comments_defect   ON comments(defect_id);
CREATE INDEX IF NOT EXISTS idx_attachments_defect ON attachments(defect_id);

-- Композитный индекс для типовых списков
CREATE INDEX IF NOT EXISTS idx_defects_filters
  ON defects(project_id, status_id, priority, assignee_id, due_date);

-- ---------------
-- Триггер на updated_at
-- ---------------
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'set_updated_at') THEN
    CREATE FUNCTION set_updated_at() RETURNS trigger AS $f$
    BEGIN
      NEW.updated_at = now();
      RETURN NEW;
    END;
    $f$ LANGUAGE plpgsql;
  END IF;
END$$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_defects_updated') THEN
    CREATE TRIGGER trg_defects_updated
    BEFORE UPDATE ON defects
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  END IF;
END$$;
