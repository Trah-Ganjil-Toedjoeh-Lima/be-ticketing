package validation

type LoginValidation struct {
	Name  string `json:"name" binding:"required_without=Email,min=0,max=36"`
	Email string `json:"email" binding:"required_without=Name,min=0,max=36"`
	Phone string `json:"phone" binding:"required,min=7,max=15,number"`
}
