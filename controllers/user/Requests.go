package user

type NewUserRequest struct {
  User      string `json:"user" example:"someUser" gorm:"unique" binding:"required"`
  Email     string `json:"email" example:"mail@mail.com" gorm:"unique" binding:"required"`
  FirstName string `json:"firstName" example:"John" binding:"required"`
  LastName  string `json:"lastName" example:"Doe" binding:"required"`
  Password  string `json:"password" example:"Password123" binding:"required"`
}

type LoginRequest struct {
  Email    string `json:"email" example:"mail@mail.com" gorm:"unique" binding:"required"`
  Password string `json:"password" example:"SomePass" binding:"required"`
}
