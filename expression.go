package scim_filter

import "fmt"

type Expr interface {
	String() string
	isScimExpr()
}

type AndExpr struct {
	Left        Expr
	Right       Expr
	encompassed bool
}

func (e AndExpr) String() string {
	if e.encompassed {
		return fmt.Sprintf("(%s and %s)", e.Left, e.Right)
	}
	return fmt.Sprintf("%s and %s", e.Left, e.Right)
}

func (e AndExpr) isScimExpr() {}

type OrExpr struct {
	Left        Expr
	Right       Expr
	encompassed bool
}

func (e OrExpr) String() string {
	if e.encompassed {
		return fmt.Sprintf("(%s or %s)", e.Left, e.Right)
	}
	return fmt.Sprintf("%s or %s", e.Left, e.Right)
}

func (e OrExpr) isScimExpr() {}

type PresentExpr struct {
	Attr        Attr
	encompassed bool
}

func (e PresentExpr) String() string {
	return fmt.Sprintf("%s pr", e.Attr)
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

type EqualsExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (e EqualsExpr) String() string {
	return fmt.Sprintf("%s eq %s", e.Attr, e.Value)
}

func (e EqualsExpr) isScimExpr() {}

type NotEqualsExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (n NotEqualsExpr) String() string {
	return fmt.Sprintf("%s ne %s", n.Attr, n.Value)
}

func (n NotEqualsExpr) isScimExpr() {}

type ContainsExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (c ContainsExpr) String() string {
	return fmt.Sprintf("%s co %s", c.Attr, c.Value)
}

func (c ContainsExpr) isScimExpr() {}

type StartsWithExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (s StartsWithExpr) String() string {
	return fmt.Sprintf("%s sw %s", s.Attr, s.Value)
}

func (s StartsWithExpr) isScimExpr() {}

type EndsWithExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (e EndsWithExpr) String() string {
	return fmt.Sprintf("%s ew %s", e.Attr, e.Value)
}

func (e EndsWithExpr) isScimExpr() {}

type GreaterThanExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (g GreaterThanExpr) String() string {
	return fmt.Sprintf("%s gt %s", g.Attr, g.Value)
}

func (g GreaterThanExpr) isScimExpr() {}

type GreaterThanOrEqualsExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (g GreaterThanOrEqualsExpr) String() string {
	return fmt.Sprintf("%s ge %s", g.Attr, g.Value)
}

func (g GreaterThanOrEqualsExpr) isScimExpr() {}

type LessThanExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (l LessThanExpr) String() string {
	return fmt.Sprintf("%s lt %s", l.Attr, l.Value)
}

func (l LessThanExpr) isScimExpr() {}

type LessThanOrEqualsExpr struct {
	Attr        Attr
	Value       Value
	encompassed bool
}

func (l LessThanOrEqualsExpr) String() string {
	return fmt.Sprintf("%s le %s", l.Attr, l.Value)
}

func (l LessThanOrEqualsExpr) isScimExpr() {}

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
