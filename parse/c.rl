package parse // import "docc.io/source/parse"

import "docc.io/source"

%%{
	machine c_blocks;
	write data;
}%%

// C sends all declarations and definitions from a C source code file to found.
func C(file []byte, found chan<- *source.Decl) {
	// Ragel setup
	var (
		data = file      // input
		cs   int         // state
		p    int         // data index
		pe   = len(data) // data end
	)

	// comment byte indices
	var commentOffset, commentEnd int

	// source code line number
	lineNo := 1
	// first byte index
	lineOffset := 0
	// source code block nesting
	blockLevel := 0

	for p < pe {
%%{

new_line = '\n' @{
	lineNo++
	lineOffset = p + 1
} ;


action comment_found {
	if commentEnd == 0 {
		commentOffset = p - 1
	}
}

action comment_set {
	commentEnd = p
}

action comment_clear {
	commentEnd = 0
}

action submit {
	if blockLevel == 0 {
		d := source.Decl{
			LineNo: lineNo,
			Source: string(data[lineOffset:p+1]),
		}
		if commentEnd > commentOffset {
			d.Comment = string(data[commentOffset:commentEnd])
		}
		found <- &d
	}
}

comment = '//' @comment_found ^new_line* new_line @comment_set
        space* new_line? @comment_clear
        ;

comment_block = '/*' @comment_found (new_line | any)* :>> '*/'
        space* new_line? @comment_set <: space* new_line? @comment_clear
        ;


# code block nesting
block_start = '{' @submit @comment_clear @{ blockLevel++ } ;
block_end = '}' @{ blockLevel-- } ;

line_end = ';' @submit @comment_clear ;


# parse quotations because they may contain special characters
qstring = '"' ([^"\\] | new_line | ('\\' any))* '"' ;
qchar = '\'' ('\\' any)? ^['\\]* '\'' ;

# skip preprocessor instructions
prep = '#' ^new_line* new_line ;

# faster than stepping out of machine
other = ! ('\n' | '/' ^[*/] | '{' | '}' | ';' | '"' | '\'' | '#') ;


main := (new_line | comment | comment_block
	| block_start | block_end | line_end
	| qstring | qchar | prep | other)* ;

write init;
write exec;

}%%

		// skip next byte for syntax mismatch case
		p++
	}

	close(found)
}
