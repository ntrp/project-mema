package episodepath

import "testing"

func TestNumbersTreatsSeriesChildrenAsSeasonsAndEpisodes(t *testing.T) {
	tests := []struct {
		path            string
		season, episode int32
		ok              bool
	}{
		{"Show/Season 2/03.mkv", 2, 3, true},
		{"Show/S03/Episode 004 - Title.mkv", 3, 4, true},
		{"Show/Specials/01.mkv", 0, 1, true},
		{"Show/Extras/01.mkv", 0, 0, false},
		{"Show/Season 2/trailer.mkv", 0, 0, false},
		{"Show/Nested/Show.S04E005.mkv", 4, 5, true},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			season, episode, ok := Numbers(test.path)
			if season != test.season || episode != test.episode || ok != test.ok {
				t.Fatalf("Numbers(%q) = S%dE%d, %v; want S%dE%d, %v", test.path, season, episode, ok, test.season, test.episode, test.ok)
			}
		})
	}
}
