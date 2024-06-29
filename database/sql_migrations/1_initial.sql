-- +migrate Up
CREATE TABLE person
(
  id BIGINT NOT NULL,
  first_name VARCHAR(256),
  last_name VARCHAR(256)
);

-- +migrate Down
DROP TABLE person;

