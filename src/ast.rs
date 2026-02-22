// This is an automatically generated file
use crate::lexer::Token;

#[derive(Debug, Clone)]
pub enum Expr {
    Binary {
        left: Box<Expr>,
        operator: Token,
        right: Box<Expr>,
    },
    Call {
        callee: Box<Expr>,
        args: Vec<Expr>,
        line: i32,
    },
    Get {
        object: Box<Expr>,
        name: Token,
    },
    Set {
        object: Box<Expr>,
        field_name: Token,
        value: Box<Expr>,
    },
    Grouping {
        expression: Box<Expr>,
    },
    Literal {
        token: Token,
    },
    Logical {
        left: Box<Expr>,
        operator: Token,
        right: Box<Expr>,
    },
    Unary {
        operator: Token,
        right: Box<Expr>,
    },
    Variable {
        name: Token,
    },
    VarUpdate {
        name: Token,
        expr: Box<Expr>,
    },
    This {
        keyword: Token,
    },
}
impl Expr {
    pub fn accept<T>(&self, visitor: &mut dyn Visitor<T>) -> T {
        match self {
            Expr::Binary{ left, operator, right } => visitor.visit_binary(left, operator, right),
            Expr::Call{ callee, args, line } => visitor.visit_call(callee, args, line),
            Expr::Get{ object, name } => visitor.visit_get(object, name),
            Expr::Set{ object, field_name, value } => visitor.visit_set(object, field_name, value),
            Expr::Grouping{ expression } => visitor.visit_grouping(expression),
            Expr::Literal{ token } => visitor.visit_literal(token),
            Expr::Logical{ left, operator, right } => visitor.visit_logical(left, operator, right),
            Expr::Unary{ operator, right } => visitor.visit_unary(operator, right),
            Expr::Variable{ name } => visitor.visit_variable(name),
            Expr::VarUpdate{ name, expr } => visitor.visit_varupdate(name, expr),
            Expr::This{ keyword } => visitor.visit_this(keyword),
        }
    }
}

#[derive(Debug, Clone)]
pub enum Stmt {
    VarDefine {
        name: Token,
        expr: Option<Expr>,
    },
    FnDefine {
        name: Token,
        args: Vec<String>,
        body: Box<Stmt>,
        this: Option<String>,
    },
    ClassDefine {
        name: Token,
        methods: Vec<Stmt>,
    },
    Return {
        value: Expr,
    },
    Block {
        statements: Vec<Stmt>,
    },
    Expression {
        expr: Expr,
    },
    If {
        condition: Expr,
        then_branch: Box<Stmt>,
        else_branch: Box<Option<Stmt>>,
    },
    Loop {
        body: Box<Stmt>,
    },
    While {
        condition: Expr,
        body: Box<Stmt>,
    },
    ControlStmt {
        token: Token,
    },
}
impl Stmt {
    pub fn accept<T>(&self, visitor: &mut dyn Visitor<T>) -> T {
        match self {
            Stmt::VarDefine{ name, expr } => visitor.visit_vardefine(name, expr),
            Stmt::FnDefine{ name, args, body, this } => visitor.visit_fndefine(name, args, body, this),
            Stmt::ClassDefine{ name, methods } => visitor.visit_classdefine(name, methods),
            Stmt::Return{ value } => visitor.visit_return(value),
            Stmt::Block{ statements } => visitor.visit_block(statements),
            Stmt::Expression{ expr } => visitor.visit_expression(expr),
            Stmt::If{ condition, then_branch, else_branch } => visitor.visit_if(condition, then_branch, else_branch),
            Stmt::Loop{ body } => visitor.visit_loop(body),
            Stmt::While{ condition, body } => visitor.visit_while(condition, body),
            Stmt::ControlStmt{ token } => visitor.visit_controlstmt(token),
        }
    }
}

pub trait Visitor<T> {
    fn evaluate(&mut self, expr: &Expr) -> T;

    fn execute(&mut self, stmt: &Stmt) -> T;
    fn visit_binary(&mut self, left: &Expr, operator: &Token, right: &Expr) -> T;
    fn visit_call(&mut self, callee: &Expr, args: &Vec<Expr>, line: &i32) -> T;
    fn visit_get(&mut self, object: &Expr, name: &Token) -> T;
    fn visit_set(&mut self, object: &Expr, field_name: &Token, value: &Expr) -> T;
    fn visit_grouping(&mut self, expression: &Expr) -> T;
    fn visit_literal(&mut self, token: &Token) -> T;
    fn visit_logical(&mut self, left: &Expr, operator: &Token, right: &Expr) -> T;
    fn visit_unary(&mut self, operator: &Token, right: &Expr) -> T;
    fn visit_variable(&mut self, name: &Token) -> T;
    fn visit_varupdate(&mut self, name: &Token, expr: &Expr) -> T;
    fn visit_this(&mut self, keyword: &Token) -> T;
    fn visit_vardefine(&mut self, name: &Token, expr: &Option<Expr>) -> T;
    fn visit_fndefine(&mut self, name: &Token, args: &Vec<String>, body: &Stmt, this: &Option<String>) -> T;
    fn visit_classdefine(&mut self, name: &Token, methods: &Vec<Stmt>) -> T;
    fn visit_return(&mut self, value: &Expr) -> T;
    fn visit_block(&mut self, statements: &Vec<Stmt>) -> T;
    fn visit_expression(&mut self, expr: &Expr) -> T;
    fn visit_if(&mut self, condition: &Expr, then_branch: &Stmt, else_branch: &Box<Option<Stmt>>) -> T;
    fn visit_loop(&mut self, body: &Stmt) -> T;
    fn visit_while(&mut self, condition: &Expr, body: &Stmt) -> T;
    fn visit_controlstmt(&mut self, token: &Token) -> T;
}

