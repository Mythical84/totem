use crate::ast::{Expr, Stmt};
use crate::error::ParseError;
use crate::lexer::{Token, TokenType, TokenType::*};

pub struct Parser {
    tokens: Vec<Token>,
    index: usize,
}

macro_rules! unwrap_return {
    ( $e:expr ) => {
        match $e {
            Ok(o) => o,
            Err(e) => return Err(e),
        }
    };
}

impl Parser {
    pub fn new(tokens: Vec<Token>) -> Parser {
        return Parser {
            tokens,
            index: 0,
        };
    }

    pub fn parse(&mut self) -> Vec<Stmt> {
        let mut stmts: Vec<Stmt> = Vec::new();
        let mut is_error = false;

        while !self.at_end() {
            match self.statement() {
                Ok(o) => stmts.push(o),
                Err(e) => {
                    let errors: Vec<ParseError> = e.into();
                    for err in errors {
                        println!("{}", err.get_message());
                    }
                    is_error = true;
                    self.synchronize();
                }
            }
        }

        if is_error {
            return vec![];
        }

        return stmts;
    }

    fn statement(&mut self) -> Result<Stmt, Vec<ParseError>> {
        if self.match_tokens(vec![IF]) { return self.if_statement(); }
        else if self.match_tokens(vec![WHILE]) { return self.while_loop(); }
        else if self.match_tokens(vec![LOOP]) { return self.loop_stmt(); }
        else if self.match_tokens(vec![LEFT_BRACE]) { return self.block(); }
        else if self.match_tokens(vec![LET]) { return self.declaration(); }
        else if self.match_tokens(vec![FN]) { return self.function_declaration(); }
        else if self.match_tokens(vec![RETURN]) { return self.return_stmt(); }
        else if self.match_tokens(vec![CONTINUE, BREAK]) { return self.control(); }
        else if self.match_tokens(vec![CLASS]) { return self.class_declaration(); }
        else { return self.expression_statement(); }
    }

    fn class_declaration(&mut self) -> Result<Stmt, Vec<ParseError>> {
        let class_name = unwrap_return!(self.consume(VAR, "Expected class identifier"));
        unwrap_return!(self.consume(LEFT_BRACE, "Expected '{' before class body"));

        let mut methods = Vec::new();
        while !self.check(RIGHT_BRACE) && !self.at_end() {
            let stmt = unwrap_return!(self.statement());
            match stmt {
                Stmt::FnDefine { .. } => { 
                    methods.push(stmt); 
                },
                _ => { 
                    return Err(vec![ParseError::new(
                        "Expected function declaration",
                        self.previous(), 
                        self.previous().line)]) 
                }
            }
        }
        
        unwrap_return!(self.consume(RIGHT_BRACE, "expected '{' after class body"));
    
        return Ok(Stmt::ClassDefine { name: class_name, methods });
    }

    fn function_declaration(&mut self) -> Result<Stmt, Vec<ParseError>> {
        let name = unwrap_return!(self.consume(VAR, "Expected function identifier"));

        unwrap_return!(self.consume(LEFT_PAREN, "Expected ( before function arguments"));

        let mut args = Vec::new();
        if !self.check(RIGHT_PAREN) {
            args.push(unwrap_return!(self.consume(VAR, "Expected function argument")).value);
            while self.match_tokens(vec![COMMA]) {
                let arg = unwrap_return!(self.consume(VAR, "Expected function argument"));
                args.push(arg.value);
                if args.len() > 255 {
                    return Err(vec![ParseError::new(
                        "Function cannot have more than 255 arguments",
                        name.clone(),
                        name.line
                    )])
                }
            }
        }

        unwrap_return!(self.consume(RIGHT_PAREN, "Expected ) after function arguments"));

        let body = unwrap_return!(self.statement());

        return Ok(Stmt::FnDefine { name, args, body: Box::new(body), this: None });
    }

    fn return_stmt(&mut self) -> Result<Stmt, Vec<ParseError>> {
        let value = unwrap_return!(self.expression());
        unwrap_return!(self.consume(SEMI, "Expected ; after return statement"));
        return Ok(Stmt::Return { value });
    }

    fn if_statement(&mut self) -> Result<Stmt, Vec<ParseError>> {
        unwrap_return!(self.consume(LEFT_PAREN, "Expected ( after 'if'"));
        let condition = unwrap_return!(self.expression());
        unwrap_return!(self.consume(RIGHT_PAREN, "Expected ) after if condition"));

        let then = unwrap_return!(self.statement());

        let mut else_statement = None;
        if self.match_tokens(vec![ELSE]) {
            else_statement = Some(unwrap_return!(self.statement()));
        }

        return Ok(Stmt::If { 
            condition: condition, 
            then_branch: Box::new(then),
            else_branch: Box::new(else_statement)
        })
    }

    fn while_loop(&mut self) -> Result<Stmt, Vec<ParseError>> {
        unwrap_return!(self.consume(LEFT_PAREN, "Expected ( after 'if'"));
        let condition = unwrap_return!(self.expression());
        unwrap_return!(self.consume(RIGHT_PAREN, "Expected ) after if condition"));

        let body = unwrap_return!(self.statement());

        return Ok(Stmt::While { condition, body: Box::new(body) })
    }

    fn loop_stmt(&mut self) -> Result<Stmt, Vec<ParseError>> {
        let body = unwrap_return!(self.statement());

        return Ok(Stmt::Loop { body: Box::new(body) });
    }

    fn block(&mut self) -> Result<Stmt, Vec<ParseError>> {
        let mut stmts = vec![];
        let mut errors = vec![];

        while !self.check(RIGHT_BRACE) && !self.at_end() {
            let stmt = self.statement();
            match stmt {
                Ok(o) => {stmts.push(o);},
                Err(e) => {
                    for err in e {
                        errors.push(err);
                    }
                    self.synchronize();
                }
            }
        }

        unwrap_return!(self.consume(RIGHT_BRACE, "Expected '}' after block"));

        if errors.len() == 0 {
            return Ok(Stmt::Block{ statements: stmts });
        } else {
            return Err(errors);
        }
    }

    fn declaration(&mut self) -> Result<Stmt, Vec<ParseError>> {
        let name = unwrap_return!(self.consume(VAR, "Expected Identifier"));
        if self.match_tokens(vec![EQUAL]) {
            let expr = unwrap_return!(self.expression());
            unwrap_return!(self.consume(SEMI, "Expected ';' after variable declaration"));
            return Ok(Stmt::VarDefine { name, expr: Some(expr) });
        } else {
            unwrap_return!(self.consume(SEMI, "Expected ';' after variable declaration"));
            return Ok(Stmt::VarDefine { name, expr: None })
        }
    }

    fn control(&mut self) -> Result<Stmt, Vec<ParseError>> {
        let token = self.previous();
        unwrap_return!(self.consume(SEMI, "Expected semi colon after control statement"));
        return Ok(Stmt::ControlStmt { token });
    }

    fn expression_statement(&mut self) -> Result<Stmt, Vec<ParseError>> {
        let expr = unwrap_return!(self.expression());
        unwrap_return!(self.consume(SEMI, "Expected semi colon after statement"));
        return Ok(Stmt::Expression { expr });
    }

    fn expression(&mut self) -> Result<Expr, Vec<ParseError>> {
        let mut expr = unwrap_return!(self.or());

        if self.match_tokens(vec![EQUAL]) {
            match expr {
                Expr::Get { object, name } => {
                        let value = unwrap_return!(self.or());
                        expr = Expr::Set { object, field_name: name, value: Box::new(value) }
                },
                _ => {}
            }
        }

        return Ok(expr);
    }

    fn or(&mut self) -> Result<Expr, Vec<ParseError>> {
        let mut expr: Result<Expr, _> = Ok(unwrap_return!(self.and()));

        while self.match_tokens(vec![DOUBLE_LINE]) {
            let operator = self.previous();
            let right = unwrap_return!(self.and());
            expr = Ok(Expr::Logical { left: Box::new(expr.unwrap()), operator, right: Box::new(right) })
        }

        return expr;
    }

    fn and(&mut self) -> Result<Expr, Vec<ParseError>> {
        let mut expr: Result<Expr, _> = Ok(unwrap_return!(self.equality()));

        while self.match_tokens(vec![DOUBLE_AMPERSAND]) {
            let operator = self.previous();
            let right = unwrap_return!(self.equality());
            expr = Ok(Expr::Logical { left: Box::new(expr.unwrap()), operator, right: Box::new(right) })
        }

        return expr;
    }

    fn equality(&mut self) -> Result<Expr, Vec<ParseError>> {
        let mut expr = self.comparison();

        while self.match_tokens(vec![DOUBLE_EQUAL, BANG_EQUAL]) {
            let operator = self.previous();
            let right = unwrap_return!(self.comparison());
            expr = Ok(Expr::Binary {
                left: Box::new(unwrap_return!(expr)),
                operator, right: Box::new(right),
            });
        }

        return expr;
    }

    fn comparison(&mut self) -> Result<Expr, Vec<ParseError>> {
        let mut expr = self.term();

        while self.match_tokens(vec![GREATER, GREATER_EQUAL, LESS, LESS_EQUAL]) {
            let operator = self.previous();
            let right = unwrap_return!(self.term());
            expr = Ok(Expr::Binary {
                left: Box::new(unwrap_return!(expr)),
                operator,
                right: Box::new(right),
            });
        }

        return expr;
    }

    fn term(&mut self) -> Result<Expr, Vec<ParseError>> {
        let mut expr = self.factor();

        while self.match_tokens(vec![SUB, ADD]) {
            let operator = self.previous();
            let right = unwrap_return!(self.factor());
            expr = Ok(Expr::Binary {
                left: Box::new(unwrap_return!(expr)),
                operator,
                right: Box::new(right),
            });
        }

        return expr;
    }

    fn factor(&mut self) -> Result<Expr, Vec<ParseError>> {
        let mut expr = self.unary();

        while self.match_tokens(vec![MUL, DIV, MOD]) {
            let operator = self.previous();
            let right = unwrap_return!(self.unary());
            expr = Ok(Expr::Binary {
                left: Box::new(unwrap_return!(expr)),
                operator,
                right: Box::new(right),
            });
        }

        return expr;
    }

    fn unary(&mut self) -> Result<Expr, Vec<ParseError>> {
        if self.match_tokens(vec![BANG, SUB]) {
            let operator = self.previous();
            let expr = unwrap_return!(self.unary());
            return Ok(Expr::Unary {
                operator,
                right: Box::new(expr),
            });
        }

        return self.call();
    }

    fn call(&mut self) -> Result<Expr, Vec<ParseError>> {
        let mut expr = unwrap_return!(self.atom());

        loop {
            if self.match_tokens(vec![LEFT_PAREN]) {
                expr = unwrap_return!(self.finish_call(expr));
            } else if self.match_tokens(vec![DOT]) {
                let name = unwrap_return!(self.consume(VAR, "Expected identifier after '.'"));
                expr = Expr::Get { object: Box::new(expr) , name };
            } else {
                break;
            }
        }

        return Ok(expr);
    }

    fn finish_call(&mut self, expr: Expr) -> Result<Expr, Vec<ParseError>> {
        let mut args = Vec::<Expr>::new();
        if !self.check(RIGHT_PAREN) {
            args.push(unwrap_return!(self.expression()));
            while self.match_tokens(vec![COMMA]) {
                args.push(unwrap_return!(self.expression()));
            }
        }


        if args.len() > 255 {
            return Err(vec![ParseError::new(
                "Function does not support more than 255 arguments", 
                self.previous(),
                self.previous().line
            )]);
        }

        unwrap_return!(self.consume(RIGHT_PAREN, "Expected ) after function call"));

        return Ok(Expr::Call { callee: Box::new(expr), args, line: self.previous().line })
        
    }

    fn atom(&mut self) -> Result<Expr, Vec<ParseError>> {
        if self.match_tokens(vec![TRUE, FALSE, NULL, STR, INT]) {
            return Ok(Expr::Literal {
                token: self.previous(),
            });
        } else if self.match_tokens(vec![VAR]) {
            return self.identifier();
        } else if self.match_tokens(vec![SELF]) {
            return Ok(Expr::This { keyword: self.previous() });
        } else if self.match_tokens(vec![LEFT_PAREN]) {
            let expr = self.expression().unwrap();
            self.consume(RIGHT_PAREN, "Expected ) after expression").unwrap();
            return Ok(Expr::Grouping {
                expression: Box::new(expr),
            });
        }

        let token = self.peek();

        return Err(vec![ParseError::new("Expected expression", token.clone(), token.line)]);
    }

    fn identifier(&mut self) -> Result<Expr, Vec<ParseError>> {
        let name = self.previous();
        if self.match_tokens(vec![EQUAL]) {
            return Ok(Expr::VarUpdate {
                name: name,
                expr: Box::new(unwrap_return!(self.expression())),
            });
        } else {
            return Ok(Expr::Variable {
                name: name,
            });
        }
    }

    fn match_tokens(&mut self, tokens: Vec<TokenType>) -> bool {
        for token_type in tokens {
            if self.check(token_type) {
                self.advance();
                return true;
            }
        }

        return false;
    }

    fn check(&self, token_type: TokenType) -> bool {
        return self.peek().token_type == token_type;
    }

    fn at_end(&self) -> bool {
        match self.peek().token_type {
            EOF => true,
            _ => false,
        }
    }

    fn peek(&self) -> &Token {
        return self.tokens.get(self.index).unwrap();
    }

    fn previous(&self) -> Token {
        return self.tokens.get(self.index - 1).unwrap().clone();
    }

    fn advance(&mut self) -> Token {
        if !self.at_end() {
            self.index += 1;
        }
        return self.previous();
    }

    fn consume(&mut self, token_type: TokenType, message: &str) -> Result<Token, Vec<ParseError>> {
        if self.check(token_type) {
            return Ok(self.advance());
        }

        let token = self.previous();

        return Err(vec![ParseError::new(message, token.clone(), token.line)]);
    }

    fn synchronize(&mut self) {
        self.advance();
        while !self.at_end() && 
            !(self.previous().token_type == SEMI || self.previous().token_type == RIGHT_BRACE) {
            self.advance();
        }
    }
}
