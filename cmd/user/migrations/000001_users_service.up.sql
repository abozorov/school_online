CREATE TABLE IF NOT EXISTS "users" (
  "id" SERIAL NOT NULL,
  "name" VARCHAR(100),
  "username" VARCHAR(100),
  "email" VARCHAR(100) UNIQUE,
  "verify_emmail" BOOLEAN DEFAULT false,
  "phone_number" VARCHAR(12) UNIQUE,
  "password_hash" TEXT,
  "refresh_token" TEXT,
  "role" VARCHAR(100),
  "birthday" TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "students" (
  "id" SERIAL NOT NULL,
  "classroom_id" INT,
  "user_id" INT UNIQUE,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS "parrents" (
  "id" SERIAL NOT NULL,
  "students_id" INT[],
  "user_id" INT,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON UPDATE CASCADE
);

-- COMMENT ON COLUMN "parrents"."students_id" IS 'INT ARRAY';   

CREATE TABLE IF NOT EXISTS "stuffs" (
  "id" SERIAL NOT NULL,
  "position" VARCHAR(100),
  "user_id" INT,
  "experience" INT, -- храним месяцы
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS "teachers" (
  "id" SERIAL NOT NULL,
  "subjects_id" INT[],
  "user_id" INT,
  "experience" INT, -- храним месяцы
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON UPDATE CASCADE
);

-- COMMENT ON COLUMN "teacher"."subjects_id" IS 'ARRAY items';

CREATE TABLE IF NOT EXISTS "subjects" (
  "id" SERIAL NOT NULL,
  "name" VARCHAR(100),
  "description" VARCHAR(100),
  PRIMARY KEY ("id")
);

