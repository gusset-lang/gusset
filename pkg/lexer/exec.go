package lexer

import (
	"io"
	"unicode"
)

func Exec(l *Lexer) {
	for {
		start := l.pos
		r, err := l.next()
		if err != nil {
			if err == io.EOF {
				l.sendEOF()
			}
			break
		}
		if r == EOF_RUNE {
			l.sendEOF()
			break
		}

		if r == '\n' {
			l.sendNewLine(start)
			continue
		}

		if unicode.IsSpace(r) {
			continue
		}

		matchedSymbolic, err := l.matchRuneSequence(start, r)
		if err != nil {
			break
		}
		if matchedSymbolic {
			continue
		}

		if r == '"' {
			if err := l.collectStringLiteral(start); err != nil {
				break
			}
			continue
		}

		if r == ':' {
			if err := l.collectSymbol(start); err != nil {
				break
			}
			continue
		}

		// match template
		if r == '`' {
			if err := l.collectTemplateLiteral(start); err != nil {
				break
			}
			continue
		}

		if unicode.IsDigit(r) {
			item, err := l.itemFromNumeric(start, r)
			if err != nil {
				break
			}
			l.sendItem(item)
			continue
		}

		item, err := l.itemFromAlphanum(start, r)
		if err != nil {
			break
		}
		if item != nil {
			l.sendItem(item)
			continue
		}
	}

	close(l.result)
}

// func Exec(l *Lexer) {
// exec:
// 	for {
// 		r, err := l.nextRune()
// 		if err != nil {
// 			break
// 		}

// 		switch {
// 		case r == EOF_RUNE:
// 			l.sendEOF()
// 			break exec
// 		case r == '\n':
// 			l.sendNewLine()
// 			l.nextLine()
// 			continue exec
// 		case unicode.IsSpace(r):
// 			l.nextCol()
// 			continue exec
// 		}

// 		item := l.itemFromRune(r)
// 		if item != nil {
// 			l.sendItem(item)
// 			l.nextCol()
// 			continue
// 		}

// 		if r == '=' {
// 			second, err := l.nextRune()
// 			if err != nil || second == EOF_RUNE {
// 				break
// 			}
// 			switch second {
// 			case '>':
// 				l.sendArrow()
// 				l.advance(2)
// 				continue exec
// 			default:
// 				if err := l.backup(); err != nil {
// 					break exec
// 				}
// 				l.sendAssign()
// 				l.nextCol()
// 				continue exec
// 			}
// 		}

// 		if r == ':' {
// 			second, err := l.nextRune()
// 			if err != nil || second == EOF_RUNE {
// 				break
// 			}
// 			switch {
// 			case second == '=':
// 				l.sendShortVar()
// 				l.advance(2)
// 				continue exec
// 			case second == '_' || unicode.IsLetter(second):
// 				if err := l.backup(); err != nil {
// 					break exec
// 				}
// 				item, runeCount, err := l.itemFromSymbol(l.pos)
// 				if err != nil {
// 					break exec
// 				}
// 				l.sendItem(item)
// 				l.advance(runeCount)
// 				continue exec
// 			default:
// 				l.sendIllegal(r)
// 				break exec
// 			}
// 		}

// 		if err := l.backup(); err != nil {
// 			break
// 		}

// 		if r == '`' {
// 			item, newPos, err := l.itemFromTemplateLiteral()
// 			if err != nil {
// 				break
// 			}
// 			l.sendItem(item)
// 			l.setPosition(newPos)
// 			continue
// 		}

// 		if r == '"' {
// 			item, runeCount, err := l.itemFromStringLiteral(l.pos)
// 			if err != nil {
// 				break
// 			}
// 			l.sendItem(item)
// 			l.advance(runeCount)
// 			continue
// 		}

// 		if unicode.IsDigit(r) {
// 			item, runeCount, err := l.itemFromNumeric(l.pos)
// 			if err != nil {
// 				break
// 			}
// 			l.sendItem(item)
// 			l.advance(runeCount)
// 			continue
// 		}

// 		item, runeCount, err := l.itemFromAlphanum(l.pos)
// 		if err != nil {
// 			break
// 		}
// 		if item != nil {
// 			l.sendItem(item)
// 			l.advance(runeCount)
// 			continue
// 		}

// 		l.sendIllegal(r)
// 	}

// 	close(l.result)
// }
