package main

import (
	"context"
	"fmt"
)

// RetryPolicy 控制重试次数。
type RetryPolicy struct {
	MaxAttempts int
}

// RunWithRetry 按退避策略重试 fn；领域错误（冲突/不存在/非法/状态迁移/策略）立即停止。
func RunWithRetry(ctx context.Context, policy RetryPolicy, fn func() error) error {
	attempts := policy.MaxAttempts
	if attempts < 1 {
		attempts = 3
	}
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		switch opsCode(lastErr) {
		case "not_found", "conflict", "invalid", "transition", "policy":
			return lastErr
		}
		if attempt < attempts {
			if err := opsDelay(ctx, opsBackoff(attempt)); err != nil {
				return err
			}
		}
	}
	return fmt.Errorf("retry exhausted: %w", lastErr)
}
