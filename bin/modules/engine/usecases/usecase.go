package usecases

import (
	"context"
	dbsModels "engine/bin/modules/dbs/models/domain"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/pkg/utils"
)

type CommandUsecase interface {
	Insert(ctx context.Context, dbs dbsModels.Dbs, payload *models.EngineRequest) (result utils.Result)
	Update(ctx context.Context, dbs dbsModels.Dbs, payload *models.EngineRequest) (result utils.Result)
	Delete(ctx context.Context, dbs dbsModels.Dbs, payload *models.EngineRequest) (result utils.Result)
}

type QueryUsecase interface {
	Get(ctx context.Context, dbs dbsModels.Dbs, table string, payload *models.GetList) (result utils.Result)
}
