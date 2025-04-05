package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"task/internal/domain"
	"task/internal/ports/httpServer/common"
	"task/internal/ports/httpServer/request"
	"task/internal/ports/httpServer/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error)
	GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetTasks(ctx context.Context) ([]*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
	CreateAssignments(ctx context.Context, taskAssignments *domain.TaskAsignments) ([]domain.Assignment, error)
	GetTaskByClass(ctx context.Context, ckass string) ([]*domain.LessonTask, error)
	SetTaskResultsByUsers(ctx context.Context, taskResults *domain.TaskResult) error
	DeleteAssignment(ctx context.Context, assignmentID uuid.UUID) error
	CreateTaskWithAssignments(ctx context.Context, assignments *domain.TaskWithAsignment) (uuid.UUID, error)
	UpdateAssignment(ctx context.Context, task *domain.TaskAsignment) error
}

type Handler struct {
	taskService TaskService
	logger      *slog.Logger
}

func NewHandler(logger *slog.Logger, taskService TaskService) *Handler {
	return &Handler{
		logger:      logger,
		taskService: taskService,
	}
}

// CreateTask godoc
// @Summary Создание шаблона задачи(без назначения на классы и уроки)
// @Description Создает шаблон/задачу(без назначения на классы и уроки), но этот шаблон может использоваться для создания назначения
// @tags tasks
// @Accept json
// @Param tasks body request.Task true "Данные задачи"
// @Produce json
// @Success 201 {object} response.TaskID
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task [post].
func (h *Handler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()
	var input request.Task

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("failed to bind body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	task, err := h.taskService.CreateTask(ctx, input.ToDomain())
	if err != nil {
		h.logger.Error("failed to create task", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	c.JSON(http.StatusCreated, response.NewTaskIDResponse(task))
}

// GetTask godoc
// @Summary Поучить задачу
// @Description Получить задачу
// @tags tasks
// @Accept json
// @Param id path string true "ID задачи"
// @Produce json
// @Success 200 {object} response.Task
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/{id} [get].
func (h *Handler) GetTask(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("failed to parse id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	task, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		h.logger.Error("failed to create task", slog.String("error", err.Error()))
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewTaskResponse(task))
}

// GetTask godoc
// @Summary Поучить все шаблоны задач
// @Description Получить все шаблоны задач
// @tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} response.Task
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/all [get].
func (h *Handler) GetTasks(c *gin.Context) {
	ctx := c.Request.Context()

	task, err := h.taskService.GetTasks(ctx)
	if err != nil {
		h.logger.Error("failed to create task", slog.String("error", err.Error()))
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewTasksResponse(task))
}

// UpdateTask godoc
// @Summary Обновить шаблон задачи
// @Description Обновить шаблон задачи
// @tags tasks
// @Accept json
// @Param id path string true "ID задачи"
// @Param tasks body request.Task true "Данные задачи"
// @Produce json
// @Success 200 {object} response.TaskID
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/{id}/update [put].
func (h *Handler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("failed to parse id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	var input request.Task

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("failed to bind body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	task, err := h.taskService.UpdateTask(ctx, input.ToDomainWithID(taskID))
	if err != nil {
		h.logger.Error("failed to update task", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewTaskIDResponse(task))
}

// DeleteTask godoc
// @Summary Удалить шаблон задачи
// @Description Удалить шаблон задачи(удалятся все назначения, которые были созданы по задаче)
// @tags tasks
// @Accept json
// @Param id path string true "ID задачи"
// @Produce json
// @Success 200 {object} response.TaskID
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/{id}/delete [delete].
func (h *Handler) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("failed to parse id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	err = h.taskService.DeleteTask(ctx, taskID)
	if err != nil {
		h.logger.Error("failed to create task", slog.String("error", err.Error()))
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewTaskIDResponse(taskID))
}

// GetTasksByClass godoc
// @Summary Поучить задачи класса
// @Description Поучить задачи класса
// @tags tasks
// @Accept json
// @Param class query string true "название класса"
// @Produce json
// @Success 200 {object} response.ClassTasks
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/get-by-class [get].
func (h *Handler) GetTaskByClass(c *gin.Context) {
	ctx := c.Request.Context()
	var class request.Class

	if err := c.BindQuery(&class); err != nil {
		h.logger.Error("failed to bind query class", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	tasks, err := h.taskService.GetTaskByClass(ctx, class.Class)
	if err != nil {
		h.logger.Error("failed to get tasks by class", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewClassTasksResponse(class.Class, tasks))
}

// DeleteAssignment godoc
// @Summary Удалить задачу с класса и урока
// @Description Удалить задачу с класса и урока
// @tags tasks
// @Accept json
// @Param task-assign body request.TaskAsignmentID true "id задачи для удаления"
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/assignment-delete [delete].
func (h *Handler) DeleteAssignment(c *gin.Context) {
	ctx := c.Request.Context()
	var input request.TaskAsignmentID

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("failed to bind body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	assignment, err := input.ToUUID()
	if err != nil {
		h.logger.Error("failed covert to uuid", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	err = h.taskService.DeleteAssignment(ctx, assignment)
	if err != nil {
		h.logger.Error("failed to delete assignment", slog.String("error", err.Error()))
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.String(http.StatusOK, "OK")
}

// AssignTaskToClasses godoc
// @Summary Назначить задачу классу и уроку
// @Description Назначить задачу классу и уроку
// @tags tasks
// @Accept json
// @Param task-assign body request.TaskAsignments true "Данные для назначения"
// @Produce json
// @Success 200 {object} response.TaskAssignments
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/assignment [post].
func (h *Handler) AssignTaskToClasses(c *gin.Context) {
	ctx := c.Request.Context()
	var input request.TaskAsignments

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("failed to bind body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	domainAssignments, err := input.ToDomain()
	if err != nil {
		h.logger.Error("failed assert to domain", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	assignments, err := h.taskService.CreateAssignments(ctx, domainAssignments)
	if err != nil {
		h.logger.Error("failed to create assignment", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, response.NewAssignmentsResponse(input.TaskID, assignments))
}

// CreateTaskWithAssignment godoc
// @Summary Создать задачу для класса
// @Description Создать задачу для класса
// @tags tasks
// @Accept json
// @Param task-assign body request.TaskWithAsignment true "Данные для создания задачи с назначением классу и уроку"
// @Produce json
// @Success 200 {object} response.AssignmentID
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/create-with-assignment [post].
func (h *Handler) CreateTaskWithAssignment(c *gin.Context) {
	ctx := c.Request.Context()
	var input request.TaskWithAsignment

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("failed to bind body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	domainAssignments, err := input.ToDomain()
	if err != nil {
		h.logger.Error("failed assert to domain", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	assignment, err := h.taskService.CreateTaskWithAssignments(ctx, domainAssignments)
	if err != nil {
		h.logger.Error("failed to create assignment", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, response.NewAssignmentIDResponse(assignment.String()))
}

// UpdateTaskAssignment godoc
// @Summary Обновить задачу для класса
// @Description Обновить задачу для класса
// @tags tasks
// @Accept json
// @Param task-assign body request.TaskAsignment true "Данные для назначения"
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/assignment-update [put].
func (h *Handler) UpdateTaskAssignment(c *gin.Context) {
	ctx := c.Request.Context()
	var input request.TaskAsignment

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("failed to bind body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	domainAssignment, err := input.ToDomain()
	if err != nil {
		h.logger.Error("failed assert to domain", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	err = h.taskService.UpdateAssignment(ctx, domainAssignment)
	if err != nil {
		h.logger.Error("failed to create assignment", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.String(http.StatusOK, "OK")
}

// TaskResults godoc
// @Summary Поставить результаты за задачу ученикам
// @Description Поставить результаты за задачу ученикам
// @tags tasks
// @Accept json
// @Param task-results body request.TaskResult true "Оценки пользователей за задачу"
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/task/result [post].
func (h *Handler) TaskResult(c *gin.Context) {
	ctx := c.Request.Context()
	var input request.TaskResult
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("failed to bind body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	taskResults, err := input.ToDomain()
	if err != nil {
		h.logger.Error("failed assert to domain", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	err = h.taskService.SetTaskResultsByUsers(ctx, taskResults)
	if err != nil {
		h.logger.Error("failed to set result", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), http.StatusInternalServerError))
		return
	}

	c.String(http.StatusOK, "OK")
}
