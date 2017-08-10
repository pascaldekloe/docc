package docc

// Decl is a C declaration.
type Decl struct {
	// LineNo is the first line number.
	LineNo int
	// Source is the original text.
	Source string
	// Comment is the original comment text.
	Comment string
}

// File is a source code file.
type File struct {
	// Path is the file location.
	Path string
	// Vars are the variable declarations.
	Vars []*Decl
	// Funcs are the function declarations.
	Funcs []*Decl
}
