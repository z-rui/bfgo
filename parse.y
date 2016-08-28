%{
package main
%}

%union {
	intVal int
	stmt Statement
	block Block
}

%start prog

%type <intVal> "+" ">"
%type <stmt> stmt
%type <block> block

%%

prog	: block { yylex.(*lexState).ast = $1 }

block	: { $$ = nil }
	| block stmt { $$ = append($1, $2) }
	;

stmt	: "+" { $$ = IncInst($1) }
	| ">" { $$ = ShiftInst($1) }
	| "," { $$ = InInst{} }
	| "." { $$ = OutInst{} }
	| "[" block "]" { $$ = Loop($2) }
	;

%%

type (
	IncInst   int
	ShiftInst int
	InInst    struct{}
	OutInst   struct{}
	Statement interface{}
	Block     []Statement
	Loop      Block
)

