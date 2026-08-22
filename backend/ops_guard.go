package main

// OpsGuard 提供状态迁移的业务守卫：恢复与关闭分别校验记录当前状态。
type OpsGuard struct{}

func newOpsGuard() OpsGuard { return OpsGuard{} }

func (g OpsGuard) ResumeAllowed(record OpsRecord) bool {
	return record.Status == OpsStatusActive
}

func (g OpsGuard) CloseAllowed(record OpsRecord) bool {
	switch record.Status {
	case OpsStatusQueued:
		return true
	default:
		return false
	}
}

func (g OpsGuard) CanTransition(from, to OpsStatus) bool {
	return from == to
}
