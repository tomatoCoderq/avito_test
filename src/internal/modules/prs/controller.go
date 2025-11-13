package prs

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tomatoCoderq/avito_task/src/models"
	"gorm.io/gorm"
)

type ServiceMethods interface {
	CreatePR(prID, prName, authorID string) (*models.PR, error)
	GetPRByID(prID string) (*models.PR, error)
	MergePR(prID string) (*models.PR, error)
	ReassignReviewer(prID, oldUserID string) (*models.PR, string, error)
}

type Controller struct {
	service ServiceMethods
}

func RegisterController(service ServiceMethods) *Controller {
	return &Controller{
		service: service,
	}
}

// Create создает PR и назначает ревьюверов
func (c *Controller) Create(ctx *gin.Context) {
	var req struct {
		PullRequestID   string `json:"pull_request_id" binding:"required"`
		PullRequestName string `json:"pull_request_name" binding:"required"`
		AuthorID        string `json:"author_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
		return
	}

	pr, err := c.service.CreatePR(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		// Автор или команда не найдены
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "author has no team") {
			ctx.JSON(404, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "author or team not found",
				},
			})
			return
		}
		// Дубликат PR
		if strings.Contains(err.Error(), "duplicate") || 
		   strings.Contains(err.Error(), "already exists") ||
		   strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		   strings.Contains(err.Error(), "violates unique constraint") {
			ctx.JSON(409, gin.H{
				"error": gin.H{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
			return
		}
		// Остальные ошибки
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create PR",
			},
		})
		return
	}

	// Формируем ответ
	reviewerIDs := make([]string, 0, len(pr.Reviewers))
	for _, reviewer := range pr.Reviewers {
		reviewerIDs = append(reviewerIDs, reviewer.ID)
	}

	ctx.JSON(201, gin.H{
		"pr": gin.H{
			"pull_request_id":     pr.ID,
			"pull_request_name":   pr.Name,
			"author_id":           pr.AuthorID,
			"status":              pr.Status,
			"assigned_reviewers":  reviewerIDs,
		},
	})
}

// Merge помечает PR как MERGED
func (c *Controller) Merge(ctx *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
		return
	}

	pr, err := c.service.MergePR(req.PullRequestID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(404, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "PR not found",
			},
		})
		return
	}
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to merge PR",
			},
		})
		return
	}

	// Формируем ответ
	reviewerIDs := make([]string, 0, len(pr.Reviewers))
	for _, reviewer := range pr.Reviewers {
		reviewerIDs = append(reviewerIDs, reviewer.ID)
	}

	ctx.JSON(200, gin.H{
		"pr": gin.H{
			"pull_request_id":     pr.ID,
			"pull_request_name":   pr.Name,
			"author_id":           pr.AuthorID,
			"status":              pr.Status,
			"assigned_reviewers":  reviewerIDs,
		},
	})
}

// Reassign переназначает ревьювера
func (c *Controller) Reassign(ctx *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id" binding:"required"`
		OldUserID     string `json:"old_reviewer_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
		return
	}

	pr, replacedBy, err := c.service.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "PR or user not found",
				},
			})
			return
		}
		
		errMsg := err.Error()
		if strings.Contains(errMsg, "PR_MERGED") {
			ctx.JSON(409, gin.H{
				"error": gin.H{
					"code":    "PR_MERGED",
					"message": "cannot reassign on merged PR",
				},
			})
			return
		}
		if strings.Contains(errMsg, "NOT_ASSIGNED") {
			ctx.JSON(409, gin.H{
				"error": gin.H{
					"code":    "NOT_ASSIGNED",
					"message": "reviewer is not assigned to this PR",
				},
			})
			return
		}
		if strings.Contains(errMsg, "NO_CANDIDATE") {
			ctx.JSON(409, gin.H{
				"error": gin.H{
					"code":    "NO_CANDIDATE",
					"message": "no active replacement candidate in team",
				},
			})
			return
		}

		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to reassign reviewer",
			},
		})
		return
	}

	// Формируем ответ
	reviewerIDs := make([]string, 0, len(pr.Reviewers))
	for _, reviewer := range pr.Reviewers {
		reviewerIDs = append(reviewerIDs, reviewer.ID)
	}

	ctx.JSON(200, gin.H{
		"pr": gin.H{
			"pull_request_id":     pr.ID,
			"pull_request_name":   pr.Name,
			"author_id":           pr.AuthorID,
			"status":              pr.Status,
			"assigned_reviewers":  reviewerIDs,
		},
		"replaced_by": replacedBy,
	})
}

// GetByID получает PR по ID с информацией о ревьюверах
func (c *Controller) GetByID(ctx *gin.Context) {
	prID := ctx.Query("pull_request_id")
	if prID == "" {
		ctx.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "pull_request_id query parameter is required",
			},
		})
		return
	}

	pr, err := c.service.GetPRByID(prID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(404, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "PR not found",
			},
		})
		return
	}
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get PR",
			},
		})
		return
	}

	// Формируем ответ
	reviewerIDs := make([]string, 0, len(pr.Reviewers))
	for _, reviewer := range pr.Reviewers {
		reviewerIDs = append(reviewerIDs, reviewer.ID)
	}

	ctx.JSON(200, gin.H{
		"pull_request_id":    pr.ID,
		"pull_request_name":  pr.Name,
		"author_id":          pr.AuthorID,
		"status":             pr.Status,
		"assigned_reviewers": reviewerIDs,
	})
}
