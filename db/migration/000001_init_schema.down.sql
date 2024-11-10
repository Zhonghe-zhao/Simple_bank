ALTER TABLE transfers DROP CONSTRAINT IF EXISTS "transfers_from_account_id_fkey";

ALTER TABLE transfers DROP CONSTRAINT IF EXISTS "transfers_to_account_id_fkey";

DROP TABLE IF EXISTS "entries";

DROP TABLE IF EXISTS "accounts";

DROP TABLE IF EXISTS "transfers";