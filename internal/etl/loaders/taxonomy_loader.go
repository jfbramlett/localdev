package loaders

import (
	"context"
	"fmt"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl/datbase"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
)

func NewTaxonomyLoader(srcDB datbase.Database, tgtDB datbase.Database) Loader {
	return &TaxonomyLoader{
		srcDB: srcDB,
		tgtDB: tgtDB,
	}
}

type TaxonomyLoader struct {
	srcDB datbase.Database
	tgtDB datbase.Database
}

func (t *TaxonomyLoader) Load(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)
	logger.Info("starting taxonomy load")
	defer logger.Info("completed taxonomy load")
	logger.Info("deleting taxonomy tables")

	if err := DeleteTables(ctx, t.tgtDB, "tagged_entities", "tag_path_enums", "tag_paths", "tag_trees",
		"tags", "taxonomies", "taxonomy_types", "sync_versions", "tag_statuses"); err != nil {
		return err
	}

	logger.Info("loading tag statuses")
	if err := t.loadTagStatuses(ctx); err != nil {
		return err
	}

	logger.Info("loading sync versions")
	if err := t.loadSyncVersions(ctx); err != nil {
		return err
	}

	logger.Info("loading taxonomy types")
	if err := t.loadTaxonomyTypes(ctx); err != nil {
		return err
	}

	logger.Info("loading taxonomies")
	if err := t.loadTaxonomies(ctx); err != nil {
		return err
	}

	logger.Info("loading tags")
	if err := t.loadTags(ctx); err != nil {
		return err
	}

	logger.Info("loading tag trees")
	if err := t.loadTagTrees(ctx); err != nil {
		return err
	}

	logger.Info("loading tag paths")
	if err := t.loadTagPaths(ctx); err != nil {
		return err
	}

	logger.Info("loading tag path enums")
	if err := t.loadTagPathEnums(ctx); err != nil {
		return err
	}

	logger.Info("loading tagged entities")
	if err := t.loadTaggedEntities(ctx); err != nil {
		return err
	}

	return nil
}

func (t *TaxonomyLoader) loadTagStatuses(ctx context.Context) error {
	return LoadAll(ctx, t.srcDB, t.tgtDB, "tag_statuses")
}

func (t *TaxonomyLoader) loadSyncVersions(ctx context.Context) error {
	return LoadAll(ctx, t.srcDB, t.tgtDB, "sync_versions")
}

func (t *TaxonomyLoader) loadTaxonomyTypes(ctx context.Context) error {
	return LoadAll(ctx, t.srcDB, t.tgtDB, "taxonomy_types")
}

func (t *TaxonomyLoader) loadTaxonomies(ctx context.Context) error {
	return LoadAll(ctx, t.srcDB, t.tgtDB, "taxonomies")
}

func (t *TaxonomyLoader) loadTags(ctx context.Context) error {
	return LoadAll(ctx, t.srcDB, t.tgtDB, "tags")
}

func (t *TaxonomyLoader) loadTagTrees(ctx context.Context) error {
	return LoadAll(ctx, t.srcDB, t.tgtDB, "tag_trees")
}

func (t *TaxonomyLoader) loadTagPaths(ctx context.Context) error {
	return LoadAll(ctx, t.srcDB, t.tgtDB, "tag_paths")
}

func (t *TaxonomyLoader) loadTagPathEnums(ctx context.Context) error {
	return LoadAll(ctx, t.srcDB, t.tgtDB, "tag_path_enums")
}

func (t *TaxonomyLoader) loadTaggedEntities(ctx context.Context) error {
	srcIDs, err := t.tgtDB.QueryIDs(ctx, "SELECT uuid FROM assets")
	if err != nil {
		return err
	}

	batches := BatchIds(srcIDs)
	for _, batch := range batches {
		stmt := fmt.Sprintf("SELECT * from tagged_entities WHERE entity_type = 'asset' AND entity_uuid in (%s)", datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, t.srcDB, t.tgtDB, "tagged_entities", stmt, batch...); err != nil {
			return err
		}
	}

	return nil
}
