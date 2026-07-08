package dlna

import (
	"net/http"
	"strings"
)

const (
	seekModeByte          = "byte"
	seekModeTime          = "time"
	seekModeBoth          = "both"
	seekModeTimeExclusive = "time_exclusive"
	seekModeNone          = "none"
)

func rejectUnsupportedSeek(w http.ResponseWriter, r *http.Request, profile RendererProfile) bool {
	if isSeekRange(r.Header.Get("Range")) && !profileAllowsByteSeek(profile) {
		http.Error(w, "DLNA byte seeking is disabled for this renderer profile", http.StatusRequestedRangeNotSatisfiable)
		return true
	}
	if strings.TrimSpace(r.Header.Get("TimeSeekRange.dlna.org")) != "" && !profileAllowsTimeSeek(profile) {
		http.Error(w, "DLNA time seeking is disabled for this renderer profile", http.StatusRequestedRangeNotSatisfiable)
		return true
	}
	return false
}

func profileAllowsByteSeek(profile RendererProfile) bool {
	switch normalizedSeekMode(profile) {
	case seekModeTime, seekModeTimeExclusive, seekModeNone:
		return false
	default:
		return true
	}
}

func profileAllowsTimeSeek(profile RendererProfile) bool {
	switch normalizedSeekMode(profile) {
	case seekModeTime, seekModeBoth, seekModeTimeExclusive:
		return true
	default:
		return false
	}
}

func normalizedSeekMode(profile RendererProfile) string {
	switch strings.ToLower(strings.TrimSpace(profile.DeliveryRules.SeekMode)) {
	case seekModeTime, seekModeBoth, seekModeTimeExclusive, seekModeNone:
		return strings.ToLower(strings.TrimSpace(profile.DeliveryRules.SeekMode))
	default:
		return seekModeByte
	}
}
