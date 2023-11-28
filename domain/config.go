package domain

type Config struct {
	Version        string   `json:"version"`
	DryRun         bool     `json:"dryRun"`
	NotifyPlan     int      `json:"notifyPlan"`
	NotifyTransfer int      `json:"notifyTransfer"`
	ReservedAmount int64    `json:"reservedAmount"`
	ReservedUnit   string   `json:"reservedUnit"`
	RsyncArgs      []string `json:"rsyncArgs"`
	Verbosity      int      `json:"verbosity"`
	CheckForUpdate int      `json:"checkForUpdate"`
	RefreshRate    int      `json:"refreshRate"`
}
