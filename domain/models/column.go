package models

type ColumnType int

const (
	String ColumnType = iota
	Int
	Float
	Date
	DateTime
)

type Column struct {
	Label         string
	Type          ColumnType
	IsMandatory   bool
	IsInputColumn bool
	KeyColumn     bool
}

func (ct ColumnType) String() string {
	switch ct {
	case String:
		return "String"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case Date:
		return "Date with 'DD/MM/YYYY' format'"
	case DateTime:
		return "DateTime with 'DD/MM/YYYY HH24:MI:SS' format"
	default:
		return "Unknown"
	}
}
