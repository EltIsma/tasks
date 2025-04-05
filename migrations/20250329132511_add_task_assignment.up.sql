BEGIN;

CREATE TABLE IF NOT EXISTS task(
    id uuid PRIMARY KEY,
    payload TEXT NOT NULL,
    deadline timestamptz
);

CREATE TABLE IF NOT EXISTS assignment(
   id uuid PRIMARY KEY,
   class TEXT NOT NULL,
   task_id uuid NOT NULL,
   lesson_id uuid NOT NULL,
   task_payload TEXT NOT NULL,
   deadline timestamptz,

   FOREIGN KEY (task_id) REFERENCES task (id) ON DELETE CASCADE,
   UNIQUE(lesson_id, task_id),
   UNIQUE(class, lesson_id, task_id)
);

CREATE TABLE IF NOT EXISTS usersMark(
   id uuid PRIMARY KEY,
   user_id uuid NOT NULL,
   task_id uuid NOT NULL,
   lesson_id uuid NOT NULL,
   mark int,

   FOREIGN KEY (task_id) REFERENCES assignment (id) ON DELETE CASCADE,
   UNIQUE(user_id, task_id, lesson_id)
);

CREATE INDEX assignment_task_id_user_id_idx on assignment (lesson_id, task_id);

END;