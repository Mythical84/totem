# This is a script that automatically generates a visitor pattern in ast.rs
# Honestly it's probably unnecessary but I made it so I'm gonna use it

expr_types = [
	["Binary", "Expr left, Token operator, Expr right"],
	["Call", "Expr callee, Vec<Expr> args, i32 line"],
	["Get", "Expr object, Token name"],
	["Set", "Expr object, Token field_name, Expr value"],
	["Grouping", "Expr expression"],
	["Literal", "Token token"],
	["Logical", "Expr left, Token operator, Expr right"],
	["Unary", "Token operator, Expr right"],
	["Variable", "Token name"],
	["VarUpdate", "Token name, Expr expr"],
	["This", "Token keyword"]
]

stmt_types = [
	["VarDefine", "Token name, Option<Expr> expr"],
	["FnDefine", "Token name, Vec<String> args, Stmt body, Option<String> this"],
	["ClassDefine", "Token name, Vec<Stmt> methods"],
	["Return", "Expr value"],
	["Block", "Vec<Stmt> statements"],
	["Expression", "Expr expr"],
	["If", "Expr condition, Stmt then_branch, Box<Option<Stmt>> else_branch"],
	["Loop", "Stmt body"],
	["While", "Expr condition, Stmt body"],
	["ControlStmt", "Token token"],
]

header = "// This is an automatically generated file\nuse crate::lexer::Token;\n\n"
enums = []
impls = []
visitor = """pub trait Visitor<T> {
    fn evaluate(&mut self, expr: &Expr) -> T;\n
    fn execute(&mut self, stmt: &Stmt) -> T;\n"""

def def_ast(class_name, types):
	global visitor
	enums.append(f"#[derive(Debug, Clone)]\npub enum {class_name} {"{"}\n")
	impls.append(f"""impl {class_name} {"{"}
    pub fn accept<T>(&self, visitor: &mut dyn Visitor<T>) -> T {"{"}
        match self {"{"}\n""")

	for t in types:
		name = t[0]
		fields = t[1].split(', ')
		enums[-1] += "    " + name + " {\n"
		impls[-1] += f"            {class_name}::{name}{"{"} "
		visitor += f"    fn visit_{name.lower()}{"("}&mut self, "
		names = ""
		for f in fields:
			split = f.split(" ")
			var_type = split[0]
			var_type_box = split[0]
			var_name = split[1]
			if var_type_box == class_name: var_type_box = f"Box<{class_name}>"
			enums[-1] += f"        {var_name}: {var_type_box},\n"
			visitor += f"{var_name}: &{var_type}"
			names += f"{var_name}"
			if not f == fields[-1]: 
				visitor += ", "
				names += ", "
		impls[-1] += names + " } => visitor.visit_" + name.lower() + "(" + names + ")" + ",\n" 
		enums[-1] += "    },\n"
		visitor += ") -> T;\n"
	enums[-1] += "}\n"
	impls[-1] += "        }\n    }\n}\n\n"


def_ast("Expr", expr_types)
def_ast("Stmt", stmt_types)

visitor += "}\n\n"

with open("ast.rs", "w") as file:
	file.write(header)
	for e, i in zip(enums, impls):
		file.write(e)
		file.write(i)
	file.write(visitor)
