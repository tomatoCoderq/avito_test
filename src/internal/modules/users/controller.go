package users

import (
	"github.com/gin-gonic/gin"
	"github.com/tomatoCoderq/avito_task/src/models"
	"gorm.io/gorm"
)

type ServiceMethods interface {
	SetIsActive(userID string, isActive bool) (*models.User, error)
	GetUserReviews(userID string) ([]models.PR, error)
}

type Controller struct {
	service ServiceMethods
}

func RegisterController(service ServiceMethods) *Controller {
	return &Controller{
		service: service,
	}
}

// SetIsActive устанавливает флаг активности пользователя
func (c *Controller) SetIsActive(ctx *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		IsActive bool   `json:"is_active"`
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

	user, err := c.service.SetIsActive(req.UserID, req.IsActive)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(404, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "user not found",
			},
		})
		return
	}
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update user",
			},
		})
		return
	}

	// Формируем ответ согласно OpenAPI спецификации
	teamName := ""
	if len(user.Teams) > 0 {
		teamName = user.Teams[0].Name
	}

	ctx.JSON(200, gin.H{
		"user": gin.H{
			"user_id":   user.ID,
			"username":  user.Name,
			"team_name": teamName,
			"is_active": user.IsActive,
		},
	})
}

// GetReview получает список PR'ов где пользователь назначен ревьювером
func (c *Controller) GetReview(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "user_id query parameter is required",
			},
		})
		return
	}

	prs, err := c.service.GetUserReviews(userID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user reviews",
			},
		})
		return
	}

	// Формируем ответ согласно OpenAPI спецификации
	pullRequests := make([]gin.H, 0, len(prs))
	for _, pr := range prs {
		pullRequests = append(pullRequests, gin.H{
			"pull_request_id":   pr.ID,
			"pull_request_name": pr.Name,
			"author_id":         pr.AuthorID,
			"status":            pr.Status,
		})
	}

	ctx.JSON(200, gin.H{
		"user_id":        userID,
		"pull_requests": pullRequests,
	})
}
