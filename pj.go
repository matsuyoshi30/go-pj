package pj

import ()

type TokenType int

const (
	TK_INT TokenType = iota
	TK_STR

	TK_SBRACES // {
	TK_EBRACES // }
	TK_COLON   // :
	TK_DQUOTE  // "

	TK_ILLEGAL
)

type Token struct {
	Type   TokenType
	Name   string
	Length int
}

type Tokenizer struct {
	input string
	pos   int // current position
}

func NewTokenizer(str string) *Tokenizer {
	return &Tokenizer{input: str}
}

func (t *Tokenizer) Tokenize() []Token {
	tokens := make([]Token, 0)

	for t.pos = 0; t.pos < len(t.input); t.pos++ {
		ch := t.input[t.pos]
		if isWhiteSpace(ch) {
			continue
		}

		switch ch {
		case '{':
			tokens = append(tokens, newToken(TK_SBRACES, `{`, 1))
		case '}':
			tokens = append(tokens, newToken(TK_EBRACES, `}`, 1))
		case ':':
			tokens = append(tokens, newToken(TK_COLON, `:`, 1))
		case '"':
			tokens = append(tokens, newToken(TK_DQUOTE, `"`, 1))
		default:
			if isLetter(ch) {
				p := t.pos
				for isLetter(t.input[t.pos]) {
					t.pos++
				}
				tokens = append(tokens, newToken(TK_STR, t.input[p:t.pos], t.pos-p))
				t.pos--
			} else if isDigit(ch) {
				p := t.pos
				for isDigit(t.input[t.pos]) {
					t.pos++
				}
				tokens = append(tokens, newToken(TK_INT, t.input[p:t.pos], t.pos-p))
				t.pos--
			} else {
				tokens = append(tokens, newToken(TK_ILLEGAL, "", 0))
			}
		}
	}

	return tokens
}

func newToken(ty TokenType, name string, length int) Token {
	return Token{
		Type:   ty,
		Name:   name,
		Length: length,
	}
}

func isWhiteSpace(c byte) bool {
	if c == ' ' {
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
