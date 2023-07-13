ALTER TABLE "sessions" DROP CONSTRAINT sessions_user_id_fkey;

DROP TABLE IF EXISTS sessions;

DROP TABLE IF EXISTS users;

SET TIME ZONE "UTC";