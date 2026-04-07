package scim_filter

import (
	"errors"
	"fmt"
	"strings"
)

type Operator string

const (
	AND Operator = "and"
	OR  Operator = "or"

	PRESENT Operator = "pr"

	EQUALS                 Operator = "eq"
	NOT_EQUALS             Operator = "ne"
	CONTAINS               Operator = "co"
	STARTS_WITH            Operator = "sw"
	ENDS_WITH              Operator = "ew"
	GREATER_THAN           Operator = "gt"
	GREATER_THAN_OR_EQUALS Operator = "ge"
	LESS_THAN              Operator = "lt"
	LESS_THAN_OR_EQUALS    Operator = "le"
)

var unexpectedEof = errors.New("unexpected end of input")

type parser struct {
	input      string
	tokens     []Token
	stack      []int
	parenDepth int
}

func Parse(filterExpr string) (Expr, error) {
	tokens, err := Tokenize(filterExpr)
	if err != nil {
		return nil, err
	}

	p := parser{
		input:  filterExpr,
		tokens: tokens,
		stack:  []int{0},
	}

	expr, err := p.parseFilter(false)
	return expr, err
}

func (p *parser) parseFilter(encompassed bool) (Expr, error) {
	//      FILTER    = attrExp / logExp / valuePath / *1"not" "(" FILTER ")"
	expr, err := p.parseSingleFilter(encompassed)
	if err != nil {
		return nil, err
	}

	if _, ok := p.peekNext(); !ok {
		return expr, nil
	}

	expr, err = p.parseLogical(expr, encompassed)
	if err == nil {
		return expr, nil
	}
	return nil, fmt.Errorf("error parsing filter: %w", err)
}

func (p *parser) parseSingleFilter(encompassed bool) (Expr, error) {
	//      FILTER    = attrExp / logExp / valuePath / *1"not" "(" FILTER ")"

	// Possible parenthesis
	t, ok := p.peekNext()
	if !ok {
		return nil, unexpectedEof
	}
	if t.Type == typeOpenParen {
		return p.parseParenthesis(p.parseFilter)
	}

	// Possible bareword expression
	if t.Type != typeBareword {
		return nil, fmt.Errorf("expected bareword, got %s", t.Type)
	}

	// Deal with "not"
	if t.Value == "not" {
		p.consumeNext()
		return p.parseParenthesis(p.parseFilter)
	}

	// try valuePath
	p.snapshotPush()
	expr, err := p.parseValuePath(encompassed)
	if err == nil {
		return expr, nil
	}
	p.snapshotPop()

	// try attrExpr
	p.snapshotPush()
	expr, err = p.parseAttrExpr(encompassed)
	if err == nil {
		return expr, nil
	}
	p.snapshotPop()

	return nil, fmt.Errorf("error parsing filter: %w", err)
}

func (p *parser) parseLogical(initial Expr, encompassed bool) (Expr, error) {
	// some shunting yard algorithm here.
	operands := []Expr{initial}
	operators := []string{}

	for {
		operator, ok := p.peekNext()
		if !ok {
			return nil, unexpectedEof
		}
		if operator.Type != typeBareword || (operator.Value != "and" && operator.Value != "or") {
			break
		}
		p.consumeNext()

		expr, err := p.parseFilter(encompassed)
		if err != nil {
			return nil, err
		}

		if operator.Value == "and" {
			i := len(operands) - 1
			operands[i] = LogicalExpr{Left: operands[i], Right: expr, Type: LogicalOpTypeAnd}
		} else {
			operators = append(operators, operator.Value)
			operands = append(operands, expr)
		}

		if _, ok := p.peekNext(); !ok {
			break
		}
	}

	if len(operands) == 1 {
		return operands[0], nil
	}

	if len(operands) < 2 || len(operands)%2 != 0 || len(operators) != len(operands)-1 {
		panic("invalid number of operands for logical expression")
	}

	for len(operators) > 0 {
		i := len(operators) - 1
		switch operators[i] {
		case "and":
			operands[i] = LogicalExpr{Left: operands[i], Right: operands[i+1], Type: LogicalOpTypeAnd}
		case "or":
			operands[i] = LogicalExpr{Left: operands[i], Right: operands[i+1], Type: LogicalOpTypeOr}
		}
		operators = operators[:i]
		operands = operands[:i+1]
	}

	if len(operands) != 1 || len(operators) != 0 {
		panic(fmt.Sprintf("invalid number of operands for logical expression"))
	}

	return operands[0], nil
}

func (p *parser) parseParenthesis(f func(bool) (Expr, error)) (Expr, error) {
	p.parenDepth++
	defer func() { p.parenDepth-- }()

	t, ok := p.consumeNext()
	if !ok {
		return nil, unexpectedEof
	}
	if t.Type != typeOpenParen {
		return nil, fmt.Errorf("expected '(', got %s", t)
	}

	res, err := f(true)
	if err != nil {
		return nil, err
	}

	t, ok = p.consumeNext()
	if !ok {
		return nil, unexpectedEof
	}
	if t.Type != typeCloseParen {
		return nil, fmt.Errorf("expected ')', got %s", t)
	}
	return res, nil
}

func (p *parser) parseValuePath(encompassed bool) (Expr, error) {
	// valuePath = attrPath "[" valFilter "]"
	//                 ; FILTER uses sub-attributes of a parent attrPath
	attr, err := p.parseAttrPath()
	if err != nil {
		return nil, err
	}

	t, ok := p.peekNext()
	if !ok {
		return nil, unexpectedEof
	}
	if t.Type != typeOpenSquareBrace {
		return nil, fmt.Errorf("expected '[', got %s", t)
	}

	p.consumeNext()
	expr, err := p.parseFilter(encompassed)
	if err != nil {
		return nil, err
	}

	t, ok = p.consumeNext()
	if !ok {
		return nil, unexpectedEof
	}
	if t.Type != typeCloseSquareBrace {
		return nil, fmt.Errorf("expected ']', got %s", t)
	}

	return PathExpr{Attr: attr, SubAttrExpr: expr, encompassed: encompassed}, nil
}

func (p *parser) parseAttrExpr(encompassed bool) (Expr, error) {
	// attrExp   = (attrPath SP "pr") /
	//                 (attrPath SP compareOp SP compValue)
	attr, err := p.parseAttrPath()
	if err != nil {
		return nil, err
	}
	operator, ok := p.consumeNext()
	if !ok {
		return nil, unexpectedEof
	}
	if operator.Type != typeBareword {
		return nil, fmt.Errorf("expected operator, got %s", operator)
	}
	if operator.Value == "pr" {
		return PresentExpr{Attr: attr, encompassed: encompassed}, nil
	}

	operand, ok := p.consumeNext()
	if !ok {
		return nil, unexpectedEof
	}

	var operandValue Value
	switch operand.Type {
	case typeQuotedStr:
		operandValue = Value{Value: operand.Value, Type: ValueTypeString}
	case typeBool:
		operandValue = Value{Value: operand.Value, Type: ValueTypeBoolean}
	case typeNumber:
		operandValue = Value{Value: operand.Value, Type: ValueTypeNumber}
	case typeNull:
		operandValue = Value{Value: operand.Value, Type: ValueTypeNull}
	default:
		return nil, fmt.Errorf("unexpected operand type: %s", operand)
	}

	var op BinaryOpType
	switch operator.Value {
	case "eq":
		op = BinaryOpTypeEquals
	case "ne":
		op = BinaryOpTypeNotEquals
	case "co":
		op = BinaryOpTypeContains
	case "sw":
		op = BinaryOpTypeStartsWith
	case "ew":
		op = BinaryOpTypeEndsWith
	case "gt":
		op = BinaryOpTypeGreaterThan
	case "ge":
		op = BinaryOpTypeGreaterThanOr
	case "lt":
		op = BinaryOpTypeLessThan
	case "le":
		op = BinaryOpTypeLessThanOr
	default:
		return nil, fmt.Errorf("unknown operator: %s", operator)
	}
	return BinaryExpr{Attr: attr, Type: op, Value: operandValue, encompassed: encompassed}, nil
}

func (p *parser) parseAttrPath() (Attr, error) {
	// attrPath  = [URI ":"] ATTRNAME *1subAttr
	//                 ; SCIM attribute name
	//                 ; URI is SCIM "schema" URI
	t, ok := p.consumeNext()
	if !ok {
		return Attr{}, unexpectedEof
	}
	if t.Type != typeBareword {
		return Attr{}, fmt.Errorf("expected bareword, got %s", t)
	}

	// todo validate attr name schema, etc...
	// attrPath  = [URI ":"] ATTRNAME *1subAttr
	//                 ; SCIM attribute name
	//                 ; URI is SCIM "schema" URI
	//
	//     ATTRNAME  = ALPHA *(nameChar)
	//
	//     nameChar  = "-" / "_" / DIGIT / ALPHA
	//
	//     subAttr   = "." ATTRNAME
	//                 ; a sub-attribute of a complex attribute
	a := strings.Split(t.Value, ".")
	switch len(a) {
	case 1:
		return Attr{Attr: a[0]}, nil
	case 2:
		return Attr{Attr: a[0], SubAttr: a[1]}, nil
	default:
		return Attr{}, fmt.Errorf("invalid attribute path: %s", t)
	}
}

func (p *parser) snapshotPush() {
	p.stack = append(p.stack, p.stack[len(p.stack)-1])
}

func (p *parser) snapshotPop() {
	if len(p.stack) == 1 {
		panic("stack underflow")
	}
	p.stack = p.stack[:len(p.stack)-1]
}

func (p *parser) consumeNext() (Token, bool) {
	res, ok := p.peekNext()
	if !ok {
		return Token{}, false
	}
	p.stack[len(p.stack)-1]++
	return res, ok
}

func (p *parser) peekNext() (Token, bool) {
	i := p.stack[len(p.stack)-1]
	if i >= len(p.tokens) {
		return Token{}, false
	}
	return p.tokens[i], true
}
