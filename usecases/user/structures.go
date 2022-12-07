package user

type LoginUser struct {
  Email    string
  Password string
}

type NewUser struct {
  User      string `example:"User" gorm:"unique"`
  Email     string `example:"some@mail.com" gorm:"unique"`
  FirstName string `example:"John"`
  LastName  string `example:"Doe"`
  Password  string `example:"SomeHashPass"`
}
