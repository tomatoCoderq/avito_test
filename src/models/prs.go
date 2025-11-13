package models

type PR struct {
	ID                string `gorm:"type:varchar(255);primaryKey"`
	Name              string
	AuthorID          string `gorm:"type:varchar(255)"`
	Author            User   `gorm:"foreignKey:AuthorID"`
	Status            string `gorm:"type:varchar(50);default:'OPEN'"`
	Reviewers         []User `gorm:"many2many:pr_reviewers;constraint:OnDelete:CASCADE;"`
	NeedMoreReviewers bool   // in task small letter but go requires capital
}
