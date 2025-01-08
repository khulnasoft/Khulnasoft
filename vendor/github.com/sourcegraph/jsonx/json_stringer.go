// Code generated by "stringer -type=ParseErrorCode,ScanErrorCode,SyntaxKind,NodeType -output=json_stringer.go"; DO NOT EDIT.

package jsonx

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InvalidSymbol-0]
	_ = x[InvalidNumberFormat-1]
	_ = x[PropertyNameExpected-2]
	_ = x[ValueExpected-3]
	_ = x[ColonExpected-4]
	_ = x[CommaExpected-5]
	_ = x[CloseBraceExpected-6]
	_ = x[CloseBracketExpected-7]
	_ = x[EndOfFileExpected-8]
	_ = x[InvalidCommentToken-9]
	_ = x[ParseErrorUnexpectedEndOfComment-10]
	_ = x[ParseErrorUnexpectedEndOfString-11]
	_ = x[ParseErrorUnexpectedEndOfNumber-12]
	_ = x[ParseErrorInvalidUnicode-13]
	_ = x[ParseErrorInvalidEscapeCharacter-14]
	_ = x[ParseErrorInvalidCharacter-15]
	_ = x[InvalidScanErrorCode-16]
}

const _ParseErrorCode_name = "InvalidSymbolInvalidNumberFormatPropertyNameExpectedValueExpectedColonExpectedCommaExpectedCloseBraceExpectedCloseBracketExpectedEndOfFileExpectedInvalidCommentTokenParseErrorUnexpectedEndOfCommentParseErrorUnexpectedEndOfStringParseErrorUnexpectedEndOfNumberParseErrorInvalidUnicodeParseErrorInvalidEscapeCharacterParseErrorInvalidCharacterInvalidScanErrorCode"

var _ParseErrorCode_index = [...]uint16{0, 13, 32, 52, 65, 78, 91, 109, 129, 146, 165, 197, 228, 259, 283, 315, 341, 361}

func (i ParseErrorCode) String() string {
	if i < 0 || i >= ParseErrorCode(len(_ParseErrorCode_index)-1) {
		return "ParseErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ParseErrorCode_name[_ParseErrorCode_index[i]:_ParseErrorCode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[UnexpectedEndOfComment-1]
	_ = x[UnexpectedEndOfString-2]
	_ = x[UnexpectedEndOfNumber-3]
	_ = x[InvalidUnicode-4]
	_ = x[InvalidEscapeCharacter-5]
	_ = x[InvalidCharacter-6]
}

const _ScanErrorCode_name = "NoneUnexpectedEndOfCommentUnexpectedEndOfStringUnexpectedEndOfNumberInvalidUnicodeInvalidEscapeCharacterInvalidCharacter"

var _ScanErrorCode_index = [...]uint8{0, 4, 26, 47, 68, 82, 104, 120}

func (i ScanErrorCode) String() string {
	if i < 0 || i >= ScanErrorCode(len(_ScanErrorCode_index)-1) {
		return "ScanErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ScanErrorCode_name[_ScanErrorCode_index[i]:_ScanErrorCode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[OpenBraceToken-1]
	_ = x[CloseBraceToken-2]
	_ = x[OpenBracketToken-3]
	_ = x[CloseBracketToken-4]
	_ = x[CommaToken-5]
	_ = x[ColonToken-6]
	_ = x[NullKeyword-7]
	_ = x[TrueKeyword-8]
	_ = x[FalseKeyword-9]
	_ = x[StringLiteral-10]
	_ = x[NumericLiteral-11]
	_ = x[LineCommentTrivia-12]
	_ = x[BlockCommentTrivia-13]
	_ = x[LineBreakTrivia-14]
	_ = x[Trivia-15]
	_ = x[EOF-16]
}

const _SyntaxKind_name = "UnknownOpenBraceTokenCloseBraceTokenOpenBracketTokenCloseBracketTokenCommaTokenColonTokenNullKeywordTrueKeywordFalseKeywordStringLiteralNumericLiteralLineCommentTriviaBlockCommentTriviaLineBreakTriviaTriviaEOF"

var _SyntaxKind_index = [...]uint8{0, 7, 21, 36, 52, 69, 79, 89, 100, 111, 123, 136, 150, 167, 185, 200, 206, 209}

func (i SyntaxKind) String() string {
	if i < 0 || i >= SyntaxKind(len(_SyntaxKind_index)-1) {
		return "SyntaxKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SyntaxKind_name[_SyntaxKind_index[i]:_SyntaxKind_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Object-0]
	_ = x[Array-1]
	_ = x[Property-2]
	_ = x[String-3]
	_ = x[Number-4]
	_ = x[Boolean-5]
	_ = x[Null-6]
}

const _NodeType_name = "ObjectArrayPropertyStringNumberBooleanNull"

var _NodeType_index = [...]uint8{0, 6, 11, 19, 25, 31, 38, 42}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
