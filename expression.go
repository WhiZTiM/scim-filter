package scim_filter

import "fmt"

type Expr interface {
	String() string
	isScimExpr()
}

type AndExpr struct {
	Left  Expr
	Right Expr
}

func (e AndExpr) String() string {
	return fmt.Sprintf("(%s and %s)", e.Left, e.Right)
}

func (e AndExpr) isScimExpr() {}

type OrExpr struct {
	Expr
	Left  Expr
	Right Expr
}

func (e OrExpr) String() string {
	return fmt.Sprintf("(%s and %s)", e.Left, e.Right)
}

func (e OrExpr) isScimExpr() {}

type PresentExpr struct {
	Attr Attr
}

func (e PresentExpr) String() string {
	return fmt.Sprintf("%s pr", e.Attr)
}

func (e PresentExpr) isScimExpr() {}

type NotExpr struct {
	Expr Attr
}

func (e NotExpr) isScimExpr() {}

func (e NotExpr) String() string {
	return fmt.Sprintf("not(%s)", e.Expr)
}

type EqualsExpr struct {
	Attr  Attr
	Value Value
}

func (e EqualsExpr) String() string {
	return fmt.Sprintf("%s eq %s", e.Attr, e.Value)
}

func (e EqualsExpr) isScimExpr() {}

type NotEqualsExpr struct {
	Attr  Attr
	Value Value
}

func (n NotEqualsExpr) String() string {
	return fmt.Sprintf("%s ne %s", n.Attr, n.Value)
}

func (n NotEqualsExpr) isScimExpr() {}

type ContainsExpr struct {
	Attr  Attr
	Value Value
}

func (c ContainsExpr) String() string {
	return fmt.Sprintf("%s co %s", c.Attr, c.Value)
}

func (c ContainsExpr) isScimExpr() {}

type StartsWithExpr struct {
	Attr  Attr
	Value Value
}

func (s StartsWithExpr) String() string {
	return fmt.Sprintf("%s sw %s", s.Attr, s.Value)
}

func (s StartsWithExpr) isScimExpr() {}

type EndsWithExpr struct {
	Attr  Attr
	Value Value
}

func (e EndsWithExpr) String() string {
	return fmt.Sprintf("%s ew %s", e.Attr, e.Value)
}

func (e EndsWithExpr) isScimExpr() {}

type GreaterThanExpr struct {
	Attr  Attr
	Value Value
}

func (g GreaterThanExpr) String() string {
	return fmt.Sprintf("%s gt %s", g.Attr, g.Value)
}

func (g GreaterThanExpr) isScimExpr() {}

type GreaterThanOrEqualsExpr struct {
	Attr  Attr
	Value Value
}

func (g GreaterThanOrEqualsExpr) String() string {
	return fmt.Sprintf("%s ge %s", g.Attr, g.Value)
}

func (g GreaterThanOrEqualsExpr) isScimExpr() {}

type LessThanExpr struct {
	Attr  Attr
	Value Value
}

func (l LessThanExpr) String() string {
	return fmt.Sprintf("%s lt %s", l.Attr, l.Value)
}

func (l LessThanExpr) isScimExpr() {}

type LessThanOrEqualsExpr struct {
	Attr  Attr
	Value Value
}

func (l LessThanOrEqualsExpr) String() string {
	return fmt.Sprintf("%s le %s", l.Attr, l.Value)
}

func (l LessThanOrEqualsExpr) isScimExpr() {}

type PathExpr struct {
	Attr Attr
}

func (p PathExpr) String() string {
	return p.Attr.String()
}

func (p PathExpr) isScimExpr() {}

type Attr struct {
	Value       string
	SubAttrExpr Expr
}

func (a Attr) String() string {
	if a.SubAttrExpr == nil {
		return a.Value
	}
	return fmt.Sprintf("%s[%s]", a.Value, a.SubAttrExpr)
}

func (a Attr) HasSubAttrExpr() bool {
	return a.SubAttrExpr != nil
}

type Value struct {
	Value string
}

func (v Value) String() string {
	return v.Value
}
