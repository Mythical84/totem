use std::rc::Rc;
use std::time::{SystemTime, UNIX_EPOCH};
use crate::env::{Callable, Environment};
use crate::error::runtime_error;
use crate::interpreter::{Interpreter, Return, ReturnType};

pub fn define_native(mut env: Environment) -> Environment {
    create_function("println", Rc::new(Println {}), &mut env);
    create_function("print", Rc::new(Print {}), &mut env);
    create_function("f", Rc::new(F {}), &mut env);
    create_function("int", Rc::new(Int {}), &mut env);
    create_function("float", Rc::new(Float {}), &mut env);
    create_function("str", Rc::new(Str {}), &mut env);
    create_function("time", Rc::new(Time {}), &mut env);
    create_function("typeof", Rc::new(Typeof {}), &mut env);
    return env;
}

fn create_function(name: &str, callable: Rc<dyn Callable>, env: &mut Environment) {
    env.define(name.to_string(),
        Return::callable(
            format!("<built-in function {}>", name),
            ReturnType::Callable,
            Some(callable))
    );
}

struct Println{}
impl Callable for Println {
    fn call(&self, arguments: Vec<Return>, _inter: &mut Interpreter, _line: i32) -> Return {
        println!("{}", arguments.get(0).unwrap().value);
        return Return::null();
    }

    fn arity(&self) -> usize { 1 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}

struct Print{}
impl Callable for Print {
    fn call(&self, arguments: Vec<Return>, _inter: &mut Interpreter, _line: i32) -> Return {
        print!("{}", arguments.get(0).unwrap().value);
        return Return::null();
    }

    fn arity(&self) -> usize { 1 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}

struct F{}
impl Callable for F {
    fn call(&self, arguments: Vec<Return>, _inter: &mut Interpreter, line: i32) -> Return {
        let mut str = arguments.get(0).unwrap().value.clone();
        if arguments.get(0).unwrap().return_type != ReturnType::String {
            runtime_error("First argument must be a string", line, &str);
        }

        for i in 1..arguments.len() {
            str = str.replacen("{}", arguments.get(i).unwrap().value.as_str(), 1);
        }

        return Return::new(str.to_string(), ReturnType::String); 
    }

    fn arity(&self) -> usize { 256 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}

struct Int{}
impl Callable for Int {
    fn call(&self, arguments: Vec<Return>, inter: &mut Interpreter, line: i32) -> Return {
        let val = arguments.get(0).unwrap();
        if !inter.is_num(&val.value) {
            runtime_error("Argument must be a number", line, &val.value);
        }

        return Return::new(val.value.split(".").next().unwrap().to_string(), ReturnType::Int);

    }

    fn arity(&self) -> usize { 1 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}

struct Float{}
impl Callable for Float {
    fn call(&self, arguments: Vec<Return>, inter: &mut Interpreter, line: i32) -> Return {
        let val = arguments.get(0).unwrap();
        if !inter.is_num(&val.value) {
            runtime_error("Argument must be a number", line, &val.value);
        }

        return Return::new(val.value.clone(), ReturnType::String);

    }

    fn arity(&self) -> usize { 1 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}

struct Str{}
impl Callable for Str {
    fn call(&self, arguments: Vec<Return>, _inter: &mut Interpreter, _line: i32) -> Return {
        return Return::new(arguments.get(0).unwrap().value.clone(), ReturnType::String);
    }

    fn arity(&self) -> usize { 1 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}

struct Time{}
impl Callable for Time {
    fn call(&self, _arguments: Vec<Return>, _inter: &mut Interpreter, _line: i32) -> Return {
        let time = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs();

        return Return::new(time.to_string(), ReturnType::Int);
    }

    fn arity(&self) -> usize { 0 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}

struct Typeof{}
impl Callable for Typeof {
    fn call(&self, arguments: Vec<Return>, _inter: &mut Interpreter, _line: i32) -> Return {
        return Return::new(
            format!("{:?}", arguments.get(0).as_ref().unwrap().return_type).to_lowercase(), 
            ReturnType::String
        );
    }

    fn arity(&self) -> usize { 1 }

    fn as_any(&self) -> &dyn std::any::Any { return self; }
}


