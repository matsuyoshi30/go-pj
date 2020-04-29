package pj

import (
	"fmt"
)

// tokenize

type TokenType int

const (
	TK_INT TokenType = iota
	TK_STR

	TK_SBRACE // {
	TK_EBRACE // }
	TK_COLON  // :

	TK_EOF

	TK_ILLEGAL
)

type Token struct {
	Type   TokenType
	Name   string
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
		if isDigit(l.ch) {
			p := l.pos
			for isDigit(l.ch) {
				l.readChar()
			}
			token = newToken(TK_INT, l.ch, l.pos-p)
			token.Name = l.input[p:l.pos]
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

type Value interface{}

type State int

const (
	ObjectStart State = iota
	ObjectOpen
	ObjectProperty

	PropertyStart
	PropertyKey
	PropertyValue
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
	} else {
		val = nil
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
		case ObjectProperty:
			if p.cur.Type == TK_EBRACE {
				p.nextToken()
				return obj
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
				p.nextToken()
			}
		case PropertyKey:
			if p.cur.Type == TK_COLON {
				propState = PropertyValue
				p.nextToken()
			}
		case PropertyValue:
			prop.val = p.cur.Name
			return prop
		}
	}

	return prop
}

func (r *Root) PrintFromRoot() {
	rootValue := *r.val
	switch rootValue.(type) {
	case Object:
		obj, _ := rootValue.(Object)
		for _, c := range obj.children {
			fmt.Printf("KEY: %s, VALUE: %s\n", c.key, c.val)
		}
	}
}
