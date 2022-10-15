package usecases_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	modelsRoles "engine/bin/modules/roles/models/domain"
	mocksRoles "engine/bin/modules/roles/models/mocks"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/modules/users/models/mocks"
	"engine/bin/modules/users/usecases"
	"engine/bin/pkg/token"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

func TestUsersCommandUsecase(t *testing.T) {
	suite.Run(t, new(UsersCommandUsecaseTest))
}

type UsersCommandUsecaseTest struct {
	suite.Suite
	usecases            *usecases.UsersCommandUsecase
	usersPostgreCommand *mocks.UsersPostgreCommand
	usersPostgreQuery   *mocks.UsersPostgreQuery
	rolesPostgreQuery   *mocksRoles.RolesPostgreQuery
}

func (s *UsersCommandUsecaseTest) SetupTest() {
	s.rolesPostgreQuery = new(mocksRoles.RolesPostgreQuery)
	s.usersPostgreQuery = new(mocks.UsersPostgreQuery)
	s.usersPostgreCommand = new(mocks.UsersPostgreCommand)
	s.usecases = usecases.NewUsersCommandUsecase(s.usersPostgreCommand, s.usersPostgreQuery, s.rolesPostgreQuery)
}

func (s *UsersCommandUsecaseTest) TestLogin() {
	params := models.ReqLogin{
		Username: "testuser",
		Password: "testpassword",
		Duration: 3600,
	}
	mockUserResp := &models.Users{
		ID:       uuid.NewString(),
		RoleID:   "1",
		Username: "testuser",
	}
	mockRolesResp := &modelsRoles.Roles{
		ID:   "1",
		Name: "Admin",
	}
	passHas, _ := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	mockUserResp.Password = string(passHas)
	s.Run("Success", func() {
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(mockUserResp, nil)
		s.rolesPostgreQuery.On("FindOneByID", mock.Anything, mockUserResp.RoleID).Return(mockRolesResp, nil)
		res := s.usecases.Login(context.TODO(), params)
		s.Nil(res.Error)
		s.usersPostgreQuery.AssertExpectations(s.T())
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(nil, driver.ErrBadConn)
		res := s.usecases.Login(context.TODO(), params)
		s.NotNil(res.Error)
		s.usersPostgreQuery.AssertCalled(s.T(), "FindOneByUsername", mock.Anything, params.Username)
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(mockUserResp, nil)
		s.rolesPostgreQuery.On("FindOneByID", mock.Anything, mockUserResp.RoleID).Return(nil, driver.ErrBadConn)
		res := s.usecases.Login(context.TODO(), params)
		s.NotNil(res.Error)
		s.rolesPostgreQuery.AssertCalled(s.T(), "FindOneByID", mock.Anything, mockUserResp.RoleID)
	})
	s.Run("username-not-found", func() {
		s.SetupTest()
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(nil, sql.ErrNoRows)
		res := s.usecases.Login(context.TODO(), params)
		s.NotNil(res.Error)
		s.usersPostgreQuery.AssertCalled(s.T(), "FindOneByUsername", mock.Anything, params.Username)
	})
	s.Run("roles-not-found", func() {
		s.SetupTest()
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(mockUserResp, nil)
		s.rolesPostgreQuery.On("FindOneByID", mock.Anything, mockUserResp.RoleID).Return(nil, sql.ErrNoRows)
		res := s.usecases.Login(context.TODO(), params)
		s.NotNil(res.Error)
		s.rolesPostgreQuery.AssertCalled(s.T(), "FindOneByID", mock.Anything, mockUserResp.RoleID)
	})
	s.Run("error-token-generate", func() {
		s.SetupTest()
		params.Duration = 10000000000
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(mockUserResp, nil)
		s.rolesPostgreQuery.On("FindOneByID", mock.Anything, mockUserResp.RoleID).Return(mockRolesResp, nil)
		res := s.usecases.Login(context.TODO(), params)
		s.NotNil(res.Error)
		s.usersPostgreQuery.AssertExpectations(s.T())
	})
	s.Run("invalid-username-password", func() {
		s.SetupTest()
		params.Password = "wrongpassword"
		params.Duration = 3600
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(mockUserResp, nil)
		s.rolesPostgreQuery.On("FindOneByID", mock.Anything, mockUserResp.RoleID).Return(mockRolesResp, nil)
		res := s.usecases.Login(context.TODO(), params)
		s.NotNil(res.Error)
		s.usersPostgreQuery.AssertCalled(s.T(), "FindOneByUsername", mock.Anything, params.Username)
	})
}

func (s *UsersCommandUsecaseTest) TestRegisterUser() {
	params := models.ReqUser{
		Username: "usertest",
		Password: "usertest",
		RoleID:   "1",
		Opts: token.Claim{
			UserID:        "38be333f-4e7e-471e-91b8-057ae7b45563",
			RoleName:      "Admin",
			Authorization: "6c438854-4df4-48a9-bc64-3ac8802e1695",
		},
	}
	s.Run("Success", func() {
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(nil, nil)
		s.usersPostgreCommand.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Users")).Return(nil)
		res := s.usecases.RegisterUser(context.TODO(), params)
		s.Nil(res.Error)
		s.usersPostgreQuery.AssertExpectations(s.T())
	})
	s.Run("existing-username", func() {
		s.SetupTest()
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(&models.Users{ID: uuid.NewString(), Username: params.Username}, nil).Once()
		res := s.usecases.RegisterUser(context.TODO(), params)
		s.NotNil(res.Error)
		s.usersPostgreQuery.AssertCalled(s.T(), "FindOneByUsername", mock.Anything, params.Username)
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(nil, driver.ErrBadConn)
		res := s.usecases.RegisterUser(context.TODO(), params)
		s.NotNil(res.Error)
		s.usersPostgreQuery.AssertCalled(s.T(), "FindOneByUsername", mock.Anything, params.Username)
	})
	s.Run("Error", func() {
		s.SetupTest()
		s.usersPostgreQuery.On("FindOneByUsername", mock.Anything, params.Username).Return(nil, nil)
		s.usersPostgreCommand.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Users")).Return(errors.New("Failed insert user"))
		res := s.usecases.RegisterUser(context.TODO(), params)
		s.NotNil(res.Error)
		s.usersPostgreCommand.AssertCalled(s.T(), "InsertOne", mock.Anything, mock.AnythingOfType("*models.Users"))
	})
	s.Run("unauthorized", func() {
		s.SetupTest()
		params.Opts.RoleName = "operator"
		res := s.usecases.RegisterUser(context.TODO(), params)
		s.NotNil(res.Error)
		s.usersPostgreCommand.AssertNotCalled(s.T(), "InsertOne")
	})
}
