-- Роли
INSERT INTO roles(name, description) VALUES
  ('manager','project manager'),
  ('engineer','developer/QA'),
  ('viewer','readonly user'),
  ('lead','tech lead')
ON CONFLICT (name) DO NOTHING;

-- Статусы
INSERT INTO statuses(name, description) VALUES
  ('new','newly created'),
  ('in_work','in progress'),
  ('review','in review'),
  ('closed','closed'),
  ('canceled','canceled')
ON CONFLICT (name) DO NOTHING;

-- Демопроект
INSERT INTO projects(name, customer)
VALUES ('Demo project','Acme Co')
ON CONFLICT DO NOTHING;
