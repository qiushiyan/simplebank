CREATE TABLE "friendships" (
    "id" bigserial PRIMARY KEY,
    "from_account_id" BIGINT NOT NULL,
    "to_account_id" BIGINT NOT NULL,
    "pending" BOOLEAN NOT NULL DEFAULT TRUE,
    "accepted" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE ("from_account_id", "to_account_id")
);
ALTER TABLE "friendships"
ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "friendships"
ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
-- connect admin and user
INSERT INTO "friendships" (
        "from_account_id",
        "to_account_id",
        "pending",
        "accepted"
    )
VALUES (1, 3, FALSE, TRUE);