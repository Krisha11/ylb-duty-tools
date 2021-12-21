package collectors

import (
	"context"
	"database/sql"
	"fmt"

	common "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty"
)

// Index stores data whose consistency needs to be checked
type Index struct {
	baseTable  string
	indexTable string
	baseKeys   []string
	indexKeys  []string
}

var (
	INDEXES = []Index{
		Index{
			baseTable:  "healthcheck_targets",
			indexTable: "healthcheck_targets__healthcheck_id__addr_id__target_id__index",
			baseKeys:   []string{"id", "healthcheck_id", "addr_id"},
			indexKeys:  []string{"target_id", "healthcheck_id", "addr_id"},
		},
		Index{
			baseTable:  "healthcheck_targets",
			indexTable: "healthcheck_targets__zone_id__target_id__index",
			baseKeys:   []string{"id", "zone_id"},
			indexKeys:  []string{"target_id", "zone_id"},
		},
		Index{
			baseTable:  "healthcheck_group_healthchecks",
			indexTable: "healthcheck_group_healthchecks__healthcheck_id__group_external_id__index",
			baseKeys:   []string{"group_external_id", "healthcheck_id"},
			indexKeys:  []string{"group_external_id", "healthcheck_id"},
		},
		Index{
			baseTable:  "healthcheck_forwarding_rules",
			indexTable: "healthcheck_forwarding_rules_denormalization",
			baseKeys:   []string{"id"},
			indexKeys:  []string{"rule_id"},
		},
	}
)

// IndexDiffCollector is collector for diff and index commands.
type IndexDiffCollector struct {
	config   *common.Config
	database *sql.DB
}

// NewIndexDiffCollector creates IndexDiffCollector.
func NewIndexDiffCollector(c *common.Config) (*IndexDiffCollector, error) {
	db, err := sql.Open("ydb", "ydb://"+(*c).Endpoint+(*c).Database+"?auth-token=secret")

	if err != nil {
		return nil, err
	}

	return &IndexDiffCollector{config: c, database: db}, nil
}

// GetDanglins checks the consistency of secondary indexes in tables from INDEXES.
func (differ *IndexDiffCollector) GetDanglings(ctx context.Context) ([]Dangling, error) {
	danglings := []Dangling{}

	var (
		allCounts  []int
		baseCount  int
		indexCount int
		count       int
	)

	for index := range INDEXES {
		query := fmt.Sprintf(`
				--!syntax_v1
				PRAGMA TablePathPrefix("%s");

				SELECT COUNT(*) AS cnt
				FROM %s
				UNION ALL SELECT COUNT(*) AS cnt
				FROM %s;
				`,
			differ.config.Prefix,
			INDEXES[index].baseTable,
			INDEXES[index].indexTable)

		counts, err := differ.database.QueryContext(ctx, query)
		defer counts.Close()

		if err != nil {
			return nil, err
		}

		for counts.Next() {
			err = counts.Scan(&count)
			if err != nil {
				return nil, err
			}

			allCounts = append(allCounts, count)
		}


		if len(allCounts) == 1 {
			count = allCounts[0]
			danglings = append(danglings, map[string]int{
				fmt.Sprintf(`One of %s and %s is empty`, INDEXES[index].baseTable, INDEXES[index].baseTable) :  count})
		}

		if len(allCounts) < 2 {
			continue
		}

		baseCount, indexCount = allCounts[0], allCounts[1]

		if baseCount != indexCount {
			danglings = append(danglings, map[string]int{
				INDEXES[index].baseTable:  baseCount,
				INDEXES[index].indexTable: indexCount})

		}
	}

	return danglings, nil
}
