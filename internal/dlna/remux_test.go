package dlna

import (
	"strings"
	"testing"
	"time"
)

func TestRemuxCachePathIncludesFileIdentityAndModTime(t *testing.T) {
	now := time.Unix(100, 0)
	left := remuxCachePath("cache", "/media/a.mkv", now, 10)
	right := remuxCachePath("cache", "/media/a.mkv", now.Add(time.Second), 10)
	other := remuxCachePath("cache", "/media/a.mkv", now, 11)

	if left == right || left == other {
		t.Fatalf("cache paths not unique: %q %q %q", left, right, other)
	}
	if !strings.HasSuffix(left, ".mkv") {
		t.Fatalf("cache path = %q", left)
	}
}
