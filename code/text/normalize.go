package text

import (
	// "io"
	// 	"strings"
	"unicode"
	"unicode/utf8"
)

// text Normalization
// any multibyte character is decomposed to it's base rune(see the baseRune variable) if possible, otherwise it is ignored
// each rune is lowercased
// punctuations are removed (see the ignorePunctuation variable)
// whitespace is reduced to a single ' '
func Normalize(input *[]rune) {
	var (
		nextRune     rune
		previousRune rune
		originalLen  = len(*input)
		targetIndex  = 0
	)

	for index := 0; index < originalLen; index++ {
		nextRune = (*input)[index]
		if unicode.IsSpace(nextRune) { // whitespace
			if unicode.IsSpace(previousRune) {
				continue
			}
			nextRune = ' '
		}

		if utf8.RuneLen(nextRune) > 1 { // multibyte characters
			if base, ok := baseRune[nextRune]; !ok {
				continue
			} else {
				nextRune = base
			}
		}

		if ignorePunctuation[nextRune] {
			continue
		}

		// lowercase
		if unicode.IsUpper(nextRune) {
			nextRune = unicode.ToLower(nextRune)
		}

		previousRune = nextRune
		(*input)[targetIndex] = nextRune
		targetIndex++
	}
	*input = (*input)[:targetIndex]
}

var ignorePunctuation = map[rune]bool{
	// brackets and braces
	'(': true, ')': true, '[': true, ']': true, '{': true, '}': true,

	// sentence punctuation
	',': true, '.': true, ':': true, ';': true, '?': true, '!': true,

	// quotes and apostrophes
	'\'': true, '"': true, '`': true,

	// mathematical and slashes
	'+': true, '-': true, '*': true, '/': true, '\\': true, '=': true, '<': true, '>': true,

	// financial and commercial
	'$': true, '%': true, '@': true, '&': true, '_': true,

	// connectors and miscellaneous
	'#': true, '^': true, '~': true, '|': true,
}

var baseRune = map[rune]rune{
	// A
	'À': 'A', 'Á': 'A', 'Â': 'A', 'Ã': 'A', 'Ä': 'A', 'Å': 'A',
	'Ā': 'A', 'Ă': 'A', 'Ą': 'A', 'Ǎ': 'A',
	'à': 'a', 'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a',
	'ā': 'a', 'ă': 'a', 'ą': 'a', 'ǎ': 'a',

	// C
	'Ç': 'C', 'Ć': 'C', 'Ĉ': 'C', 'Ċ': 'C', 'Č': 'C',
	'ç': 'c', 'ć': 'c', 'ĉ': 'c', 'ċ': 'c', 'č': 'c',

	// D
	'Ď': 'D',
	'ď': 'd',

	// E
	'È': 'E', 'É': 'E', 'Ê': 'E', 'Ë': 'E',
	'Ē': 'E', 'Ĕ': 'E', 'Ė': 'E', 'Ę': 'E', 'Ě': 'E',
	'è': 'e', 'é': 'e', 'ê': 'e', 'ë': 'e',
	'ē': 'e', 'ĕ': 'e', 'ė': 'e', 'ę': 'e', 'ě': 'e',

	// G
	'Ĝ': 'G', 'Ğ': 'G', 'Ġ': 'G', 'Ģ': 'G',
	'ĝ': 'g', 'ğ': 'g', 'ġ': 'g', 'ģ': 'g',

	// H
	'Ĥ': 'H',
	'ĥ': 'h',

	// I
	'Ì': 'I', 'Í': 'I', 'Î': 'I', 'Ï': 'I',
	'Ĩ': 'I', 'Ī': 'I', 'Ĭ': 'I', 'Į': 'I', 'Ǐ': 'I',
	'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i',
	'ĩ': 'i', 'ī': 'i', 'ĭ': 'i', 'į': 'i', 'ǐ': 'i',

	// N
	'Ñ': 'N', 'Ń': 'N', 'Ņ': 'N', 'Ň': 'N',
	'ñ': 'n', 'ń': 'n', 'ņ': 'n', 'ň': 'n',

	// O
	'Ò': 'O', 'Ó': 'O', 'Ô': 'O', 'Õ': 'O', 'Ö': 'O',
	'Ō': 'O', 'Ŏ': 'O', 'Ő': 'O', 'Ǒ': 'O',
	'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o',
	'ō': 'o', 'ŏ': 'o', 'ő': 'o', 'ǒ': 'o',

	// R
	'Ŕ': 'R', 'Ŗ': 'R', 'Ř': 'R',
	'ŕ': 'r', 'ŗ': 'r', 'ř': 'r',

	// S
	'Ś': 'S', 'Ŝ': 'S', 'Ş': 'S', 'Š': 'S',
	'ś': 's', 'ŝ': 's', 'ş': 's', 'š': 's',

	// T
	'Ţ': 'T', 'Ť': 'T',
	'ţ': 't', 'ť': 't',

	// U
	'Ù': 'U', 'Ú': 'U', 'Û': 'U', 'Ü': 'U',
	'Ũ': 'U', 'Ū': 'U', 'Ŭ': 'U', 'Ů': 'U', 'Ű': 'U', 'Ų': 'U', 'Ǔ': 'U',
	'ù': 'u', 'ú': 'u', 'û': 'u', 'ü': 'u',
	'ũ': 'u', 'ū': 'u', 'ŭ': 'u', 'ů': 'u', 'ű': 'u', 'ų': 'u', 'ǔ': 'u',

	// W
	'Ẁ': 'W', 'Ẃ': 'W', 'Ŵ': 'W', 'Ẅ': 'W',
	'ẁ': 'w', 'ẃ': 'w', 'ŵ': 'w', 'ẅ': 'w',

	// Y
	'Ỳ': 'Y', 'Ý': 'Y', 'Ŷ': 'Y', 'Ÿ': 'Y', 'Ỹ': 'Y',
	'ỳ': 'y', 'ý': 'y', 'ŷ': 'y', 'ÿ': 'y', 'ỹ': 'y',

	// Z
	'Ź': 'Z', 'Ż': 'Z', 'Ž': 'Z',
	'ź': 'z', 'ż': 'z', 'ž': 'z',
}
