package etl

import (
	"context"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl/datbase"
	loaders2 "github.com/splice/platform/localdev/v2/localdev/internal/etl/loaders"
)

type ManualDriver struct {
	srcDB datbase.Database
	tgtDB datbase.Database
}

func NewManualDriver(srcDSN string, tgtDSN string) (*ManualDriver, error) {
	srcDB, err := datbase.NewDatabase(srcDSN)
	if err != nil {
		return nil, err
	}
	tgtDB, err := datbase.NewDatabase(tgtDSN)
	if err != nil {
		return nil, err
	}
	return &ManualDriver{
		srcDB: srcDB,
		tgtDB: tgtDB,
	}, nil
}

func (d *ManualDriver) Run(ctx context.Context) error {
	loadersToRun := []loaders2.Loader{
		loaders2.NewCatalogLoader(d.srcDB, d.tgtDB),
		loaders2.NewMerchandisingLoader(d.srcDB, d.tgtDB),
		loaders2.NewTaxonomyLoader(d.srcDB, d.tgtDB),
		loaders2.NewAPILoader(d.srcDB, d.tgtDB),
		loaders2.NewUserLoader(d.srcDB, d.tgtDB),
	}

	if err := d.tgtDB.DisableForeignKeys(ctx); err != nil {
		return err
	}
	defer func() {
		_ = d.tgtDB.EnableForeignKeys(ctx)
	}()

	for _, load := range loadersToRun {
		if err := load.Load(ctx); err != nil {
			return err
		}
	}
	return nil
}
