package usecases

import (
	"context"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/modules/users/repositories/commands"
	"engine/bin/modules/users/repositories/queries"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/token"

	"engine/bin/pkg/utils"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

type UsersCommandUsecase struct {
	UsersPostgreCommand commands.UsersPostgre
	UsersPostgreQuery   queries.UsersPostgre
}

func NewCommandUsecase(UsersPostgreCommand commands.UsersPostgre, UsersPostgreQuery queries.UsersPostgre) *UsersCommandUsecase {
	return &UsersCommandUsecase{
		UsersPostgreCommand: UsersPostgreCommand,
		UsersPostgreQuery:   UsersPostgreQuery,
	}
}

func (u UsersCommandUsecase) Login(ctx context.Context, params models.ReqLogin) utils.Result {
	var result utils.Result

	users, err := u.UsersPostgreQuery.GetByEmail(ctx, params.Email)
	if err != nil {
		log.Println(err)
		errObj := httpError.NewNotFound()
		errObj.Message = "Email tidak ditemukan"
		result.Error = errObj
		return result
	}

	err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(params.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println(err)
		errObj := httpError.NewBadRequest()
		errObj.Message = "Login gagal"
		result.Error = errObj
		return result
	}

	token, err := token.Generate(users.ID)
	if err != nil {
		log.Println(err)
		errObj := httpError.NewBadRequest()
		errObj.Message = "Login gagal"
		result.Error = errObj
		return result
	}

	result.Data = token
	return result
}

func (u UsersCommandUsecase) RegisterUser(ctx context.Context, params models.ReqUser) utils.Result {
	var result utils.Result

	validatebyEmail, _ := u.UsersPostgreQuery.GetByEmail(ctx, params.Email)
	if validatebyEmail != nil {
		log.Println("Email sudah terpakai")
		errObj := httpError.NewBadRequest()
		errObj.Message = "Email sudah terpakai"
		result.Error = errObj
		return result
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		errObj := httpError.NewBadRequest()
		errObj.Message = "Login gagal"
		result.Error = errObj
		return result
	}
	password := string(hashedPassword)

	users := &models.Users{
		Username:  params.Username,
		Email:     params.Email,
		Password:  password,
		CreatedBy: params.UserId,
		CreatedAt: null.TimeFrom(time.Now()),
	}

	err = u.UsersPostgreCommand.InsertOne(ctx, users)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "Failed insert user"
		result.Error = errObj
		return result
	}

	result.Data = users
	return result
}
