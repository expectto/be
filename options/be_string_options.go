package options

type StringOption Option

const (
	// considering here 0 as no option
	_ = 0

	// Alpha option represents the presence of alphabetical characters.
	Alpha StringOption = 1 << (iota - 1)

	// Numeric option represents the presence of numeric characters.
	Numeric

	// Whitespace option represents the presence of whitespace characters.
	Whitespace

	// Dots option represents the presence of dot characters.
	Dots

	// Punctuation option represents the presence of punctuation characters.
	Punctuation

	// SpecialCharacters option represents the presence of special characters.
	SpecialCharacters
)

func (f StringOption) String() string {
	switch f {
	case Alpha:
		return "alpha"
	case Numeric:
		return "numeric"
	case Whitespace:
		return "whitespace"
	case Dots:
		return "dots"
	case Punctuation:
		return "punctuation"
	case SpecialCharacters:
		return "special characters"
	default:
		return "unknown"
	}
}

// ExtractStringOptions extracts individual options from a combined string option
func ExtractStringOptions(combined StringOption) []StringOption {
	var options []StringOption

	// Iterate over each bit position to check if the option is set
	for flag := StringOption(1); flag <= combined; flag <<= 1 {
		if combined&flag != 0 {
			options = append(options, flag)
		}
	}

	return options
}
