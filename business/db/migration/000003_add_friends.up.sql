CREATE TABLE "friendships" (
    "id" bigserial PRIMARY KEY,
    "from_account_id" BIGINT NOT NULL,
    "to_account_id" BIGINT NOT NULL,
    "status" VARCHAR NOT NULL DEFAULT 'pending',
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
        "status"
    )
VALUES (1, 3, 'accepted');
INSERT INTO "friendships" (
        "from_account_id",
        "to_account_id",
        "status"
    )
VALUES (4, 2, 'pending');