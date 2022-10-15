package usecases

import (
	"context"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type CommandUsecase interface {
	SetupDataset(ctx context.Context, dialect string, db *sqlx.DB) (result utils.Result)
	Insert(ctx context.Context, engineConfig models.EngineConfig, payload *models.EngineRequest) (result utils.Result)
	Update(ctx context.Context, engineConfig models.EngineConfig, payload *models.EngineRequest) (result utils.Result)
	Delete(ctx context.Context, engineConfig models.EngineConfig, payload *models.EngineRequest) (result utils.Result)
	Patch(ctx context.Context, engineConfig models.EngineConfig, payload *models.EngineRequest) (result utils.Result)
}

type QueryUsecase interface {
	Get(ctx context.Context, engineConfig models.EngineConfig, table string, payload *models.GetList) (result utils.Result)
}
