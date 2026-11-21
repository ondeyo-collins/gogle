/*
type eventType int
	errors I'm worried About:
		a < in a character string
*/
package gogle

import (
	"bufio"
	"fmt"
	"errors"
	"io"
	"unicode"
)

/*
	consider:
	<html>
		<p> my name is collins </p>
		<a href="example.com"> go to example.com </a>
	</html>
*/
type eventType int

const (
	EventTypeTextNode       eventType = iota // "my", "name", "is", "collins" - an eventTextNode is triggered for each word that is not in an end tag or a start tag
	EventTypeEndDocument                     // triggered right after the > symbol of </html> is read
	EventTypeOpeningTag                      // triggered right after the h in <html> is read (and other tags too)
	EventTypeClosingTag                      // triggered right after the / in </html> is read (and other tags too)
	EventTypeAttributeKey                    // triggered right after the h in <a href="example.com"> is read
	EventTypeAttributeValue                  // triggered right after the = in <a href="example.com"> is read
	EventTypeComment                         //
	EventTypeUnknown                         // used to signal that an event.type is not set
)

var (
	eventErrorInvalidWhiteSpace = errors.New("Invalid whitespace")
	eventErrorInvalidTag        = errors.New("Invalid tag")
	eventErrorInvalidEndOfFile  = errors.New("Invalid end of file")
	eventErrorInvalidNewLine    = errors.New("Invalid new line")
)

type ParserEvent struct {
	Type   eventType
	Buffer []rune
	Error  error
}
type parserState byte

const (
	parserStateOpeningTag parserState = iota
	parserStateUnknown
	parserStateAttributeKey
	parserStateAttributeValue
)

func Parse(input io.Reader) func(*ParserEvent) {
	var (
		nextRune              rune
		reader                            = bufio.NewReader(input)
		currentState          parserState = parserStateUnknown
		numberOfSpacesSkipped int         = 0
		lineNumber            int         = 1
	)

	// skip whitespace and read the first character
	// at the end of input, returns 0
	var skipWhiteSpace = func() rune {
		numberOfSpacesSkipped = 0
	read:
		character, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return 0
			} else {
				panic(fmt.Errorf("Unknown Error, [DEBUGGING]: file: index/htmlparser.go, function: skipWhiteSpace: %w", err))
			}
		}
		if nextRune == '\n' {
			lineNumber++
		}
		if !unicode.IsSpace(character) {
			return character
		} else {
			numberOfSpacesSkipped++
		}
		goto read
	}

	return func(event *ParserEvent) {
		event.Buffer = event.Buffer[:0]
		event.Error = nil

	functionStart:
		nextRune = skipWhiteSpace() //skipWhiteSpace between tokens
		if nextRune == 0 {
			event.Type = EventTypeEndDocument
			return
		}

		// these states are used to make sure we don't break out of an opening tag in between calls untill we've read the entire tag
		if currentState == parserStateOpeningTag {
			reader.UnreadRune()
			goto attributeKey
		}

		if currentState == parserStateAttributeKey {
			if nextRune == '=' {
				goto attributeValue
			}
			goto attributeKey
		}

		if currentState == parserStateAttributeValue {
			reader.UnreadRune()
			goto attributeKey
		}

		// possible tag
		if nextRune == '<' {
			nextRune = skipWhiteSpace() // allow <a < a
			if nextRune == '/' {
				goto closingTag
			}

			if unicode.IsLetter(nextRune) {
				reader.UnreadRune()
				goto openingTag
			}

			if nextRune == '!' {
				// discard comment or document declaration
				event.Type = EventTypeComment
				var secondRune = skipWhiteSpace()
				for {
					nextRune = skipWhiteSpace()
					if nextRune == 0 {
						event.Error = eventErrorInvalidEndOfFile
						return
					}

					if nextRune == '-' {
						nextRune, _, _ = reader.ReadRune()
						if nextRune == '-' {
							nextRune, _, _ = reader.ReadRune()
							if nextRune == '>' && secondRune == '-' {
								break
							}
						}
					}
					if nextRune == '>' && secondRune != '-' {
						break
					}
				}
				return
			}

			event.Error = fmt.Errorf("%w, line: %v", eventErrorInvalidTag, lineNumber)
			return
		}
		reader.UnreadRune()
		goto textNode;

	openingTag:
		{
			currentState = parserStateOpeningTag
			event.Type = EventTypeOpeningTag
			// read up to first space or >
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if numberOfSpacesSkipped > 0 {
					if nextRune == '>' { // <a >
						currentState = parserStateUnknown
						break
					}
					reader.UnreadRune() // we have a rune, but we skipped a few spaces to get it, it is an attributeKey
					break
				}

				if nextRune == '>' { // <a>
					currentState = parserStateUnknown
					break
				}
				event.Buffer = append(event.Buffer, nextRune)
			}
			return
		}

	attributeKey:
		{
			currentState = parserStateAttributeKey
			event.Type = EventTypeAttributeKey
			// read up until a space or '=' or >
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if numberOfSpacesSkipped > 0 && nextRune != '=' && nextRune != '>' { // a second key
					currentState = parserStateAttributeKey
					break
				}

				// NOTE: numberOfSpacesSkipped makes no difference here
				if nextRune == '=' { // <a href = ..
					reader.UnreadRune() // NOTE: '=' is used by the next call to go to attributeValue state (see functionStart label)
					currentState = parserStateAttributeKey
					break
				}
				if nextRune == '>' { // <a blackButton  >,
					currentState = parserStateUnknown
					break
				}
				event.Buffer = append(event.Buffer, nextRune)
			}
			return
		}

	attributeValue:
		{
			currentState = parserStateAttributeValue
			event.Type = EventTypeAttributeValue
			// read up until '>' or space, or quotes (if it started with a quote)
			var firstRune = skipWhiteSpace()
			if !(firstRune == '\'' || firstRune == '"') {
				reader.UnreadRune()
			}
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if nextRune == '>' {
					currentState = parserStateUnknown
					break
				}

				if numberOfSpacesSkipped > 0 && !(firstRune == '\'' || firstRune == '"') {
					reader.UnreadRune();
					break
				}
				if nextRune == firstRune && (firstRune == '\'' || firstRune == '"') {
					break
				}
				event.Buffer = append(event.Buffer, nextRune)
			}
			return
		}

	textNode:
		{
			event.Type = EventTypeTextNode
			// read up to a space or <
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if numberOfSpacesSkipped > 0 {
					reader.UnreadRune()
					if len(event.Buffer) < 1 {
						// NOTE: so we dont return empty space
						goto functionStart
					}
					break
				}
				if nextRune == '<' {
					reader.UnreadRune();
					break;
				}

				event.Buffer = append(event.Buffer, nextRune)
			}
			if len(event.Buffer) < 1 {
				// NOTE: so we dont return empty space
				goto functionStart
			}
			return
		}

	closingTag:
		{
			event.Type = EventTypeClosingTag
			// read up to >
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if nextRune == '>' {
					break
				}
				event.Buffer = append(event.Buffer, nextRune)
			}
			return
		}
	}
}
