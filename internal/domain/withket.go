package domain

type ProjectInfo struct {
	ID                    uint     `gorm:"primaryKey"` // Primary key field
	Title                 string   `json:"title"`
	URL                   string   `json:"url"`
	StatusMarks           string   `json:"status_marks"`
	EstimatedAmount       string   `json:"estimated_amount"`
	EstimatedDuration     string   `json:"estimated_duration"`
	WorkStartDate         string   `json:"work_start_date"`
	NumberOfApplicants    string   `json:"number_of_applicants"`
	ProjectCategoryOrRole string   `json:"project_category_or_role"`
	Location              string   `json:"location"`
	Skills                []string `gorm:"-"` // Excluded from DB, as GORM doesn’t support slices by default
}

func (ProjectInfo) TableName() string {
	return "wishket"
}
