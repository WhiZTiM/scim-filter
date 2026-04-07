package scim_filter

import "fmt"

var _ = Expr(LogicalExpr{})
var _ = Expr(BinaryExpr{})
var _ = Expr(PresentExpr{})
var _ = Expr(NotExpr{})
var _ = Expr(PathExpr{})

type Expr interface {
	String() string
	Visit(v Visitor) error
	isScimExpr()
}

type LogicalOpType string

const (
	LogicalOpTypeAnd LogicalOpType = "and"
	LogicalOpTypeOr  LogicalOpType = "or"
)

type LogicalExpr struct {
	Left        Expr
	Right       Expr
	Type        LogicalOpType
	encompassed bool
}

func (e LogicalExpr) String() string {
	if e.encompassed {
		return fmt.Sprintf("(%s %s %s)", e.Left, e.Type, e.Right)
	}
	return fmt.Sprintf("%s %s %s", e.Left, e.Type, e.Right)
}

func (e LogicalExpr) Visit(v Visitor) error {
	return v.VisitOpLogical(e)
}

func (e LogicalExpr) isScimExpr() {}

type PresentExpr struct {
	Attr        Attr
	encompassed bool
}

func (e PresentExpr) String() string {
	return fmt.Sprintf("%s pr", e.Attr)
}

func (e PresentExpr) Visit(v Visitor) error {
	return v.VisitPresent(e)
}

func (e PresentExpr) isScimExpr() {}

type NotExpr struct {
	Expr        Attr
	encompassed bool
}

func (e NotExpr) isScimExpr() {}

func (e NotExpr) String() string {
	return fmt.Sprintf("not(%s)", e.Expr)
}

func (e NotExpr) Visit(v Visitor) error {
	return v.VisitNot(e)
}

type BinaryOpType string

const (
	BinaryOpTypeEquals        BinaryOpType = "eq"
	BinaryOpTypeNotEquals     BinaryOpType = "ne"
	BinaryOpTypeContains      BinaryOpType = "co"
	BinaryOpTypeStartsWith    BinaryOpType = "sw"
	BinaryOpTypeEndsWith      BinaryOpType = "ew"
	BinaryOpTypeGreaterThan   BinaryOpType = "gt"
	BinaryOpTypeGreaterThanOr BinaryOpType = "ge"
	BinaryOpTypeLessThan      BinaryOpType = "lt"
	BinaryOpTypeLessThanOr    BinaryOpType = "le"
)

type BinaryExpr struct {
	Attr        Attr
	Value       Value
	Type        BinaryOpType
	encompassed bool
}

func (c BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", c.Attr, c.Type, c.Value)
}

func (c BinaryExpr) isScimExpr() {}

func (c BinaryExpr) Visit(v Visitor) error {
	return v.VisitOpBinary(c)
}

type PathExpr struct {
	Attr        Attr
	SubAttrExpr Expr
	encompassed bool
}

func (p PathExpr) String() string {
	if p.SubAttrExpr == nil {
		return p.Attr.String()
	}
	return fmt.Sprintf("%s[%s]", p.Attr, p.SubAttrExpr)
}

func (a PathExpr) HasSubAttrExpr() bool {
	return a.SubAttrExpr != nil
}

func (a PathExpr) Visit(v Visitor) error {
	return v.VisitPath(a)
}

func (p PathExpr) isScimExpr() {}

type Attr struct {
	Attr    string
	SubAttr string
}

func (a Attr) String() string {
	if a.SubAttr == "" {
		return fmt.Sprintf("%s", a.Attr)
	}
	return fmt.Sprintf("%s.%s", a.Attr, a.SubAttr)
}

type ValueType string

const (
	ValueTypeString  ValueType = "string"
	ValueTypeNumber  ValueType = "number"
	ValueTypeBoolean ValueType = "boolean"
	ValueTypeNull    ValueType = "null"
)

type Value struct {
	Value string
	Type  ValueType
}

func (v Value) String() string {
	switch v.Type {
	case ValueTypeString:
		return fmt.Sprintf("\"%s\"", v.Value)
	case ValueTypeNumber:
		return v.Value
	case ValueTypeBoolean:
		return fmt.Sprintf("%t", v.Value == "true")
	case ValueTypeNull:
		return "null"
	}
	panic(fmt.Sprintf("unknown value type: %s", v.Type))
}
