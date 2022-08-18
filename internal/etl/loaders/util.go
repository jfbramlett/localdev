package loaders

import (
	"context"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl/datbase"
)

type Loader interface {
	Load(ctx context.Context) error
}

func LoadAll(ctx context.Context, srcDB datbase.Database, tgtDB datbase.Database, table string) error {
	queryResult, err := srcDB.QueryAll(ctx, table)
	if err != nil {
		return err
	}

	err = tgtDB.Insert(ctx, table, queryResult.Cols, queryResult.Vals)
	if err != nil {
		return err
	}
	return nil
}

func LoadByQuery(ctx context.Context, srcDB datbase.Database, tgtDB datbase.Database, table string, query string, args ...interface{}) error {
	var queryResult *datbase.QueryResult
	var err error
	if len(args) > 0 {
		queryResult, err = srcDB.Query(ctx, query, args...)
	} else {
		queryResult, err = srcDB.Query(ctx, query)
	}
	if err != nil {
		return err
	}

	err = tgtDB.Insert(ctx, table, queryResult.Cols, queryResult.Vals)
	if err != nil {
		return err
	}
	return nil
}

func LoadByQueryIgnoreCol(ctx context.Context, srcDB datbase.Database, tgtDB datbase.Database, table string, colsToIgnore []string, query string, args ...interface{}) error {
	var queryResult *datbase.QueryResult
	var err error
	if len(args) > 0 {
		queryResult, err = srcDB.Query(ctx, query, args...)
	} else {
		queryResult, err = srcDB.Query(ctx, query)
	}
	if err != nil {
		return err
	}

	for _, c := range colsToIgnore {
		idx := -1
		for i, col := range queryResult.Cols {
			if col == c {
				idx = i
				break
			}
		}
		if idx >= 0 {
			if idx == 0 {
				queryResult.Cols = queryResult.Cols[1:]
				for i, v := range queryResult.Vals {
					queryResult.Vals[i] = v[1:]
				}
			} else {
				queryResult.Cols = append(queryResult.Cols[:idx], queryResult.Cols[idx+1:]...)
				for i, v := range queryResult.Vals {
					queryResult.Vals[i] = append(v[:idx], v[idx+1:]...)
				}
			}
		}
	}

	err = tgtDB.Insert(ctx, table, queryResult.Cols, queryResult.Vals)
	if err != nil {
		return err
	}
	return nil
}

func BatchIds(vals []interface{}) [][]interface{} {
	var batches [][]interface{}

	chunkSize := 500

	for i := 0; i < len(vals); i += chunkSize {
		end := i + chunkSize

		if end > len(vals) {
			end = len(vals)
		}

		batches = append(batches, vals[i:end])
	}

	return batches
}

func DeleteTables(ctx context.Context, tgtDB datbase.Database, tables ...string) error {
	for _, t := range tables {
		if err := tgtDB.DeleteAll(ctx, t); err != nil {
			return err
		}
	}
	return nil
}
