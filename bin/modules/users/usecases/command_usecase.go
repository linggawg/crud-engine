package usecases

import (
	"context"
	"database/sql"
	rolesQueries "engine/bin/modules/roles/repositories/queries"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/modules/users/repositories/commands"
	"engine/bin/modules/users/repositories/queries"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/token"

	"github.com/google/uuid"

	"engine/bin/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

type UsersCommandUsecase struct {
	UsersPostgreCommand commands.UsersPostgre
	UsersPostgreQuery   queries.UsersPostgre
	RolesPostgreQuery   rolesQueries.RolesPostgre
}

func NewUsersCommandUsecase(UsersPostgreCommand commands.UsersPostgre, UsersPostgreQuery queries.UsersPostgre, RolesPostgreQuery rolesQueries.RolesPostgre) *UsersCommandUsecase {
	return &UsersCommandUsecase{
		UsersPostgreCommand: UsersPostgreCommand,
		UsersPostgreQuery:   UsersPostgreQuery,
		RolesPostgreQuery:   RolesPostgreQuery,
	}
}

func (u UsersCommandUsecase) Login(ctx context.Context, params models.ReqLogin) utils.Result {
	var result utils.Result

	users, err := u.UsersPostgreQuery.FindOneByUsername(ctx, params.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			errObj := httpError.NewNotFound()
			errObj.Message = "username / password salah"
			result.Error = errObj
		} else {
			errObj := httpError.NewInternalServerError()
			errObj.Message = err.Error()
			result.Error = errObj
		}
		return result
	}

	roles, err := u.RolesPostgreQuery.FindOneByID(ctx, users.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			errObj := httpError.NewNotFound()
			errObj.Message = "roles tidak ditemukan"
			result.Error = errObj
		} else {
			errObj := httpError.NewInternalServerError()
			errObj.Message = err.Error()
			result.Error = errObj
		}
		return result
	}

	err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(params.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		errObj := httpError.NewBadRequest()
		errObj.Message = "login gagal"
		result.Error = errObj
		return result
	}

	token, err := token.Generate(users.ID, roles.Name, params.Duration)
	if err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	result.Data = token
	return result
}

func (u UsersCommandUsecase) RegisterUser(ctx context.Context, params models.ReqUser) utils.Result {
	var result utils.Result

	validatebyUsername, err := u.UsersPostgreQuery.FindOneByUsername(ctx, params.Username)
	if err != nil && err != sql.ErrNoRows {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	} else if validatebyUsername != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = "username sudah terpakai"
		result.Error = errObj
		return result
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	password := string(hashedPassword)

	users := &models.Users{
		ID:         uuid.NewString(),
		Username:   params.Username,
		RoleID:     params.RoleID,
		Password:   password,
		CreatedBy:  &params.Opts.UserID,
		CreatedAt:  null.TimeFrom(time.Now()),
		ModifiedAt: null.TimeFrom(time.Now()),
		ModifiedBy: &params.Opts.UserID,
	}

	err = u.UsersPostgreCommand.InsertOne(ctx, users)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	result.Data = users
	return result
}
