package teams

import (
	"github.com/gin-gonic/gin"
	"github.com/tomatoCoderq/avito_task/src/models"
	"gorm.io/gorm"
)

type ServiceMethods interface {
	TeamCreate(team *models.Team) (*models.Team, error)
	TeamGetByName(name string) (*models.Team, error)
}

type Controller struct {
	service ServiceMethods
}

func RegisterController(service ServiceMethods) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) TeamCreate(ctx *gin.Context) {
	var req struct {
		TeamName string `json:"team_name" binding:"required"`
		Members  []struct {
			UserID   string `json:"user_id" binding:"required"`
			Username string `json:"username" binding:"required"`
			IsActive bool   `json:"is_active"`
		} `json:"members" binding:"required"`
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

	// Создаем команду с пользователями
	team := &models.Team{
		Name:  req.TeamName,
		Users: make([]models.User, len(req.Members)),
	}

	for i, member := range req.Members {
		team.Users[i] = models.User{
			ID:       member.UserID,
			Name:     member.Username,
			IsActive: member.IsActive,
		}
	}

	createdTeam, err := c.service.TeamCreate(team)
	if err != nil {
		// Проверяем на дубликат команды
		if err.Error() == "team already exists" {
			ctx.JSON(400, gin.H{
				"error": gin.H{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		}
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create team",
			},
		})
		return
	}

	// Формируем ответ согласно OpenAPI
	members := make([]gin.H, len(createdTeam.Users))
	for i, user := range createdTeam.Users {
		members[i] = gin.H{
			"user_id":   user.ID,
			"username":  user.Name,
			"is_active": user.IsActive,
		}
	}

	ctx.JSON(201, gin.H{
		"team": gin.H{
			"team_name": createdTeam.Name,
			"members":   members,
		},
	})
}

func (c *Controller) TeamGetByName(ctx *gin.Context) {
	name := ctx.Query("team_name")
	if name == "" {
		ctx.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "team_name query parameter is required",
			},
		})
		return
	}

	team, err := c.service.TeamGetByName(name)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(404, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get team",
			},
		})
		return
	}

	// Формируем ответ согласно OpenAPI
	members := make([]gin.H, len(team.Users))
	for i, user := range team.Users {
		members[i] = gin.H{
			"user_id":   user.ID,
			"username":  user.Name,
			"is_active": user.IsActive,
		}
	}

	ctx.JSON(200, gin.H{
		"team_name": team.Name,
		"members":   members,
	})
}
