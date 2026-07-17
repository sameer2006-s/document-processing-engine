package auth

import( 
    "time"
    "github.com/google/uuid"
)
type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`

	FirstName string
	LastName  string

	Active bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time

	SessionToken string `gorm:"uniqueIndex"`
}