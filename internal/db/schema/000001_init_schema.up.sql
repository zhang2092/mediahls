CREATE TABLE "users" (
  "id" varchar NOT NULL PRIMARY KEY,
  "username" varchar NOT NULL UNIQUE,
  "hashed_password" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);