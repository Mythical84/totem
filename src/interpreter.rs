use std::rc::Rc;

use crate::ast::Expr;
use crate::ast::Stmt;
use crate::ast::Visitor;
use crate::env::Callable;
use crate::env::Environment;
use crate::env::Function;
use crate::error::runtime_error;
use crate::lexer::Token;
use crate::lexer::TokenType;
use crate::lexer::TokenType::*;
use crate::stdlib;

#[derive(Clone)]
pub struct Return {
    pub value: String,
    pub return_type: ReturnType,
    pub callable: Option<Rc<dyn Callable>>,
    pub instance: Option<ClassInstance>
}
impl Return {
    pub fn new(value: String, return_type: ReturnType) -> Return {
        return Return { value, return_type, callable: None, instance: None };
    }

    pub fn callable(value: String, return_type: ReturnType, callable: Option<Rc<dyn Callable>>) -> Return {
        return Return { value, return_type, callable, instance: None };
    }

    pub fn instance(value: String, instance: Option<ClassInstance>) -> Return {
        return Return { value, return_type: ReturnType::Instance, callable: None, instance };
    }

    pub fn null() -> Return {
        return Return::new("null".to_string(), ReturnType::Null);
    }
}

#[derive(PartialEq, Debug, Clone)]
pub enum ReturnType {
    String,
    Int,
    Float,
    Bool,
    Null,
    Continue,
    Break,
    Return(Box<ReturnType>),
    Callable,
    Instance
}

#[derive(Clone)]
pub struct Class {
    pub name: String,
    env: Environment
}

impl Callable for Class {
    fn call(&self, arguments: Vec<Return>, inter: &mut Interpreter, line: i32) -> Return {
        return Return::instance(
           format!("<{} instance>", self.name),
           Some(ClassInstance::new(self.env.clone(), false))
        );
    }

    fn arity(&self) -> usize { 0 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}

#[derive(Clone)]
pub struct ClassInstance {
    env: Environment,
    readonly: bool
}
impl ClassInstance {
    fn new(env: Environment, readonly: bool) -> ClassInstance {
        ClassInstance {
            env: Environment::new(env),
            readonly
        }
    }
}

pub struct Interpreter {
    pub env: Environment,
    pub global: Environment
}
impl Visitor<Return> for Interpreter {
    fn evaluate(&mut self, expr: &Expr) -> Return {
        return expr.accept(self);
    }

    fn execute(&mut self, stmt: &Stmt) -> Return {
        return stmt.accept(self);
    }

    // Statements
    fn visit_if(&mut self, condition: &Expr, then_branch: &Stmt, else_branch: &Box<Option<Stmt>>) -> Return {
        let c = self.evaluate(condition);
        if self.is_truthy(&c) {
            return self.execute(then_branch);
        } else {
            match else_branch.as_ref() {
                None => {
                    return Return::null();
                },
                Some(block) => {
                    return self.execute(block);
                }
            }
        }
    }

    fn visit_loop(&mut self, body: &Stmt) -> Return {
        loop {
            let b = self.execute(body);
            if b.return_type == ReturnType::Break { break; }
        }

        return Return::null();
    }

    fn visit_while(&mut self, condition: &Expr, body: &Stmt) -> Return {
        let mut t = self.evaluate(condition);

        while self.is_truthy(&t) {
            let b = self.execute(body);
            if b.return_type == ReturnType::Break { break; }
            t = self.evaluate(condition);
        }
        
        return Return::null();
    }
    
    fn visit_block(&mut self, statements: &Vec<Stmt>) -> Return {
        let temp = Environment::new(self.env.clone());
        self.env = temp;

        if statements.len() == 0 {
            return Return::null();
        }

        let mut out = Return::null();
        for i in 0..statements.len() {
            out = self.execute(statements.get(i).unwrap());
            match out.return_type {
                ReturnType::Break | ReturnType::Continue => { break; }
                ReturnType::Return(t) => { out = Return { 
                        value: out.value,
                        return_type: *t,
                        callable: out.callable,
                        instance: out.instance
                    };
                    break; 
                }
                _ => {}
            };

        }

        self.env = self.env.enclosing.clone().unwrap();
        return out;
    }

    fn visit_vardefine(&mut self, name: &Token, expr: &Option<Expr>) -> Return {
        match expr {
            Some(expr) => {
                let expression = self.evaluate(expr);
                self.env.define(name.value.clone(), expression.clone());
            },
            None => {
                self.env.define(name.value.clone(), Return::null());
            }
        }
        return Return::null();
    }

    fn visit_fndefine(&mut self, name: &Token, args: &Vec<String>, body: &Stmt, this: &Option<String>) -> Return {
        let args = args.to_vec();
        let func = Rc::new(Function {
            body: body.clone(),
            args,
            this: this.clone()
        });
        let ret = Return::callable(format!("<function {}>", name.value), ReturnType::Callable, Some(func));
        self.env.define(name.value.clone(), ret);
        return Return::null();
    }

    fn visit_classdefine(&mut self, name: &Token, methods: &Vec<Stmt>) -> Return {
        let mut env = Environment::new(self.global.clone());
        for stmt in methods {
            match stmt {
                Stmt::FnDefine { name, args, body, .. } => {
                    env.define(name.value.clone(), Return::callable(
                        format!("<method {}>", name.value), 
                        ReturnType::Callable, 
                        Some(Rc::new(Function {
                            body: *body.clone(),
                            args: args.to_vec(),
                            this: None 
                        } )
                    )));
                },
                _ => { panic!("This shouldn't fire"); }
            }
        }
        let ret = Return::callable(
            format!("<class {}>", name.value),
            ReturnType::Callable,
            Some(Rc::new(Class { name: name.value.clone(), env }))
        );
        self.env.define(name.value.clone(), ret);
        return Return::null();
    }

    fn visit_return(&mut self, value: &Expr) -> Return {
        let expr = self.evaluate(value);
        return Return {
            value: expr.value, 
            return_type: ReturnType::Return(Box::new(expr.return_type)),
            callable: expr.callable,
            instance: expr.instance
        }

    }

    fn visit_controlstmt(&mut self, token: &Token) -> Return {
        let return_type;

        match token.token_type {
            CONTINUE => { return_type = ReturnType::Continue; }
            BREAK => { return_type = ReturnType::Break; }
            // Should never fire
            _ => { panic!("Unknown control statement") }
        }
        
        return Return::new("".to_string(), return_type);
    }

    fn visit_expression(&mut self, expr: &Expr) -> Return {
        return self.evaluate(expr);
    }

    // Expressions
    fn visit_logical(&mut self, left: &Expr, operator: &Token, right: &Expr) -> Return {
        let left = self.evaluate(left);

        if operator.token_type == DOUBLE_LINE {
            if self.is_truthy(&left) {return left;}
        } else if operator.token_type == DOUBLE_AMPERSAND {
            if !self.is_truthy(&left) {return left;}
        }
        return self.evaluate(right);
    }
    
    fn visit_binary(&mut self, left: &Expr, operator: &Token, right: &Expr) -> Return {
        let left = self.evaluate(left);
        let right = self.evaluate(right);

        let i_left;
        let i_right;
        let mut float = false;

        if (left.return_type == ReturnType::String || right.return_type == ReturnType::String)
            && operator.token_type == ADD {
            return Return::new(left.value + &right.value, ReturnType::String);
        } else if operator.token_type == DOUBLE_EQUAL {
            return Return::new(self.is_equal(left, right).to_string(), ReturnType::Bool);
        } else if operator.token_type == BANG_EQUAL {
            return Return::new((!self.is_equal(left, right)).to_string(), ReturnType::Bool);
        } else {
            if !self.is_num(&left.value) {
                runtime_error("Token must be int", operator.line, &left.value);
            }
            if !self.is_num(&right.value) {
                runtime_error("Token must be int", operator.line, &right.value);
            }
            i_left = left.value.parse::<f32>().unwrap();
            i_right = right.value.parse::<f32>().unwrap();
            if left.value.contains(".") || right.value.contains(".") {
                float = true;
            }
        }

        match operator.token_type {
            ADD => {
                if float { return Return::new(((i_left + i_right)).to_string(), ReturnType::Float); }
                else {  return Return::new(((i_left + i_right) as i32).to_string(), ReturnType::Int);}
            }
            SUB => {
                if float { return Return::new((i_left - i_right).to_string(), ReturnType::Float); }
                else {  return Return::new(((i_left - i_right) as i32).to_string(), ReturnType::Int);}
            }
            DIV => {
                if float { return Return::new((i_left / i_right).to_string(), ReturnType::Float); }
                else {  return Return::new(((i_left / i_right) as i32).to_string(), ReturnType::Int);}
            }
            MUL => {
                if float { return Return::new((i_left * i_right).to_string(), ReturnType::Float); }
                else {  return Return::new(((i_left * i_right) as i32).to_string(), ReturnType::Int);}
            }
            MOD => {
                if float { return Return::new((i_left % i_right).to_string(), ReturnType::Float); }
                else {  return Return::new(((i_left % i_right) as i32).to_string(), ReturnType::Int);}
            }
            GREATER => {
                return Return::new((i_left > i_right).to_string(), ReturnType::Bool);
            }
            LESS => {
                return Return::new((i_left < i_right).to_string(), ReturnType::Bool);
            }
            GREATER_EQUAL => {
                return Return::new((i_left >= i_right).to_string(), ReturnType::Bool);
            }
            LESS_EQUAL => {
                return Return::new((i_left <= i_right).to_string(), ReturnType::Bool);
            }
            _ => {
                panic!("Unimplemented binary operator token");
            }
        }
    }

    fn visit_grouping(&mut self, expression: &Expr) -> Return {
        return self.evaluate(expression);
    }

    fn visit_literal(&mut self, token: &Token) -> Return {
        let return_type = self.token_type_to_return_type(&token.token_type);
        return Return::new(token.value.clone(), return_type);
    }

    fn visit_unary(&mut self, operator: &Token, right: &Expr) -> Return {
        let right = self.evaluate(right);

        match operator.token_type {
            SUB => {
                if !self.is_num(&right.value) {
                    runtime_error("Token must be int", operator.line, &right.value);
                }
                let int = right.value.parse::<i32>().unwrap();
                return Return::new((-int).to_string(), ReturnType::Int);
            }
            BANG => {
                let val = !self.is_truthy(&right);
                return Return::new((val).to_string(), ReturnType::Bool);
            }
            // Unreachable
            _ => {
                panic!("Unknown unary operator token");
            }
        }
    }

    fn visit_call(&mut self, callee: &Expr, args: &Vec<Expr>, line: &i32) -> Return {
        let callee = self.evaluate(callee);
        let mut arguments = Vec::new();

        for arg in args {
            arguments.push(self.evaluate(arg));
        }

        match callee.callable {
            Some(callable) => { 
                if callable.arity() < 256 &&  arguments.len() != callable.arity() { 
                    runtime_error("Incorrect number of function arguments", *line, &callee.value);
                }

                return callable.call(arguments, self, *line);
            },
            None => {runtime_error("Object is not callable", *line, &callee.value); }
        }

        // Should never trigger
        panic!("Should never trigger");
    }

    fn visit_get(&mut self, object: &Expr, name: &Token) -> Return {
        let object = self.evaluate(object);

        match object.return_type {
            ReturnType::Instance => {
                return object.instance.unwrap().env.get_value(&name.value, name.line);
            },
            _ => { runtime_error("Only instances have properties", name.line, &name.value); }
        }

        panic!("Should never be reached");
    }

    fn visit_this(&mut self, keyword: &Token) -> Return {
        let name = self.env.enclosing.clone().unwrap().get_value(&"!self".to_string(), keyword.line);
        println!("{}", name.value);
        return self.env.enclosing.clone().unwrap().get_value(&"test".to_string(), keyword.line);
    }

    fn visit_set(&mut self, object: &Expr, field_name: &Token, value: &Expr) -> Return {
        let mut instance = self.evaluate(object);

        if instance.instance.is_none() {
            runtime_error("Only instances have fields", field_name.line, &field_name.value);
        }

        if instance.instance.as_ref().unwrap().readonly {
            runtime_error(&format!("{} object is readonly", instance.value), field_name.line, &field_name.value);
        }

        println!(": {:?}", object);

        match object {
            Expr::Variable { name } => {
                let val = self.evaluate(value);
                instance.instance.as_mut().unwrap().env.define(field_name.value.clone(), val);
                self.env.define(name.value.clone(), instance);
            },
            Expr::This { .. } => {
                let val = self.evaluate(value);
                instance.instance.as_mut().unwrap().env.define(field_name.value.clone(), val.clone());
                
            }
            _ => {panic!("Honestly idk but this probably won't fire")}
        }
        
        return Return::null();
    }

    fn visit_variable(&mut self, name: &Token) -> Return {
        return self.env.get_value(&name.value, name.line);
    }

    fn visit_varupdate(&mut self, name: &Token, expr: &Expr) -> Return {
        let expression = self.evaluate(expr);
        self.env.update(&name.value, expression.clone(), name.line);
        return expression;
    }

}
impl Interpreter {
    pub fn new() -> Interpreter {
        let mut global = Environment::parent();
        global = stdlib::define_native(global);
        return Interpreter {
            global: global.clone(),
            env: Environment::new(global),
        };
    }

    pub fn interpret(&mut self, stmts: Vec<Stmt>) {
        for stmt in stmts {
            self.execute(&stmt);
        }
    }

    fn is_truthy(&self, val: &Return) -> bool {
        if val.return_type == ReturnType::Null {
            return false;
        } else if val.return_type == ReturnType::Bool {
            return val.value.parse::<bool>().unwrap();
        } else if val.return_type == ReturnType::Int || val.return_type == ReturnType::Float {
            return val.value.parse::<f32>().unwrap() != 0.0;
        }
        return true;
    }

    fn is_equal(&self, val1: Return, val2: Return) -> bool {
        if val1.return_type != val2.return_type {
            return false;
        }
        return val1.value == val2.value;
    }

    pub fn is_num(&self, val: &String) -> bool {
        let mut decimal = false;
        for c in val.chars() {
            if !decimal && c == '.' {
                decimal = true;
            } else if !c.is_digit(10) {
                return false;
            }
        }
        return true;
    }

    fn token_type_to_return_type(&self, token_type: &TokenType) -> ReturnType {
        if token_type == &INT {
            return ReturnType::Int;
        } else if token_type == &STR {
            return ReturnType::String;
        } else if token_type == &TRUE || token_type == &FALSE {
            return ReturnType::Bool;
        } else if token_type == &NULL {
            return ReturnType::Null;
        } else if token_type == &DOUBLE {
            return ReturnType::Float;
        } else {
            panic!("Unimplemented literal type");
        }
    }
}
