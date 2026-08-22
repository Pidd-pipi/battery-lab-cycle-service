package main

// buildReportRows 把记录按「站点 + 状态」分组，生成报告行。
func buildReportRows(items []OpsRecord) []OpsReportRow {
	rows := []OpsReportRow{}
	index := map[string]int{}
	for _, item := range items {
		key := item.LabelValue("site") + "|" + string(item.Status)
		i, ok := index[key]
		if !ok {
			rows = append(rows, OpsReportRow{Site: item.LabelValue("site"), Status: item.Status, IDs: []string{}})
			index[key] = len(rows) - 1
			i = len(rows) - 1
		}
		rows[i].Count++
		rows[i].IDs = append(rows[i].IDs, item.ID)
	}
	return rows
}
