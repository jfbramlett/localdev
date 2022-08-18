package loaders

import (
	"context"
	"fmt"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl/datbase"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
)

func NewAPILoader(srcDB datbase.Database, tgtDB datbase.Database) Loader {
	return &APILoader{
		srcDB: srcDB,
		tgtDB: tgtDB,
	}
}

type APILoader struct {
	srcDB datbase.Database
	tgtDB datbase.Database
}

func (api *APILoader) Load(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)
	logger.Info("starting api load")
	defer logger.Info("completed api load")
	logger.Info("deleting api tables")

	if err := DeleteTables(ctx, api.tgtDB, "audio_tags", "bundled_gear", "cancellation_categories", "cancellation_ui_options",
		"file_types", "genre_locations", "genre_origins", "genre_influences", "genres", "genre_gear_connections", "gear_connections",
		"leasable_components", "leasable_gear_installation_details", "leasable_plugin_installers", "leasable_plugins",
		"plan_feature", "plan_group", "plans", "plugin_manufacturers", "publishing_event_types",
		"publishing_states", "submission_types", "world_library_tags", "desktop_versions", "sample_packs", "content_group_entries",
		"pack_extras", "pack_grouping_items", "premium_sample_replacement", "sample_pack_metadata",
		"synth_patch_collections", "companion_pack_groups", "companion_pack_group_memberships", "pack_groupings", "content_groups",
		"plugin_descriptions", "revision_plugins", "plugin_versions", "plugin_links", "plugin_installers", "premium_samples",
		"audio_samples", "sample_premium_sets", "audio_sample_mlt_samples", "premium_sample_audio_tags", "premium_presets",
		"world_library_sample_tags", "world_library_sample_attributes", "world_library_paths", "world_library_audio_science_tags",
		"world_library_audio_science_attributes", "world_library_path_tags", "sample_providers"); err != nil {
		return err
	}

	logger.Info("loading audio_tags")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "audio_tags"); err != nil {
		return err
	}

	logger.Info("loading bundled_gear")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "bundled_gear"); err != nil {
		return err
	}

	logger.Info("loading cancellation_categories")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "cancellation_categories"); err != nil {
		return err
	}

	logger.Info("loading cancellation_ui_options")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "cancellation_ui_options"); err != nil {
		return err
	}

	logger.Info("loading file_types")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "file_types"); err != nil {
		return err
	}

	logger.Info("loading genre_locations")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "genre_locations"); err != nil {
		return err
	}

	logger.Info("loading genre_origins")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "genre_origins"); err != nil {
		return err
	}

	logger.Info("loading genre_influences")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "genre_influences"); err != nil {
		return err
	}

	logger.Info("loading genres")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "genres"); err != nil {
		return err
	}

	logger.Info("loading genre_gear_connections")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "genre_gear_connections"); err != nil {
		return err
	}

	logger.Info("loading gear_connections")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "gear_connections"); err != nil {
		return err
	}

	logger.Info("loading leasable_components")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "leasable_components"); err != nil {
		return err
	}

	logger.Info("loading leasable_gear_installation_details")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "leasable_gear_installation_details"); err != nil {
		return err
	}

	logger.Info("loading leasable_plugin_installers")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "leasable_plugin_installers"); err != nil {
		return err
	}

	logger.Info("loading leasable_plugins")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "leasable_plugins"); err != nil {
		return err
	}

	logger.Info("loading plan_feature")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "plan_feature"); err != nil {
		return err
	}

	logger.Info("loading plan_group")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "plan_group"); err != nil {
		return err
	}

	logger.Info("loading plans")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "plans"); err != nil {
		return err
	}

	logger.Info("loading plugin_manufacturers")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "plugin_manufacturers"); err != nil {
		return err
	}

	logger.Info("loading publishing_event_types")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "publishing_event_types"); err != nil {
		return err
	}

	logger.Info("loading publishing_states")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "publishing_states"); err != nil {
		return err
	}

	logger.Info("loading submission_types")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "submission_types"); err != nil {
		return err
	}

	logger.Info("loading world_library_tags")
	if err := LoadAll(ctx, api.srcDB, api.tgtDB, "world_library_tags"); err != nil {
		return err
	}

	logger.Info("loading desktop_versions")
	stmt := `SELECT * FROM desktop_versions
				ORDER BY created_at desc LIMIT 50`
	if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "desktop_versions", stmt); err != nil {
		return err
	}

	logger.Info("loading sample_providers")
	if err := api.loadSampleProviders(ctx); err != nil {
		return err
	}

	logger.Info("loading packs")
	if err := api.loadPacks(ctx); err != nil {
		return err
	}

	logger.Info("loading samples")
	if err := api.loadSamples(ctx); err != nil {
		return err
	}

	logger.Info("loading presets")
	if err := api.loadPresets(ctx); err != nil {
		return err
	}

	logger.Info("loading plugins")
	if err := api.loadPlugins(ctx); err != nil {
		return err
	}

	logger.Info("loading world_samples")
	if err := api.loadWorldSamples(ctx); err != nil {
		return err
	}

	return nil
}

func (api *APILoader) loadPacks(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)

	ids, err := api.tgtDB.QueryIDs(ctx, `SELECT p.uuid FROM assets p JOIN asset_types ats on p.asset_type_id = ats.id WHERE ats.slug = 'pack'`)
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch sample_packs")
		stmt := fmt.Sprintf(`
SELECT *
FROM sample_packs
WHERE asset_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "sample_packs", stmt, batch...); err != nil {
			return err
		}
	}

	ids, err = api.tgtDB.QueryIDs(ctx, `SELECT p.uuid FROM sample_packs p `)
	if err != nil {
		return err
	}

	batchIDs = BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch content_group_entries")
		stmt := fmt.Sprintf(`
SELECT *
FROM content_group_entries
WHERE sample_pack_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "content_group_entries", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch pack_extras")
		stmt = fmt.Sprintf(`
SELECT *
FROM pack_extras
WHERE sample_pack_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "pack_extras", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch pack_grouping_items")
		stmt = fmt.Sprintf(`
SELECT *
FROM pack_grouping_items
WHERE sample_pack_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "pack_grouping_items", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch premium_sample_replacement")
		stmt = fmt.Sprintf(`
SELECT *
FROM premium_sample_replacement
WHERE pack_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "premium_sample_replacement", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch sample_pack_metadata")
		stmt = fmt.Sprintf(`
SELECT *
FROM sample_pack_metadata
WHERE sample_pack_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "sample_pack_metadata", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch synth_patch_collections")
		stmt = fmt.Sprintf(`
SELECT *
FROM synth_patch_collections
WHERE sample_pack_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "synth_patch_collections", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch companion_pack_groups")
		stmt = fmt.Sprintf(`
SELECT *
FROM companion_pack_groups
WHERE parent_pack_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "companion_pack_groups", stmt, batch...); err != nil {
			return err
		}

	}

	ids, err = api.tgtDB.QueryIDs(ctx, `SELECT p.uuid FROM companion_pack_groups p `)
	if err != nil {
		return err
	}

	batchIDs = BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch companion_pack_group_memberships")
		stmt := fmt.Sprintf(`
SELECT *
FROM companion_pack_group_memberships
WHERE companion_pack_group_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "companion_pack_group_memberships", stmt, batch...); err != nil {
			return err
		}
	}

	ids, err = api.tgtDB.QueryIDs(ctx, `SELECT pack_grouping_uuid FROM pack_grouping_items p`)
	if err != nil {
		return err
	}

	batchIDs = BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch pack_groupings")
		stmt := fmt.Sprintf(`
SELECT *
FROM pack_groupings
WHERE uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "pack_groupings", stmt, batch...); err != nil {
			return err
		}
	}

	ids, err = api.tgtDB.QueryIDs(ctx, `SELECT uuid FROM content_group_entries p`)
	if err != nil {
		return err
	}

	batchIDs = BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch content_groups")
		stmt := fmt.Sprintf(`
SELECT *
FROM content_groups
WHERE uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "content_groups", stmt, batch...); err != nil {
			return err
		}
	}

	return nil
}

func (api *APILoader) loadPlugins(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)

	if err := LoadByQueryIgnoreCol(ctx, api.srcDB, api.tgtDB, "plugin_descriptions", []string{"pushable"}, `select * from plugin_descriptions where published = 1 AND purchasable = 1`); err != nil {
		return err
	}

	ids, err := api.tgtDB.QueryIDs(ctx, `SELECT p.id FROM plugin_descriptions p`)
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch revision_plugins")
		stmt := fmt.Sprintf(`
SELECT *
FROM revision_plugins
WHERE plugin_description_id in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "revision_plugins", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch plugin_versions")
		stmt = fmt.Sprintf(`
SELECT *
FROM plugin_versions
WHERE plugin_description_id in (%s)`, datbase.ValuesToParameters(batch))
		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "plugin_versions", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch plugin_links")
		stmt = fmt.Sprintf(`
SELECT *
FROM plugin_links
WHERE plugin_description_id in (%s)`, datbase.ValuesToParameters(batch))
		if err := LoadByQueryIgnoreCol(ctx, api.srcDB, api.tgtDB, "plugin_links", []string{"screenshot"}, stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch plugin_installers")
		stmt = fmt.Sprintf(`
SELECT *
FROM plugin_installers
WHERE plugin_description_id in (%s)`, datbase.ValuesToParameters(batch))
		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "plugin_installers", stmt, batch...); err != nil {
			return err
		}

	}
	return nil
}

func (api *APILoader) loadSamples(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)

	ids, err := api.tgtDB.QueryIDs(ctx, `SELECT p.uuid FROM assets p JOIN asset_types ats on p.asset_type_id = ats.id WHERE ats.slug = 'sample'`)
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch premium_samples")
		stmt := fmt.Sprintf(`
SELECT *
FROM premium_samples
WHERE asset_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "premium_samples", stmt, batch...); err != nil {
			return err
		}
	}

	ids, err = api.tgtDB.QueryIDs(ctx, `SELECT p.file_hash FROM premium_samples p`)
	if err != nil {
		return err
	}
	batchIDs = BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch audio_samples")
		stmt := fmt.Sprintf(`
SELECT *
FROM audio_samples
WHERE file_hash in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQueryIgnoreCol(ctx, api.srcDB, api.tgtDB, "audio_samples", []string{"bit_depth"}, stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch sample_premium_sets")
		stmt = fmt.Sprintf(`
SELECT *
FROM sample_premium_sets
WHERE file_hash in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "sample_premium_sets", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch audio_sample_mlt_samples")
		stmt = fmt.Sprintf(`
SELECT *
FROM audio_sample_mlt_samples
WHERE child_sample_filehash in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "audio_sample_mlt_samples", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch premium_sample_audio_tags")
		stmt = fmt.Sprintf(`
SELECT *
FROM premium_sample_audio_tags
WHERE file_hash in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "premium_sample_audio_tags", stmt, batch...); err != nil {
			return err
		}
	}

	return nil
}

func (api *APILoader) loadPresets(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)

	ids, err := api.tgtDB.QueryIDs(ctx, `SELECT p.uuid FROM assets p JOIN asset_types ats on p.asset_type_id = ats.id WHERE ats.slug = 'preset'`)
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch premium_presets")
		stmt := fmt.Sprintf(`
SELECT *
FROM premium_presets
WHERE asset_uuid in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "premium_presets", stmt, batch...); err != nil {
			return err
		}
	}
	return nil
}

func (api *APILoader) loadWorldSamples(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)
	ids, err := api.tgtDB.QueryIDs(ctx, `SELECT p.id FROM audio_samples p`)
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch world_library_sample_tags")
		stmt := fmt.Sprintf(`
SELECT *
FROM world_library_sample_tags
WHERE audio_sample_id in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "world_library_sample_tags", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch world_library_sample_attributes")
		stmt = fmt.Sprintf(`
SELECT *
FROM world_library_sample_attributes
WHERE audio_sample_id in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "world_library_sample_attributes", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch world_library_paths")
		stmt = fmt.Sprintf(`
SELECT *
FROM world_library_paths
WHERE audio_sample_id in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "world_library_paths", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch world_library_audio_science_tags")
		stmt = fmt.Sprintf(`
SELECT *
FROM world_library_audio_science_tags
WHERE audio_sample_id in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "world_library_audio_science_tags", stmt, batch...); err != nil {
			return err
		}

		logger.Info("loading batch world_library_audio_science_attributes")
		stmt = fmt.Sprintf(`
SELECT *
FROM world_library_audio_science_attributes
WHERE audio_sample_id in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "world_library_audio_science_attributes", stmt, batch...); err != nil {
			return err
		}

	}

	ids, err = api.tgtDB.QueryIDs(ctx, `SELECT p.id FROM world_library_paths p`)
	if err != nil {
		return err
	}

	batchIDs = BatchIds(ids)
	for _, batch := range batchIDs {
		logger.Info("loading batch world_library_path_tags")
		stmt := fmt.Sprintf(`
SELECT *
FROM world_library_path_tags
WHERE world_library_path_id in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "world_library_path_tags", stmt, batch...); err != nil {
			return err
		}
	}
	return nil
}

func (api *APILoader) loadSampleProviders(ctx context.Context) error {
	ids, err := api.tgtDB.QueryIDs(ctx, `SELECT p.name FROM providers p`)
	if err != nil {
		return err
	}

	batchIDs := BatchIds(ids)
	for _, batch := range batchIDs {
		stmt := fmt.Sprintf(`
SELECT *
FROM sample_providers
WHERE name in (%s)`, datbase.ValuesToParameters(batch))

		if err := LoadByQuery(ctx, api.srcDB, api.tgtDB, "sample_providers", stmt, batch...); err != nil {
			return err
		}
	}

	return nil
}
