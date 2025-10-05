-- Аккуратный откат сидов (пройдёт на чистой БД; если есть данные — возможны ограничения FK)
DELETE FROM projects WHERE name = 'Demo project';
DELETE FROM statuses WHERE name IN ('new','in_work','review','closed','canceled');
DELETE FROM roles    WHERE name IN ('manager','engineer','viewer','lead');
