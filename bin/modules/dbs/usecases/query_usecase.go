package usecases

import (
	"context"
	"engine/bin/modules/dbs/repositories/queries"
	serviceModels "engine/bin/modules/services/models/domain"
	dbsServices "engine/bin/modules/services/repositories/queries"
	usersQuery "engine/bin/modules/users/repositories/queries"
	userServiceQuery "engine/bin/modules/userservice/repositories/queries"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"net/url"
)

type DbsQueryUsecase struct {
	dbsPostgreQuery      queries.DbsPostgre
	servicesPostgreQuery dbsServices.ServicesPostgre
	userServiceQuery     userServiceQuery.UserServicePostgre
	usersQuery           usersQuery.UsersPostgre
}

func NewQueryUsecase(dbsPostgreQuery queries.DbsPostgre, servicesPostgreQuery dbsServices.ServicesPostgre, userServiceQuery userServiceQuery.UserServicePostgre, usersQuery usersQuery.UsersPostgre) *DbsQueryUsecase {
	return &DbsQueryUsecase{
		dbsPostgreQuery:      dbsPostgreQuery,
		servicesPostgreQuery: servicesPostgreQuery,
		userServiceQuery:     userServiceQuery,
		usersQuery:           usersQuery,
	}
}

func (u DbsQueryUsecase) GetByID(ctx context.Context, id string) utils.Result {
	var result utils.Result

	dbs, err := u.dbsPostgreQuery.GetByID(ctx, id)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Data dbs tidak ditemukan"
		result.Error = errObj
		return result
	}

	result.Data = dbs
	return result
}

func (h *DbsQueryUsecase) GetDbsConnection(ctx context.Context, userId, method, serviceUrl string, isQuery bool) utils.Result {
	var (
		result  utils.Result
		service *serviceModels.Services
		err     error
	)

	_, err = h.usersQuery.GetByID(ctx, userId)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "User tidak ditemukan"
		result.Error = errObj
		return result
	}

	if isQuery {
		serviceUrl, _ = url.QueryUnescape(serviceUrl)
		service, err = h.servicesPostgreQuery.GetByServiceDefinitionAndMethod(ctx, serviceUrl, method)
		if err != nil {
			errObj := httpError.NewNotFound()
			errObj.Message = "Data service definition tidak ditemukan"
			result.Error = errObj
			return result
		}
	} else {
		service, err = h.servicesPostgreQuery.GetByServiceUrlAndMethod(ctx, serviceUrl, method)
		if err != nil {
			errObj := httpError.NewNotFound()
			errObj.Message = "Data service url tidak ditemukan"
			result.Error = errObj
			return result
		}
	}

	_, err = h.userServiceQuery.GetByServiceIDAndUserId(ctx, service.ID, userId)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Data user service tidak ditemukan"
		result.Error = errObj
		return result
	}
	dbs, err := h.dbsPostgreQuery.GetByID(ctx, service.DbID)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Data dbs tidak ditemukan"
		result.Error = errObj
		return result
	}

	result.Data = dbs
	return result
}
