package types

type UserRegisterInput struct {
	FirstName  string `form:"first_name" binding:"required"`
	SecondName string `form:"second_name" binding:"required"`
	LastName   string `form:"last_name" binding:"required"`
	Email      string `form:"email" binding:"required,email"`
	Phone      string `form:"phone" binding:"required"`
	Password   string `form:"password" binding:"required,min=6"`
	Role       string `form:"role" binding:"required"`
}
