use std::process::exit;

use crate::lexer::Token;

#[derive(Debug)]
pub struct ParseError {
    message: String,
    line: i32,
    token: Token,
}
impl ParseError {
    pub fn new(message: &str, token: Token, line: i32) -> ParseError {
        return ParseError {
            message: message.to_string(),
            line,
            token,
        };
    }

    pub fn get_message(&self) -> String {
        let mut temp = String::from("Syntax error on line ");
        temp += &self.line.to_string();
        temp += " at token \"";
        temp += &self.token.value;
        temp += "\": ";
        temp += &self.message;

        return temp;
    }
}

// This needs to be changed if I create destructors, otherwise it's fine
pub fn runtime_error(message: &str, line: i32, value: &String) {
    println!(
        "Runtime error on line {0} at token \"{1}\": {2}",
        line, value, message
    );
    exit(1);
}
