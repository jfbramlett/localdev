package loaders

import (
	"context"
	"fmt"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl/datbase"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
)

const (
	maxPublished   = 20
	maxUnpublished = 10
)

func NewCatalogLoader(srcDB datbase.Database, tgtDB datbase.Database) Loader {
	return &CatalogLoader{
		srcDB: srcDB,
		tgtDB: tgtDB,
	}
}

type CatalogLoader struct {
	srcDB datbase.Database
	tgtDB datbase.Database
}

func (c *CatalogLoader) Load(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)
	logger.Info("starting catalog load")
	defer logger.Info("completed catalog load")
	logger.Info("deleting catalog tables")
	if err := DeleteTables(ctx, c.tgtDB, "asset_relationships", "asset_devices", "asset_errors", "asset_files",
		"asset_histories", "asset_instructors", "asset_jobs", "asset_price_histories", "asset_prices",
		"samples", "lessons", "videos", "midis", "assets",
		"providers", "device_versions", "devices", "marketplace_statuses",
		"asset_categories", "asset_types", "asset_statuses", "activity_types", "asset_data_sources", "asset_error_levels",
		"asset_error_types", "asset_file_types", "device_types", "instructors", "provider_statuses"); err != nil {
		return err
	}

	logger.Info("loading asset statuses")
	if err := c.loadAssetStatuses(ctx); err != nil {
		return err
	}

	logger.Info("loading asset types")
	if err := c.loadAssetTypes(ctx); err != nil {
		return err
	}

	logger.Info("loading asset categories")
	if err := c.loadAssetCategories(ctx); err != nil {
		return err
	}

	logger.Info("loading marketplace statuses")
	if err := c.loadMarketplaceStatuses(ctx); err != nil {
		return err
	}

	logger.Info("loading activity types")
	if err := c.loadActivityTypes(ctx); err != nil {
		return err
	}

	logger.Info("loading asset error levels")
	if err := c.loadAssetErrorLevels(ctx); err != nil {
		return err
	}

	logger.Info("loading asset error types")
	if err := c.loadAssetErrorTypes(ctx); err != nil {
		return err
	}

	logger.Info("loading asset file types")
	if err := c.loadAssetFileTypes(ctx); err != nil {
		return err
	}

	logger.Info("loading device types")
	if err := c.loadDeviceTypes(ctx); err != nil {
		return err
	}

	logger.Info("loading instructors")
	if err := c.loadInstructors(ctx); err != nil {
		return err
	}

	logger.Info("loading provider statuses")
	if err := c.loadProviderStatuses(ctx); err != nil {
		return err
	}

	logger.Info("loading providers")
	if err := c.loadProviders(ctx); err != nil {
		return err
	}

	logger.Info("loading packs")
	if err := c.loadPacks(ctx); err != nil {
		return err
	}

	logger.Info("loading samples")
	if err := c.loadSamples(ctx); err != nil {
		return err
	}

	logger.Info("loading lessons")
	if err := c.loadLessons(ctx); err != nil {
		return err
	}

	logger.Info("loading videos")
	if err := c.loadVideos(ctx); err != nil {
		return err
	}

	logger.Info("loading midis")
	if err := c.loadMidis(ctx); err != nil {
		return err
	}

	logger.Info("loading presets")
	if err := c.loadPresets(ctx); err != nil {
		return err
	}

	logger.Info("loading plugins")
	if err := c.loadPlugins(ctx); err != nil {
		return err
	}

	logger.Info("loading devices")
	if err := c.loadDevices(ctx); err != nil {
		return err
	}

	logger.Info("loading device versions")
	if err := c.loadDeviceVersions(ctx); err != nil {
		return err
	}

	logger.Info("loading asset devices")
	if err := c.loadAssetDevices(ctx); err != nil {
		return err
	}

	logger.Info("loading asset errors")
	if err := c.loadAssetErrors(ctx); err != nil {
		return err
	}

	logger.Info("loading asset files")
	if err := c.loadAssetFiles(ctx); err != nil {
		return err
	}

	logger.Info("loading asset histories - skipped")
	//if err := c.loadAssetHistories(ctx); err != nil {
	//	return err
	//}

	logger.Info("loading asset instructors")
	if err := c.loadAssetInstructors(ctx); err != nil {
		return err
	}

	logger.Info("loading asset jobs")
	if err := c.loadAssetJobs(ctx); err != nil {
		return err
	}

	logger.Info("loading asset price histories")
	if err := c.loadAssetPriceHistories(ctx); err != nil {
		return err
	}

	logger.Info("loading asset prices")
	if err := c.loadAssetPrices(ctx); err != nil {
		return err
	}

	return nil
}

func (c *CatalogLoader) loadAssetStatuses(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "asset_statuses")
}

func (c *CatalogLoader) loadAssetTypes(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "asset_types")
}

func (c *CatalogLoader) loadAssetCategories(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "asset_categories")
}

func (c *CatalogLoader) loadMarketplaceStatuses(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "marketplace_statuses")
}

func (c *CatalogLoader) loadActivityTypes(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "activity_types")
}

func (c *CatalogLoader) loadAssetErrorLevels(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "asset_error_levels")
}

func (c *CatalogLoader) loadAssetErrorTypes(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "asset_error_types")
}

func (c *CatalogLoader) loadAssetFileTypes(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "asset_file_types")
}

func (c *CatalogLoader) loadDeviceTypes(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "device_types")
}

func (c *CatalogLoader) loadInstructors(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "instructors")
}

func (c *CatalogLoader) loadProviderStatuses(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "provider_statuses")
}

func (c *CatalogLoader) loadProviders(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "providers")
}

func (c *CatalogLoader) loadPacks(ctx context.Context) error {
	if err := c.loadAssets(ctx, "pack", "published", maxPublished); err != nil {
		return err
	}
	if err := c.loadAssets(ctx, "pack", "unpublished", maxUnpublished); err != nil {
		return err
	}
	return nil
}

func (c *CatalogLoader) loadLessons(ctx context.Context) error {
	if err := c.loadAssets(ctx, "lesson", "published", maxPublished); err != nil {
		return err
	}
	if err := c.loadAssets(ctx, "lesson", "unpublished", maxUnpublished); err != nil {
		return err
	}

	return c.loadAssetTypeData(ctx, "lessons", "lesson")
}

func (c *CatalogLoader) loadSamples(ctx context.Context) error {
	if err := c.loadRelatedAssets(ctx, "pack", "sample"); err != nil {
		return err
	}
	return c.loadAssetTypeRelatedData(ctx, "samples", "sample", "id")
}

func (c *CatalogLoader) loadVideos(ctx context.Context) error {
	err := c.loadRelatedAssets(ctx, "lesson", "video")
	if err != nil {
		return err
	}

	if err := c.loadAssets(ctx, "video", "unpublished", maxUnpublished); err != nil {
		return err
	}
	return c.loadAssetTypeData(ctx, "videos", "video")

}

func (c *CatalogLoader) loadMidis(ctx context.Context) error {
	if err := c.loadAssets(ctx, "midi", "published", maxPublished); err != nil {
		return err
	}
	if err := c.loadAssets(ctx, "midi", "unpublished", maxUnpublished); err != nil {
		return err
	}
	return c.loadAssetTypeRelatedData(ctx, "midis", "midi", "uuid")
}

func (c *CatalogLoader) loadPresets(ctx context.Context) error {
	if err := c.loadAssets(ctx, "preset", "published", maxPublished); err != nil {
		return err
	}
	return c.loadAssets(ctx, "preset", "unpublished", maxUnpublished)
}

func (c *CatalogLoader) loadPlugins(ctx context.Context) error {
	if err := c.loadAssets(ctx, "plugin", "published", maxPublished); err != nil {
		return err
	}
	return c.loadAssets(ctx, "plugin", "unpublished", maxUnpublished)
}

func (c *CatalogLoader) loadDevices(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "devices")
}

func (c *CatalogLoader) loadDeviceVersions(ctx context.Context) error {
	return LoadAll(ctx, c.srcDB, c.tgtDB, "device_versions")
}

func (c *CatalogLoader) loadAssetDevices(ctx context.Context) error {
	return c.loadAssetContent(ctx, "asset_devices", "uuid")
}

func (c *CatalogLoader) loadAssetErrors(ctx context.Context) error {
	return c.loadAssetContent(ctx, "asset_errors", "id")
}

func (c *CatalogLoader) loadAssetFiles(ctx context.Context) error {
	return c.loadAssetContent(ctx, "asset_files", "id")
}

//func (c *CatalogLoader) loadAssetHistories(ctx context.Context) error {
//	return c.loadAssetContent(ctx, "asset_histories", "id")
//}

func (c *CatalogLoader) loadAssetInstructors(ctx context.Context) error {
	return c.loadAssetContent(ctx, "asset_instructors", "uuid")
}

func (c *CatalogLoader) loadAssetJobs(ctx context.Context) error {
	return c.loadAssetContent(ctx, "asset_jobs", "uuid")
}

func (c *CatalogLoader) loadAssetPriceHistories(ctx context.Context) error {
	return c.loadAssetContent(ctx, "asset_price_histories", "uuid")
}

func (c *CatalogLoader) loadAssetPrices(ctx context.Context) error {
	return c.loadAssetContent(ctx, "asset_prices", "uuid")
}

func (c *CatalogLoader) loadAssetContent(ctx context.Context, table string, identifier string) error {
	ids, err := c.tgtDB.QueryIDs(ctx, fmt.Sprintf(`SELECT p.%s FROM assets p`, identifier))
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		stmt := fmt.Sprintf(`
SELECT *
FROM %s
WHERE asset_%s in (%s)`, table, identifier, datbase.ValuesToParameters(batch))
		if err := LoadByQuery(ctx, c.srcDB, c.tgtDB, table, stmt, batch...); err != nil {
			return err
		}
	}

	return nil

}

func (c *CatalogLoader) loadAssetTypeRelatedData(ctx context.Context, table string, assetType string, identifier string) error {
	ids, err := c.tgtDB.QueryIDs(ctx, fmt.Sprintf(`SELECT p.%s FROM assets p JOIN asset_types ats on p.asset_type_id = ats.id WHERE ats.slug = '%s'`, identifier, assetType))
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		stmt := fmt.Sprintf(`
SELECT *
FROM %s
WHERE asset_%s in (%s)`, table, identifier, datbase.ValuesToParameters(batch))
		if err := LoadByQuery(ctx, c.srcDB, c.tgtDB, table, stmt, batch...); err != nil {
			return err
		}
	}

	return nil

}

func (c *CatalogLoader) loadRelatedAssets(ctx context.Context, parentType string, childType string) error {
	ids, err := c.tgtDB.QueryIDs(ctx, fmt.Sprintf(`SELECT p.id FROM assets p JOIN asset_types ats on p.asset_type_id = ats.id WHERE ats.slug = '%s'`, parentType))
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		stmt := fmt.Sprintf(`
SELECT p.*
FROM assets p
JOIN asset_relationships ar on p.id = ar.child_asset_id
JOIN asset_types ats on p.asset_type_id = ats.id
WHERE ats.slug = '%s' AND ar.parent_asset_id in (%s)`, childType, datbase.ValuesToParameters(batch))

		if err = LoadByQuery(ctx, c.srcDB, c.tgtDB, "assets", stmt, batch...); err != nil {
			return err
		}

		stmt = fmt.Sprintf(`
SELECT *
FROM asset_relationships
WHERE parent_asset_id in (%s)`, datbase.ValuesToParameters(batch))

		if err = LoadByQuery(ctx, c.srcDB, c.tgtDB, "asset_relationships", stmt, batch...); err != nil {
			return err
		}
	}
	return nil
}

func (c *CatalogLoader) loadAssets(ctx context.Context, assetType string, assetStatus string, limit int) error {
	stmt := fmt.Sprintf(`
SELECT p.*
FROM assets p
JOIN asset_statuses ast on p.asset_status_id = ast.id
JOIN asset_types ats on p.asset_type_id = ats.id
WHERE ast.slug = '%s' AND ats.slug = '%s'
LIMIT %d`, assetStatus, assetType, limit)

	return LoadByQuery(ctx, c.srcDB, c.tgtDB, "assets", stmt)
}

func (c *CatalogLoader) loadAssetTypeData(ctx context.Context, table string, assetType string) error {
	ids, err := c.tgtDB.QueryIDs(ctx, fmt.Sprintf(`SELECT p.uuid FROM assets p JOIN asset_types ats on p.asset_type_id = ats.id WHERE ats.slug = '%s'`, assetType))
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {

		stmt := fmt.Sprintf(`
SELECT *
FROM %s
WHERE uuid in (%s)`, table, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, c.srcDB, c.tgtDB, table, stmt, batch...); err != nil {
			return err
		}
	}
	return nil
}
