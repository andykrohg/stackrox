// Code generated by pg-bindings generator. DO NOT EDIT.
package n41ton42

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_41_to_n_42_postgres_process_baselines/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_41_to_n_42_postgres_process_baselines/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNumWithoutPostgres() + 41,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNumWithoutPostgres()) + 42},
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving process_baselines from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.ProcessBaselinesSchema
	log       = loghelper.LogWrapper{}
)

func move(gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var processBaselines []*storage.ProcessBaseline
	err := walk(ctx, legacyStore, func(obj *storage.ProcessBaseline) error {
		processBaselines = append(processBaselines, obj)
		if len(processBaselines) == batchSize {
			if err := store.UpsertMany(ctx, processBaselines); err != nil {
				log.WriteToStderrf("failed to persist process_baselines to store %v", err)
				return err
			}
			processBaselines = processBaselines[:0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(processBaselines) > 0 {
		if err = store.UpsertMany(ctx, processBaselines); err != nil {
			log.WriteToStderrf("failed to persist process_baselines to store %v", err)
			return err
		}
	}
	return nil
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.ProcessBaseline) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
