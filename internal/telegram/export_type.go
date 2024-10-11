package telegram

//go:generate go run github.com/dmarkham/enumer -type ExportType -trimprefix Type -transform snake -text -output export_type_string.go

type ExportType uint8

const (
	TypePersonalChat  ExportType = iota
	TypePrivateGroup  ExportType = iota
	TypeSavedMessages ExportType = iota
)
