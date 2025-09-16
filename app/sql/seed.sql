INSERT INTO users (email, password_hash, role) VALUES
('manager@example.com',  'hash1', 'manager'),
('engineer@example.com', 'hash2', 'engineer'),
('viewer@example.com',   'hash3', 'viewer');

INSERT INTO projects (name, customer) VALUES
('Demo Project', 'Customer A');

INSERT INTO defects (project_id, title, description, priority, status, created_by)
VALUES
(1, 'Ошибка входа', 'Не работает вход при неверном пароле', 1, 'new', 1);

INSERT INTO comments (defect_id, author_id, body) VALUES
(1, 2, 'Воспроизвёл баг, буду чинить');
