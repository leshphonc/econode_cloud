package device

// RegisterRequest 注册
type RegisterRequest struct {
	SerialNo string         `json:"serial_no" binding:"required"`
	Meta     map[string]any `json:"meta"`
}

type RegisterResponse struct {
	SerialNo  string `json:"serial_no"`
	ClaimCode string `json:"claim_code"`
}

type ActivateRequest struct {
	SerialNo  string         `json:"serial_no" binding:"required"`
	Model     string         `json:"model" binding:"required"`
	PowerMode int16          `json:"power_mode" binding:"required,oneof=1 2 3"`
	HWVersion string         `json:"hw_version" binding:"required"`
	FWVersion string         `json:"fw_version" binding:"required"`
	ClaimCode string         `json:"claim_code" binding:"required"`
	Meta      map[string]any `json:"meta"`
}

type ActivateResponse struct {
	DeviceUID    string         `json:"device_uid"`
	Name         string         `json:"name"`
	Model        string         `json:"model"`
	Status       int16          `json:"status"`
	PowerMode    int16          `json:"power_mode"`
	HWVersion    string         `json:"hw_version"`
	FWVersion    string         `json:"fw_version"`
	ClaimAt      int64          `json:"claim_at"`
	ActiveErrors []string       `json:"active_errors"`
	Meta         map[string]any `json:"meta"`
}

type HeartbeatRequest struct {
	ReportedAtMs *int64         `json:"reported_at_ms"`
	Meta         map[string]any `json:"meta"`
}
