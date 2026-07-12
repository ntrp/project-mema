package providers

import "strings"

func alpha3Language(languageID string) string {
	key := strings.ToLower(strings.TrimSpace(languageID))
	switch key {
	case "ar", "ara", "arabic":
		return "ara"
	case "bg", "bul", "bulgarian":
		return "bul"
	case "cs", "ces", "cze", "czech":
		return "ces"
	case "da", "dan", "danish":
		return "dan"
	case "de", "deu", "ger", "german":
		return "deu"
	case "el", "ell", "gre", "greek":
		return "ell"
	case "en", "eng", "english":
		return "eng"
	case "fi", "fin", "finnish":
		return "fin"
	case "fr", "fra", "fre", "french":
		return "fra"
	case "he", "heb", "hebrew":
		return "heb"
	case "hu", "hun", "hungarian":
		return "hun"
	case "id", "ind", "indonesian":
		return "ind"
	case "it", "ita", "italian":
		return "ita"
	case "ja", "jpn", "japanese":
		return "jpn"
	case "ko", "kor", "korean":
		return "kor"
	case "nl", "nld", "dut", "dutch":
		return "nld"
	case "pl", "pol", "polish":
		return "pol"
	case "pt-br", "pob", "por-br", "brazilian portuguese":
		return "pob"
	case "pt", "por", "portuguese":
		return "por"
	case "ro", "ron", "rum", "romanian":
		return "ron"
	case "ru", "rus", "russian":
		return "rus"
	case "es", "spa", "spanish":
		return "spa"
	case "sv", "swe", "swedish":
		return "swe"
	case "th", "tha", "thai":
		return "tha"
	case "tr", "tur", "turkish":
		return "tur"
	case "uk", "ukr", "ukrainian":
		return "ukr"
	case "vi", "vie", "vietnamese":
		return "vie"
	case "zh", "zho", "chi", "chinese":
		return "zho"
	default:
		return key
	}
}
