package dto

type AppConfigRequest struct {
}

type AppConfigResponse struct {
	HydraPermissions []string `json:"hydraPermissions"`
	GezuPermissions  []string `json:"gezuPermissions"`
}
