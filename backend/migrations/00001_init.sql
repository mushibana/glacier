-- +goose Up
CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT
);

-- +goose Down
DROP TABLE tasks;
