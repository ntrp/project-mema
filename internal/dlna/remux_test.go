package dlna

import (
	"strings"
	"testing"
	"time"
)

func TestRemuxCachePathIncludesFileIdentityAndModTime(t *testing.T) {
	now := time.Unix(100, 0)
	left := remuxCachePath("cache", "/media/a.mkv", now, 10, ".mkv")
	right := remuxCachePath("cache", "/media/a.mkv", now.Add(time.Second), 10, ".mkv")
	other := remuxCachePath("cache", "/media/a.mkv", now, 11, ".mkv")

	if left == right || left == other {
		t.Fatalf("cache paths not unique: %q %q %q", left, right, other)
	}
	if !strings.HasSuffix(left, ".mkv") {
		t.Fatalf("cache path = %q", left)
	}
}

func TestDLNAOutputArgsUseStructuredContainers(t *testing.T) {
	decision := remuxDecision()
	args := dlnaOutputArgs("/media/a.mkv", "pipe:1", decision, mpegtsOutputTarget())

	if !hasArgPair(args, "-c:v", "copy") || !hasArgPair(args, "-c:a", "copy") {
		t.Fatalf("args = %#v", args)
	}
	if !hasArgPair(args, "-f", "mpegts") {
		t.Fatalf("args = %#v", args)
	}
}

func hasArgPair(args []string, key string, value string) bool {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == key && args[i+1] == value {
			return true
		}
	}
	return false
}
