package filetypes

type CopyConfig struct {
	Src, Dst          string
	Overwrite, Unpack bool
}

type BackupConfig struct {
	CopyConfig     CopyConfig
	AddDate, Force bool
	TypeArch       string
}
