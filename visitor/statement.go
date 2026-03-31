package visitor

import (
	. "main/lexer"
)

type Stmt interface {
	Accept(visitor Visitor) error
}

type Expression struct {
	Expr Expr
}

func (self Expression) Accept(visitor Visitor) error {
	return visitor.VisitExpression(self)
}

type TryCatch struct {
	TryBody   Stmt
	CatchBody Stmt
	ErrName   Token
}

func (self TryCatch) Accept(visitor Visitor) error {
	return visitor.VisitTryCatch(self)
}

type Block struct {
	Body []Stmt
}

func (self Block) Accept(visitor Visitor) error {
	return visitor.VisitBlock(self)
}

type If struct {
	Condition Expr
	IfBody    Stmt
	ElseBody  Stmt
}

func (self If) Accept(visitor Visitor) error {
	return visitor.VisitIf(self)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (self While) Accept(visitor Visitor) error {
	return visitor.VisitWhile(self)
}

type Loop struct {
	Body Stmt
}

func (self Loop) Accept(visitor Visitor) error {
	return visitor.VisitLoop(self)
}

type Func struct {
	Name Token
	Args []string
	Body Stmt
}

func (self Func) Accept(visitor Visitor) error {
	return visitor.VisitFunc(self)
}

type ClassDecl struct {
	Name       Token
	Body       []Stmt
	Superclass Expr
}

func (self ClassDecl) Accept(visitor Visitor) error {
	return visitor.VisitClassDecl(self)
}

type Return struct {
	Expr Expr
	Line int
}

func (self Return) Accept(visitor Visitor) error {
	return visitor.VisitReturn(self)
}

type Break struct{
	Line int
}

func (self Break) Accept(visitor Visitor) error {
	return visitor.VisitBreak(self)
}

type Continue struct {
	Line int
}

func (self Continue) Accept(visitor Visitor) error {
	return visitor.VisitContinue(self)
}

type Import struct {
	Path Token
}

func (self Import) Accept(visitor Visitor) error {
	return visitor.VisitImport(self)
}

type Throw struct {
	Expr Expr
	Line int
}

func (self Throw) Accept(visitor Visitor) error {
	return visitor.VisitThrow(self)
}
