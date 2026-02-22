use std::any::Any;
use crate::error::runtime_error;
use crate::interpreter::{Interpreter, Return};
use crate::ast::{Stmt, Visitor};
use std::collections::HashMap;

pub trait Callable {
    fn call(&self, arguments: Vec<Return>, inter: &mut Interpreter, line: i32) -> Return;
    fn arity(&self) -> usize;
    fn as_any(&self) -> &dyn Any;
}

#[derive(Clone)]
pub struct Function {
    pub body: Stmt,
    pub args: Vec<String>,
}
impl Callable for Function {
    fn call(&self, arguments: Vec<Return>, inter: &mut Interpreter, _line: i32) -> Return {
        let temp = Environment::new(inter.env.clone());
        inter.env = temp;
        for i in 0..arguments.len() {
            inter.env.define(self.args.get(i).unwrap().clone(), 
                arguments.get(i).unwrap().clone());
        }

        let out = inter.execute(&self.body);
        println!("{}", out.callable.is_none());
        inter.env = inter.env.enclosing.clone().unwrap();
        return out;
    }

    fn arity(&self) -> usize {
        return self.args.len();
    }

    fn as_any(&self) -> &dyn Any {
        return self;
    }

}

#[derive(Clone)]
pub struct Environment {
    pub env: HashMap<String, Return>,
    pub enclosing: Box<Option<Environment>>
}
impl Environment {
    pub fn parent() -> Environment {
        return Environment { 
            env: HashMap::<String, Return>::new(),
            enclosing: Box::new(None)
        };
    }

    pub fn new(env: Environment) -> Environment {
        return Environment {
            env: HashMap::<String, Return>::new(),
            enclosing: Box::new(Some(env))
        }
    }

    pub fn define(&mut self, name: String, value: Return) {
        self.env.insert(name, value);
    }

    pub fn update(&mut self, name: &String, value: Return, line: i32) {
        if self.env.contains_key(name) {
            self.env.insert(name.clone(), value);
        } else {
            match self.enclosing.as_mut() {
                None => { runtime_error("Variable does not exist", line, name); },
                Some(e) => {e.update(name, value, line);}
            }
        }

    }

    pub fn get_value(&self, name: &String, line: i32) -> Return {
        if self.env.contains_key(name) {
            return self.env.get(name).unwrap().clone();
        } else {
            match self.enclosing.as_ref() {
                None => { runtime_error("Variable does not exist", line, name); },
                Some(e) => { return e.get_value(name, line); }
            }
        }

        // panic function in case runtime_error get's changed and I forget to update this
        panic!("This should never be reached");
    }
}
