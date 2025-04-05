package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"task/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryPG struct {
	conn *pgxpool.Pool
}

func NewRepositoruPG(conn *pgxpool.Pool) *RepositoryPG {
	return &RepositoryPG{
		conn: conn,
	}
}

func (pg *RepositoryPG) CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {
	var id uuid.UUID
	err := pg.conn.QueryRow(ctx, "INSERT INTO task (id, payload, deadline) VALUES($1, $2, $3) RETURNING id", task.ID, task.Payload, task.Deadline).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("can't create new task records:%w", err)
	}

	return id, nil
}

func (pg *RepositoryPG) GetTaskByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var task domain.Task
	err := pg.conn.QueryRow(ctx, "SELECT  id, payload, deadline FROM task WHERE id = $1", id).Scan(&task.ID, &task.Payload, &task.Deadline)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (pg *RepositoryPG) GetTasks(ctx context.Context) ([]*domain.Task, error) {
	rows, err := pg.conn.Query(ctx, "SELECT id, payload, deadline FROM task")
	if err != nil {
		return nil, fmt.Errorf("error executing prepared statement: %w", err)
	}

	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID,
			&task.Payload,
			&task.Deadline,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning task row: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows: %w", err)
	}

	return tasks, nil
}

func (pg *RepositoryPG) UpdateTask(ctx context.Context, task *domain.Task) error {
	_, err := pg.conn.Exec(ctx, "UPDATE task SET payload = $1, deadline = $2 WHERE id = $3", task.Payload, task.Deadline, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (pg *RepositoryPG) DeleteTask(ctx context.Context, id uuid.UUID) error {

	tag, err := pg.conn.Exec(ctx, "DELETE FROM task WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func (pg *RepositoryPG) CreateAssignments(ctx context.Context, task *domain.TaskAsignments) ([]domain.Assignment, error) {
	taskDetails, err := pg.GetTaskByID(ctx, task.TaskID)
	if err != nil {
		return nil, fmt.Errorf("can't get task details")
	}

	assignments := make([]domain.Assignment, 0, len(task.ToAssign))

	sql := "INSERT INTO assignment (id, class, task_id, lesson_id, task_payload, deadline) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING"
	batch := &pgx.Batch{}
	for _, cl := range task.ToAssign {
		assignmentID := uuid.New()

		assignments = append(assignments, domain.Assignment{
			AssignmentID: assignmentID,
			Class:        cl.Class,
			LessonID:     cl.LessonID,
		})

		batch.Queue(sql, assignmentID, cl.Class, task.TaskID, cl.LessonID, taskDetails.Payload, taskDetails.Deadline)
	}

	results := pg.conn.SendBatch(ctx, batch)
	defer results.Close()

	for range task.ToAssign {
		_, err := results.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				fmt.Println("error %w", pgErr.Message)
				if pgErr.Code == pgerrcode.UniqueViolation {
					continue
				}
			}

			return nil, fmt.Errorf("unnable create assignment %w", err)

		}

	}

	return assignments, nil
}

func (pg *RepositoryPG) UpdateAssignment(ctx context.Context, task *domain.TaskAsignment) error {
	_, err := pg.conn.Exec(ctx, "UPDATE assignment SET task_payload = $1, class = $2 WHERE id = $3", task.Payload, task.Class, task.AssignmentID)
	if err != nil {
		return err
	}

	return nil
}

func (pg *RepositoryPG) GetTaskByClass(ctx context.Context, class string) ([]*domain.LessonTask, error) {
	rows, err := pg.conn.Query(ctx, "SELECT id, lesson_id, task_id, task_payload, deadline FROM assignment where class = $1", class)
	if err != nil {
		return nil, fmt.Errorf("error executing prepared statement: %w", err)
	}

	defer rows.Close()

	var tasks []*domain.LessonTask
	for rows.Next() {
		var task domain.LessonTask
		err := rows.Scan(
			&task.TaskID,
			&task.LessonID,
			&task.TaskTemplateID,
			&task.Payload,
			&task.Deadline,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning task row: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows: %w", err)
	}

	return tasks, nil
}

func (pg *RepositoryPG) SetTaskResultsByUsers(ctx context.Context, taskResults *domain.TaskResult) error {
	sql := "INSERT INTO usersMark (id, user_id, task_id, lesson_id, mark) VALUES($1, $2, $3, $4, $5) ON CONFLICT (user_id, task_id, lesson_id) DO UPDATE SET mark = $5"
	batch := &pgx.Batch{}
	for _, userResults := range taskResults.UsersResult {
		batch.Queue(sql, uuid.New(), userResults.UserID, taskResults.TaskID, taskResults.LessonID, userResults.Mark)
	}

	results := pg.conn.SendBatch(ctx, batch)
	defer results.Close()

	for range taskResults.UsersResult {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("unable to update row: %w", err)
		}
	}

	return nil
}

func (pg *RepositoryPG) DeleteAssignment(ctx context.Context, assignmentID uuid.UUID) error {
	tag, err := pg.conn.Exec(ctx, "DELETE FROM assignment WHERE id=$1", assignmentID)
	if err != nil {
		return err
	}

	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrAssignmentNotFound
	}

	return nil
}

func (pg *RepositoryPG) CreateTaskWithAssignments(ctx context.Context, assignment *domain.TaskWithAsignment) (uuid.UUID, error) {
	tx, err := pg.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("starting transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO task (id, payload, deadline) VALUES($1, $2, $3)", assignment.TaskID, assignment.Payload, assignment.Deadline)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create new assignment records:%w", err)
	}

	var id uuid.UUID
	err = tx.QueryRow(ctx, "INSERT INTO assignment (id, class, task_id, lesson_id, task_payload, deadline) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		uuid.New(), assignment.Class, assignment.TaskID, assignment.LessonID, assignment.Payload, assignment.Deadline).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("can't create new assignment records:%w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("committing transaction: %w", err)
	}

	return id, nil
}
