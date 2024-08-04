package telegram

//go:generate enumer -type ExportType -trimprefix Type -transform snake -text -output export_type_string.go

type ExportType uint8

const (
	TypePersonalChat ExportType = iota
	TypePrivateGroup ExportType = iota
)
