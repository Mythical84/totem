package visitor

type Visitor interface {
	// Statements
	VisitExpression(expression Expression) error
	VisitBlock(block Block) error
	VisitTryCatch(try_catch TryCatch) error
	VisitIf(if_stmt If) error
	VisitWhile(while_stmt While) error
	VisitLoop(loop Loop) error
	VisitFunc(func_stmt Func) error
	VisitReturn(return_stmt Return) error
	VisitBreak(break_stmt Break) error
	VisitContinue(continue_stmt Continue) error
	VisitClassDecl(class_decl ClassDecl) error
	VisitImport(import_stmt Import) error
	VisitThrow(throw Throw) error

	// Expressions
	VisitBinary(binary Binary) (any, error)
	VisitTruth(truth Truth) (any, error)
	VisitLogical(logical Logical) (any, error)
	VisitGrouping(grouping Grouping) (any, error)
	VisitLiteral(literal Literal) (any, error)
	VisitUnary(unary Unary) (any, error)
	VisitVarGet(var_get VarGet) (any, error)
	VisitVarDef(var_def VarDef) (any, error)
	VisitCall(call Call) (any, error)
	VisitGet(get Get) (any, error)
	VisitSet(set Set) (any, error)
	VisitSelf(self Self) (any, error)
	VisitTuple(tuple Tuple) (any, error)
	VisitTupleDef(tuple_def TupleDef) (any, error)
	VisitListDef(list_def ListDef) (any, error)
	VisitIterGet(iter_get IterGet) (any, error)
	VisitIterSet(iter_set IterSet) (any, error)
}
