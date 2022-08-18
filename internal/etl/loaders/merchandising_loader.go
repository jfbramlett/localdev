package loaders

import (
	"context"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl/datbase"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
)

func NewMerchandisingLoader(srcDB datbase.Database, tgtDB datbase.Database) Loader {
	return &MerchandisingLoader{
		srcDB: srcDB,
		tgtDB: tgtDB,
	}
}

type MerchandisingLoader struct {
	srcDB datbase.Database
	tgtDB datbase.Database
}

func (m *MerchandisingLoader) Load(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)
	logger.Info("starting merchandising load")
	defer logger.Info("completed merchandising load")
	logger.Info("deleting merchandising tables")
	if err := DeleteTables(ctx, m.tgtDB, "page_content_blocks", "content_block_content_items", "content_items", "content_blocks",
		"pages", "content_types", "content_block_statuses"); err != nil {
		return err
	}

	logger.Info("loading content block statuses")
	if err := m.loadContentBlockStatuses(ctx); err != nil {
		return err
	}

	logger.Info("loading content types")
	if err := m.loadContentTypes(ctx); err != nil {
		return err
	}

	logger.Info("loading content pages")
	if err := m.loadPages(ctx); err != nil {
		return err
	}

	logger.Info("loading content blocks")
	if err := m.loadContentBlocks(ctx); err != nil {
		return err
	}

	logger.Info("loading content items")
	if err := m.loadContentItems(ctx); err != nil {
		return err
	}

	logger.Info("loading content block content items")
	if err := m.loadContentBlockContentItems(ctx); err != nil {
		return err
	}

	logger.Info("loading page content blocks")
	if err := m.loadPageContentBlocks(ctx); err != nil {
		return err
	}

	return nil
}

func (m *MerchandisingLoader) loadContentBlockStatuses(ctx context.Context) error {
	return LoadAll(ctx, m.srcDB, m.tgtDB, "content_block_statuses")
}

func (m *MerchandisingLoader) loadContentTypes(ctx context.Context) error {
	return LoadAll(ctx, m.srcDB, m.tgtDB, "content_types")
}

func (m *MerchandisingLoader) loadPages(ctx context.Context) error {
	return LoadAll(ctx, m.srcDB, m.tgtDB, "pages")
}

func (m *MerchandisingLoader) loadContentBlocks(ctx context.Context) error {
	return LoadAll(ctx, m.srcDB, m.tgtDB, "content_blocks")
}

func (m *MerchandisingLoader) loadContentItems(ctx context.Context) error {
	return LoadAll(ctx, m.srcDB, m.tgtDB, "content_items")
}

func (m *MerchandisingLoader) loadContentBlockContentItems(ctx context.Context) error {
	return LoadAll(ctx, m.srcDB, m.tgtDB, "content_block_content_items")
}

func (m *MerchandisingLoader) loadPageContentBlocks(ctx context.Context) error {
	return LoadAll(ctx, m.srcDB, m.tgtDB, "page_content_blocks")
}
