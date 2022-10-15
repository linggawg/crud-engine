package usecases_test

import (
	"context"
	"database/sql"
	models "engine/bin/modules/services/models/domain"
	mocks2 "engine/bin/modules/services/models/mocks"
	"engine/bin/modules/services/usecases"
	"engine/bin/pkg/token"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestUsersCommandUsecase(t *testing.T) {
	suite.Run(t, new(ServicesCommandUsecaseTest))
}

type ServicesCommandUsecaseTest struct {
	suite.Suite
	usecases               *usecases.ServicesCommandUsecase
	servicesPostgreCommand *mocks2.ServicesPostgreCommand
	servicesPostgreQuery   *mocks2.ServicesPostgreQuery
}

func (s *ServicesCommandUsecaseTest) SetupTest() {
	s.servicesPostgreCommand = new(mocks2.ServicesPostgreCommand)
	s.servicesPostgreQuery = new(mocks2.ServicesPostgreQuery)
	s.usecases = usecases.NewServicesCommandUsecase(s.servicesPostgreCommand, s.servicesPostgreQuery)
}

func (s *ServicesCommandUsecaseTest) TestDeleteByServiceUrl() {
	payload := models.ServicesRequest{
		ServiceUrl: "users",
		Opts: token.Claim{
			UserID:        "",
			RoleName:      "Admin",
			Authorization: "",
		},
	}
	s.Run("Success", func() {
		s.servicesPostgreCommand.On("DeleteByServiceUrl", mock.Anything, payload.ServiceUrl).Return(nil)
		res := s.usecases.DeleteByServiceUrl(context.TODO(), payload)
		s.Nil(res.Error)
		s.servicesPostgreCommand.AssertExpectations(s.T())
	})
	s.Run("Error", func() {
		s.SetupTest()
		s.servicesPostgreCommand.On("DeleteByServiceUrl", mock.Anything, payload.ServiceUrl).Return(sql.ErrNoRows)
		res := s.usecases.DeleteByServiceUrl(context.TODO(), payload)
		s.NotNil(res.Error)
		s.servicesPostgreCommand.AssertExpectations(s.T())
	})
	s.Run("Unauthorized", func() {
		s.SetupTest()
		payload.Opts.RoleName = "operator"
		res := s.usecases.DeleteByServiceUrl(context.TODO(), payload)
		s.NotNil(res.Error)
		s.servicesPostgreCommand.AssertNotCalled(s.T(), "DeleteByServiceUrl")
	})
}
