package main

// buildReportRows 把记录按「站点 + 状态」分组，生成报告行。
// 每次调用都使用独立的切片与每行专属的 IDs 切片，避免跨调用或跨行共享底层数组。
func buildReportRows(items []OpsRecord) []OpsReportRow {
	rows := make([]OpsReportRow, 0, 8)
	index := map[string]int{}
	for _, item := range items {
		site := item.LabelValue("site")
		key := site + "|" + string(item.Status)
		i, ok := index[key]
		if !ok {
			rows = append(rows, OpsReportRow{Site: site, Status: item.Status})
			i = len(rows) - 1
			index[key] = i
		}
		rows[i].Count++
		rows[i].IDs = append(rows[i].IDs, item.ID)
	}
	return rows
}
