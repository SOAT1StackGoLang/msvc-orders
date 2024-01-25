package endpoint

type (
	GetCategoryRequest struct {
		ID string `json:"id"`
	}

	GetCategoryResponse struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at,omitempty"`
		Name      string `json:"name"`
	}
)
