CREATE TABLE "users" (
    "username" varchar PRIMARY KEY,
    "email" VARCHAR UNIQUE,
    "hashed_password" varchar NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "accounts"
ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts"
ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
-- create two users, one admin and one normal user
INSERT INTO "users" ("username", "email", "hashed_password")
VALUES (
        'admin',
        'admin@gmail.com',
        '$2a$10$SZw5eHfb.e94cBVUbrZ0BOTmXOFtINl7ZOv9Ac9WbYJdOnNZ9Gk6a'
    ),
    (
        'user',
        'user@gmail.com',
        '$2a$10$SZw5eHfb.e94cBVUbrZ0BOTmXOFtINl7ZOv9Ac9WbYJdOnNZ9Gk6a'
    );
-- create accounts for the two users
INSERT INTO "accounts" ("owner", "name", "balance", "currency")
VALUES ('admin', 'Investment 1', 100, 'USD'),
    ('admin', 'Investment 2', 100, 'EUR'),
    ('user', 'Personal', 100, 'USD'),
    ('user', 'Family', 100, 'EUR');