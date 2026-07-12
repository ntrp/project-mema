package subtitles

import "testing"

func TestClusterAProvidersAreRuntimeSupported(t *testing.T) {
	for _, key := range []string{"assrt", "betaseries", "gestdown", "jimaku", "regielive", "subdl", "subsarr", "subsource", "subsro", "subtis", "subx", "whisperai"} {
		if !RuntimeSupported(key) {
			t.Fatalf("expected %s to be runtime supported", key)
		}
	}
}
