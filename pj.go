package pj

import (
	"fmt"
	"strconv"
)

// tokenize

type TokenType int

const (
	TK_INT TokenType = iota
	TK_STR

	TK_SBRACE   // {
	TK_EBRACE   // }
	TK_SBRACKET // [
	TK_EBRACKET // ]
	TK_COMMA    // ,
	TK_COLON    // :

	TK_TRUE
	TK_FALSE

	TK_EOF

	TK_ILLEGAL
)

type Token struct {
	Type   TokenType
	Name   string
	Value  int
	Length int
}

type Lexer struct {
	input   string
	pos     int // current position
	readPos int
	ch      byte // current char
}

func NewLexer(str string) *Lexer {
	l := &Lexer{input: str}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) nextChar() byte {
	return l.input[l.readPos]
}

func (l *Lexer) NextToken() Token {
	var token Token

	for isWhiteSpace(l.ch) {
		l.readChar()
	}

	switch l.ch {
	case '{':
		token = newToken(TK_SBRACE, l.ch, 1)
	case '}':
		token = newToken(TK_EBRACE, l.ch, 1)
	case '[':
		token = newToken(TK_SBRACKET, l.ch, 1)
	case ']':
		token = newToken(TK_EBRACKET, l.ch, 1)
	case ',':
		token = newToken(TK_COMMA, l.ch, 1)
	case ':':
		token = newToken(TK_COLON, l.ch, 1)
	case '"':
		l.readChar()
		p := l.pos
		for !isQuote(l.ch) {
			l.readChar()
		}
		token = newToken(TK_STR, l.ch, l.pos-p)
		token.Name = l.input[p:l.pos]
	case 0:
		token = newToken(TK_EOF, l.ch, 1)
	default:
		if isLetter(l.ch) {
			p := l.pos
			for isLetter(l.nextChar()) {
				l.readChar()
			}
			token = newToken(LookUpIdent(l.input[p:l.pos+1]), l.ch, l.pos-p)
			token.Name = l.input[p : l.pos+1]
		} else if isDigit(l.ch) {
			p := l.pos
			for isDigit(l.nextChar()) {
				l.readChar()
			}
			token = newToken(TK_INT, l.ch, l.pos-p)
			token.Name = l.input[p : l.pos+1]
			num, err := strconv.Atoi(token.Name)
			if err != nil {
				token.Type = TK_ILLEGAL
			}
			token.Value = num
		} else {
			token = newToken(TK_ILLEGAL, l.ch, -1)
		}
	}

	l.readChar()
	return token
}

func newToken(ty TokenType, ch byte, length int) Token {
	return Token{
		Type:   ty,
		Name:   string(ch),
		Length: length,
	}
}

var keywords = map[string]TokenType{
	"true":  TK_TRUE,
	"false": TK_FALSE,
}

func LookUpIdent(str string) TokenType {
	if tok, ok := keywords[str]; ok {
		return tok
	}
	return TK_ILLEGAL
}

func isWhiteSpace(c byte) bool {
	if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
		return true
	}
	return false
}

func isQuote(c byte) bool {
	if c == '"' {
		return true
	}
	return false
}

func isLetter(c byte) bool {
	if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
		return true
	}
	return false
}

func isDigit(c byte) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	return false
}

// parse

type Parser struct {
	lexer *Lexer
	cur   Token
	peek  Token
	err   []error
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{lexer: l}
	// cur and peek are both set
	p.nextToken()
	p.nextToken()
	return p
}

type NodeType int

const (
	RootNode NodeType = iota
	ObjectNode
	PropertyNode
	ArrayNode
)

type Root struct {
	ty  NodeType
	val *Value
}

type Object struct {
	ty       NodeType
	children []Property
}

type Property struct {
	ty  NodeType
	key string
	val Value
}

type Array struct {
	ty       NodeType
	children []Value
}

type Value interface{}

type State int

const (
	ObjectStart State = iota
	ObjectOpen

	PropertyStart
	PropertyKey
	PropertyValue

	ArrayStart
	ArrayOpen
)

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) Parse() (*Root, error) {
	val := p.parseValue()

	return &Root{ty: RootNode, val: &val}, nil
}

func (p *Parser) parseValue() Value {
	var val Value

	if p.cur.Type == TK_SBRACE {
		val = p.parseObject()
	} else if p.cur.Type == TK_SBRACKET {
		val = p.parseArray()
	} else {
		switch p.cur.Type {
		case TK_INT:
			val = p.cur.Value
		case TK_STR:
			val = p.cur.Name
		case TK_TRUE:
			val = true
		case TK_FALSE:
			val = false
		default:
			val = nil
		}
	}

	return val
}

func (p *Parser) parseObject() Value {
	obj := Object{ty: ObjectNode}
	objState := ObjectStart

	for p.cur.Type != TK_EOF {
		switch objState {
		case ObjectStart:
			if p.cur.Type == TK_SBRACE {
				objState = ObjectOpen
				p.nextToken()
			}
		case ObjectOpen:
			if p.cur.Type == TK_EBRACE {
				p.nextToken()
				return obj
			} else {
				prop := p.parseProperty()
				obj.children = append(obj.children, prop)
				p.nextToken()
			}
		}
	}

	return obj
}

func (p *Parser) parseProperty() Property {
	prop := Property{ty: PropertyNode}
	propState := PropertyStart

	for p.cur.Type != TK_EOF {
		switch propState {
		case PropertyStart:
			if p.cur.Type == TK_STR {
				prop.key = p.cur.Name
				propState = PropertyKey
			}
			p.nextToken()
		case PropertyKey:
			if p.cur.Type == TK_COLON {
				propState = PropertyValue
				p.nextToken()
			}
		case PropertyValue:
			switch p.cur.Type {
			case TK_SBRACE:
				prop.val = p.parseObject()
			case TK_SBRACKET:
				prop.val = p.parseArray()
			case TK_INT:
				prop.val = p.cur.Value
			case TK_STR:
				prop.val = p.cur.Name
			case TK_TRUE:
				prop.val = true
			case TK_FALSE:
				prop.val = false
			}
			return prop
		}
	}

	return prop
}

func (p *Parser) parseArray() Value {
	arr := Array{ty: ArrayNode}
	arrState := ArrayStart

	for p.cur.Type != TK_EOF {
		switch arrState {
		case ArrayStart:
			if p.cur.Type == TK_SBRACKET {
				arrState = ArrayOpen
				p.nextToken()
			}
		case ArrayOpen:
			if p.cur.Type == TK_EBRACKET {
				p.nextToken()
				return arr
			} else if p.cur.Type == TK_COMMA {
				p.nextToken()
			} else {
				val := p.parseValue()
				arr.children = append(arr.children, val)
				p.nextToken()
			}
		}
	}

	return arr
}

func (r *Root) PrintFromRoot() {
	rootValue := *r.val
	switch rootValue.(type) {
	case Object:
		obj, _ := rootValue.(Object)
		for _, c := range obj.children {
			fmt.Printf("KEY: %s, ", c.key)
			if cObj, ok := c.val.(Object); ok {
				fmt.Printf("VALUE: { ")
				cObj.printObject()
				fmt.Printf(" }")
			} else if cArr, ok := c.val.(Array); ok {
				fmt.Printf("VALUE: [")
				cArr.printArray()
				fmt.Printf(" ]")
			} else {
				fmt.Printf("VALUE: %s\n", c.val)
			}
		}
	case Array:
		arr, _ := rootValue.(Array)
		fmt.Printf("ARRAY: [")
		arr.printArray()
		fmt.Printf(" ]\n")
	}
}

func (o *Object) printObject() {
	for _, c := range o.children {
		if valObj, ok := c.val.(Object); ok {
			valObj.printObject()
		} else {
			fmt.Printf("KEY: %s, VALUE: %s", c.key, c.val)
		}
	}
}

func (a *Array) printArray() {
	for _, c := range a.children {
		if valArr, ok := c.(Array); ok {
			valArr.printArray()
		} else {
			fmt.Printf(" %v ", c)
		}
	}
}
