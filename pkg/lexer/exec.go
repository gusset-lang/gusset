package lexer

import (
	"unicode"
)

func Exec(l *Lexer) {
exec:
	for {
		r, err := l.nextRune()
		if err != nil {
			break
		}

		switch {
		case r == EOF_RUNE:
			l.sendEOF()
			break exec
		case r == '\n':
			l.sendNewLine()
			l.nextLine()
			continue exec
		case unicode.IsSpace(r):
			l.nextCol()
			continue exec
		}

		item := l.itemFromRune(r)
		if item != nil {
			l.sendItem(item)
			l.nextCol()
			continue
		}

		if r == ':' {
			second, err := l.nextRune()
			if err != nil || second == EOF_RUNE {
				break
			}
			switch second {
			case '=':
				l.sendShortVar()
				l.advance(2)
				continue exec
			default:
				if err := l.backup(); err != nil {
					break exec
				}
				l.sendColon()
				l.nextCol()
				continue exec
			}
		}

		if err := l.backup(); err != nil {
			break
		}

		if r == '"' {
			item, runeCount, err := l.itemFromStringLiteral(l.pos)
			if err != nil {
				break
			}
			l.sendItem(item)
			l.advance(runeCount)
			continue
		}

		if unicode.IsDigit(r) {
			item, runeCount, err := l.itemFromNumeric(l.pos)
			if err != nil {
				break
			}
			l.sendItem(item)
			l.advance(runeCount)
			continue
		}

		item, runeCount, err := l.itemFromAlphanum(l.pos)
		if err != nil {
			break
		}
		if item != nil {
			l.sendItem(item)
			l.advance(runeCount)
			continue
		}

		l.sendIllegal(r)
	}

	close(l.result)
}
