package usecases_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	models2 "engine/bin/modules/services/models/domain"
	mocks2 "engine/bin/modules/services/models/mocks"
	models "engine/bin/modules/users-services/models/domain"
	"engine/bin/modules/users-services/models/mocks"
	"engine/bin/modules/users-services/usecases"
	"engine/bin/pkg/token"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestUsersCommandUsecase(t *testing.T) {
	suite.Run(t, new(UsersServicesCommandUsecaseTest))
}

type UsersServicesCommandUsecaseTest struct {
	suite.Suite
	usersServicesPostgreCommand *mocks.UsersServicesPostgreCommand
	usersServicesPostgreQuery   *mocks.UsersServicesPostgreQuery
	servicesPostgreCommand      *mocks2.ServicesPostgreCommand
	servicesPostgreQuery        *mocks2.ServicesPostgreQuery
	usecases                    *usecases.UsersServicesCommandUsecase
}

func (s *UsersServicesCommandUsecaseTest) SetupTest() {
	s.usersServicesPostgreCommand = new(mocks.UsersServicesPostgreCommand)
	s.usersServicesPostgreQuery = new(mocks.UsersServicesPostgreQuery)
	s.servicesPostgreCommand = new(mocks2.ServicesPostgreCommand)
	s.servicesPostgreQuery = new(mocks2.ServicesPostgreQuery)
	s.usecases = usecases.NewUsersServicesCommandUsecase(s.usersServicesPostgreCommand, s.usersServicesPostgreQuery, s.servicesPostgreCommand, s.servicesPostgreQuery)
}

func (s *UsersServicesCommandUsecaseTest) TestInsertAllByServices() {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	payload := models.UsersServicesRequest{
		ServiceUrl: "users",
		UserID:     "38be333f-4e7e-471e-91b8-057ae7b45563",
		DbID:       "6c438854-4df4-48a9-bc64-3ac8802e1695",
		Opts: token.Claim{
			UserID:        "38be333f-4e7e-471e-91b8-057ae7b45563",
			RoleName:      "Admin",
			Authorization: "6c438854-4df4-48a9-bc64-3ac8802e1695",
		},
	}
	services := &models2.Services{
		DbID:       payload.DbID,
		QueryID:    nil,
		ServiceUrl: &payload.ServiceUrl,
	}
	usersServices := &models.UsersServices{
		UserID:    payload.UserID,
		ServiceID: services.ID,
	}
	s.Run("success", func() {
		for _, method := range methods {
			s.servicesPostgreQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, method).Return(nil, nil)
			s.servicesPostgreCommand.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Services")).Return(nil)
			s.usersServicesPostgreQuery.On("FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, payload.UserID, method).Return(nil, nil)
			s.usersServicesPostgreCommand.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.UsersServices")).Return(nil)
		}
		res := s.usecases.InsertAllByServices(context.TODO(), payload)
		s.Nil(res.Error)
		s.servicesPostgreQuery.AssertExpectations(s.T())
	})

	s.Run("success-data-exists", func() {
		s.SetupTest()
		for _, method := range methods {
			s.servicesPostgreQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, method).Return(services, nil)
			s.usersServicesPostgreQuery.On("FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, payload.UserID, method).Return(usersServices, nil)
		}
		res := s.usecases.InsertAllByServices(context.TODO(), payload)
		s.Nil(res.Error)
		s.servicesPostgreQuery.AssertExpectations(s.T())
		s.servicesPostgreCommand.AssertNotCalled(s.T(), "InsertOne")
		s.usersServicesPostgreQuery.AssertNotCalled(s.T(), "InsertOne")
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		for _, method := range methods {
			s.servicesPostgreQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, method).Return(nil, driver.ErrBadConn)
		}
		res := s.usecases.InsertAllByServices(context.TODO(), payload)
		s.NotNil(res.Error)
		s.servicesPostgreQuery.AssertNotCalled(s.T(), "InsertOne")
		s.usersServicesPostgreQuery.AssertNotCalled(s.T(), "FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull")
		s.usersServicesPostgreCommand.AssertNotCalled(s.T(), "InsertOne")
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		for _, method := range methods {
			s.servicesPostgreQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, method).Return(nil, nil)
			s.servicesPostgreCommand.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Services")).Return(nil)
			s.usersServicesPostgreQuery.On("FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, payload.UserID, method).Return(nil, driver.ErrBadConn)
		}
		res := s.usecases.InsertAllByServices(context.TODO(), payload)
		s.NotNil(res.Error)
		s.usersServicesPostgreCommand.AssertNotCalled(s.T(), "InsertOne")
	})
	s.Run("error-services-InsertOne", func() {
		s.SetupTest()
		for _, method := range methods {
			s.servicesPostgreQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, method).Return(nil, nil)
			s.servicesPostgreCommand.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Services")).Return(sql.ErrTxDone)
		}
		res := s.usecases.InsertAllByServices(context.TODO(), payload)
		s.NotNil(res.Error)
		s.usersServicesPostgreQuery.AssertNotCalled(s.T(), "FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull")
		s.usersServicesPostgreCommand.AssertNotCalled(s.T(), "InsertOne")
	})
	s.Run("error-user-services-InsertOne", func() {
		s.SetupTest()
		for _, method := range methods {
			s.servicesPostgreQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, method).Return(nil, nil)
			s.servicesPostgreCommand.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Services")).Return(nil)
			s.usersServicesPostgreQuery.On("FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull", mock.Anything, payload.ServiceUrl, payload.UserID, method).Return(nil, nil)
			s.usersServicesPostgreCommand.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.UsersServices")).Return(sql.ErrTxDone)
		}
		res := s.usecases.InsertAllByServices(context.TODO(), payload)
		s.NotNil(res.Error)
	})
	s.Run("unauthorized", func() {
		s.SetupTest()
		payload.Opts.RoleName = "operator"
		res := s.usecases.InsertAllByServices(context.TODO(), payload)
		s.NotNil(res.Error)
		s.servicesPostgreQuery.AssertNotCalled(s.T(), "FindOneByServiceUrlAndMethodAndQueryIsNull")
	})
}

func (s *UsersServicesCommandUsecaseTest) TestDeleteByUsersIdAndServiceUrl() {
	payload := models.UsersServicesRequest{
		ServiceUrl: "users",
		UserID:     "38be333f-4e7e-471e-91b8-057ae7b45563",
		DbID:       "6c438854-4df4-48a9-bc64-3ac8802e1695",
		Opts: token.Claim{
			UserID:        "38be333f-4e7e-471e-91b8-057ae7b45563",
			RoleName:      "Admin",
			Authorization: "6c438854-4df4-48a9-bc64-3ac8802e1695",
		},
	}
	s.Run("success", func() {
		s.usersServicesPostgreCommand.On("DeleteByUsersIdAndServiceUrl", mock.Anything, payload.UserID, payload.ServiceUrl).Return(nil)
		res := s.usecases.DeleteByUsersIdAndServiceUrl(context.TODO(), payload)
		s.Nil(res.Error)
		s.usersServicesPostgreCommand.AssertExpectations(s.T())
	})
	s.Run("error-delete", func() {
		s.SetupTest()
		s.usersServicesPostgreCommand.On("DeleteByUsersIdAndServiceUrl", mock.Anything, payload.UserID, payload.ServiceUrl).Return(sql.ErrTxDone)
		res := s.usecases.DeleteByUsersIdAndServiceUrl(context.TODO(), payload)
		s.NotNil(res.Error)
		s.usersServicesPostgreCommand.AssertExpectations(s.T())
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		s.usersServicesPostgreCommand.On("DeleteByUsersIdAndServiceUrl", mock.Anything, payload.UserID, payload.ServiceUrl).Return(driver.ErrBadConn)
		res := s.usecases.DeleteByUsersIdAndServiceUrl(context.TODO(), payload)
		s.NotNil(res.Error)
		s.usersServicesPostgreCommand.AssertNotCalled(s.T(), "DeleteByUsersIdAndServiceUrl")
	})
	s.Run("unauthorized", func() {
		s.SetupTest()
		payload.Opts.RoleName = "operator"
		res := s.usecases.DeleteByUsersIdAndServiceUrl(context.TODO(), payload)
		s.NotNil(res.Error)
		s.usersServicesPostgreCommand.AssertNotCalled(s.T(), "DeleteByUsersIdAndServiceUrl")
	})
}
