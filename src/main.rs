use std::cell::RefCell;
use std::fs;
use std::sync::Arc;

use crate::interpreter::Interpreter;
use crate::lexer::Lexer;
use crate::parser::Parser;

mod lexer;
mod parser;
mod ast;
mod interpreter;
mod env;
mod error;
mod stdlib;

const VERSION: &str = "v1.0.0";

fn main() {
    if cfg!(debug_assertions) {
        println!("Totem debug build {}\n", VERSION);
    }
    
    let args: Vec<String> = std::env::args().collect();
    #[cfg(debug_assertions)]
    println!("args: {:?}\n", args);

    if args.len() == 1 { panic!("Specify an input file") }
    let code = fs::read_to_string(args[1].clone()).expect(&format!("File does not exist: {}", args[1]));
    let lex = Lexer::new(code);
    let tokens = lex.tokenize();

    #[cfg(debug_assertions)]
    println!("token stream: {:?}\n", tokens);
    let mut parser = Parser::new(tokens);
    let expression = parser.parse();

    if expression.len() > 0 {
        #[cfg(debug_assertions)]
        println!("ast: {:?}\n", expression);
        let mut interpreter = Interpreter::new();
        #[cfg(debug_assertions)]
        println!("------------Program Output--------------");
        interpreter.interpret(expression);
    }
}
