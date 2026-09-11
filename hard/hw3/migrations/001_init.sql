-- +migrate Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    result TEXT NOT NULL DEFAULT '',
    translator TEXT NOT NULL,
    code TEXT NOT NULL
);

-- +migrate Down

DROP TABLE tasks;
DROP TABLE users;
