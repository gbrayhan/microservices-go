package user

import (
  "fmt"
  model "github.com/gbrayhan/microservices-go/models/user"
  "golang.org/x/crypto/bcrypt"
)

func Login(user LoginUser) () {
  userMap := map[string]interface{}{"email": user.Email}
  userModel, err := model.GetUserByMap(userMap)
  if err != nil {
    return
  }
  fmt.Println(userModel)

  return
}

func CreateUser(user *NewUser) (userModel model.User, err error) {
  userModel = model.User{User: user.User, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}

  // Generate "hash" to store from user password
  hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
  if err != nil {
    return
  }
  userModel.HashPassword = string(hash)
  userModel.Status = true
  err = model.CreateUser(&userModel)

  return
}
