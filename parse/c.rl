package parse

%%{
	machine c_blocks;
	write data;
}%%


func C(data []byte) (lineCount int, comments []string) {
	// first byte index
	var commentOffset int
	// source code lines
	lineNo := 1
	// source code block nesting
	blockLevel := 0

	// Ragel state
	var cs int
	// Ragel data index
	var p int = 0
	// Ragel data end
	var pe int = len(data)

	for p < pe {
%%{
action comment_found {
	commentOffset = p - 1
}

action comment_submit {
	if blockLevel == 0 {
		comments = append(comments, string(data[commentOffset:p+1]))
	}
}


# source code is line oriented
new_line = '\n' @{ lineNo++ } ;

# preprocessor instruction
prep = '#' ^new_line* new_line ;

# comment line
comment = '//' @comment_found ^new_line* new_line @comment_submit ;

# block comment
comment_block = '/*' @comment_found (new_line | any)* :>> '*/' @comment_submit ;

# double quoted string
qstring = '"' ([^"\\] | new_line | ('\\' any))* '"' ;

# single quoted character
qchar = '\'' ('\\' any)? ^['\\]* '\'' ;

# code block nesting
block_start = '{' @{ blockLevel++ } ;
block_end = '}' @{ blockLevel-- } ;

main := (new_line | comment | comment_block | qstring | qchar | block_start | block_end)* ;

write init;
write exec;
}%%

		p++	// skip next byte
	}

	lineCount = lineNo
	return
}
