package validation

type RegisterValidation struct {
	Name  string `json:"name" binding:"required,min=3,max=36"`
	Email string `json:"email" binding:"required,email,min=5,max=36"`
	Phone string `json:"phone" binding:"required,min=7,max=15,number"`
}
