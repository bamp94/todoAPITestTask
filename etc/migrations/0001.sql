CREATE TABLE if NOT EXISTS todo_tasks
(
  id BIGSERIAL NOT NULL PRIMARY KEY,
  task TEXT NOT NULL,
  authorization_token TEXT
);