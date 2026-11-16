-- +goose Up
CREATE TABLE streams (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  is_pinned BOOLEAN NOT NULL DEFAULT FALSE,

  created_at TEXT NOT NULL,
  updated_at TEXT,
  deleted_at TEXT
);

CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  parent_id INTEGER,
  order_num INTEGER NOT NULL,

  name TEXT NOT NULL,
  description TEXT,
  stream_id INTEGER,
  start_datetime TEXT,
  deadline_datetime TEXT,
  timelimit TEXT,
  is_onetime BOOLEAN NOT NULL DEFAULT FALSE,
  repeatable BOOLEAN NOT NULL DEFAULT FALSE,
  finished BOOLEAN NOT NULL DEFAULT FALSE,

  created_at TEXT NOT NULL,
  updated_at TEXT,
  deleted_at TEXT,

  FOREIGN KEY (parent_id) REFERENCES tasks(id) ON DELETE SET NULL,
  FOREIGN KEY (stream_id) REFERENCES streams(id) ON DELETE SET NULL
);

CREATE TABLE timetable (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  for_date TEXT NOT NULL,
  task_id INTEGER NOT NULL UNIQUE,
  start_time TEXT NOT NULL,
  end_time TEXT NOT NULL,

  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE streams;
DROP TABLE tasks;
DROP TABLE timetable;
