CREATE TABLE IF NOT EXISTS "journals" (
  "id" SERIAL NOT NULL,
  "classroom_id" INT,
  "subject_id" INT,
  "teacher_id" INT,
  "date" DATE,
  "lesson_number" INT,
  "room" INT,
  "attendence" BOOLEAN DEFAULT true,
  "student_id" INT,
  "grade" INT,
  "homework" TEXT,
  PRIMARY KEY ("id")
);