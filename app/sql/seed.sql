-- ========================
--  Начальные данные
-- ========================

-- Роли
INSERT INTO roles (name, description) VALUES
('manager',  'Руководитель/менеджер проектов'),
('engineer', 'Инженер-исполнитель'),
('viewer',   'Только просмотр'),
('lead',     'Руководитель/наблюдатель')
ON CONFLICT (name) DO NOTHING;

-- Статусы
INSERT INTO statuses (name, description) VALUES
('new',      'Новый дефект'),
('in_work',  'В работе'),
('review',   'На проверке'),
('closed',   'Закрыт'),
('canceled', 'Отменён')
ON CONFLICT (name) DO NOTHING;

-- Пользователи (примерные данные)
-- bcrypt-хеш пароля "password123":
-- $2b$12$KIXQ4S6W7YQk4e6x4blJeOkGDZsBoBPXaGNsq4CPGK/c/pX9nuS3u
INSERT INTO users (email, password_hash, role_id)
VALUES
('manager@example.com',
 '$2b$12$KIXQ4S6W7YQk4e6x4blJeOkGDZsBoBPXaGNsq4CPGK/c/pX9nuS3u',
 (SELECT id FROM roles WHERE name='manager')),
('engineer@example.com',
 '$2b$12$KIXQ4S6W7YQk4e6x4blJeOkGDZsBoBPXaGNsq4CPGK/c/pX9nuS3u',
 (SELECT id FROM roles WHERE name='engineer')),
('viewer@example.com',
 '$2b$12$KIXQ4S6W7YQk4e6x4blJeOkGDZsBoBPXaGNsq4CPGK/c/pX9nuS3u',
 (SELECT id FROM roles WHERE name='viewer'))
ON CONFLICT (email) DO NOTHING;

-- Проекты
INSERT INTO projects (name, customer)
VALUES
('Demo Project', 'Demo Customer')
ON CONFLICT DO NOTHING;

-- Дефекты (используем подзапросы, чтобы взять id статусов/ролей)
INSERT INTO defects (project_id, title, description, priority, assignee_id, status_id, due_date, created_by)
VALUES
(
  (SELECT id FROM projects WHERE name='Demo Project'),
  'Ошибка входа',
  'При неверном пароле не показывается сообщение об ошибке',
  1,
  (SELECT id FROM users WHERE email='engineer@example.com'),
  (SELECT id FROM statuses WHERE name='in_work'),
  CURRENT_DATE + INTERVAL '7 days',
  (SELECT id FROM users WHERE email='manager@example.com')
),
(
  (SELECT id FROM projects WHERE name='Demo Project'),
  'Кнопка уезжает на мобильном',
  'Неверная вёрстка на экране авторизации (iPhone 13, Safari)',
  3,
  NULL,
  (SELECT id FROM statuses WHERE name='new'),
  NULL,
  (SELECT id FROM users WHERE email='manager@example.com')
);

-- Комментарии
INSERT INTO comments (defect_id, author_id, body)
VALUES
(
  (SELECT id FROM defects WHERE title='Ошибка входа' LIMIT 1),
  (SELECT id FROM users WHERE email='engineer@example.com'),
  'Воспроизвёл баг, начал фиксить'
),
(
  (SELECT id FROM defects WHERE title='Ошибка входа' LIMIT 1),
  (SELECT id FROM users WHERE email='manager@example.com'),
  'Ок, жду обновление'
);

-- Пример вложения
INSERT INTO attachments (defect_id, file_path, mime, uploaded_by)
VALUES
(
  (SELECT id FROM defects WHERE title='Ошибка входа' LIMIT 1),
  'uploads/login-error.png',
  'image/png',
  (SELECT id FROM users WHERE email='engineer@example.com')
);
