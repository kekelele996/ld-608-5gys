package constants

var LogTemplates = map[string][]string{
	"FlightTurnaround": {"创建航班过站：%s", "更新航班过站：%s", "航班过站状态变更：%s %s -> %s", "航班过站放行成功：航班 %s"},
	"GroundTask":       {"创建地勤任务：%s", "更新地勤任务：%s", "地勤任务状态变更：%s %s -> %s", "地勤任务同步至放行时间线：航班 %s"},
	"GroundResource":   {"创建保障资源：%s", "更新保障资源：%s", "保障资源状态变更：%s %s -> %s", "保障资源台账导出：%s"},
	"ResourceBooking":  {"创建资源预约：%s", "更新资源预约：%s", "资源预约状态变更：%s %s -> %s", "资源预约批量释放：航班 %s，%d 条"},
	"DelayEvent":       {"登记延误事件：%s", "更新延误事件：%s", "关闭延误事件：%s", "未关闭延误阻塞放行：航班 %s，%d 条"},
}
