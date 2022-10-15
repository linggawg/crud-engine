package queries_test

import (
	"context"
	"engine/bin/modules/resources-mapping/repositories/queries"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestResourceMappingPostgreQuery(t *testing.T) {
	suite.Run(t, new(ResourceMappingPostgreQueryTest))
}

type ResourceMappingPostgreQueryTest struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository queries.ResourcesMappingPostgre
}

func (s *ResourceMappingPostgreQueryTest) SetupTest() {
	var (
		err error
	)
	mockDb, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	db := sqlx.NewDb(mockDb, "sqlmock")
	s.mock = mock
	s.repository = queries.NewResourcesMappingQuery(db)
}

func (s *ResourceMappingPostgreQueryTest) TestSelectInformationSchema() {
	servicesId := "168a3654-aa3d-4f2e-a5b7-745b20d2da78"
	rows := sqlmock.NewRows([]string{"id", "service_id", "source_origin", "source_alias"})
	rows.AddRow("2b99f877-a089-484f-ade9-7660419f1e05", "168a3654-aa3d-4f2e-a5b7-745b20d2da78", "id", "apps_id")
	rows.AddRow("37a758eb-e17f-41dc-a4b3-06b27a391385", "168a3654-aa3d-4f2e-a5b7-745b20d2da78", "name", "apps_name")
	query := regexp.QuoteMeta("SELECT id, service_id, source_origin, source_alias FROM resources_mapping WHERE service_id = $1;")
	s.Run("Success", func() {
		s.mock.ExpectQuery(query).WithArgs(servicesId).WillReturnRows(rows)
		data, err := s.repository.FindByServiceId(context.TODO(), servicesId)
		s.NoError(err)
		s.Len(data, 2)
	})
	s.Run("RowsWillBeClosed", func() {
		s.mock.ExpectQuery(query).WithArgs(servicesId).RowsWillBeClosed()
		data, err := s.repository.FindByServiceId(context.TODO(), servicesId)
		s.Error(err)
		s.Empty(data)
	})
}

func TestFindByServiceId(t *testing.T) {
	mockDb, _, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db := sqlx.NewDb(mockDb, "sqlmock")

	db.Close()
	entity, err := queries.NewResourcesMappingQuery(db).FindByServiceId(context.TODO(), "")
	assert.Empty(t, entity)
	assert.Error(t, err)
	assert.EqualError(t, err, "error establishing a database connection")
}
