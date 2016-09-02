%{

/*
 *  R : A Computer Langage for Statistical Data Analysis
 *  Copyright (C) 1995, 1996, 1997  Robert Gentleman and Ross Ihaka
 *  Copyright (C) 1997--2000  Robert Gentleman, Ross Ihaka and the
 *                            R Development Core Team
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program; if not, write to the Free Software
 *  Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */

%}

%token		END_OF_INPUT ERROR
%token		STR_CONST NUM_CONST NULL_CONST SYMBOL FUNCTION
%token		LEFT_ASSIGN RIGHT_ASSIGN LBB
%token		FOR IN IF ELSE WHILE NEXT BREAK REPEAT
%token		GT GE LT LE EQ NE AND OR

%left		'?'
%left		LOW WHILE FOR REPEAT
%right		IF
%left		ELSE
%right		LEFT_ASSIGN
%left		RIGHT_ASSIGN
%left		'~' TILDE
%left		OR
%left		AND
%left		UNOT NOT
%left   	GT GE LT LE EQ NE
%left		'+' '-'
%left		'*' '/' '%'
%left		SPECIAL
%left		':'
%left		UMINUS UPLUS
%right		'^'
%left		'$'
%nonassoc	'(' '[' LBB

%%

prog	:	END_OF_INPUT
	|	'\n'
	|	expr '\n'
	|	expr ';'
	|	error
	;

expr	: 	NUM_CONST
	|	STR_CONST
	|	NULL_CONST
	|	SYMBOL
	|	'{' exprlist '}'
	|	'(' expr ')'
	|	'-' expr %prec UMINUS
	|	'+' expr %prec UMINUS
	|	'!' expr %prec UNOT
	|	'~' expr %prec TILDE
	|	'?' expr
	|	expr ':'  expr
	|	expr '+'  expr
	|	expr '-' expr
	|	expr '*' expr
	|	expr '/' expr
	|	expr '^' expr
	|	expr SPECIAL expr
	|	expr '%' expr
	|	expr '~' expr
	|	expr LT expr
	|	expr LE expr
	|	expr EQ expr
	|	expr NE expr
	|	expr GE expr
	|	expr GT expr
	|	expr AND expr
	|	expr OR expr
	|	expr LEFT_ASSIGN expr
	|	expr RIGHT_ASSIGN expr
	|	FUNCTION '(' formlist ')' cr expr %prec LOW
	|	expr '(' sublist ')'
	|	IF ifcond expr
	|	IF ifcond expr ELSE expr
	|	FOR forcond expr %prec FOR
	|	WHILE cond expr
	|	REPEAT expr
	|	expr LBB sublist ']' ']'
	|	expr '[' sublist ']'
	|	expr '$' SYMBOL
	|	expr '$' STR_CONST
	|	NEXT
	|	BREAK
	;

cond	:	'(' expr ')'
	;

ifcond	:	'(' expr ')'
	;

forcond :	'(' SYMBOL IN expr ')'
	;

exprlist:
	|	expr
	|	exprlist ';' expr
	|	exprlist ';'
	|	exprlist '\n' expr
	|	exprlist '\n'
	;

sublist	:	sub
	|	sublist cr ',' sub
	;

sub	:
	|	expr
	|	SYMBOL '='
	|	SYMBOL '=' expr
	|	STR_CONST '='
	|	STR_CONST '=' expr
	|	NULL_CONST '='
	|	NULL_CONST '=' expr
	;

formlist:
	|	SYMBOL
	|	SYMBOL '=' expr
	|	formlist ',' SYMBOL
	|	formlist ',' SYMBOL '=' expr
	;

cr	:
	;
%%
