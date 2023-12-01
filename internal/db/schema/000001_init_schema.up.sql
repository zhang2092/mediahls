CREATE TABLE "users" (
  "id" varchar NOT NULL PRIMARY KEY,
  "username" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "users" ADD CONSTRAINT "username_key" UNIQUE ("username");
ALTER TABLE "users" ADD CONSTRAINT "email_key" UNIQUE ("email");
CREATE INDEX ON "users" ("username");
CREATE INDEX ON "users" ("email");

CREATE TABLE "videos" (
  "id" varchar NOT NULL PRIMARY KEY,
  "title" varchar NOT NULL,
  "description" varchar NOT NULL,
  "images" varchar NOT NULL,
  "origin_link" varchar NOT NULL,
  "play_link" varchar NOT NULL,
  -- 100: 下架
  -- 0: 添加视频
  -- 1: 视频转码中
  -- 2: 视频转码失败
  -- 200: 视频正常显示播放
  "status" int NOT NULL DEFAULT (0),
  "is_deleted" boolean NOT NULL DEFAULT false, -- 删除
  "user_id" varchar NOT NULL,
  "create_at" timestamptz NOT NULL DEFAULT (now()),
  "create_by" varchar NOT NULL,
  "update_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
  "update_by" varchar NOT NULL DEFAULT ('')
);

ALTER TABLE "videos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
CREATE INDEX ON "videos" ("status");
CREATE INDEX ON "videos" ("user_id", "status", "create_at");
CREATE INDEX ON "videos" ("title");
CREATE INDEX ON "videos" ("user_id", "title");