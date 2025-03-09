package baisl

type Token struct {
	TType    TokenType
	Location SourceLocation
	// Nil unless HasValue is true
	Value    string
	HasValue bool
}
