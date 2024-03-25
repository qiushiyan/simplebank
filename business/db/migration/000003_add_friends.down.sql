ALTER TABLE IF EXISTS "friendships" DROP CONSTRAINT IF EXISTS "friendships_from_account_id_fkey";
ALTER TABLE IF EXISTS "friendships" DROP CONSTRAINT IF EXISTS "friendships_to_account_id_fkey";
DROP TABLE IF EXISTS "friendships";