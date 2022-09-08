CREATE TABLE IF NOT EXISTS "user" (
  id            uuid                                 DEFAULT uuid_generate_v4(),
  username      VARCHAR(25)                 NOT NULL,
  password      VARCHAR(60)                 NOT NULL,
  created_at    timestamp WITHOUT TIME ZONE NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  updated_at    timestamp WITHOUT TIME ZONE NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  deleted_at    timestamp WITHOUT TIME ZONE     NULL,
  PRIMARY KEY (id)
);