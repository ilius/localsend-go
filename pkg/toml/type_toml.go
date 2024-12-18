package toml

// tomlType represents any Go type that corresponds to a TOML type.
// While the first draft of the TOML spec has a simplistic type system that
// probably doesn't need this level of sophistication, we seem to be militating
// toward adding real composite types.
type tomlType interface {
	typeString() string
}

type tomlBaseType string

func (btype tomlBaseType) typeString() string { return string(btype) }
func (btype tomlBaseType) String() string     { return btype.typeString() }

var (
	tomlInteger   tomlBaseType = "Integer"
	tomlFloat     tomlBaseType = "Float"
	tomlDatetime  tomlBaseType = "Datetime"
	tomlString    tomlBaseType = "String"
	tomlBool      tomlBaseType = "Bool"
	tomlArray     tomlBaseType = "Array"
	tomlHash      tomlBaseType = "Hash"
	tomlArrayHash tomlBaseType = "ArrayHash"
)

// typeOfPrimitive returns a tomlType of any primitive value in TOML.
// Primitive values are: Integer, Float, Datetime, String and Bool.
//
// Passing a lexer item other than the following will cause a BUG message
// to occur: itemString, itemBool, itemInteger, itemFloat, itemDatetime.
func (p *parser) typeOfPrimitive(lexItem item) tomlType {
	switch lexItem.typ {
	case itemInteger:
		return tomlInteger
	case itemFloat:
		return tomlFloat
	case itemDatetime:
		return tomlDatetime
	case itemString, itemStringEsc:
		return tomlString
	case itemMultilineString:
		return tomlString
	case itemRawString:
		return tomlString
	case itemRawMultilineString:
		return tomlString
	case itemBool:
		return tomlBool
	}
	p.bug("Cannot infer primitive type of lex item '%s'.", lexItem)
	panic("unreachable")
}
