package stats

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceMethods interface {
	GetUserStats(userID string) (*UserStats, error)
	GetOverviewStats() (*OverviewStats, error)
	GetTeamStats(teamName string) (*TeamStats, error)
}

type Controller struct {
	service ServiceMethods
}

func RegisterController(service ServiceMethods) *Controller {
	return &Controller{
		service: service,
	}
}

// GetUserStats возвращает статистику по пользователю
func (c *Controller) GetUserStats(ctx *gin.Context) {
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

	stats, err := c.service.GetUserStats(userID)
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
				"message": "Failed to get user statistics",
			},
		})
		return
	}

	ctx.JSON(200, gin.H{
		"user_id":  stats.UserID,
		"username": stats.Username,
		"statistics": gin.H{
			"authored_prs": gin.H{
				"total":  stats.AuthoredTotal,
				"open":   stats.AuthoredOpen,
				"merged": stats.AuthoredMerged,
			},
			"reviewing_prs": gin.H{
				"total":  stats.ReviewingTotal,
				"open":   stats.ReviewingOpen,
				"merged": stats.ReviewingMerged,
			},
			"team_name": stats.TeamName,
		},
	})
}

// GetOverview возвращает общую статистику системы
func (c *Controller) GetOverview(ctx *gin.Context) {
	stats, err := c.service.GetOverviewStats()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get overview statistics",
			},
		})
		return
	}

	ctx.JSON(200, gin.H{
		"total_users":                stats.TotalUsers,
		"active_users":               stats.ActiveUsers,
		"total_teams":                stats.TotalTeams,
		"total_prs":                  stats.TotalPRs,
		"open_prs":                   stats.OpenPRs,
		"merged_prs":                 stats.MergedPRs,
		"top_reviewers":              stats.TopReviewers,
	})
}

// GetTeamStats возвращает статистику по команде
func (c *Controller) GetTeamStats(ctx *gin.Context) {
	teamName := ctx.Query("team_name")
	if teamName == "" {
		ctx.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "team_name query parameter is required",
			},
		})
		return
	}

	stats, err := c.service.GetTeamStats(teamName)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(404, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "team not found",
			},
		})
		return
	}
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get team statistics",
			},
		})
		return
	}

	ctx.JSON(200, gin.H{
		"team_name": stats.TeamName,
		"team_statistics": gin.H{
			"total_members":    stats.TotalMembers,
			"active_members":   stats.ActiveMembers,
			"total_prs":        stats.TotalPRs,
			"open_prs":         stats.OpenPRs,
			"merged_prs":       stats.MergedPRs,
			"top_contributors": stats.TopContributors,
		},
	})
}