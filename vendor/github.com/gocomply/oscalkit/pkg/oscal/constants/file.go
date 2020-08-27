package constants

// Representing OSCAL file format. XML, JSON, YAML, ...
type DocumentFormat int

const (
	UnknownFormat DocumentFormat = iota
	XmlFormat
	JsonFormat
	YamlFormat
)

type DocumentType int

const (
	UnknownDocument = iota
	CatalogDocument
	ProfileDocument
	SSPDocument
	ComponentDocument
)
