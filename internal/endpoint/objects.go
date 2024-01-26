package endpoint

type (
	// CATEGORIES
	CoreCategory struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at,omitempty"`
		Name      string `json:"name"`
	}
	GetCategoryRequest struct {
		ID string `json:"id"`
	}

	GetCategoryResponse struct {
		CoreCategory
	}

	InsertCategoryRequest struct {
		Name string `json:"name" description:"Nome da categoria de produto"`
	}

	InsertCategoryResponse struct {
		CoreCategory
	}

	DeleteCategoryRequest struct {
		ID string `json:"id"`
	}

	DeleteCategoryResponse struct {
		ID string `json:"id"`
	}

	ListCategoriesRequest struct {
		Limit  int `json:"limit" default:"10" description:"Quantidade de registros"`
		Offset int `json:"offset"`
	}

	ListCategoriesResponse struct {
		Categories []CoreCategory `json:"categories"`
		Limit      int            `json:"limit" default:"10"`
		Offset     int            `json:"offset"`
		Total      int64          `json:"total"`
	}
)

type (
	// PRODUCTS

	Product struct {
		ID          string `json:"id,omitempty"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Price       string `json:"price"`
		CreatedAt   string `json:"created_at,omitempty" readOnly:"true"`
		UpdatedAt   string `json:"updated_at,omitempty" readOnly:"true"`
		DeletedAt   string `json:"deleted_at,omitempty" readOnly:"true"`
	}

	ProductList struct {
		Products []Product `json:"products"`
		Limit    int       `json:"limit" default:"10"`
		Offset   int       `json:"offset"`
		Total    int64     `json:"total"`
	}
)
