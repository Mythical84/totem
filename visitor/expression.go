package visitor

import (
	. "main/lexer"
)

type Expr interface {
	Accept(visitor Visitor) (any, error)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (self Binary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitBinary(self)
}

type Truth struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (self Truth) Accept(visitor Visitor) (any, error) {
	return visitor.VisitTruth(self)
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (self Logical) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLogical(self)
}

type Unary struct {
	Operator Token
	Expr     Expr
}

func (self Unary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitUnary(self)
}

type Grouping struct {
	Expr Expr
}

func (self Grouping) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGrouping(self)
}

type Literal struct {
	Value any
}

func (self Literal) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLiteral(self)
}

type VarGet struct {
	Name Token
}

func (self VarGet) Accept(visitor Visitor) (any, error) {
	return visitor.VisitVarGet(self)
}

type VarDef struct {
	Name  Token
	Expr  Expr
	Local bool
}

func (self VarDef) Accept(visitor Visitor) (any, error) {
	return visitor.VisitVarDef(self)
}

type List struct {
	Vals []Expr
}

type Call struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

func (self Call) Accept(visitor Visitor) (any, error) {
	return visitor.VisitCall(self)
}

type Get struct {
	Expr Expr
	Name Token
}

func (self Get) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGet(self)
}

type Set struct {
	Expr Expr
	Name Token
	Val  Expr
}

func (self Set) Accept(visitor Visitor) (any, error) {
	return visitor.VisitSet(self)
}

type Self struct {
	Line int
}

func (self Self) Accept(visitor Visitor) (any, error) {
	return visitor.VisitSelf(self)
}

type Tuple struct {
	Vals []Expr
}

func (self Tuple) Accept(visitor Visitor) (any, error) {
	return visitor.VisitTuple(self)
}

type TupleDef struct {
	Vars  []string
	Vals  []Expr
	Local bool
}

func (self TupleDef) Accept(visitor Visitor) (any, error) {
	return visitor.VisitTupleDef(self)
}

type ListDef struct {
	Vals []Expr
}

func (self ListDef) Accept(visitor Visitor) (any, error) {
	return visitor.VisitListDef(self)
}

type IterGet struct {
	Expr  Expr
	Index Expr
	Line  int
}

func (self IterGet) Accept(visitor Visitor) (any, error) {
	return visitor.VisitIterGet(self)
}

type IterSet struct {
	Expr  Expr
	Val   Expr
	Index Expr
	Line  int
}

func (self IterSet) Accept(visitor Visitor) (any, error) {
	return visitor.VisitIterSet(self)
}
