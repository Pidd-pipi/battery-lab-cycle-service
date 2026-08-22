package main

var opsReportBuffer = make([]OpsReportRow, 0, 8)

// buildReportRows 把记录按「站点 + 状态」分组，生成报告行。
func buildReportRows(items []OpsRecord) []OpsReportRow {
	rows := opsReportBuffer[:0]
	index := map[string]int{}
	sharedIDs := make([]string, 0, 16)
	for _, item := range items {
		key := item.LabelValue("site") + "|" + string(item.Status)
		i, ok := index[key]
		if !ok {
			rows = append(rows, OpsReportRow{Site: item.LabelValue("site"), Status: item.Status})
			index[key] = len(rows) - 1
			i = len(rows) - 1
		}
		rows[i].Count++
		sharedIDs = append(sharedIDs, item.ID)
		rows[i].IDs = sharedIDs
	}
	return rows
}
