package event

import "fmt"

// TODO err 需要加到 bizerr 里
func validateEventTypeAction(t, action int16) error {
	switch t {

	// 2 = 错误（error）
	// 3 = 告警（alert）
	case 2, 3:
		// 错误 / 告警只允许 raise / fix
		if action != 1 && action != 2 {
			return fmt.Errorf(
				"action %d is not allowed for event type %d",
				action, t,
			)
		}

	// 1 = 状态
	// 4 = 投递
	case 1, 4:
		// 状态 / 投递允许 raise / update
		if action != 1 && action != 3 {
			return fmt.Errorf(
				"action %d is not allowed for event type %d",
				action, t,
			)
		}

	// 5 = 维护
	case 5:
		// 维护事件一般只记录一次
		if action != 1 {
			return fmt.Errorf(
				"action %d is not allowed for maintenance event",
				action,
			)
		}

	// 6 = 调试
	case 6:
		// 调试只允许 update（采样/日志流）
		if action != 3 {
			return fmt.Errorf(
				"action %d is not allowed for debug event",
				action,
			)
		}
	}

	return nil
}
func validateEventSeverity(t, severity int16) error {
	switch t {

	// 调试事件：只能 debug / info
	case 6:
		if severity > 2 {
			return fmt.Errorf(
				"severity %d is too high for debug event",
				severity,
			)
		}

	// 状态 / 投递 / 维护：不应该是 critical
	case 1, 4, 5:
		if severity > 4 {
			return fmt.Errorf(
				"severity %d is not allowed for event type %d",
				severity, t,
			)
		}

	// 错误 / 告警：允许全量（1~5）
	case 2, 3:
		return nil
	}

	return nil
}
