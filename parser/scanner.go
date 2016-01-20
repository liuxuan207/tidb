// CAUTION: Generated file - DO NOT EDIT.

// Copyright 2013 The ql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSES/QL-LICENSE file.

// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/mysql"
	"github.com/pingcap/tidb/util/charset"
	"github.com/pingcap/tidb/util/stringutil"
)

type lexer struct {
	c            int
	col          int
	errs         []error
	expr         ast.ExprNode
	i            int
	inj          int
	lcol         int
	line         int
	list         []ast.StmtNode
	ncol         int
	nline        int
	sc           int
	src          string
	val          []byte
	ungetBuf     []byte
	root         bool
	prepare      bool
	stmtStartPos int
	stringLit    []byte

	// record token's offset of the input
	tokenEndOffset   int
	tokenStartOffset int

	// Charset information
	charset   string
	collation string
}

// NewLexer builds a new lexer.
func NewLexer(src string) (l *lexer) {
	l = &lexer{
		src:   src,
		nline: 1,
		ncol:  0,
	}
	l.next()
	return
}

func (l *lexer) Errors() []error {
	return l.errs
}

func (l *lexer) Stmts() []ast.StmtNode {
	return l.list
}

func (l *lexer) Expr() ast.ExprNode {
	return l.expr
}

func (l *lexer) Inj() int {
	return l.inj
}

func (l *lexer) SetInj(inj int) {
	l.inj = inj
}

func (l *lexer) SetPrepare() {
	l.prepare = true
}

func (l *lexer) IsPrepare() bool {
	return l.prepare
}

func (l *lexer) Root() bool {
	return l.root
}

func (l *lexer) SetRoot(root bool) {
	l.root = root
}

func (l *lexer) SetCharsetInfo(charset, collation string) {
	l.charset = charset
	l.collation = collation
}

func (l *lexer) GetCharsetInfo() (string, string) {
	return l.charset, l.collation
}

// The select statement is not at the end of the whole statement, if the last
// field text was set from its offset to the end of the src string, update
// the last field text.
func (l *lexer) SetLastSelectFieldText(st *ast.SelectStmt, lastEnd int) {
	lastField := st.Fields.Fields[len(st.Fields.Fields)-1]
	if lastField.Offset+len(lastField.Text()) >= len(l.src)-1 {
		lastField.SetText(l.src[lastField.Offset:lastEnd])
	}
}

func (l *lexer) startOffset(offset int) int {
	offset--
	for unicode.IsSpace(rune(l.src[offset])) {
		offset++
	}
	return offset
}

func (l *lexer) endOffset(offset int) int {
	offset--
	for offset > 0 && unicode.IsSpace(rune(l.src[offset-1])) {
		offset--
	}
	return offset
}

func (l *lexer) unget(b byte) {
	l.ungetBuf = append(l.ungetBuf, b)
	l.i--
	l.ncol--
	l.tokenEndOffset--
}

func (l *lexer) next() int {
	if un := len(l.ungetBuf); un > 0 {
		nc := l.ungetBuf[0]
		l.ungetBuf = l.ungetBuf[1:]
		l.c = int(nc)
		return l.c
	}

	if l.c != 0 {
		l.val = append(l.val, byte(l.c))
	}
	l.c = 0
	if l.i < len(l.src) {
		l.c = int(l.src[l.i])
		l.i++
	}
	switch l.c {
	case '\n':
		l.lcol = l.ncol
		l.nline++
		l.ncol = 0
	default:
		l.ncol++
	}
	l.tokenEndOffset++
	return l.c
}

func (l *lexer) err0(ln, c int, arg interface{}) {
	var argStr string
	if arg != nil {
		argStr = fmt.Sprintf(" %v", arg)
	}

	err := fmt.Errorf("line %d column %d near \"%s\"%s", ln, c, l.val, argStr)
	l.errs = append(l.errs, err)
}

func (l *lexer) err(arg interface{}) {
	l.err0(l.line, l.col, arg)
}

func (l *lexer) errf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	l.err0(l.line, l.col, s)
}

func (l *lexer) Error(s string) {
	// Notice: ignore origin error info.
	l.err(nil)
}

func (l *lexer) stmtText() string {
	endPos := l.i
	if l.src[l.i-1] == '\n' {
		endPos = l.i - 1 // trim new line
	}
	if l.src[l.stmtStartPos] == '\n' {
		l.stmtStartPos++
	}

	text := l.src[l.stmtStartPos:endPos]

	l.stmtStartPos = l.i
	return text
}

func (l *lexer) Lex(lval *yySymType) (r int) {
	defer func() {
		lval.line, lval.col, lval.offset = l.line, l.col, l.tokenStartOffset
		l.tokenStartOffset = l.tokenEndOffset
	}()
	const (
		INITIAL = iota
		S1
		S2
		S3
		S4
	)

	if n := l.inj; n != 0 {
		l.inj = 0
		return n
	}

	c0, c := 0, l.c

yystate0:

	l.val = l.val[:0]
	c0, l.line, l.col = l.c, l.nline, l.ncol

	switch yyt := l.sc; yyt {
	default:
		panic(fmt.Errorf(`invalid start condition %d`, yyt))
	case 0: // start condition: INITIAL
		goto yystart1
	case 1: // start condition: S1
		goto yystart1214
	case 2: // start condition: S2
		goto yystart1220
	case 3: // start condition: S3
		goto yystart1226
	case 4: // start condition: S4
		goto yystart1229
	}

	goto yystate0 // silence unused label error
	goto yystate1 // silence unused label error
yystate1:
	c = l.next()
yystart1:
	switch {
	default:
		goto yystate3 // c >= '\x01' && c <= '\b' || c == '\v' || c == '\f' || c >= '\x0e' && c <= '\x1f' || c == '$' || c == '%%' || c >= '(' && c <= ',' || c == ':' || c == ';' || c >= '[' && c <= '^' || c == '{' || c >= '}' && c <= 'ÿ'
	case c == '!':
		goto yystate6
	case c == '"':
		goto yystate8
	case c == '#':
		goto yystate9
	case c == '&':
		goto yystate11
	case c == '-':
		goto yystate15
	case c == '.':
		goto yystate17
	case c == '/':
		goto yystate22
	case c == '0':
		goto yystate27
	case c == '<':
		goto yystate36
	case c == '=':
		goto yystate41
	case c == '>':
		goto yystate42
	case c == '?':
		goto yystate45
	case c == '@':
		goto yystate46
	case c == 'A' || c == 'a':
		goto yystate65
	case c == 'B' || c == 'b':
		goto yystate118
	case c == 'C' || c == 'c':
		goto yystate155
	case c == 'D' || c == 'd':
		goto yystate285
	case c == 'E' || c == 'e':
		goto yystate423
	case c == 'F' || c == 'f':
		goto yystate461
	case c == 'G' || c == 'g':
		goto yystate502
	case c == 'H' || c == 'h':
		goto yystate523
	case c == 'I' || c == 'i':
		goto yystate566
	case c == 'J' || c == 'j':
		goto yystate615
	case c == 'K' || c == 'k':
		goto yystate619
	case c == 'L' || c == 'l':
		goto yystate633
	case c == 'M' || c == 'm':
		goto yystate693
	case c == 'N' || c == 'n':
		goto yystate760
	case c == 'O' || c == 'o':
		goto yystate784
	case c == 'P' || c == 'p':
		goto yystate806
	case c == 'Q' || c == 'q':
		goto yystate842
	case c == 'R' || c == 'r':
		goto yystate852
	case c == 'S' || c == 's':
		goto yystate900
	case c == 'T' || c == 't':
		goto yystate1015
	case c == 'U' || c == 'u':
		goto yystate1078
	case c == 'V' || c == 'v':
		goto yystate1124
	case c == 'W' || c == 'w':
		goto yystate1153
	case c == 'X' || c == 'x':
		goto yystate1182
	case c == 'Y' || c == 'y':
		goto yystate1188
	case c == 'Z' || c == 'z':
		goto yystate1202
	case c == '\'':
		goto yystate14
	case c == '\n':
		goto yystate5
	case c == '\t' || c == '\r' || c == ' ':
		goto yystate4
	case c == '\x00':
		goto yystate2
	case c == '_':
		goto yystate1210
	case c == '`':
		goto yystate1211
	case c == '|':
		goto yystate1212
	case c >= '1' && c <= '9':
		goto yystate34
	}

yystate2:
	c = l.next()
	goto yyrule1

yystate3:
	c = l.next()
	goto yyrule310

yystate4:
	c = l.next()
	switch {
	default:
		goto yyrule2
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate5
	}

yystate5:
	c = l.next()
	switch {
	default:
		goto yyrule2
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate5
	}

yystate6:
	c = l.next()
	switch {
	default:
		goto yyrule310
	case c == '=':
		goto yystate7
	}

yystate7:
	c = l.next()
	goto yyrule33

yystate8:
	c = l.next()
	goto yyrule13

yystate9:
	c = l.next()
	switch {
	default:
		goto yyrule3
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate10
	}

yystate10:
	c = l.next()
	switch {
	default:
		goto yyrule3
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate10
	}

yystate11:
	c = l.next()
	switch {
	default:
		goto yyrule310
	case c == '&':
		goto yystate12
	case c == '^':
		goto yystate13
	}

yystate12:
	c = l.next()
	goto yyrule27

yystate13:
	c = l.next()
	goto yyrule28

yystate14:
	c = l.next()
	goto yyrule14

yystate15:
	c = l.next()
	switch {
	default:
		goto yyrule310
	case c == '-':
		goto yystate16
	}

yystate16:
	c = l.next()
	goto yyrule6

yystate17:
	c = l.next()
	switch {
	default:
		goto yyrule310
	case c >= '0' && c <= '9':
		goto yystate18
	}

yystate18:
	c = l.next()
	switch {
	default:
		goto yyrule10
	case c == 'E' || c == 'e':
		goto yystate19
	case c >= '0' && c <= '9':
		goto yystate18
	}

yystate19:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '+' || c == '-':
		goto yystate20
	case c >= '0' && c <= '9':
		goto yystate21
	}

yystate20:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate21
	}

yystate21:
	c = l.next()
	switch {
	default:
		goto yyrule10
	case c >= '0' && c <= '9':
		goto yystate21
	}

yystate22:
	c = l.next()
	switch {
	default:
		goto yyrule310
	case c == '*':
		goto yystate23
	case c == '/':
		goto yystate26
	}

yystate23:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate24
	case c >= '\x01' && c <= ')' || c >= '+' && c <= 'ÿ':
		goto yystate23
	}

yystate24:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate24
	case c == '/':
		goto yystate25
	case c >= '\x01' && c <= ')' || c >= '+' && c <= '.' || c >= '0' && c <= 'ÿ':
		goto yystate23
	}

yystate25:
	c = l.next()
	goto yyrule5

yystate26:
	c = l.next()
	switch {
	default:
		goto yyrule4
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate26
	}

yystate27:
	c = l.next()
	switch {
	default:
		goto yyrule9
	case c == '.':
		goto yystate18
	case c == '8' || c == '9':
		goto yystate29
	case c == 'B' || c == 'b':
		goto yystate30
	case c == 'E' || c == 'e':
		goto yystate19
	case c == 'X' || c == 'x':
		goto yystate32
	case c >= '0' && c <= '7':
		goto yystate28
	}

yystate28:
	c = l.next()
	switch {
	default:
		goto yyrule9
	case c == '.':
		goto yystate18
	case c == '8' || c == '9':
		goto yystate29
	case c == 'E' || c == 'e':
		goto yystate19
	case c >= '0' && c <= '7':
		goto yystate28
	}

yystate29:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate18
	case c == 'E' || c == 'e':
		goto yystate19
	case c >= '0' && c <= '9':
		goto yystate29
	}

yystate30:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '0' || c == '1':
		goto yystate31
	}

yystate31:
	c = l.next()
	switch {
	default:
		goto yyrule12
	case c == '0' || c == '1':
		goto yystate31
	}

yystate32:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate33
	}

yystate33:
	c = l.next()
	switch {
	default:
		goto yyrule11
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate33
	}

yystate34:
	c = l.next()
	switch {
	default:
		goto yyrule9
	case c == '.':
		goto yystate18
	case c == 'E' || c == 'e':
		goto yystate19
	case c >= '0' && c <= '9':
		goto yystate35
	}

yystate35:
	c = l.next()
	switch {
	default:
		goto yyrule9
	case c == '.':
		goto yystate18
	case c == 'E' || c == 'e':
		goto yystate19
	case c >= '0' && c <= '9':
		goto yystate35
	}

yystate36:
	c = l.next()
	switch {
	default:
		goto yyrule310
	case c == '<':
		goto yystate37
	case c == '=':
		goto yystate38
	case c == '>':
		goto yystate40
	}

yystate37:
	c = l.next()
	goto yyrule29

yystate38:
	c = l.next()
	switch {
	default:
		goto yyrule30
	case c == '>':
		goto yystate39
	}

yystate39:
	c = l.next()
	goto yyrule37

yystate40:
	c = l.next()
	goto yyrule34

yystate41:
	c = l.next()
	goto yyrule31

yystate42:
	c = l.next()
	switch {
	default:
		goto yyrule310
	case c == '=':
		goto yystate43
	case c == '>':
		goto yystate44
	}

yystate43:
	c = l.next()
	goto yyrule32

yystate44:
	c = l.next()
	goto yyrule36

yystate45:
	c = l.next()
	goto yyrule39

yystate46:
	c = l.next()
	switch {
	default:
		goto yyrule38
	case c == '@':
		goto yystate47
	case c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate64
	}

yystate47:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == 'G' || c == 'g':
		goto yystate49
	case c == 'L' || c == 'l':
		goto yystate56
	case c == 'S' || c == 's':
		goto yystate58
	case c >= 'A' && c <= 'F' || c >= 'H' && c <= 'K' || c >= 'M' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'k' || c >= 'm' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate48
	}

yystate48:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate48
	}

yystate49:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate48
	case c == 'L' || c == 'l':
		goto yystate50
	}

yystate50:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate48
	case c == 'O' || c == 'o':
		goto yystate51
	}

yystate51:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate48
	case c == 'B' || c == 'b':
		goto yystate52
	}

yystate52:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate48
	case c == 'A' || c == 'a':
		goto yystate53
	}

yystate53:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate48
	case c == 'L' || c == 'l':
		goto yystate54
	}

yystate54:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate48
	case c == '.':
		goto yystate55
	}

yystate55:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate48
	}

yystate56:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate48
	case c == 'O' || c == 'o':
		goto yystate57
	}

yystate57:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate48
	case c == 'C' || c == 'c':
		goto yystate52
	}

yystate58:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate48
	case c == 'E' || c == 'e':
		goto yystate59
	}

yystate59:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate48
	case c == 'S' || c == 's':
		goto yystate60
	}

yystate60:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate48
	case c == 'S' || c == 's':
		goto yystate61
	}

yystate61:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate48
	case c == 'I' || c == 'i':
		goto yystate62
	}

yystate62:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate48
	case c == 'O' || c == 'o':
		goto yystate63
	}

yystate63:
	c = l.next()
	switch {
	default:
		goto yyrule216
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate48
	case c == 'N' || c == 'n':
		goto yystate54
	}

yystate64:
	c = l.next()
	switch {
	default:
		goto yyrule217
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate64
	}

yystate65:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'C' || c == 'E' || c >= 'G' && c <= 'K' || c == 'M' || c >= 'O' && c <= 'R' || c == 'T' || c >= 'W' && c <= 'Z' || c == '_' || c == 'a' || c == 'c' || c == 'e' || c >= 'g' && c <= 'k' || c == 'm' || c >= 'o' && c <= 'r' || c == 't' || c >= 'w' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate67
	case c == 'D' || c == 'd':
		goto yystate69
	case c == 'F' || c == 'f':
		goto yystate78
	case c == 'L' || c == 'l':
		goto yystate82
	case c == 'N' || c == 'n':
		goto yystate87
	case c == 'S' || c == 's':
		goto yystate90
	case c == 'U' || c == 'u':
		goto yystate92
	case c == 'V' || c == 'v':
		goto yystate105
	}

yystate66:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate67:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate68
	}

yystate68:
	c = l.next()
	switch {
	default:
		goto yyrule40
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate69:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate70
	case c == 'M' || c == 'm':
		goto yystate75
	}

yystate70:
	c = l.next()
	switch {
	default:
		goto yyrule41
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate71
	}

yystate71:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate72
	}

yystate72:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate73
	}

yystate73:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate74
	}

yystate74:
	c = l.next()
	switch {
	default:
		goto yyrule42
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate75:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate76
	}

yystate76:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate77
	}

yystate77:
	c = l.next()
	switch {
	default:
		goto yyrule43
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate78:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate79
	}

yystate79:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate80
	}

yystate80:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate81
	}

yystate81:
	c = l.next()
	switch {
	default:
		goto yyrule44
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate82:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate83
	case c == 'T' || c == 't':
		goto yystate84
	}

yystate83:
	c = l.next()
	switch {
	default:
		goto yyrule45
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate84:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate85
	}

yystate85:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate86
	}

yystate86:
	c = l.next()
	switch {
	default:
		goto yyrule46
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate87:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate88
	case c == 'Y' || c == 'y':
		goto yystate89
	}

yystate88:
	c = l.next()
	switch {
	default:
		goto yyrule47
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate89:
	c = l.next()
	switch {
	default:
		goto yyrule48
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate90:
	c = l.next()
	switch {
	default:
		goto yyrule50
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate91
	}

yystate91:
	c = l.next()
	switch {
	default:
		goto yyrule49
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate92:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate93
	}

yystate93:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate94
	}

yystate94:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate95
	}

yystate95:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate96
	}

yystate96:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate97
	}

yystate97:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate98
	}

yystate98:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate99
	}

yystate99:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate100
	}

yystate100:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate101
	}

yystate101:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate102
	}

yystate102:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate103
	}

yystate103:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate104
	}

yystate104:
	c = l.next()
	switch {
	default:
		goto yyrule51
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate105:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate106
	}

yystate106:
	c = l.next()
	switch {
	default:
		goto yyrule52
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate107
	}

yystate107:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate108
	}

yystate108:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate109
	}

yystate109:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate110
	}

yystate110:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate111
	}

yystate111:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate112
	}

yystate112:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate113
	}

yystate113:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate114
	}

yystate114:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate115
	}

yystate115:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate116
	}

yystate116:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate117
	}

yystate117:
	c = l.next()
	switch {
	default:
		goto yyrule53
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate118:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'H' || c == 'J' || c == 'K' || c == 'M' || c == 'N' || c >= 'P' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'h' || c == 'j' || c == 'k' || c == 'm' || c == 'n' || c >= 'p' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate122
	case c == 'I' || c == 'i':
		goto yystate131
	case c == 'L' || c == 'l':
		goto yystate141
	case c == 'O' || c == 'o':
		goto yystate144
	case c == 'Y' || c == 'y':
		goto yystate152
	case c == '\'':
		goto yystate119
	}

yystate119:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '0' || c == '1':
		goto yystate120
	}

yystate120:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '0' || c == '1':
		goto yystate120
	case c == '\'':
		goto yystate121
	}

yystate121:
	c = l.next()
	goto yyrule12

yystate122:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate123
	case c == 'T' || c == 't':
		goto yystate126
	}

yystate123:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate124
	}

yystate124:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate125
	}

yystate125:
	c = l.next()
	switch {
	default:
		goto yyrule54
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate126:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate127
	}

yystate127:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate128
	}

yystate128:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate129
	}

yystate129:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate130
	}

yystate130:
	c = l.next()
	switch {
	default:
		goto yyrule55
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate131:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'M' || c >= 'O' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'm' || c >= 'o' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate132
	case c == 'N' || c == 'n':
		goto yystate136
	case c == 'T' || c == 't':
		goto yystate140
	}

yystate132:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate133
	}

yystate133:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate134
	}

yystate134:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate135
	}

yystate135:
	c = l.next()
	switch {
	default:
		goto yyrule280
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate136:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate137
	}

yystate137:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate138
	}

yystate138:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate139
	}

yystate139:
	c = l.next()
	switch {
	default:
		goto yyrule294
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate140:
	c = l.next()
	switch {
	default:
		goto yyrule275
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate141:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate142
	}

yystate142:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate143
	}

yystate143:
	c = l.next()
	switch {
	default:
		goto yyrule297
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate144:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate145
	case c == 'T' || c == 't':
		goto yystate150
	}

yystate145:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate146
	}

yystate146:
	c = l.next()
	switch {
	default:
		goto yyrule304
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate147
	}

yystate147:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate148
	}

yystate148:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate149
	}

yystate149:
	c = l.next()
	switch {
	default:
		goto yyrule305
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate150:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate151
	}

yystate151:
	c = l.next()
	switch {
	default:
		goto yyrule56
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate152:
	c = l.next()
	switch {
	default:
		goto yyrule57
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate153
	}

yystate153:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate154
	}

yystate154:
	c = l.next()
	switch {
	default:
		goto yyrule306
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate155:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'G' || c >= 'I' && c <= 'N' || c == 'P' || c == 'Q' || c == 'S' || c == 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'g' || c >= 'i' && c <= 'n' || c == 'p' || c == 'q' || c == 's' || c == 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate156
	case c == 'H' || c == 'h':
		goto yystate160
	case c == 'O' || c == 'o':
		goto yystate177
	case c == 'R' || c == 'r':
		goto yystate245
	case c == 'U' || c == 'u':
		goto yystate253
	}

yystate156:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate157
	}

yystate157:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate158
	case c == 'T' || c == 't':
		goto yystate159
	}

yystate158:
	c = l.next()
	switch {
	default:
		goto yyrule58
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate159:
	c = l.next()
	switch {
	default:
		goto yyrule59
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate160:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate161
	case c == 'E' || c == 'e':
		goto yystate171
	}

yystate161:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate162
	}

yystate162:
	c = l.next()
	switch {
	default:
		goto yyrule292
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate163
	case c == 'S' || c == 's':
		goto yystate168
	}

yystate163:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate164
	}

yystate164:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate165
	}

yystate165:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate166
	}

yystate166:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate167
	}

yystate167:
	c = l.next()
	switch {
	default:
		goto yyrule60
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate168:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate169
	}

yystate169:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate170
	}

yystate170:
	c = l.next()
	switch {
	default:
		goto yyrule61
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate171:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate172
	}

yystate172:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate173
	}

yystate173:
	c = l.next()
	switch {
	default:
		goto yyrule62
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate174
	}

yystate174:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate175
	}

yystate175:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate176
	}

yystate176:
	c = l.next()
	switch {
	default:
		goto yyrule63
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate177:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'K' || c >= 'O' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'k' || c >= 'o' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate178
	case c == 'L' || c == 'l':
		goto yystate184
	case c == 'M' || c == 'm':
		goto yystate196
	case c == 'N' || c == 'n':
		goto yystate214
	case c == 'U' || c == 'u':
		goto yystate242
	}

yystate178:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate179
	}

yystate179:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate180
	}

yystate180:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate181
	}

yystate181:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate182
	}

yystate182:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate183
	}

yystate183:
	c = l.next()
	switch {
	default:
		goto yyrule64
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate184:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate185
	case c == 'U' || c == 'u':
		goto yystate192
	}

yystate185:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate186
	}

yystate186:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate187
	}

yystate187:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate188
	case c == 'I' || c == 'i':
		goto yystate189
	}

yystate188:
	c = l.next()
	switch {
	default:
		goto yyrule65
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate189:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate190
	}

yystate190:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate191
	}

yystate191:
	c = l.next()
	switch {
	default:
		goto yyrule66
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate192:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate193
	}

yystate193:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate194
	}

yystate194:
	c = l.next()
	switch {
	default:
		goto yyrule67
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate195
	}

yystate195:
	c = l.next()
	switch {
	default:
		goto yyrule68
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate196:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c == 'N' || c == 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c == 'n' || c == 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate197
	case c == 'P' || c == 'p':
		goto yystate206
	}

yystate197:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate198
	case c == 'I' || c == 'i':
		goto yystate201
	}

yystate198:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate199
	}

yystate199:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate200
	}

yystate200:
	c = l.next()
	switch {
	default:
		goto yyrule69
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate201:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate202
	}

yystate202:
	c = l.next()
	switch {
	default:
		goto yyrule70
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate203
	}

yystate203:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate204
	}

yystate204:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate205
	}

yystate205:
	c = l.next()
	switch {
	default:
		goto yyrule71
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate206:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate207
	}

yystate207:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate208
	}

yystate208:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate209
	}

yystate209:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate210
	}

yystate210:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate211
	}

yystate211:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate212
	}

yystate212:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate213
	}

yystate213:
	c = l.next()
	switch {
	default:
		goto yyrule72
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate214:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'M' || c >= 'O' && c <= 'R' || c == 'T' || c == 'U' || c >= 'W' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'm' || c >= 'o' && c <= 'r' || c == 't' || c == 'u' || c >= 'w' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate215
	case c == 'N' || c == 'n':
		goto yystate221
	case c == 'S' || c == 's':
		goto yystate231
	case c == 'V' || c == 'v':
		goto yystate238
	}

yystate215:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate216
	}

yystate216:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate217
	}

yystate217:
	c = l.next()
	switch {
	default:
		goto yyrule73
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate218
	}

yystate218:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate219
	}

yystate219:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate220
	}

yystate220:
	c = l.next()
	switch {
	default:
		goto yyrule74
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate221:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate222
	}

yystate222:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate223
	}

yystate223:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate224
	}

yystate224:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate225
	}

yystate225:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate226
	}

yystate226:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate227
	}

yystate227:
	c = l.next()
	switch {
	default:
		goto yyrule75
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate228
	}

yystate228:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate229
	}

yystate229:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate230
	}

yystate230:
	c = l.next()
	switch {
	default:
		goto yyrule76
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate231:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate232
	}

yystate232:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate233
	}

yystate233:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate234
	}

yystate234:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate235
	}

yystate235:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate236
	}

yystate236:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate237
	}

yystate237:
	c = l.next()
	switch {
	default:
		goto yyrule77
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate238:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate239
	}

yystate239:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate240
	}

yystate240:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate241
	}

yystate241:
	c = l.next()
	switch {
	default:
		goto yyrule78
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate242:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate243
	}

yystate243:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate244
	}

yystate244:
	c = l.next()
	switch {
	default:
		goto yyrule79
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate245:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate246
	case c == 'O' || c == 'o':
		goto yystate250
	}

yystate246:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate247
	}

yystate247:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate248
	}

yystate248:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate249
	}

yystate249:
	c = l.next()
	switch {
	default:
		goto yyrule80
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate250:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate251
	}

yystate251:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate252
	}

yystate252:
	c = l.next()
	switch {
	default:
		goto yyrule81
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate253:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate254
	}

yystate254:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Q' || c == 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'q' || c == 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate255
	case c == 'R' || c == 'r':
		goto yystate259
	case c == 'T' || c == 't':
		goto yystate281
	}

yystate255:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate256
	}

yystate256:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate257
	}

yystate257:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate258
	}

yystate258:
	c = l.next()
	switch {
	default:
		goto yyrule82
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate259:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate260
	}

yystate260:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate261
	}

yystate261:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate262
	}

yystate262:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate263
	}

yystate263:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'S' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 's' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate264
	case c == 'T' || c == 't':
		goto yystate268
	case c == 'U' || c == 'u':
		goto yystate277
	}

yystate264:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate265
	}

yystate265:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate266
	}

yystate266:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate267
	}

yystate267:
	c = l.next()
	switch {
	default:
		goto yyrule83
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate268:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate269
	}

yystate269:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate270
	}

yystate270:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate271
	}

yystate271:
	c = l.next()
	switch {
	default:
		goto yyrule85
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate272
	}

yystate272:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate273
	}

yystate273:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate274
	}

yystate274:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate275
	}

yystate275:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate276
	}

yystate276:
	c = l.next()
	switch {
	default:
		goto yyrule271
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate277:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate278
	}

yystate278:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate279
	}

yystate279:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate280
	}

yystate280:
	c = l.next()
	switch {
	default:
		goto yyrule86
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate281:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate282
	}

yystate282:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate283
	}

yystate283:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate284
	}

yystate284:
	c = l.next()
	switch {
	default:
		goto yyrule84
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate285:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'B' || c == 'C' || c >= 'F' && c <= 'H' || c >= 'J' && c <= 'N' || c == 'P' || c == 'Q' || c == 'S' || c == 'T' || c >= 'V' && c <= 'Z' || c == '_' || c == 'b' || c == 'c' || c >= 'f' && c <= 'h' || c >= 'j' && c <= 'n' || c == 'p' || c == 'q' || c == 's' || c == 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate286
	case c == 'D' || c == 'd':
		goto yystate352
	case c == 'E' || c == 'e':
		goto yystate354
	case c == 'I' || c == 'i':
		goto yystate397
	case c == 'O' || c == 'o':
		goto yystate405
	case c == 'R' || c == 'r':
		goto yystate410
	case c == 'U' || c == 'u':
		goto yystate413
	}

yystate286:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate287
	case c == 'Y' || c == 'y':
		goto yystate306
	}

yystate287:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate288
	case c == 'E' || c == 'e':
		goto yystate294
	}

yystate288:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate289
	}

yystate289:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate290
	}

yystate290:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate291
	}

yystate291:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate292
	}

yystate292:
	c = l.next()
	switch {
	default:
		goto yyrule87
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate293
	}

yystate293:
	c = l.next()
	switch {
	default:
		goto yyrule88
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate294:
	c = l.next()
	switch {
	default:
		goto yyrule287
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate295
	case c == '_':
		goto yystate299
	}

yystate295:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate296
	}

yystate296:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate297
	}

yystate297:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate298
	}

yystate298:
	c = l.next()
	switch {
	default:
		goto yyrule290
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate299:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate300
	case c == 'S' || c == 's':
		goto yystate303
	}

yystate300:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate301
	}

yystate301:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate302
	}

yystate302:
	c = l.next()
	switch {
	default:
		goto yyrule89
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate303:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate304
	}

yystate304:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate305
	}

yystate305:
	c = l.next()
	switch {
	default:
		goto yyrule90
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate306:
	c = l.next()
	switch {
	default:
		goto yyrule91
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'P' && c <= 'Z' || c >= 'a' && c <= 'm' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate307
	case c == 'O' || c == 'o':
		goto yystate311
	case c == '_':
		goto yystate326
	}

yystate307:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate308
	}

yystate308:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate309
	}

yystate309:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate310
	}

yystate310:
	c = l.next()
	switch {
	default:
		goto yyrule92
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate311:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'E' || c >= 'G' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	case c == 'F' || c == 'f':
		goto yystate312
	}

yystate312:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'V' || c == 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'v' || c == 'x' || c == 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate313
	case c == 'W' || c == 'w':
		goto yystate318
	case c == 'Y' || c == 'y':
		goto yystate322
	}

yystate313:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate314
	}

yystate314:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate315
	}

yystate315:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate316
	}

yystate316:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate317
	}

yystate317:
	c = l.next()
	switch {
	default:
		goto yyrule94
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate318:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate319
	}

yystate319:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate320
	}

yystate320:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate321
	}

yystate321:
	c = l.next()
	switch {
	default:
		goto yyrule93
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate322:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate323
	}

yystate323:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate324
	}

yystate324:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate325
	}

yystate325:
	c = l.next()
	switch {
	default:
		goto yyrule95
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate326:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'L' || c >= 'N' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'l' || c >= 'n' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate327
	case c == 'M' || c == 'm':
		goto yystate331
	case c == 'S' || c == 's':
		goto yystate346
	}

yystate327:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate328
	}

yystate328:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate329
	}

yystate329:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate330
	}

yystate330:
	c = l.next()
	switch {
	default:
		goto yyrule96
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate331:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate332
	}

yystate332:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate333
	case c == 'N' || c == 'n':
		goto yystate342
	}

yystate333:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate334
	}

yystate334:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate335
	}

yystate335:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate336
	}

yystate336:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate337
	}

yystate337:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate338
	}

yystate338:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate339
	}

yystate339:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate340
	}

yystate340:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate341
	}

yystate341:
	c = l.next()
	switch {
	default:
		goto yyrule97
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate342:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate343
	}

yystate343:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate344
	}

yystate344:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate345
	}

yystate345:
	c = l.next()
	switch {
	default:
		goto yyrule98
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate346:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate347
	}

yystate347:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate348
	}

yystate348:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate349
	}

yystate349:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate350
	}

yystate350:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate351
	}

yystate351:
	c = l.next()
	switch {
	default:
		goto yyrule99
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate352:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate353
	}

yystate353:
	c = l.next()
	switch {
	default:
		goto yyrule100
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate354:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'B' || c == 'D' || c == 'E' || c >= 'G' && c <= 'K' || c >= 'M' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c == 'b' || c == 'd' || c == 'e' || c >= 'g' && c <= 'k' || c >= 'm' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate355
	case c == 'C' || c == 'c':
		goto yystate363
	case c == 'F' || c == 'f':
		goto yystate368
	case c == 'L' || c == 'l':
		goto yystate373
	case c == 'S' || c == 's':
		goto yystate391
	}

yystate355:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate356
	}

yystate356:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate357
	}

yystate357:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate358
	}

yystate358:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate359
	}

yystate359:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate360
	}

yystate360:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate361
	}

yystate361:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate362
	}

yystate362:
	c = l.next()
	switch {
	default:
		goto yyrule101
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate363:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate364
	}

yystate364:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate365
	}

yystate365:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate366
	}

yystate366:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate367
	}

yystate367:
	c = l.next()
	switch {
	default:
		goto yyrule281
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate368:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate369
	}

yystate369:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate370
	}

yystate370:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate371
	}

yystate371:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate372
	}

yystate372:
	c = l.next()
	switch {
	default:
		goto yyrule102
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate373:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate374
	case c == 'E' || c == 'e':
		goto yystate388
	}

yystate374:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate375
	}

yystate375:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate376
	case c == '_':
		goto yystate378
	}

yystate376:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate377
	}

yystate377:
	c = l.next()
	switch {
	default:
		goto yyrule103
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate378:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate379
	}

yystate379:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate380
	}

yystate380:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate381
	}

yystate381:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate382
	}

yystate382:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate383
	}

yystate383:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate384
	}

yystate384:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate385
	}

yystate385:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate386
	}

yystate386:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate387
	}

yystate387:
	c = l.next()
	switch {
	default:
		goto yyrule104
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate388:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate389
	}

yystate389:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate390
	}

yystate390:
	c = l.next()
	switch {
	default:
		goto yyrule105
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate391:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate392
	}

yystate392:
	c = l.next()
	switch {
	default:
		goto yyrule106
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate393
	}

yystate393:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate394
	}

yystate394:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate395
	}

yystate395:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate396
	}

yystate396:
	c = l.next()
	switch {
	default:
		goto yyrule107
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate397:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c == 'T' || c == 'U' || c >= 'W' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c == 't' || c == 'u' || c >= 'w' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate398
	case c == 'V' || c == 'v':
		goto yystate404
	}

yystate398:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate399
	}

yystate399:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate400
	}

yystate400:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate401
	}

yystate401:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate402
	}

yystate402:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate403
	}

yystate403:
	c = l.next()
	switch {
	default:
		goto yyrule109
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate404:
	c = l.next()
	switch {
	default:
		goto yyrule110
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate405:
	c = l.next()
	switch {
	default:
		goto yyrule111
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate406
	}

yystate406:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate407
	}

yystate407:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate408
	}

yystate408:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate409
	}

yystate409:
	c = l.next()
	switch {
	default:
		goto yyrule284
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate410:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate411
	}

yystate411:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate412
	}

yystate412:
	c = l.next()
	switch {
	default:
		goto yyrule108
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate413:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate414
	case c == 'P' || c == 'p':
		goto yystate416
	}

yystate414:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate415
	}

yystate415:
	c = l.next()
	switch {
	default:
		goto yyrule112
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate416:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate417
	}

yystate417:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate418
	}

yystate418:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate419
	}

yystate419:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate420
	}

yystate420:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate421
	}

yystate421:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate422
	}

yystate422:
	c = l.next()
	switch {
	default:
		goto yyrule113
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate423:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c == 'M' || c >= 'O' && c <= 'R' || c >= 'T' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'k' || c == 'm' || c >= 'o' && c <= 'r' || c >= 't' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate424
	case c == 'N' || c == 'n':
		goto yystate427
	case c == 'S' || c == 's':
		goto yystate436
	case c == 'X' || c == 'x':
		goto yystate441
	}

yystate424:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate425
	}

yystate425:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate426
	}

yystate426:
	c = l.next()
	switch {
	default:
		goto yyrule114
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate427:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c == 'E' || c == 'F' || c >= 'H' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c == 'e' || c == 'f' || c >= 'h' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate428
	case c == 'G' || c == 'g':
		goto yystate429
	case c == 'U' || c == 'u':
		goto yystate434
	}

yystate428:
	c = l.next()
	switch {
	default:
		goto yyrule115
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate429:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate430
	}

yystate430:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate431
	}

yystate431:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate432
	}

yystate432:
	c = l.next()
	switch {
	default:
		goto yyrule116
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate433
	}

yystate433:
	c = l.next()
	switch {
	default:
		goto yyrule117
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate434:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate435
	}

yystate435:
	c = l.next()
	switch {
	default:
		goto yyrule119
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate436:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate437
	}

yystate437:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate438
	}

yystate438:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate439
	}

yystate439:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate440
	}

yystate440:
	c = l.next()
	switch {
	default:
		goto yyrule120
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate441:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'H' || c >= 'J' && c <= 'O' || c >= 'Q' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'h' || c >= 'j' && c <= 'o' || c >= 'q' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate442
	case c == 'I' || c == 'i':
		goto yystate447
	case c == 'P' || c == 'p':
		goto yystate451
	case c == 'T' || c == 't':
		goto yystate456
	}

yystate442:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate443
	}

yystate443:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate444
	}

yystate444:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate445
	}

yystate445:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate446
	}

yystate446:
	c = l.next()
	switch {
	default:
		goto yyrule118
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate447:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate448
	}

yystate448:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate449
	}

yystate449:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate450
	}

yystate450:
	c = l.next()
	switch {
	default:
		goto yyrule121
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate451:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate452
	}

yystate452:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate453
	}

yystate453:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate454
	}

yystate454:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate455
	}

yystate455:
	c = l.next()
	switch {
	default:
		goto yyrule122
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate456:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate457
	}

yystate457:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate458
	}

yystate458:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate459
	}

yystate459:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate460
	}

yystate460:
	c = l.next()
	switch {
	default:
		goto yyrule123
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate461:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'H' || c == 'J' || c == 'K' || c == 'M' || c == 'N' || c == 'P' || c == 'Q' || c == 'S' || c == 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'h' || c == 'j' || c == 'k' || c == 'm' || c == 'n' || c == 'p' || c == 'q' || c == 's' || c == 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate462
	case c == 'I' || c == 'i':
		goto yystate466
	case c == 'L' || c == 'l':
		goto yystate474
	case c == 'O' || c == 'o':
		goto yystate478
	case c == 'R' || c == 'r':
		goto yystate492
	case c == 'U' || c == 'u':
		goto yystate495
	}

yystate462:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate463
	}

yystate463:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate464
	}

yystate464:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate465
	}

yystate465:
	c = l.next()
	switch {
	default:
		goto yyrule268
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate466:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate467
	case c == 'R' || c == 'r':
		goto yystate471
	}

yystate467:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate468
	}

yystate468:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate469
	}

yystate469:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate470
	}

yystate470:
	c = l.next()
	switch {
	default:
		goto yyrule124
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate471:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate472
	}

yystate472:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate473
	}

yystate473:
	c = l.next()
	switch {
	default:
		goto yyrule125
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate474:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate475
	}

yystate475:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate476
	}

yystate476:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate477
	}

yystate477:
	c = l.next()
	switch {
	default:
		goto yyrule283
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate478:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c == 'S' || c == 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c == 's' || c == 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate479
	case c == 'U' || c == 'u':
		goto yystate484
	}

yystate479:
	c = l.next()
	switch {
	default:
		goto yyrule126
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate480
	}

yystate480:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate481
	}

yystate481:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate482
	}

yystate482:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate483
	}

yystate483:
	c = l.next()
	switch {
	default:
		goto yyrule127
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate484:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate485
	}

yystate485:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate486
	}

yystate486:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate487
	}

yystate487:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate488
	}

yystate488:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate489
	}

yystate489:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate490
	}

yystate490:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate491
	}

yystate491:
	c = l.next()
	switch {
	default:
		goto yyrule128
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate492:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate493
	}

yystate493:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate494
	}

yystate494:
	c = l.next()
	switch {
	default:
		goto yyrule129
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate495:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate496
	}

yystate496:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate497
	}

yystate497:
	c = l.next()
	switch {
	default:
		goto yyrule130
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate498
	}

yystate498:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate499
	}

yystate499:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate500
	}

yystate500:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate501
	}

yystate501:
	c = l.next()
	switch {
	default:
		goto yyrule131
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate502:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate503
	case c == 'R' || c == 'r':
		goto yystate508
	}

yystate503:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate504
	}

yystate504:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate505
	}

yystate505:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate506
	}

yystate506:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate507
	}

yystate507:
	c = l.next()
	switch {
	default:
		goto yyrule207
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate508:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate509
	case c == 'O' || c == 'o':
		goto yystate513
	}

yystate509:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate510
	}

yystate510:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate511
	}

yystate511:
	c = l.next()
	switch {
	default:
		goto yyrule132
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate512
	}

yystate512:
	c = l.next()
	switch {
	default:
		goto yyrule133
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate513:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate514
	}

yystate514:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate515
	}

yystate515:
	c = l.next()
	switch {
	default:
		goto yyrule134
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate516
	}

yystate516:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate517
	}

yystate517:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate518
	}

yystate518:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate519
	}

yystate519:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate520
	}

yystate520:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate521
	}

yystate521:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate522
	}

yystate522:
	c = l.next()
	switch {
	default:
		goto yyrule135
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate523:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'H' || c >= 'J' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'h' || c >= 'j' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate524
	case c == 'I' || c == 'i':
		goto yystate529
	case c == 'O' || c == 'o':
		goto yystate541
	}

yystate524:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'U' || c >= 'W' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'u' || c >= 'w' && c <= 'z':
		goto yystate66
	case c == 'V' || c == 'v':
		goto yystate525
	}

yystate525:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate526
	}

yystate526:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate527
	}

yystate527:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate528
	}

yystate528:
	c = l.next()
	switch {
	default:
		goto yyrule136
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate529:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate530
	}

yystate530:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate531
	}

yystate531:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate532
	}

yystate532:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate533
	}

yystate533:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate534
	}

yystate534:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate535
	}

yystate535:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate536
	}

yystate536:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate537
	}

yystate537:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate538
	}

yystate538:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate539
	}

yystate539:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate540
	}

yystate540:
	c = l.next()
	switch {
	default:
		goto yyrule137
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate541:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate542
	}

yystate542:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate543
	}

yystate543:
	c = l.next()
	switch {
	default:
		goto yyrule138
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate544
	}

yystate544:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate545
	case c == 'S' || c == 's':
		goto yystate560
	}

yystate545:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate546
	}

yystate546:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate547
	case c == 'N' || c == 'n':
		goto yystate556
	}

yystate547:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate548
	}

yystate548:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate549
	}

yystate549:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate550
	}

yystate550:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate551
	}

yystate551:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate552
	}

yystate552:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate553
	}

yystate553:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate554
	}

yystate554:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate555
	}

yystate555:
	c = l.next()
	switch {
	default:
		goto yyrule139
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate556:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate557
	}

yystate557:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate558
	}

yystate558:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate559
	}

yystate559:
	c = l.next()
	switch {
	default:
		goto yyrule140
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate560:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate561
	}

yystate561:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate562
	}

yystate562:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate563
	}

yystate563:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate564
	}

yystate564:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate565
	}

yystate565:
	c = l.next()
	switch {
	default:
		goto yyrule141
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate566:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c == 'E' || c >= 'H' && c <= 'M' || c >= 'O' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c == 'e' || c >= 'h' && c <= 'm' || c >= 'o' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate567
	case c == 'F' || c == 'f':
		goto yystate576
	case c == 'G' || c == 'g':
		goto yystate581
	case c == 'N' || c == 'n':
		goto yystate586
	case c == 'S' || c == 's':
		goto yystate607
	}

yystate567:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate568
	}

yystate568:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate569
	}

yystate569:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate570
	}

yystate570:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate571
	}

yystate571:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'E' || c >= 'G' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	case c == 'F' || c == 'f':
		goto yystate572
	}

yystate572:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate573
	}

yystate573:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate574
	}

yystate574:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate575
	}

yystate575:
	c = l.next()
	switch {
	default:
		goto yyrule142
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate576:
	c = l.next()
	switch {
	default:
		goto yyrule143
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate577
	}

yystate577:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate578
	}

yystate578:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate579
	}

yystate579:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate580
	}

yystate580:
	c = l.next()
	switch {
	default:
		goto yyrule144
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate581:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate582
	}

yystate582:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate583
	}

yystate583:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate584
	}

yystate584:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate585
	}

yystate585:
	c = l.next()
	switch {
	default:
		goto yyrule145
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate586:
	c = l.next()
	switch {
	default:
		goto yyrule151
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'M' || c >= 'O' && c <= 'R' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'm' || c >= 'o' && c <= 'r' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate587
	case c == 'N' || c == 'n':
		goto yystate590
	case c == 'S' || c == 's':
		goto yystate593
	case c == 'T' || c == 't':
		goto yystate597
	}

yystate587:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate588
	}

yystate588:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate589
	}

yystate589:
	c = l.next()
	switch {
	default:
		goto yyrule146
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate590:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate591
	}

yystate591:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate592
	}

yystate592:
	c = l.next()
	switch {
	default:
		goto yyrule147
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate593:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate594
	}

yystate594:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate595
	}

yystate595:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate596
	}

yystate596:
	c = l.next()
	switch {
	default:
		goto yyrule148
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate597:
	c = l.next()
	switch {
	default:
		goto yyrule307
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate598
	case c == 'O' || c == 'o':
		goto yystate606
	}

yystate598:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate599
	case c == 'R' || c == 'r':
		goto yystate602
	}

yystate599:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate600
	}

yystate600:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate601
	}

yystate601:
	c = l.next()
	switch {
	default:
		goto yyrule308
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate602:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'U' || c >= 'W' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'u' || c >= 'w' && c <= 'z':
		goto yystate66
	case c == 'V' || c == 'v':
		goto yystate603
	}

yystate603:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate604
	}

yystate604:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate605
	}

yystate605:
	c = l.next()
	switch {
	default:
		goto yyrule149
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate606:
	c = l.next()
	switch {
	default:
		goto yyrule150
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate607:
	c = l.next()
	switch {
	default:
		goto yyrule152
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate608
	}

yystate608:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate609
	}

yystate609:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate610
	}

yystate610:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate611
	}

yystate611:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate612
	}

yystate612:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate613
	}

yystate613:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate614
	}

yystate614:
	c = l.next()
	switch {
	default:
		goto yyrule153
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate615:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate616
	}

yystate616:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate617
	}

yystate617:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate618
	}

yystate618:
	c = l.next()
	switch {
	default:
		goto yyrule154
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate619:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate620
	}

yystate620:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate621
	}

yystate621:
	c = l.next()
	switch {
	default:
		goto yyrule155
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate622
	}

yystate622:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate623
	}

yystate623:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate624
	}

yystate624:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate625
	}

yystate625:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate626
	}

yystate626:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate627
	}

yystate627:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate628
	}

yystate628:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate629
	}

yystate629:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate630
	}

yystate630:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Y' || c == '_' || c >= 'a' && c <= 'y':
		goto yystate66
	case c == 'Z' || c == 'z':
		goto yystate631
	}

yystate631:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate632
	}

yystate632:
	c = l.next()
	switch {
	default:
		goto yyrule156
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate633:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'H' || c >= 'J' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'h' || c >= 'j' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate634
	case c == 'I' || c == 'i':
		goto yystate649
	case c == 'O' || c == 'o':
		goto yystate655
	}

yystate634:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'E' || c >= 'G' && c <= 'M' || c >= 'O' && c <= 'U' || c >= 'W' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'e' || c >= 'g' && c <= 'm' || c >= 'o' && c <= 'u' || c >= 'w' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate635
	case c == 'F' || c == 'f':
		goto yystate640
	case c == 'N' || c == 'n':
		goto yystate642
	case c == 'V' || c == 'v':
		goto yystate646
	}

yystate635:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate636
	}

yystate636:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate637
	}

yystate637:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate638
	}

yystate638:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate639
	}

yystate639:
	c = l.next()
	switch {
	default:
		goto yyrule157
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate640:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate641
	}

yystate641:
	c = l.next()
	switch {
	default:
		goto yyrule158
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate642:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate643
	}

yystate643:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate644
	}

yystate644:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate645
	}

yystate645:
	c = l.next()
	switch {
	default:
		goto yyrule159
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate646:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate647
	}

yystate647:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate648
	}

yystate648:
	c = l.next()
	switch {
	default:
		goto yyrule160
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate649:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c == 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c == 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate650
	case c == 'M' || c == 'm':
		goto yystate652
	}

yystate650:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate651
	}

yystate651:
	c = l.next()
	switch {
	default:
		goto yyrule161
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate652:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate653
	}

yystate653:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate654
	}

yystate654:
	c = l.next()
	switch {
	default:
		goto yyrule162
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate655:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'M' || c >= 'O' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'm' || c >= 'o' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate656
	case c == 'N' || c == 'n':
		goto yystate671
	case c == 'W' || c == 'w':
		goto yystate681
	}

yystate656:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate657
	case c == 'K' || c == 'k':
		goto yystate670
	}

yystate657:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate658
	case c == 'T' || c == 't':
		goto yystate668
	}

yystate658:
	c = l.next()
	switch {
	default:
		goto yyrule163
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate659
	}

yystate659:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate660
	}

yystate660:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate661
	}

yystate661:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate662
	}

yystate662:
	c = l.next()
	switch {
	default:
		goto yyrule272
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate663
	}

yystate663:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate664
	}

yystate664:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate665
	}

yystate665:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate666
	}

yystate666:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate667
	}

yystate667:
	c = l.next()
	switch {
	default:
		goto yyrule273
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate668:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate669
	}

yystate669:
	c = l.next()
	switch {
	default:
		goto yyrule164
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate670:
	c = l.next()
	switch {
	default:
		goto yyrule165
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate671:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate672
	}

yystate672:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate673
	case c == 'T' || c == 't':
		goto yystate677
	}

yystate673:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate674
	}

yystate674:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate675
	}

yystate675:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate676
	}

yystate676:
	c = l.next()
	switch {
	default:
		goto yyrule299
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate677:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate678
	}

yystate678:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate679
	}

yystate679:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate680
	}

yystate680:
	c = l.next()
	switch {
	default:
		goto yyrule303
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate681:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate682
	case c == '_':
		goto yystate684
	}

yystate682:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate683
	}

yystate683:
	c = l.next()
	switch {
	default:
		goto yyrule166
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate684:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate685
	}

yystate685:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate686
	}

yystate686:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate687
	}

yystate687:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate688
	}

yystate688:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate689
	}

yystate689:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate690
	}

yystate690:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate691
	}

yystate691:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate692
	}

yystate692:
	c = l.next()
	switch {
	default:
		goto yyrule167
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate693:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c >= 'F' && c <= 'H' || c >= 'J' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'h' || c >= 'j' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate694
	case c == 'E' || c == 'e':
		goto yystate701
	case c == 'I' || c == 'i':
		goto yystate717
	case c == 'O' || c == 'o':
		goto yystate754
	}

yystate694:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate695
	}

yystate695:
	c = l.next()
	switch {
	default:
		goto yyrule168
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate696
	}

yystate696:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate697
	}

yystate697:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate698
	}

yystate698:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate699
	}

yystate699:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate700
	}

yystate700:
	c = l.next()
	switch {
	default:
		goto yyrule169
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate701:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate702
	}

yystate702:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate703
	}

yystate703:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate704
	}

yystate704:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate705
	}

yystate705:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'H' || c >= 'J' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'h' || c >= 'j' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate706
	case c == 'I' || c == 'i':
		goto yystate710
	case c == 'T' || c == 't':
		goto yystate713
	}

yystate706:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate707
	}

yystate707:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate708
	}

yystate708:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate709
	}

yystate709:
	c = l.next()
	switch {
	default:
		goto yyrule298
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate710:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate711
	}

yystate711:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate712
	}

yystate712:
	c = l.next()
	switch {
	default:
		goto yyrule279
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate713:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate714
	}

yystate714:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate715
	}

yystate715:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate716
	}

yystate716:
	c = l.next()
	switch {
	default:
		goto yyrule301
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate717:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate718
	case c == 'N' || c == 'n':
		goto yystate727
	}

yystate718:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate719
	}

yystate719:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate720
	}

yystate720:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate721
	}

yystate721:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate722
	}

yystate722:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate723
	}

yystate723:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate724
	}

yystate724:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate725
	}

yystate725:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate726
	}

yystate726:
	c = l.next()
	switch {
	default:
		goto yyrule170
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate727:
	c = l.next()
	switch {
	default:
		goto yyrule171
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate728
	case c == '_':
		goto yystate749
	}

yystate728:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate729
	}

yystate729:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate730
	}

yystate730:
	c = l.next()
	switch {
	default:
		goto yyrule172
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate731
	}

yystate731:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate732
	case c == 'S' || c == 's':
		goto yystate743
	}

yystate732:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate733
	}

yystate733:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate734
	}

yystate734:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate735
	}

yystate735:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate736
	}

yystate736:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate737
	}

yystate737:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate738
	}

yystate738:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate739
	}

yystate739:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate740
	}

yystate740:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate741
	}

yystate741:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate742
	}

yystate742:
	c = l.next()
	switch {
	default:
		goto yyrule173
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate743:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate744
	}

yystate744:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate745
	}

yystate745:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate746
	}

yystate746:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate747
	}

yystate747:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate748
	}

yystate748:
	c = l.next()
	switch {
	default:
		goto yyrule174
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate749:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate750
	}

yystate750:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate751
	}

yystate751:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate752
	}

yystate752:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate753
	}

yystate753:
	c = l.next()
	switch {
	default:
		goto yyrule175
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate754:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate755
	case c == 'N' || c == 'n':
		goto yystate757
	}

yystate755:
	c = l.next()
	switch {
	default:
		goto yyrule176
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate756
	}

yystate756:
	c = l.next()
	switch {
	default:
		goto yyrule177
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate757:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate758
	}

yystate758:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate759
	}

yystate759:
	c = l.next()
	switch {
	default:
		goto yyrule178
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate760:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'N' || c >= 'P' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'n' || c >= 'p' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate761
	case c == 'O' || c == 'o':
		goto yystate771
	case c == 'U' || c == 'u':
		goto yystate774
	}

yystate761:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate762
	case c == 'T' || c == 't':
		goto yystate765
	}

yystate762:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate763
	}

yystate763:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate764
	}

yystate764:
	c = l.next()
	switch {
	default:
		goto yyrule179
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate765:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate766
	}

yystate766:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate767
	}

yystate767:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate768
	}

yystate768:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate769
	}

yystate769:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate770
	}

yystate770:
	c = l.next()
	switch {
	default:
		goto yyrule180
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate771:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c == 'U' || c == 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c == 'u' || c == 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate772
	case c == 'W' || c == 'w':
		goto yystate773
	}

yystate772:
	c = l.next()
	switch {
	default:
		goto yyrule181
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate773:
	c = l.next()
	switch {
	default:
		goto yyrule274
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate774:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate775
	case c == 'M' || c == 'm':
		goto yystate779
	}

yystate775:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate776
	}

yystate776:
	c = l.next()
	switch {
	default:
		goto yyrule267
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate777
	}

yystate777:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'E' || c >= 'G' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	case c == 'F' || c == 'f':
		goto yystate778
	}

yystate778:
	c = l.next()
	switch {
	default:
		goto yyrule243
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate779:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate780
	}

yystate780:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate781
	}

yystate781:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate782
	}

yystate782:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate783
	}

yystate783:
	c = l.next()
	switch {
	default:
		goto yyrule282
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate784:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'E' || c >= 'G' && c <= 'M' || c == 'O' || c == 'Q' || c == 'S' || c == 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'm' || c == 'o' || c == 'q' || c == 's' || c == 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'F' || c == 'f':
		goto yystate785
	case c == 'N' || c == 'n':
		goto yystate790
	case c == 'P' || c == 'p':
		goto yystate793
	case c == 'R' || c == 'r':
		goto yystate798
	case c == 'U' || c == 'u':
		goto yystate802
	}

yystate785:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'E' || c >= 'G' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	case c == 'F' || c == 'f':
		goto yystate786
	}

yystate786:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate787
	}

yystate787:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate788
	}

yystate788:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate789
	}

yystate789:
	c = l.next()
	switch {
	default:
		goto yyrule182
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate790:
	c = l.next()
	switch {
	default:
		goto yyrule183
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate791
	}

yystate791:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate792
	}

yystate792:
	c = l.next()
	switch {
	default:
		goto yyrule184
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate793:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate794
	}

yystate794:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate795
	}

yystate795:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate796
	}

yystate796:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate797
	}

yystate797:
	c = l.next()
	switch {
	default:
		goto yyrule185
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate798:
	c = l.next()
	switch {
	default:
		goto yyrule187
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate799
	}

yystate799:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate800
	}

yystate800:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate801
	}

yystate801:
	c = l.next()
	switch {
	default:
		goto yyrule186
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate802:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate803
	}

yystate803:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate804
	}

yystate804:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate805
	}

yystate805:
	c = l.next()
	switch {
	default:
		goto yyrule188
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate806:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'N' || c == 'P' || c == 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'n' || c == 'p' || c == 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate807
	case c == 'O' || c == 'o':
		goto yystate814
	case c == 'R' || c == 'r':
		goto yystate818
	}

yystate807:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate808
	}

yystate808:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate809
	}

yystate809:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate810
	}

yystate810:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate811
	}

yystate811:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate812
	}

yystate812:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate813
	}

yystate813:
	c = l.next()
	switch {
	default:
		goto yyrule189
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate814:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate815
	}

yystate815:
	c = l.next()
	switch {
	default:
		goto yyrule190
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate816
	}

yystate816:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate817
	}

yystate817:
	c = l.next()
	switch {
	default:
		goto yyrule191
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate818:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'H' || c >= 'J' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'h' || c >= 'j' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate819
	case c == 'I' || c == 'i':
		goto yystate830
	case c == 'O' || c == 'o':
		goto yystate835
	}

yystate819:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate820
	case c == 'P' || c == 'p':
		goto yystate826
	}

yystate820:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate821
	}

yystate821:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate822
	}

yystate822:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate823
	}

yystate823:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate824
	}

yystate824:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate825
	}

yystate825:
	c = l.next()
	switch {
	default:
		goto yyrule285
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate826:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate827
	}

yystate827:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate828
	}

yystate828:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate829
	}

yystate829:
	c = l.next()
	switch {
	default:
		goto yyrule192
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate830:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate831
	}

yystate831:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate832
	}

yystate832:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate833
	}

yystate833:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate834
	}

yystate834:
	c = l.next()
	switch {
	default:
		goto yyrule193
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate835:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate836
	}

yystate836:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate837
	}

yystate837:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate838
	}

yystate838:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate839
	}

yystate839:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate840
	}

yystate840:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate841
	}

yystate841:
	c = l.next()
	switch {
	default:
		goto yyrule194
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate842:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate843
	}

yystate843:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate844
	case c == 'I' || c == 'i':
		goto yystate849
	}

yystate844:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate845
	}

yystate845:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate846
	}

yystate846:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate847
	}

yystate847:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate848
	}

yystate848:
	c = l.next()
	switch {
	default:
		goto yyrule195
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate849:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate850
	}

yystate850:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate851
	}

yystate851:
	c = l.next()
	switch {
	default:
		goto yyrule196
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate852:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c >= 'F' && c <= 'H' || c == 'J' || c == 'K' || c == 'M' || c == 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'h' || c == 'j' || c == 'k' || c == 'm' || c == 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate853
	case c == 'E' || c == 'e':
		goto yystate856
	case c == 'I' || c == 'i':
		goto yystate884
	case c == 'L' || c == 'l':
		goto yystate888
	case c == 'O' || c == 'o':
		goto yystate892
	}

yystate853:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate854
	}

yystate854:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate855
	}

yystate855:
	c = l.next()
	switch {
	default:
		goto yyrule208
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate856:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'E' || c >= 'H' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'e' || c >= 'h' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate857
	case c == 'F' || c == 'f':
		goto yystate860
	case c == 'G' || c == 'g':
		goto yystate868
	case c == 'P' || c == 'p':
		goto yystate872
	}

yystate857:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate858
	case c == 'L' || c == 'l':
		goto yystate859
	}

yystate858:
	c = l.next()
	switch {
	default:
		goto yyrule209
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate859:
	c = l.next()
	switch {
	default:
		goto yyrule286
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate860:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate861
	}

yystate861:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate862
	}

yystate862:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate863
	}

yystate863:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate864
	}

yystate864:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate865
	}

yystate865:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate866
	}

yystate866:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate867
	}

yystate867:
	c = l.next()
	switch {
	default:
		goto yyrule214
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate868:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate869
	}

yystate869:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate870
	}

yystate870:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate871
	}

yystate871:
	c = l.next()
	switch {
	default:
		goto yyrule212
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate872:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate873
	case c == 'L' || c == 'l':
		goto yystate880
	}

yystate873:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate874
	}

yystate874:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate875
	}

yystate875:
	c = l.next()
	switch {
	default:
		goto yyrule210
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate876
	}

yystate876:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate877
	}

yystate877:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate878
	}

yystate878:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate879
	}

yystate879:
	c = l.next()
	switch {
	default:
		goto yyrule211
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate880:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate881
	}

yystate881:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate882
	}

yystate882:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate883
	}

yystate883:
	c = l.next()
	switch {
	default:
		goto yyrule213
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate884:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate885
	}

yystate885:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate886
	}

yystate886:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate887
	}

yystate887:
	c = l.next()
	switch {
	default:
		goto yyrule197
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate888:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate889
	}

yystate889:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate890
	}

yystate890:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate891
	}

yystate891:
	c = l.next()
	switch {
	default:
		goto yyrule215
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate892:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate893
	case c == 'W' || c == 'w':
		goto yystate899
	}

yystate893:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate894
	}

yystate894:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate895
	}

yystate895:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate896
	}

yystate896:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate897
	}

yystate897:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate898
	}

yystate898:
	c = l.next()
	switch {
	default:
		goto yyrule198
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate899:
	c = l.next()
	switch {
	default:
		goto yyrule199
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate900:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c == 'D' || c == 'F' || c == 'G' || c >= 'J' && c <= 'L' || c == 'N' || c == 'P' || c == 'R' || c == 'S' || c >= 'V' && c <= 'X' || c == 'Z' || c == '_' || c == 'a' || c == 'b' || c == 'd' || c == 'f' || c == 'g' || c >= 'j' && c <= 'l' || c == 'n' || c == 'p' || c == 'r' || c == 's' || c >= 'v' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate901
	case c == 'E' || c == 'e':
		goto yystate907
	case c == 'H' || c == 'h':
		goto yystate944
	case c == 'I' || c == 'i':
		goto yystate950
	case c == 'M' || c == 'm':
		goto yystate955
	case c == 'O' || c == 'o':
		goto yystate962
	case c == 'Q' || c == 'q':
		goto yystate965
	case c == 'T' || c == 't':
		goto yystate983
	case c == 'U' || c == 'u':
		goto yystate990
	case c == 'Y' || c == 'y':
		goto yystate1009
	}

yystate901:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate902
	}

yystate902:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate903
	}

yystate903:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate904
	}

yystate904:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate905
	}

yystate905:
	c = l.next()
	switch {
	default:
		goto yyrule200
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate906
	}

yystate906:
	c = l.next()
	switch {
	default:
		goto yyrule201
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate907:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'K' || c >= 'M' && c <= 'Q' || c >= 'U' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'k' || c >= 'm' && c <= 'q' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate908
	case c == 'L' || c == 'l':
		goto yystate924
	case c == 'R' || c == 'r':
		goto yystate928
	case c == 'S' || c == 's':
		goto yystate938
	case c == 'T' || c == 't':
		goto yystate943
	}

yystate908:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate909
	}

yystate909:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate910
	}

yystate910:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate911
	}

yystate911:
	c = l.next()
	switch {
	default:
		goto yyrule218
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate912
	}

yystate912:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate913
	}

yystate913:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate914
	}

yystate914:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate915
	}

yystate915:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate916
	}

yystate916:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate917
	}

yystate917:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate918
	}

yystate918:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate919
	}

yystate919:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate920
	}

yystate920:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate921
	}

yystate921:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate922
	}

yystate922:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate923
	}

yystate923:
	c = l.next()
	switch {
	default:
		goto yyrule219
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate924:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate925
	}

yystate925:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate926
	}

yystate926:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate927
	}

yystate927:
	c = l.next()
	switch {
	default:
		goto yyrule220
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate928:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate929
	}

yystate929:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate930
	}

yystate930:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate931
	}

yystate931:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate932
	}

yystate932:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Y' || c == '_' || c >= 'a' && c <= 'y':
		goto yystate66
	case c == 'Z' || c == 'z':
		goto yystate933
	}

yystate933:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate934
	}

yystate934:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate935
	}

yystate935:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate936
	}

yystate936:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate937
	}

yystate937:
	c = l.next()
	switch {
	default:
		goto yyrule202
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate938:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate939
	}

yystate939:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate940
	}

yystate940:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate941
	}

yystate941:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate942
	}

yystate942:
	c = l.next()
	switch {
	default:
		goto yyrule203
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate943:
	c = l.next()
	switch {
	default:
		goto yyrule221
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate944:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate945
	case c == 'O' || c == 'o':
		goto yystate948
	}

yystate945:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate946
	}

yystate946:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate947
	}

yystate947:
	c = l.next()
	switch {
	default:
		goto yyrule222
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate948:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate949
	}

yystate949:
	c = l.next()
	switch {
	default:
		goto yyrule223
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate950:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate951
	}

yystate951:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate952
	}

yystate952:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate953
	}

yystate953:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate954
	}

yystate954:
	c = l.next()
	switch {
	default:
		goto yyrule264
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate955:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate956
	}

yystate956:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate957
	}

yystate957:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate958
	}

yystate958:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate959
	}

yystate959:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate960
	}

yystate960:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate961
	}

yystate961:
	c = l.next()
	switch {
	default:
		goto yyrule278
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate962:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate963
	}

yystate963:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate964
	}

yystate964:
	c = l.next()
	switch {
	default:
		goto yyrule204
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate965:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate966
	}

yystate966:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate967
	}

yystate967:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate968
	}

yystate968:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate969
	}

yystate969:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate970
	}

yystate970:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate971
	}

yystate971:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate972
	}

yystate972:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'E' || c >= 'G' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	case c == 'F' || c == 'f':
		goto yystate973
	}

yystate973:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate974
	}

yystate974:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate975
	}

yystate975:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate976
	}

yystate976:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate977
	}

yystate977:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate978
	}

yystate978:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate979
	}

yystate979:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate980
	}

yystate980:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate981
	}

yystate981:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate982
	}

yystate982:
	c = l.next()
	switch {
	default:
		goto yyrule270
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate983:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate984
	}

yystate984:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c == 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c == 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate985
	case c == 'T' || c == 't':
		goto yystate987
	}

yystate985:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate986
	}

yystate986:
	c = l.next()
	switch {
	default:
		goto yyrule205
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate987:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate988
	}

yystate988:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate989
	}

yystate989:
	c = l.next()
	switch {
	default:
		goto yyrule206
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate990:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate991
	case c == 'M' || c == 'm':
		goto yystate1008
	}

yystate991:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate992
	case c == 'S' || c == 's':
		goto yystate996
	}

yystate992:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate993
	}

yystate993:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate994
	}

yystate994:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate995
	}

yystate995:
	c = l.next()
	switch {
	default:
		goto yyrule224
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate996:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate997
	}

yystate997:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate998
	}

yystate998:
	c = l.next()
	switch {
	default:
		goto yyrule225
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate999
	}

yystate999:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1000
	}

yystate1000:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate1001
	}

yystate1001:
	c = l.next()
	switch {
	default:
		goto yyrule226
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate66
	case c == '_':
		goto yystate1002
	}

yystate1002:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1003
	}

yystate1003:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1004
	}

yystate1004:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate1005
	}

yystate1005:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1006
	}

yystate1006:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate1007
	}

yystate1007:
	c = l.next()
	switch {
	default:
		goto yyrule227
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1008:
	c = l.next()
	switch {
	default:
		goto yyrule228
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1009:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1010
	}

yystate1010:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate1011
	}

yystate1011:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1012
	}

yystate1012:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1013
	}

yystate1013:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1014
	}

yystate1014:
	c = l.next()
	switch {
	default:
		goto yyrule229
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1015:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c == 'F' || c == 'G' || c >= 'J' && c <= 'N' || c == 'P' || c == 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c == 'f' || c == 'g' || c >= 'j' && c <= 'n' || c == 'p' || c == 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1016
	case c == 'E' || c == 'e':
		goto yystate1021
	case c == 'H' || c == 'h':
		goto yystate1024
	case c == 'I' || c == 'i':
		goto yystate1027
	case c == 'O' || c == 'o':
		goto yystate1048
	case c == 'R' || c == 'r':
		goto yystate1049
	}

yystate1016:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate1017
	}

yystate1017:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate1018
	}

yystate1018:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1019
	}

yystate1019:
	c = l.next()
	switch {
	default:
		goto yyrule230
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1020
	}

yystate1020:
	c = l.next()
	switch {
	default:
		goto yyrule231
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1021:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate1022
	}

yystate1022:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1023
	}

yystate1023:
	c = l.next()
	switch {
	default:
		goto yyrule302
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1024:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1025
	}

yystate1025:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1026
	}

yystate1026:
	c = l.next()
	switch {
	default:
		goto yyrule232
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1027:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate1028
	case c == 'N' || c == 'n':
		goto yystate1035
	}

yystate1028:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1029
	}

yystate1029:
	c = l.next()
	switch {
	default:
		goto yyrule288
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1030
	}

yystate1030:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1031
	}

yystate1031:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1032
	}

yystate1032:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate1033
	}

yystate1033:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'P' || c == 'p':
		goto yystate1034
	}

yystate1034:
	c = l.next()
	switch {
	default:
		goto yyrule289
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1035:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate1036
	}

yystate1036:
	c = l.next()
	switch {
	default:
		goto yyrule276
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'H' || c >= 'J' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'h' || c >= 'j' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate1037
	case c == 'I' || c == 'i':
		goto yystate1041
	case c == 'T' || c == 't':
		goto yystate1044
	}

yystate1037:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate1038
	}

yystate1038:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1039
	}

yystate1039:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate1040
	}

yystate1040:
	c = l.next()
	switch {
	default:
		goto yyrule296
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1041:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1042
	}

yystate1042:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1043
	}

yystate1043:
	c = l.next()
	switch {
	default:
		goto yyrule277
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1044:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1045
	}

yystate1045:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'w' || c == 'y' || c == 'z':
		goto yystate66
	case c == 'X' || c == 'x':
		goto yystate1046
	}

yystate1046:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1047
	}

yystate1047:
	c = l.next()
	switch {
	default:
		goto yyrule300
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1048:
	c = l.next()
	switch {
	default:
		goto yyrule233
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1049:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'H' || c >= 'J' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'h' || c >= 'j' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1050
	case c == 'I' || c == 'i':
		goto yystate1064
	case c == 'U' || c == 'u':
		goto yystate1071
	}

yystate1050:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1051
	case c == 'N' || c == 'n':
		goto yystate1056
	}

yystate1051:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate1052
	}

yystate1052:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1053
	}

yystate1053:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1054
	}

yystate1054:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate1055
	}

yystate1055:
	c = l.next()
	switch {
	default:
		goto yyrule234
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1056:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1057
	}

yystate1057:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1058
	}

yystate1058:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate1059
	}

yystate1059:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1060
	}

yystate1060:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1061
	}

yystate1061:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1062
	}

yystate1062:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1063
	}

yystate1063:
	c = l.next()
	switch {
	default:
		goto yyrule235
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1064:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate1065
	case c == 'M' || c == 'm':
		goto yystate1070
	}

yystate1065:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate1066
	}

yystate1066:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1067
	}

yystate1067:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1068
	}

yystate1068:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1069
	}

yystate1069:
	c = l.next()
	switch {
	default:
		goto yyrule236
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1070:
	c = l.next()
	switch {
	default:
		goto yyrule237
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1071:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1072
	case c == 'N' || c == 'n':
		goto yystate1073
	}

yystate1072:
	c = l.next()
	switch {
	default:
		goto yyrule269
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1073:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate1074
	}

yystate1074:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1075
	}

yystate1075:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1076
	}

yystate1076:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1077
	}

yystate1077:
	c = l.next()
	switch {
	default:
		goto yyrule238
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1078:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c == 'O' || c == 'Q' || c == 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c == 'o' || c == 'q' || c == 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1079
	case c == 'P' || c == 'p':
		goto yystate1110
	case c == 'S' || c == 's':
		goto yystate1118
	}

yystate1079:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'H' || c == 'J' || c >= 'M' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'h' || c == 'j' || c >= 'm' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate1080
	case c == 'I' || c == 'i':
		goto yystate1089
	case c == 'K' || c == 'k':
		goto yystate1095
	case c == 'L' || c == 'l':
		goto yystate1100
	case c == 'S' || c == 's':
		goto yystate1104
	}

yystate1080:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1081
	}

yystate1081:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate1082
	}

yystate1082:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate1083
	}

yystate1083:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1084
	}

yystate1084:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1085
	}

yystate1085:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1086
	}

yystate1086:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1087
	}

yystate1087:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate1088
	}

yystate1088:
	c = l.next()
	switch {
	default:
		goto yyrule239
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1089:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c == 'P' || c >= 'R' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c == 'p' || c >= 'r' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1090
	case c == 'Q' || c == 'q':
		goto yystate1092
	}

yystate1090:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1091
	}

yystate1091:
	c = l.next()
	switch {
	default:
		goto yyrule240
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1092:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate1093
	}

yystate1093:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1094
	}

yystate1094:
	c = l.next()
	switch {
	default:
		goto yyrule241
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1095:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1096
	}

yystate1096:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1097
	}

yystate1097:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate1098
	}

yystate1098:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1099
	}

yystate1099:
	c = l.next()
	switch {
	default:
		goto yyrule242
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1100:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1101
	}

yystate1101:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	case c == 'C' || c == 'c':
		goto yystate1102
	}

yystate1102:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate1103
	}

yystate1103:
	c = l.next()
	switch {
	default:
		goto yyrule244
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1104:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1105
	}

yystate1105:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate1106
	}

yystate1106:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1107
	}

yystate1107:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1108
	}

yystate1108:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate1109
	}

yystate1109:
	c = l.next()
	switch {
	default:
		goto yyrule265
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1110:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate1111
	case c == 'P' || c == 'p':
		goto yystate1115
	}

yystate1111:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1112
	}

yystate1112:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1113
	}

yystate1113:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1114
	}

yystate1114:
	c = l.next()
	switch {
	default:
		goto yyrule245
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1115:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1116
	}

yystate1116:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1117
	}

yystate1117:
	c = l.next()
	switch {
	default:
		goto yyrule246
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1118:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1119
	case c == 'I' || c == 'i':
		goto yystate1121
	}

yystate1119:
	c = l.next()
	switch {
	default:
		goto yyrule247
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1120
	}

yystate1120:
	c = l.next()
	switch {
	default:
		goto yyrule248
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1121:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1122
	}

yystate1122:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate1123
	}

yystate1123:
	c = l.next()
	switch {
	default:
		goto yyrule249
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1124:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1125
	case c == 'E' || c == 'e':
		goto yystate1147
	}

yystate1125:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate1126
	case c == 'R' || c == 'r':
		goto yystate1130
	}

yystate1126:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	case c == 'U' || c == 'u':
		goto yystate1127
	}

yystate1127:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1128
	}

yystate1128:
	c = l.next()
	switch {
	default:
		goto yyrule250
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1129
	}

yystate1129:
	c = l.next()
	switch {
	default:
		goto yyrule251
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1130:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'D' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c == 'a' || c >= 'd' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate1131
	case c == 'C' || c == 'c':
		goto yystate1137
	case c == 'I' || c == 'i':
		goto yystate1141
	}

yystate1131:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1132
	}

yystate1132:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1133
	}

yystate1133:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1134
	}

yystate1134:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1135
	}

yystate1135:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate1136
	}

yystate1136:
	c = l.next()
	switch {
	default:
		goto yyrule295
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1137:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate1138
	}

yystate1138:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1139
	}

yystate1139:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1140
	}

yystate1140:
	c = l.next()
	switch {
	default:
		goto yyrule293
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1141:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1142
	}

yystate1142:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c == 'A' || c >= 'C' && c <= 'Z' || c == '_' || c == 'a' || c >= 'c' && c <= 'z':
		goto yystate66
	case c == 'B' || c == 'b':
		goto yystate1143
	}

yystate1143:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate1144
	}

yystate1144:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1145
	}

yystate1145:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1146
	}

yystate1146:
	c = l.next()
	switch {
	default:
		goto yyrule252
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1147:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1148
	}

yystate1148:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1149
	}

yystate1149:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1150
	}

yystate1150:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1151
	}

yystate1151:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1152
	}

yystate1152:
	c = l.next()
	switch {
	default:
		goto yyrule253
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1153:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c == 'F' || c == 'G' || c >= 'I' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c == 'f' || c == 'g' || c >= 'i' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1154
	case c == 'E' || c == 'e':
		goto yystate1161
	case c == 'H' || c == 'h':
		goto yystate1173
	case c == 'R' || c == 'r':
		goto yystate1178
	}

yystate1154:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1155
	}

yystate1155:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1156
	}

yystate1156:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1157
	}

yystate1157:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1158
	}

yystate1158:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	case c == 'G' || c == 'g':
		goto yystate1159
	}

yystate1159:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	case c == 'S' || c == 's':
		goto yystate1160
	}

yystate1160:
	c = l.next()
	switch {
	default:
		goto yyrule254
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1161:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1162
	}

yystate1162:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate1163
	}

yystate1163:
	c = l.next()
	switch {
	default:
		goto yyrule255
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'c' || c >= 'e' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'D' || c == 'd':
		goto yystate1164
	case c == 'O' || c == 'o':
		goto yystate1167
	}

yystate1164:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1165
	}

yystate1165:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate1166
	}

yystate1166:
	c = l.next()
	switch {
	default:
		goto yyrule256
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1167:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'E' || c >= 'G' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	case c == 'F' || c == 'f':
		goto yystate1168
	}

yystate1168:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	case c == 'Y' || c == 'y':
		goto yystate1169
	}

yystate1169:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1170
	}

yystate1170:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1171
	}

yystate1171:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1172
	}

yystate1172:
	c = l.next()
	switch {
	default:
		goto yyrule257
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1173:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1174
	}

yystate1174:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1175
	case c == 'R' || c == 'r':
		goto yystate1176
	}

yystate1175:
	c = l.next()
	switch {
	default:
		goto yyrule258
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1176:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1177
	}

yystate1177:
	c = l.next()
	switch {
	default:
		goto yyrule259
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1178:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1179
	}

yystate1179:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1180
	}

yystate1180:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1181
	}

yystate1181:
	c = l.next()
	switch {
	default:
		goto yyrule260
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1182:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1186
	case c == '\'':
		goto yystate1183
	}

yystate1183:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate1184
	}

yystate1184:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c == '\'':
		goto yystate1185
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate1184
	}

yystate1185:
	c = l.next()
	goto yyrule11

yystate1186:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1187
	}

yystate1187:
	c = l.next()
	switch {
	default:
		goto yyrule261
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1188:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1189
	}

yystate1189:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	case c == 'A' || c == 'a':
		goto yystate1190
	}

yystate1190:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1191
	}

yystate1191:
	c = l.next()
	switch {
	default:
		goto yyrule291
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'V' || c >= 'X' && c <= 'Z' || c >= 'a' && c <= 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	case c == 'W' || c == 'w':
		goto yystate1192
	case c == '_':
		goto yystate1196
	}

yystate1192:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1193
	}

yystate1193:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1194
	}

yystate1194:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	case c == 'K' || c == 'k':
		goto yystate1195
	}

yystate1195:
	c = l.next()
	switch {
	default:
		goto yyrule262
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1196:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate66
	case c == 'M' || c == 'm':
		goto yystate1197
	}

yystate1197:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1198
	}

yystate1198:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	case c == 'N' || c == 'n':
		goto yystate1199
	}

yystate1199:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	case c == 'T' || c == 't':
		goto yystate1200
	}

yystate1200:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	case c == 'H' || c == 'h':
		goto yystate1201
	}

yystate1201:
	c = l.next()
	switch {
	default:
		goto yyrule263
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1202:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	case c == 'E' || c == 'e':
		goto yystate1203
	}

yystate1203:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	case c == 'R' || c == 'r':
		goto yystate1204
	}

yystate1204:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	case c == 'O' || c == 'o':
		goto yystate1205
	}

yystate1205:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'E' || c >= 'G' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	case c == 'F' || c == 'f':
		goto yystate1206
	}

yystate1206:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	case c == 'I' || c == 'i':
		goto yystate1207
	}

yystate1207:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate1208
	}

yystate1208:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	case c == 'L' || c == 'l':
		goto yystate1209
	}

yystate1209:
	c = l.next()
	switch {
	default:
		goto yyrule266
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1210:
	c = l.next()
	switch {
	default:
		goto yyrule309
	case c == '$' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate1211:
	c = l.next()
	goto yyrule15

yystate1212:
	c = l.next()
	switch {
	default:
		goto yyrule310
	case c == '|':
		goto yystate1213
	}

yystate1213:
	c = l.next()
	goto yyrule35

	goto yystate1214 // silence unused label error
yystate1214:
	c = l.next()
yystart1214:
	switch {
	default:
		goto yyrule16
	case c == '"':
		goto yystate1216
	case c == '\\':
		goto yystate1218
	case c == '\x00':
		goto yystate2
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate1215
	}

yystate1215:
	c = l.next()
	switch {
	default:
		goto yyrule16
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate1215
	}

yystate1216:
	c = l.next()
	switch {
	default:
		goto yyrule19
	case c == '"':
		goto yystate1217
	}

yystate1217:
	c = l.next()
	goto yyrule18

yystate1218:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate1219
	}

yystate1219:
	c = l.next()
	goto yyrule17

	goto yystate1220 // silence unused label error
yystate1220:
	c = l.next()
yystart1220:
	switch {
	default:
		goto yyrule20
	case c == '\'':
		goto yystate1222
	case c == '\\':
		goto yystate1224
	case c == '\x00':
		goto yystate2
	case c >= '\x01' && c <= '&' || c >= '(' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate1221
	}

yystate1221:
	c = l.next()
	switch {
	default:
		goto yyrule20
	case c >= '\x01' && c <= '&' || c >= '(' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate1221
	}

yystate1222:
	c = l.next()
	switch {
	default:
		goto yyrule23
	case c == '\'':
		goto yystate1223
	}

yystate1223:
	c = l.next()
	goto yyrule22

yystate1224:
	c = l.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate1225
	}

yystate1225:
	c = l.next()
	goto yyrule21

	goto yystate1226 // silence unused label error
yystate1226:
	c = l.next()
yystart1226:
	switch {
	default:
		goto yystate1227 // c >= '\x01' && c <= '\b' || c >= '\n' && c <= '\x1f' || c >= '!' && c <= 'ÿ'
	case c == '\t' || c == ' ':
		goto yystate1228
	case c == '\x00':
		goto yystate2
	}

yystate1227:
	c = l.next()
	goto yyrule8

yystate1228:
	c = l.next()
	switch {
	default:
		goto yyrule7
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate1228
	}

	goto yystate1229 // silence unused label error
yystate1229:
	c = l.next()
yystart1229:
	switch {
	default:
		goto yyrule24
	case c == '\x00':
		goto yystate2
	case c == '`':
		goto yystate1231
	case c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ':
		goto yystate1230
	}

yystate1230:
	c = l.next()
	switch {
	default:
		goto yyrule24
	case c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ':
		goto yystate1230
	}

yystate1231:
	c = l.next()
	switch {
	default:
		goto yyrule26
	case c == '`':
		goto yystate1232
	}

yystate1232:
	c = l.next()
	goto yyrule25

yyrule1: // \0
	{
		return 0
	}
yyrule2: // [ \t\n\r]+

	goto yystate0
yyrule3: // #.*

	goto yystate0
yyrule4: // \/\/.*

	goto yystate0
yyrule5: // \/\*([^*]|\*+[^*/])*\*+\/

	goto yystate0
yyrule6: // --
	{
		l.sc = S3
		goto yystate0
	}
yyrule7: // [ \t]+.*
	{
		{
			l.sc = 0
		}
		goto yystate0
	}
yyrule8: // [^ \t]
	{
		{
			l.sc = 0
			l.c = '-'
			n := len(l.val)
			l.unget(l.val[n-1])
			return '-'
		}
		goto yystate0
	}
yyrule9: // {int_lit}
	{
		return l.int(lval)
	}
yyrule10: // {float_lit}
	{
		return l.float(lval)
	}
yyrule11: // {hex_lit}
	{
		return l.hex(lval)
	}
yyrule12: // {bit_lit}
	{
		return l.bit(lval)
	}
yyrule13: // \"
	{
		l.sc = S1
		goto yystate0
	}
yyrule14: // '
	{
		l.sc = S2
		goto yystate0
	}
yyrule15: // `
	{
		l.sc = S4
		goto yystate0
	}
yyrule16: // [^\"\\]*
	{
		l.stringLit = append(l.stringLit, l.val...)
		goto yystate0
	}
yyrule17: // \\.
	{
		l.stringLit = append(l.stringLit, l.val...)
		goto yystate0
	}
yyrule18: // \"\"
	{
		l.stringLit = append(l.stringLit, '"')
		goto yystate0
	}
yyrule19: // \"
	{
		l.stringLit = append(l.stringLit, '"')
		l.sc = 0
		return l.str(lval, "\"")
	}
yyrule20: // [^'\\]*
	{
		l.stringLit = append(l.stringLit, l.val...)
		goto yystate0
	}
yyrule21: // \\.
	{
		l.stringLit = append(l.stringLit, l.val...)
		goto yystate0
	}
yyrule22: // ''
	{
		l.stringLit = append(l.stringLit, '\'')
		goto yystate0
	}
yyrule23: // '
	{
		l.stringLit = append(l.stringLit, '\'')
		l.sc = 0
		return l.str(lval, "'")
	}
yyrule24: // [^`]*
	{
		l.stringLit = append(l.stringLit, l.val...)
		goto yystate0
	}
yyrule25: // ``
	{
		l.stringLit = append(l.stringLit, '`')
		goto yystate0
	}
yyrule26: // `
	{
		l.sc = 0
		lval.item = string(l.stringLit)
		l.stringLit = l.stringLit[0:0]
		return identifier
	}
yyrule27: // "&&"
	{
		return andand
	}
yyrule28: // "&^"
	{
		return andnot
	}
yyrule29: // "<<"
	{
		return lsh
	}
yyrule30: // "<="
	{
		return le
	}
yyrule31: // "="
	{
		return eq
	}
yyrule32: // ">="
	{
		return ge
	}
yyrule33: // "!="
	{
		return neq
	}
yyrule34: // "<>"
	{
		return neq
	}
yyrule35: // "||"
	{
		return oror
	}
yyrule36: // ">>"
	{
		return rsh
	}
yyrule37: // "<=>"
	{
		return nulleq
	}
yyrule38: // "@"
	{
		return at
	}
yyrule39: // "?"
	{
		return placeholder
	}
yyrule40: // {abs}
	{
		lval.item = string(l.val)
		return abs
	}
yyrule41: // {add}
	{
		return add
	}
yyrule42: // {adddate}
	{
		lval.item = string(l.val)
		return addDate
	}
yyrule43: // {admin}
	{
		lval.item = string(l.val)
		return admin
	}
yyrule44: // {after}
	{
		lval.item = string(l.val)
		return after
	}
yyrule45: // {all}
	{
		return all
	}
yyrule46: // {alter}
	{
		return alter
	}
yyrule47: // {and}
	{
		return and
	}
yyrule48: // {any}
	{
		lval.item = string(l.val)
		return any
	}
yyrule49: // {asc}
	{
		return asc
	}
yyrule50: // {as}
	{
		return as
	}
yyrule51: // {auto_increment}
	{
		lval.item = string(l.val)
		return autoIncrement
	}
yyrule52: // {avg}
	{
		lval.item = string(l.val)
		return avg
	}
yyrule53: // {avg_row_length}
	{
		lval.item = string(l.val)
		return avgRowLength
	}
yyrule54: // {begin}
	{
		lval.item = string(l.val)
		return begin
	}
yyrule55: // {between}
	{
		return between
	}
yyrule56: // {both}
	{
		return both
	}
yyrule57: // {by}
	{
		return by
	}
yyrule58: // {case}
	{
		return caseKwd
	}
yyrule59: // {cast}
	{
		return cast
	}
yyrule60: // {character}
	{
		return character
	}
yyrule61: // {charset}
	{
		lval.item = string(l.val)
		return charsetKwd
	}
yyrule62: // {check}
	{
		return check
	}
yyrule63: // {checksum}
	{
		lval.item = string(l.val)
		return checksum
	}
yyrule64: // {coalesce}
	{
		lval.item = string(l.val)
		return coalesce
	}
yyrule65: // {collate}
	{
		return collate
	}
yyrule66: // {collation}
	{
		lval.item = string(l.val)
		return collation
	}
yyrule67: // {column}
	{
		return column
	}
yyrule68: // {columns}
	{
		lval.item = string(l.val)
		return columns
	}
yyrule69: // {comment}
	{
		lval.item = string(l.val)
		return comment
	}
yyrule70: // {commit}
	{
		lval.item = string(l.val)
		return commit
	}
yyrule71: // {committed}
	{
		lval.item = string(l.val)
		return committed
	}
yyrule72: // {compression}
	{
		lval.item = string(l.val)
		return compression
	}
yyrule73: // {concat}
	{
		lval.item = string(l.val)
		return concat
	}
yyrule74: // {concat_ws}
	{
		lval.item = string(l.val)
		return concatWs
	}
yyrule75: // {connection}
	{
		lval.item = string(l.val)
		return connection
	}
yyrule76: // {connection_id}
	{
		lval.item = string(l.val)
		return connectionID
	}
yyrule77: // {constraint}
	{
		return constraint
	}
yyrule78: // {convert}
	{
		return convert
	}
yyrule79: // {count}
	{
		lval.item = string(l.val)
		return count
	}
yyrule80: // {create}
	{
		return create
	}
yyrule81: // {cross}
	{
		return cross
	}
yyrule82: // {curdate}
	{
		lval.item = string(l.val)
		return curDate
	}
yyrule83: // {current_date}
	{
		lval.item = string(l.val)
		return currentDate
	}
yyrule84: // {curtime}
	{
		lval.item = string(l.val)
		return curTime
	}
yyrule85: // {current_time}
	{
		lval.item = string(l.val)
		return currentTime
	}
yyrule86: // {current_user}
	{
		lval.item = string(l.val)
		return currentUser
	}
yyrule87: // {database}
	{
		lval.item = string(l.val)
		return database
	}
yyrule88: // {databases}
	{
		return databases
	}
yyrule89: // {date_add}
	{
		lval.item = string(l.val)
		return dateAdd
	}
yyrule90: // {date_sub}
	{
		lval.item = string(l.val)
		return dateSub
	}
yyrule91: // {day}
	{
		lval.item = string(l.val)
		return day
	}
yyrule92: // {dayname}
	{
		lval.item = string(l.val)
		return dayname
	}
yyrule93: // {dayofweek}
	{
		lval.item = string(l.val)
		return dayofweek
	}
yyrule94: // {dayofmonth}
	{
		lval.item = string(l.val)
		return dayofmonth
	}
yyrule95: // {dayofyear}
	{
		lval.item = string(l.val)
		return dayofyear
	}
yyrule96: // {day_hour}
	{
		lval.item = string(l.val)
		return dayHour
	}
yyrule97: // {day_microsecond}
	{
		lval.item = string(l.val)
		return dayMicrosecond
	}
yyrule98: // {day_minute}
	{
		lval.item = string(l.val)
		return dayMinute
	}
yyrule99: // {day_second}
	{
		lval.item = string(l.val)
		return daySecond
	}
yyrule100: // {ddl}
	{
		return ddl
	}
yyrule101: // {deallocate}
	{
		lval.item = string(l.val)
		return deallocate
	}
yyrule102: // {default}
	{
		return defaultKwd
	}
yyrule103: // {delayed}
	{
		return delayed
	}
yyrule104: // {delay_key_write}
	{
		lval.item = string(l.val)
		return delayKeyWrite
	}
yyrule105: // {delete}
	{
		return deleteKwd
	}
yyrule106: // {desc}
	{
		return desc
	}
yyrule107: // {describe}
	{
		return describe
	}
yyrule108: // {drop}
	{
		return drop
	}
yyrule109: // {distinct}
	{
		return distinct
	}
yyrule110: // {div}
	{
		return div
	}
yyrule111: // {do}
	{
		lval.item = string(l.val)
		return do
	}
yyrule112: // {dual}
	{
		return dual
	}
yyrule113: // {duplicate}
	{
		lval.item = string(l.val)
		return duplicate
	}
yyrule114: // {else}
	{
		return elseKwd
	}
yyrule115: // {end}
	{
		lval.item = string(l.val)
		return end
	}
yyrule116: // {engine}
	{
		lval.item = string(l.val)
		return engine
	}
yyrule117: // {engines}
	{
		lval.item = string(l.val)
		return engines
	}
yyrule118: // {execute}
	{
		lval.item = string(l.val)
		return execute
	}
yyrule119: // {enum}
	{
		return enum
	}
yyrule120: // {escape}
	{
		lval.item = string(l.val)
		return escape
	}
yyrule121: // {exists}
	{
		return exists
	}
yyrule122: // {explain}
	{
		return explain
	}
yyrule123: // {extract}
	{
		lval.item = string(l.val)
		return extract
	}
yyrule124: // {fields}
	{
		lval.item = string(l.val)
		return fields
	}
yyrule125: // {first}
	{
		lval.item = string(l.val)
		return first
	}
yyrule126: // {for}
	{
		return forKwd
	}
yyrule127: // {foreign}
	{
		return foreign
	}
yyrule128: // {found_rows}
	{
		lval.item = string(l.val)
		return foundRows
	}
yyrule129: // {from}
	{
		return from
	}
yyrule130: // {full}
	{
		lval.item = string(l.val)
		return full
	}
yyrule131: // {fulltext}
	{
		return fulltext
	}
yyrule132: // {grant}
	{
		return grant
	}
yyrule133: // {grants}
	{
		lval.item = string(l.val)
		return grants
	}
yyrule134: // {group}
	{
		return group
	}
yyrule135: // {group_concat}
	{
		lval.item = string(l.val)
		return groupConcat
	}
yyrule136: // {having}
	{
		return having
	}
yyrule137: // {high_priority}
	{
		return highPriority
	}
yyrule138: // {hour}
	{
		lval.item = string(l.val)
		return hour
	}
yyrule139: // {hour_microsecond}
	{
		lval.item = string(l.val)
		return hourMicrosecond
	}
yyrule140: // {hour_minute}
	{
		lval.item = string(l.val)
		return hourMinute
	}
yyrule141: // {hour_second}
	{
		lval.item = string(l.val)
		return hourSecond
	}
yyrule142: // {identified}
	{
		lval.item = string(l.val)
		return identified
	}
yyrule143: // {if}
	{
		lval.item = string(l.val)
		return ifKwd
	}
yyrule144: // {ifnull}
	{
		lval.item = string(l.val)
		return ifNull
	}
yyrule145: // {ignore}
	{
		return ignore
	}
yyrule146: // {index}
	{
		return index
	}
yyrule147: // {inner}
	{
		return inner
	}
yyrule148: // {insert}
	{
		return insert
	}
yyrule149: // {interval}
	{
		return interval
	}
yyrule150: // {into}
	{
		return into
	}
yyrule151: // {in}
	{
		return in
	}
yyrule152: // {is}
	{
		return is
	}
yyrule153: // {isolation}
	{
		lval.item = string(l.val)
		return isolation
	}
yyrule154: // {join}
	{
		return join
	}
yyrule155: // {key}
	{
		return key
	}
yyrule156: // {key_block_size}
	{
		lval.item = string(l.val)
		return keyBlockSize
	}
yyrule157: // {leading}
	{
		return leading
	}
yyrule158: // {left}
	{
		lval.item = string(l.val)
		return left
	}
yyrule159: // {length}
	{
		lval.item = string(l.val)
		return length
	}
yyrule160: // {level}
	{
		lval.item = string(l.val)
		return level
	}
yyrule161: // {like}
	{
		return like
	}
yyrule162: // {limit}
	{
		return limit
	}
yyrule163: // {local}
	{
		lval.item = string(l.val)
		return local
	}
yyrule164: // {locate}
	{
		lval.item = string(l.val)
		return locate
	}
yyrule165: // {lock}
	{
		return lock
	}
yyrule166: // {lower}
	{
		lval.item = string(l.val)
		return lower
	}
yyrule167: // {low_priority}
	{
		return lowPriority
	}
yyrule168: // {max}
	{
		lval.item = string(l.val)
		return max
	}
yyrule169: // {max_rows}
	{
		lval.item = string(l.val)
		return maxRows
	}
yyrule170: // {microsecond}
	{
		lval.item = string(l.val)
		return microsecond
	}
yyrule171: // {min}
	{
		lval.item = string(l.val)
		return min
	}
yyrule172: // {minute}
	{
		lval.item = string(l.val)
		return minute
	}
yyrule173: // {minute_microsecond}
	{
		lval.item = string(l.val)
		return minuteMicrosecond
	}
yyrule174: // {minute_second}
	{
		lval.item = string(l.val)
		return minuteSecond
	}
yyrule175: // {min_rows}
	{
		lval.item = string(l.val)
		return minRows
	}
yyrule176: // {mod}
	{
		return mod
	}
yyrule177: // {mode}
	{
		lval.item = string(l.val)
		return mode
	}
yyrule178: // {month}
	{
		lval.item = string(l.val)
		return month
	}
yyrule179: // {names}
	{
		lval.item = string(l.val)
		return names
	}
yyrule180: // {national}
	{
		lval.item = string(l.val)
		return national
	}
yyrule181: // {not}
	{
		return not
	}
yyrule182: // {offset}
	{
		lval.item = string(l.val)
		return offset
	}
yyrule183: // {on}
	{
		return on
	}
yyrule184: // {only}
	{
		lval.item = string(l.val)
		return only
	}
yyrule185: // {option}
	{
		return option
	}
yyrule186: // {order}
	{
		return order
	}
yyrule187: // {or}
	{
		return or
	}
yyrule188: // {outer}
	{
		return outer
	}
yyrule189: // {password}
	{
		lval.item = string(l.val)
		return password
	}
yyrule190: // {pow}
	{
		lval.item = string(l.val)
		return pow
	}
yyrule191: // {power}
	{
		lval.item = string(l.val)
		return power
	}
yyrule192: // {prepare}
	{
		lval.item = string(l.val)
		return prepare
	}
yyrule193: // {primary}
	{
		return primary
	}
yyrule194: // {procedure}
	{
		return procedure
	}
yyrule195: // {quarter}
	{
		lval.item = string(l.val)
		return quarter
	}
yyrule196: // {quick}
	{
		lval.item = string(l.val)
		return quick
	}
yyrule197: // {right}
	{
		return right
	}
yyrule198: // {rollback}
	{
		lval.item = string(l.val)
		return rollback
	}
yyrule199: // {row}
	{
		lval.item = string(l.val)
		return row
	}
yyrule200: // {schema}
	{
		lval.item = string(l.val)
		return schema
	}
yyrule201: // {schemas}
	{
		return schemas
	}
yyrule202: // {serializable}
	{
		lval.item = string(l.val)
		return serializable
	}
yyrule203: // {session}
	{
		lval.item = string(l.val)
		return session
	}
yyrule204: // {some}
	{
		lval.item = string(l.val)
		return some
	}
yyrule205: // {start}
	{
		lval.item = string(l.val)
		return start
	}
yyrule206: // {status}
	{
		lval.item = string(l.val)
		return status
	}
yyrule207: // {global}
	{
		lval.item = string(l.val)
		return global
	}
yyrule208: // {rand}
	{
		lval.item = string(l.val)
		return rand
	}
yyrule209: // {read}
	{
		return read
	}
yyrule210: // {repeat}
	{
		lval.item = string(l.val)
		return repeat
	}
yyrule211: // {repeatable}
	{
		lval.item = string(l.val)
		return repeatable
	}
yyrule212: // {regexp}
	{
		return regexp
	}
yyrule213: // {replace}
	{
		lval.item = string(l.val)
		return replace
	}
yyrule214: // {references}
	{
		return references
	}
yyrule215: // {rlike}
	{
		return rlike
	}
yyrule216: // {sys_var}
	{
		lval.item = string(l.val)
		return sysVar
	}
yyrule217: // {user_var}
	{
		lval.item = string(l.val)
		return userVar
	}
yyrule218: // {second}
	{
		lval.item = string(l.val)
		return second
	}
yyrule219: // {second_microsecond}
	{
		lval.item = string(l.val)
		return secondMicrosecond
	}
yyrule220: // {select}
	{
		return selectKwd
	}
yyrule221: // {set}
	{
		return set
	}
yyrule222: // {share}
	{
		return share
	}
yyrule223: // {show}
	{
		return show
	}
yyrule224: // {subdate}
	{
		lval.item = string(l.val)
		return subDate
	}
yyrule225: // {substr}
	{
		lval.item = string(l.val)
		return substring
	}
yyrule226: // {substring}
	{
		lval.item = string(l.val)
		return substring
	}
yyrule227: // {substring_index}
	{
		lval.item = string(l.val)
		return substringIndex
	}
yyrule228: // {sum}
	{
		lval.item = string(l.val)
		return sum
	}
yyrule229: // {sysdate}
	{
		lval.item = string(l.val)
		return sysDate
	}
yyrule230: // {table}
	{
		return tableKwd
	}
yyrule231: // {tables}
	{
		lval.item = string(l.val)
		return tables
	}
yyrule232: // {then}
	{
		return then
	}
yyrule233: // {to}
	{
		return to
	}
yyrule234: // {trailing}
	{
		return trailing
	}
yyrule235: // {transaction}
	{
		lval.item = string(l.val)
		return transaction
	}
yyrule236: // {triggers}
	{
		lval.item = string(l.val)
		return triggers
	}
yyrule237: // {trim}
	{
		lval.item = string(l.val)
		return trim
	}
yyrule238: // {truncate}
	{
		lval.item = string(l.val)
		return truncate
	}
yyrule239: // {uncommitted}
	{
		lval.item = string(l.val)
		return uncommitted
	}
yyrule240: // {union}
	{
		return union
	}
yyrule241: // {unique}
	{
		return unique
	}
yyrule242: // {unknown}
	{
		lval.item = string(l.val)
		return unknown
	}
yyrule243: // {nullif}
	{
		lval.item = string(l.val)
		return nullIf
	}
yyrule244: // {unlock}
	{
		return unlock
	}
yyrule245: // {update}
	{
		return update
	}
yyrule246: // {upper}
	{
		lval.item = string(l.val)
		return upper
	}
yyrule247: // {use}
	{
		return use
	}
yyrule248: // {user}
	{
		lval.item = string(l.val)
		return user
	}
yyrule249: // {using}
	{
		return using
	}
yyrule250: // {value}
	{
		lval.item = string(l.val)
		return value
	}
yyrule251: // {values}
	{
		return values
	}
yyrule252: // {variables}
	{
		lval.item = string(l.val)
		return variables
	}
yyrule253: // {version}
	{
		lval.item = string(l.val)
		return version
	}
yyrule254: // {warnings}
	{
		lval.item = string(l.val)
		return warnings
	}
yyrule255: // {week}
	{
		lval.item = string(l.val)
		return week
	}
yyrule256: // {weekday}
	{
		lval.item = string(l.val)
		return weekday
	}
yyrule257: // {weekofyear}
	{
		lval.item = string(l.val)
		return weekofyear
	}
yyrule258: // {when}
	{
		return when
	}
yyrule259: // {where}
	{
		return where
	}
yyrule260: // {write}
	{
		return write
	}
yyrule261: // {xor}
	{
		return xor
	}
yyrule262: // {yearweek}
	{
		lval.item = string(l.val)
		return yearweek
	}
yyrule263: // {year_month}
	{
		lval.item = string(l.val)
		return yearMonth

	}
yyrule264: // {signed}
	{
		lval.item = string(l.val)
		return signed
	}
yyrule265: // {unsigned}
	{
		return unsigned
	}
yyrule266: // {zerofill}
	{
		return zerofill
	}
yyrule267: // {null}
	{
		lval.item = nil
		return null
	}
yyrule268: // {false}
	{
		return falseKwd
	}
yyrule269: // {true}
	{
		return trueKwd
	}
yyrule270: // {calc_found_rows}
	{
		lval.item = string(l.val)
		return calcFoundRows
	}
yyrule271: // {current_ts}
	{
		lval.item = string(l.val)
		return currentTs
	}
yyrule272: // {localtime}
	{
		return localTime
	}
yyrule273: // {localts}
	{
		return localTs
	}
yyrule274: // {now}
	{
		lval.item = string(l.val)
		return now
	}
yyrule275: // {bit}
	{
		lval.item = string(l.val)
		return bitType
	}
yyrule276: // {tiny}
	{
		lval.item = string(l.val)
		return tinyIntType
	}
yyrule277: // {tinyint}
	{
		lval.item = string(l.val)
		return tinyIntType
	}
yyrule278: // {smallint}
	{
		lval.item = string(l.val)
		return smallIntType
	}
yyrule279: // {mediumint}
	{
		lval.item = string(l.val)
		return mediumIntType
	}
yyrule280: // {bigint}
	{
		lval.item = string(l.val)
		return bigIntType
	}
yyrule281: // {decimal}
	{
		lval.item = string(l.val)
		return decimalType
	}
yyrule282: // {numeric}
	{
		lval.item = string(l.val)
		return numericType
	}
yyrule283: // {float}
	{
		lval.item = string(l.val)
		return floatType
	}
yyrule284: // {double}
	{
		lval.item = string(l.val)
		return doubleType
	}
yyrule285: // {precision}
	{
		lval.item = string(l.val)
		return precisionType
	}
yyrule286: // {real}
	{
		lval.item = string(l.val)
		return realType
	}
yyrule287: // {date}
	{
		lval.item = string(l.val)
		return dateType
	}
yyrule288: // {time}
	{
		lval.item = string(l.val)
		return timeType
	}
yyrule289: // {timestamp}
	{
		lval.item = string(l.val)
		return timestampType
	}
yyrule290: // {datetime}
	{
		lval.item = string(l.val)
		return datetimeType
	}
yyrule291: // {year}
	{
		lval.item = string(l.val)
		return yearType
	}
yyrule292: // {char}
	{
		lval.item = string(l.val)
		return charType
	}
yyrule293: // {varchar}
	{
		lval.item = string(l.val)
		return varcharType
	}
yyrule294: // {binary}
	{
		lval.item = string(l.val)
		return binaryType
	}
yyrule295: // {varbinary}
	{
		lval.item = string(l.val)
		return varbinaryType
	}
yyrule296: // {tinyblob}
	{
		lval.item = string(l.val)
		return tinyblobType
	}
yyrule297: // {blob}
	{
		lval.item = string(l.val)
		return blobType
	}
yyrule298: // {mediumblob}
	{
		lval.item = string(l.val)
		return mediumblobType
	}
yyrule299: // {longblob}
	{
		lval.item = string(l.val)
		return longblobType
	}
yyrule300: // {tinytext}
	{
		lval.item = string(l.val)
		return tinytextType
	}
yyrule301: // {mediumtext}
	{
		lval.item = string(l.val)
		return mediumtextType
	}
yyrule302: // {text}
	{
		lval.item = string(l.val)
		return textType
	}
yyrule303: // {longtext}
	{
		lval.item = string(l.val)
		return longtextType
	}
yyrule304: // {bool}
	{
		lval.item = string(l.val)
		return boolType
	}
yyrule305: // {boolean}
	{
		lval.item = string(l.val)
		return booleanType
	}
yyrule306: // {byte}
	{
		lval.item = string(l.val)
		return byteType
	}
yyrule307: // {int}
	{
		lval.item = string(l.val)
		return intType
	}
yyrule308: // {integer}
	{
		lval.item = string(l.val)
		return integerType
	}
yyrule309: // {ident}
	{
		lval.item = string(l.val)
		return l.handleIdent(lval)
	}
yyrule310: // .
	{
		return c0
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	return int(unicode.ReplacementChar)
}

func (l *lexer) npos() (line, col int) {
	if line, col = l.nline, l.ncol; col == 0 {
		line--
		col = l.lcol + 1
	}
	return
}

func (l *lexer) str(lval *yySymType, pref string) int {
	l.sc = 0
	// TODO: performance issue.
	s := string(l.stringLit)
	l.stringLit = l.stringLit[0:0]
	if pref == "'" {
		s = strings.Replace(s, "\\'", "'", -1)
		s = strings.TrimSuffix(s, "'") + "\""
		pref = "\""
	}
	v := stringutil.RemoveUselessBackslash(pref + s)
	v, err := strconv.Unquote(v)
	if err != nil {
		v = strings.TrimSuffix(s, pref)
	}
	lval.item = v
	return stringLit
}

func (l *lexer) trimIdent(idt string) string {
	idt = strings.TrimPrefix(idt, "`")
	idt = strings.TrimSuffix(idt, "`")
	return idt
}

func (l *lexer) int(lval *yySymType) int {
	n, err := strconv.ParseUint(string(l.val), 0, 64)
	if err != nil {
		l.errf("integer literal: %v", err)
		return int(unicode.ReplacementChar)
	}

	switch {
	case n < math.MaxInt64:
		lval.item = int64(n)
	default:
		lval.item = uint64(n)
	}
	return intLit
}

func (l *lexer) float(lval *yySymType) int {
	n, err := strconv.ParseFloat(string(l.val), 64)
	if err != nil {
		l.errf("float literal: %v", err)
		return int(unicode.ReplacementChar)
	}

	lval.item = float64(n)
	return floatLit
}

// https://dev.mysql.com/doc/refman/5.7/en/hexadecimal-literals.html
func (l *lexer) hex(lval *yySymType) int {
	s := string(l.val)
	h, err := mysql.ParseHex(s)
	if err != nil {
		l.errf("hexadecimal literal: %v", err)
		return int(unicode.ReplacementChar)
	}
	lval.item = h
	return hexLit
}

// https://dev.mysql.com/doc/refman/5.7/en/bit-type.html
func (l *lexer) bit(lval *yySymType) int {
	s := string(l.val)
	b, err := mysql.ParseBit(s, -1)
	if err != nil {
		l.errf("bit literal: %v", err)
		return int(unicode.ReplacementChar)
	}
	lval.item = b
	return bitLit
}

func (l *lexer) handleIdent(lval *yySymType) int {
	s := lval.item.(string)
	// A character string literal may have an optional character set introducer and COLLATE clause:
	// [_charset_name]'string' [COLLATE collation_name]
	// See: https://dev.mysql.com/doc/refman/5.7/en/charset-literal.html
	if !strings.HasPrefix(s, "_") {
		return identifier
	}
	cs, _, err := charset.GetCharsetInfo(s[1:])
	if err != nil {
		return identifier
	}
	lval.item = cs
	return underscoreCS
}
