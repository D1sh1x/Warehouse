CREATE TABLE IF NOT EXISTS users (
  id        SERIAL PRIMARY KEY,
  username  TEXT NOT NULL UNIQUE,
  password  TEXT NOT NULL,
  role      TEXT NOT NULL CHECK (role IN ('admin','manager','viewer'))
);

CREATE TABLE IF NOT EXISTS items (
  id     SERIAL PRIMARY KEY,
  name   TEXT NOT NULL,
  count  INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS history (
  id         BIGSERIAL PRIMARY KEY,
  item_id    INTEGER,
  action     TEXT NOT NULL CHECK (action IN ('insert','update','delete')),
  changed_by TEXT,
  "timestamp" TIMESTAMPTZ NOT NULL DEFAULT now(),
  old_data   JSONB,
  new_data   JSONB,
  CONSTRAINT fk_history_item
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL
);

INSERT INTO users (username, password, role) VALUES
  ('alice','alice','admin'),
  ('bob','bob','manager'),
  ('eve','eve','viewer')
ON CONFLICT (username) DO NOTHING;
