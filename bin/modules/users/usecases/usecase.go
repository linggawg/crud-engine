package usecases

import (
	"context"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/pkg/utils"
)

type CommandUsecase interface {
	Login(ctx context.Context, params models.ReqLogin) utils.Result
	RegisterUser(ctx context.Context, params models.ReqUser) utils.Result
}
