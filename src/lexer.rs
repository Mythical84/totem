use std::process;

#[allow(non_camel_case_types)]
#[derive(Debug, PartialEq, Clone)]
pub enum TokenType {
    // Single character tokens
    ADD, SUB, MUL, DIV, MOD,
    EQUAL,LESS, GREATER, BANG,
    LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE,
    SEMI, COMMA, DOT,

    // Double character tokens
    BANG_EQUAL, DOUBLE_EQUAL, GREATER_EQUAL, LESS_EQUAL,
    DOUBLE_AMPERSAND, DOUBLE_LINE,

    // literals
    INT, STR, VAR, DOUBLE,

    // keywords
    TRUE, FALSE, NULL,
    IF, ELSE,
    LOOP, BREAK, CONTINUE, WHILE,
    CLASS, SELF, SUPER,
    FN, RETURN,
    LET, IMPORT,

    // End of file
    EOF,
}

#[derive(Debug, Clone)]
pub struct Token {
    pub token_type: TokenType,
    pub value: String,
    pub line: i32,
}

pub struct Lexer {
    current: usize,
    line: i32,
    tokens: Vec<Token>,
    chars: Vec<char>,
}

impl Lexer {
    pub fn new(code: String) -> Lexer {
        return Lexer {
            current: 0,
            line: 1,
            tokens: Vec::new(),
            chars: code.chars().collect(),
        };
    }

    pub fn tokenize(mut self) -> Vec<Token> {
        while self.current <= self.chars.len() {
            self.current += 1;
            if self.current > self.chars.len() {
                break;
            }
            let c = *self.chars.get(self.current - 1).unwrap();

            match c {
                // Single character lexemes
                ' ' => { continue; }
                '\t' => { continue; }
                '*' => { self.push_token(TokenType::MUL, c.to_string()); continue; }
                '+' => { self.push_token(TokenType::ADD, c.to_string()); continue; }
                '-' => { self.push_token(TokenType::SUB, c.to_string()); continue; }
                '/' => { self.push_token(TokenType::DIV, c.to_string()); continue; }
                '%' => { self.push_token(TokenType::MOD, c.to_string()); continue; }
                ';' => { self.push_token(TokenType::SEMI, c.to_string()); continue; }
                ',' => { self.push_token(TokenType::COMMA, c.to_string()); continue;}
                '.' => { self.push_token(TokenType::DOT, c.to_string()); continue; }
                '(' => { self.push_token(TokenType::LEFT_PAREN, c.to_string()); continue; }
                ')' => { self.push_token(TokenType::RIGHT_PAREN, c.to_string()); continue; }
                '{' => { self.push_token(TokenType::LEFT_BRACE, c.to_string()); continue; }
                '}' => { self.push_token(TokenType::RIGHT_BRACE, c.to_string()); continue; }

                // Double character lexemes
                // TODO: for code cleanliness reasons, it may be smart to write a macro for this
                '<' => {
                    if self.peek() == '=' {
                        self.push_token(TokenType::LESS_EQUAL, "<=".to_string());
                        self.current += 1;
                    } else {
                        self.push_token(TokenType::LESS, c.to_string())
                    }
                    continue;
                }
                '>' => {
                    if self.peek() == '=' {
                        self.push_token(TokenType::GREATER_EQUAL, ">=".to_string());
                        self.current += 1;
                    } else {
                        self.push_token(TokenType::GREATER, c.to_string())
                    }
                    continue;
                }
                '!' => {
                    if self.peek() == '=' {
                        self.push_token(TokenType::BANG_EQUAL, "!=".to_string());
                        self.current += 1;
                    } else {
                        self.push_token(TokenType::BANG, c.to_string())
                    }
                    continue;
                }
                '=' => {
                    if self.peek() == '=' {
                        self.push_token(TokenType::DOUBLE_EQUAL, "==".to_string());
                        self.current += 1;
                    } else {
                        self.push_token(TokenType::EQUAL, c.to_string())
                    }
                    continue;
                }
                '|' => {
                    if self.peek() == '|' {
                        self.push_token(TokenType::DOUBLE_LINE, "||".to_string());
                        self.current += 1;
                        continue;
                    }
                }
                '&' => {
                    if self.peek() == '&' {
                        self.push_token(TokenType::DOUBLE_AMPERSAND, "&&".to_string());
                        self.current += 1;
                        continue;
                    }
                }
                '\n' => {
                    self.line += 1;
                    continue;
                }

                // multi character lexemes
                _ => {
                    if c == '#' {
                        self.comment();
                    } else if c.is_digit(10) {
                        self.token_num();
                    } else if c == '"' {
                        self.token_string();
                    } else if c.is_alphabetic() || c == '-' || c == '_' {
                        self.token_identifier();
                    } else {
                        syntax_error(self.line, format!("Unexpected character: {}", c));
                    }
                }
            };
        }

        self.push_token(TokenType::EOF, "EOF".to_string());

        return self.tokens;
    }

    fn token_identifier(&mut self) {
        let c = self.chars.get(self.current - 1).unwrap();
        let mut token = c.to_string();
        while self.peek().is_alphanumeric() || self.peek() == '-' || self.peek() == '_' {
            self.current += 1;
            token += &self.chars.get(self.current - 1).unwrap().to_string();
        }

        if self.is_keyword(&token) {
            self.token_keyword(token);
        } else if token == "true".to_string() {
            self.push_token(TokenType::TRUE, token);
        } else if token == "false".to_string() {
            self.push_token(TokenType::FALSE, token);
        } else if token == "null".to_string() {
            self.push_token(TokenType::NULL, token);
        } else {
            self.push_token(TokenType::VAR, token);
        }
    }

    fn comment(&mut self) {
        while self.peek() != '\n' && self.current < self.chars.len() {
            self.current += 1;
        }
    }

    fn is_keyword(&self, token: &str) -> bool {
        let keywords: [&str; 13] = [
            "if", "else", "loop", "break", "continue", "class", "self", "super", "fn",
            "return", "let", "import", "while",
        ];

        return keywords.contains(&token);
    }

    fn token_keyword(&mut self, token: String) {
        let token_type;
        match token.as_str() {
            "if" => token_type = TokenType::IF,
            "else" => token_type = TokenType::ELSE,
            "loop" => token_type = TokenType::LOOP,
            "break" => token_type = TokenType::BREAK,
            "continue" => token_type = TokenType::CONTINUE,
            "while" => token_type = TokenType::WHILE,
            "class" => token_type = TokenType::CLASS,
            "self" => token_type = TokenType::SELF,
            "super" => token_type = TokenType::SUPER,
            "fn" => token_type = TokenType::FN,
            "return" => token_type = TokenType::RETURN,
            "let" => token_type = TokenType::LET,
            "import" => token_type = TokenType::IMPORT,
            // This will only fire if I fuck something up in the lexer
            _ => panic!("Invalid token"),
        }

        self.push_token(token_type, token);
    }

    fn token_string(&mut self) {
        let mut lexeme = String::new();
        while self.peek() != '"' {
            self.current += 1;
            lexeme += &self.chars.get(self.current - 1).unwrap().to_string();
        }
        self.current += 1;
        self.push_token(TokenType::STR, lexeme)
    }

    fn token_num(&mut self) {
        let mut float = false;
        let mut lexeme = self.chars.get(self.current - 1).unwrap().to_string();
        while self.peek().is_digit(10) || self.peek() == '.' {
            self.current += 1;
            lexeme += &self.chars.get(self.current - 1).unwrap().to_string();
            if self.peek() == '.' {
                float = true;
            }
        }

        if float {
            self.push_token(TokenType::DOUBLE, lexeme);
        } else {
            self.push_token(TokenType::INT, lexeme);
        }
    }

    fn peek(&self) -> char {
        if self.current > self.chars.len() {
            return '\0';
        }
        return *self.chars.get(self.current).unwrap();
    }

    fn push_token(&mut self, token_type: TokenType, value: String) {
        self.tokens.push(Token {
            token_type,
            value,
            line: self.line,
        });
    }
}

fn syntax_error(line: i32, message: String) {
    println!("Error on line {}: {}", line, message);
    process::exit(1);
}
