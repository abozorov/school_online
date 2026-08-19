CREATE TABLE IF NOT EXISTS "users" (
  "id" SERIAL NOT NULL,
  "name" VARCHAR(100) NOT NULL,
  "username" VARCHAR(100) NOT NULL,
  "email" VARCHAR(100) UNIQUE NOT NULL,
  "verify_email" BOOLEAN DEFAULT false,
  "phone_number" VARCHAR(12) UNIQUE,
  "password_hash" TEXT NOT NULL,
  "refresh_token" TEXT,
  "role" VARCHAR(100) CHECK ("role" IN ('user', 'staff', 'teacher', 'student', 'parent')),
  "birthday" VARCHAR(10),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "students" (
  "id" SERIAL NOT NULL,
  "classroom_id" INT,
  "user_id" INT UNIQUE,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users"("id")
);


CREATE TABLE IF NOT EXISTS "parents" (
  "id" SERIAL NOT NULL,
  "students_id" INT[],
  "user_id" INT,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users"("id")
);

-- COMMENT ON COLUMN "parents"."students_id" IS 'INT ARRAY';   

CREATE TABLE IF NOT EXISTS "staffs" (
  "id" SERIAL NOT NULL,
  "position" VARCHAR(100),
  "user_id" INT,
  "experience" INT, -- храним месяцы
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users"("id")
);


CREATE TABLE IF NOT EXISTS "teachers" (
  "id" SERIAL NOT NULL,
  "subjects_id" INT[],
  "user_id" INT,
  "experience" INT, -- храним месяцы
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users"("id")
);

-- COMMENT ON COLUMN "teacher"."subjects_id" IS 'ARRAY items';

CREATE TABLE IF NOT EXISTS "subjects" (
  "id" SERIAL NOT NULL,
  "name" VARCHAR(100),
  "description" VARCHAR(100),
  PRIMARY KEY ("id")
);

