
-- HSTORE
CREATE EXTENSION IF NOT EXISTS hstore;
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- AMPQ
-- CREATE EXTENSION IF NOT EXISTS amqp;
-- --  добавляем уникальный индекс чтобы при импорте брокер не дублировался
-- CREATE UNIQUE INDEX IF NOT EXISTS host_index on amqp.broker (host);
-- --  создаем брокера
-- INSERT INTO amqp.broker (host, username, password)
-- VALUES ('osnovi-finansov.ru', 'osnovi_finansov', 'ktulhu77')
-- ON CONFLICT (host)
--   DO NOTHING;
