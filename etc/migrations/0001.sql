CREATE TABLE if NOT EXISTS todo_lists
(
    id                  BIGSERIAL NOT NULL PRIMARY KEY,
    name                VARCHAR   NOT NULL,
    authorization_token VARCHAR   NOT NULL
);

CREATE TABLE if NOT EXISTS todo_tasks
(
    id           BIGSERIAL NOT NULL PRIMARY KEY,
    todo_list_id INT       NOT NULL REFERENCES todo_lists,
    task         TEXT      NOT NULL
);

INSERT INTO todo_lists (name, authorization_token)
VALUES ('Daily tasks', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9'),
       ('My tasks', 'zdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ikp'),
       ('Our tasks', 'iOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4');