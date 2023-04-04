package validation

type LoginValidation struct {
	Name  string `json:"name" binding:"required_without=Email,min=0,max=36"`
	Email string `json:"email" binding:"required_without=Name,min=0,max=36"`
	Phone string `json:"phone" binding:"required,min=7,max=15,number"`
}

type RegisterValidation struct {
	Name  string `json:"name" binding:"required,min=3,max=36"`
	Email string `json:"email" binding:"required,email,min=5,max=36"`
	Phone string `json:"phone" binding:"required,min=7,max=15,number"`
}

type OtpVerification struct {
	Email string `json:"email" binding:"required,email,min=5,max=36"`
	Otp   string `json:"otp" binding:"required,number,min=6,max=6"`
}

type RegisterEmailValidation struct {
	Email string `json:"email" binding:"required,email,min=5,max=36"`
}
