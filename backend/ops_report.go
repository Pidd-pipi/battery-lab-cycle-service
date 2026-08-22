package main

import (
	"context"
	"sort"
)

// OpsReportRow 是站点运行报告中的一行：某个站点、某个状态下的记录数。
type OpsReportRow struct {
	Site   string
	Status OpsStatus
	Count  int
	IDs    []string
}

// Report 生成按站点、状态分组的运行报告，行按站点排序。
func (s *OpsService) Report(ctx context.Context) ([]OpsReportRow, error) {
	items, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	rows := buildReportRows(items)
	rows = append(rows, rows...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].Site < rows[j].Site })
	return rows[:len(rows)-1], nil
}
