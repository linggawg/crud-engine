package usecases

import (
	"context"
	"database/sql"
	modelsServices "engine/bin/modules/services/models/domain"
	commandsServices "engine/bin/modules/services/repositories/commands"
	queriesServices "engine/bin/modules/services/repositories/queries"
	models "engine/bin/modules/users-services/models/domain"
	"engine/bin/modules/users-services/repositories/commands"
	"engine/bin/modules/users-services/repositories/queries"
	"strings"

	httpError "engine/bin/pkg/http-error"

	"engine/bin/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type UsersServicesCommandUsecase struct {
	UsersServicesPostgreCommand commands.UsersServicesPostgre
	UsersServicesPostgreQuery   queries.UsersServicesPostgre
	ServicesPostgreCommand      commandsServices.ServicesPostgre
	ServicesPostgreQuery        queriesServices.ServicesPostgre
}

func NewUsersServicesCommandUsecase(
	UsersServicesPostgreCommand commands.UsersServicesPostgre,
	UsersServicesPostgreQuery queries.UsersServicesPostgre,
	ServicesPostgreCommand commandsServices.ServicesPostgre,
	ServicesPostgreQuery queriesServices.ServicesPostgre) *UsersServicesCommandUsecase {
	return &UsersServicesCommandUsecase{
		UsersServicesPostgreCommand: UsersServicesPostgreCommand,
		UsersServicesPostgreQuery:   UsersServicesPostgreQuery,
		ServicesPostgreCommand:      ServicesPostgreCommand,
		ServicesPostgreQuery:        ServicesPostgreQuery,
	}
}

func (u UsersServicesCommandUsecase) InsertAllByServices(ctx context.Context, payload models.UsersServicesRequest) utils.Result {
	var result utils.Result
	if !strings.EqualFold(payload.Opts.RoleName, "admin") {
		errObj := httpError.NewUnauthorized()
		errObj.Message = "unauthorized access"
		result.Error = errObj
		return result
	}

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, method := range methods {
		var serviceID, serviceUrl string
		service, err := u.ServicesPostgreQuery.FindOneByServiceUrlAndMethodAndQueryIsNull(ctx, payload.ServiceUrl, method)
		if err != nil && err != sql.ErrNoRows {
			errObj := httpError.NewInternalServerError()
			errObj.Message = err.Error()
			result.Error = errObj
			return result
		}

		if service == nil {
			service := modelsServices.Services{
				ID:         uuid.New().String(),
				DbID:       payload.DbID,
				Method:     method,
				ServiceUrl: &payload.ServiceUrl,
				CreatedAt:  null.TimeFrom(time.Now()),
				CreatedBy:  &payload.Opts.UserID,
				ModifiedAt: null.TimeFrom(time.Now()),
				ModifiedBy: &payload.Opts.UserID,
			}

			err := u.ServicesPostgreCommand.InsertOne(ctx, &service)
			if err != nil {
				errObj := httpError.NewInternalServerError()
				errObj.Message = err.Error()
				result.Error = errObj
				return result
			}
			serviceID = service.ID
			serviceUrl = *service.ServiceUrl
		} else {
			serviceID = service.ID
			serviceUrl = *service.ServiceUrl
		}

		userservice, err := u.UsersServicesPostgreQuery.FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull(ctx, serviceUrl, payload.UserID, method)
		if err != nil && err != sql.ErrNoRows {
			errObj := httpError.NewInternalServerError()
			errObj.Message = err.Error()
			result.Error = errObj
			return result
		}

		if userservice == nil {
			userservice := models.UsersServices{
				ID:         uuid.New().String(),
				UserID:     payload.UserID,
				ServiceID:  serviceID,
				CreatedAt:  null.TimeFrom(time.Now()),
				CreatedBy:  &payload.Opts.UserID,
				ModifiedAt: null.TimeFrom(time.Now()),
				ModifiedBy: &payload.Opts.UserID,
			}

			err := u.UsersServicesPostgreCommand.InsertOne(ctx, &userservice)
			if err != nil {
				errObj := httpError.NewInternalServerError()
				errObj.Message = err.Error()
				result.Error = errObj
				return result
			}
		}
	}
	return result
}

func (u UsersServicesCommandUsecase) DeleteByUsersIdAndServiceUrl(ctx context.Context, payload models.UsersServicesRequest) utils.Result {
	var result utils.Result
	if !strings.EqualFold(payload.Opts.RoleName, "admin") {
		errObj := httpError.NewUnauthorized()
		errObj.Message = "unauthorized access"
		result.Error = errObj
		return result
	}

	err := u.UsersServicesPostgreCommand.DeleteByUsersIdAndServiceUrl(ctx, payload.UserID, payload.ServiceUrl)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	return result
}
