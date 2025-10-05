-- Снять триггер и функцию
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_defects_updated') THEN
    DROP TRIGGER trg_defects_updated ON defects;
  END IF;
END$$;

DROP FUNCTION IF EXISTS set_updated_at();

-- Удаляем таблицы в обратном порядке зависимостей
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS defects;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS statuses;
DROP TABLE IF EXISTS roles;
