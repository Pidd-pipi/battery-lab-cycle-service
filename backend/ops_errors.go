package main

import (
	"errors"
	"fmt"
)

var (
	ErrOpsNotFound   = errors.New("operations record not found")
	ErrOpsConflict   = errors.New("operations revision conflict")
	ErrOpsInvalid    = errors.New("operations request is invalid")
	ErrOpsTransition = errors.New("operations status transition is not allowed")
	ErrOpsPolicy     = errors.New("operations policy rejected the request")
)

type OpsError struct {
	Code      string
	Operation string
	Cause     error
}

func (e *OpsError) Error() string {
	if e.Cause == nil {
		return e.Code + ": " + e.Operation
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Operation, e.Cause)
}

// Unwrap 让 errors.Is/errors.As 能穿过 OpsError 找到底层 sentinel，
// 调用方据此区分冲突、不存在等分类，而非一律当作内部错误。
func (e *OpsError) Unwrap() error { return e.Cause }

// wrapOps 用 *OpsError 包裹领域错误，保留 Cause 链与 Code 信息。
func wrapOps(code, operation string, cause error) error {
	return &OpsError{Code: code, Operation: operation, Cause: cause}
}

// opsCode 把错误归类成对外稳定的分类字符串。
// 先按 sentinel 识别（兼容被 wrapOps 包裹的情况），再回退到 OpsError.Code。
func opsCode(err error) string {
	switch {
	case errors.Is(err, ErrOpsNotFound):
		return "not_found"
	case errors.Is(err, ErrOpsConflict):
		return "conflict"
	case errors.Is(err, ErrOpsInvalid):
		return "invalid"
	case errors.Is(err, ErrOpsTransition):
		return "transition"
	case errors.Is(err, ErrOpsPolicy):
		return "policy"
	}
	var opsErr *OpsError
	if errors.As(err, &opsErr) {
		return opsErr.Code
	}
	return "internal"
}

// opsIsDomain 判断是否为不可重试的领域错误（冲突/不存在/非法/状态迁移/策略）。
func opsIsDomain(err error) bool {
	return errors.Is(err, ErrOpsNotFound) ||
		errors.Is(err, ErrOpsConflict) ||
		errors.Is(err, ErrOpsInvalid) ||
		errors.Is(err, ErrOpsTransition) ||
		errors.Is(err, ErrOpsPolicy)
}

func opsIsNotFound(err error) bool   { return errors.Is(err, ErrOpsNotFound) }
func opsIsConflict(err error) bool   { return errors.Is(err, ErrOpsConflict) }
func opsIsInvalid(err error) bool    { return errors.Is(err, ErrOpsInvalid) }
func opsIsTransition(err error) bool { return errors.Is(err, ErrOpsTransition) }
func opsIsPolicy(err error) bool     { return errors.Is(err, ErrOpsPolicy) }
