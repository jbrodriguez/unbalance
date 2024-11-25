package domain

// UnassignedDevice represents the structure of the disk information in the JSON file
type UnassignedDevice struct {
	PoolName     string `json:"pool_name"`
	DiskLabel    string `json:"disk_label"`
	Serial       string `json:"serial"`
	SerialShort  string `json:"serial_short"`
	Part         string `json:"part"`
	Logfile      string `json:"logfile"`
	Label        string `json:"label"`
	ProgName     string `json:"prog_name"`
	EnableScript string `json:"enable_script"`
	FSType       string `json:"fstype"`
	Mountpoint   string `json:"mountpoint"`
	Luks         string `json:"luks"`
	FileSystem   string `json:"file_system"`
	CommandBg    string `json:"command_bg"`
	UserCommand  string `json:"user_command"`
	Command      string `json:"command"`
	Owner        string `json:"owner"`
	Target       string `json:"target"`
	Disk         string `json:"disk"`
	Device       string `json:"device"`
	UUID         string `json:"uuid"`
	Size         uint64 `json:"size"`
	Used         uint64 `json:"used"`
	Avail        uint64 `json:"avail"`
	Formatting   bool   `json:"formatting"`
	ArrayDisk    bool   `json:"array_disk"`
	NotUdev      bool   `json:"not_udev"`
	PassThrough  bool   `json:"pass_through"`
	Pool         bool   `json:"pool"`
	Preclearing  bool   `json:"preclearing"`
	Unmounting   bool   `json:"unmounting"`
	ReadOnly     bool   `json:"read_only"`
	Mounted      bool   `json:"mounted"`
	Mounting     bool   `json:"mounting"`
	NotUnmounted bool   `json:"not_unmounted"`
	DisableMount bool   `json:"disable_mount"`
	Running      bool   `json:"running"`
	UdDevice     bool   `json:"ud_device"`
	PartReadOnly bool   `json:"part_read_only"`
	Clearing     bool   `json:"clearing"`
	Shared       bool   `json:"shared"`
}
