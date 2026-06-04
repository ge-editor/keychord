package keychord

// Locale-specific heuristic for Japanese IME.
// Some IMEs immediately convert vowels into kana
// during key chord input:
//
//	Ctrl+X i → Ctrl+X い
//
// To avoid chord mismatch, normalize five vowels
// and 'n' back to ASCII.
func NormalizeLocaleASCII(r rune) rune {
	// Full-width ASCII: U+FF01 ～ U+FF5E
	if r >= '！' && r <= '～' {
		r -= 0xFEE0
	}

	// Full-width space
	if r == '　' {
		return ' '
	}

	// Japanese vowel kana normalization
	// IME heuristic:
	// Ctrl+X i → Ctrl+X い → i
	switch r {
	case 'あ', 'ア':
		return 'a'
	case 'い', 'イ':
		return 'i'
	case 'う', 'ウ':
		return 'u'
	case 'え', 'エ':
		return 'e'
	case 'お', 'オ':
		return 'o'
	case 'ん', 'ン':
		return 'n'
	}

	return r
}
