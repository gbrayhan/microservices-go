package user

import model "github.com/gbrayhan/microservices-go/models/user"

func UserModelToResponseMapper(userModel model.User) (createUserResponse *UserResponse) {
  createUserResponse = &UserResponse{ID: userModel.ID, User: userModel.User,
    Email: userModel.Email, FirstName: userModel.FirstName, LastName: userModel.LastName,
    Status: userModel.Status, CreatedAt: userModel.CreatedAt, UpdatedAt: userModel.UpdatedAt}

  return
}
