CREATE TABLE "users" (
    "username" varchar PRIMARY KEY,
    "email" VARCHAR UNIQUE NOT NULL,
    "hashed_password" varchar NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "accounts"
ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts"
ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
-- create admin user
INSERT INTO "users" ("username", "email", "hashed_password")
VALUES (
        'admin',
        'test@gmail.com',
        '$2a$10$SZw5eHfb.e94cBVUbrZ0BOTmXOFtINl7ZOv9Ac9WbYJdOnNZ9Gk6a'
    );