// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package n19ton20

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_19_to_n_20_postgres_compliance_run_results/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_19_to_n_20_postgres_compliance_run_results/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"

	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	ctx         context.Context

	legacyDB   *rocksdb.RocksDB
	postgresDB *pghelper.TestPostgres
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")
	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	var err error
	s.legacyDB, err = rocksdb.NewTemp(s.T().Name())
	s.NoError(err)

	s.Require().NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.postgresDB = pghelper.ForT(s.T(), true)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestComplianceRunResultsMigration() {
	newStore := pgStore.New(s.postgresDB.Pool)
	legacyStore, err := legacy.New(s.legacyDB)
	s.NoError(err)

	// Prepare data and write to legacy DB
	var complianceRunResultss []*storage.ComplianceRunResults
	for i := 0; i < 200; i++ {
		complianceRunResults := &storage.ComplianceRunResults{}
		s.NoError(testutils.FullInit(complianceRunResults, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		complianceRunResultss = append(complianceRunResultss, complianceRunResults)
	}

	s.NoError(legacyStore.UpsertMany(s.ctx, complianceRunResultss))

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.Pool, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(complianceRunResultss), count)
	for _, complianceRunResults := range complianceRunResultss {
		fetched, exists, err := newStore.Get(s.ctx, complianceRunResults.GetRunMetadata().GetRunId())
		s.NoError(err)
		s.True(exists)
		s.Equal(complianceRunResults, fetched)
	}
}
