package httpapi

import (
	"net/http"

	"media-manager/internal/delivery"
	"media-manager/internal/dlna"
)

func (s *Server) TraceDLNAProfileMatch(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DLNAProfileMatchTraceRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	trace := s.dlnaDiagnosticsManager().TraceRendererProfile(r.Context(), dlnaRequestFromTrace(body))
	writeJSON(w, http.StatusOK, dlnaProfileMatchTraceResponse(trace))
}

func (s *Server) TraceDLNADeliveryDecision(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DLNADeliveryTraceRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	trace := s.dlnaDiagnosticsManager().TraceDeliveryDecision(r.Context(), dlnaDeliveryTraceInput(body))
	writeJSON(w, http.StatusOK, dlnaDeliveryTraceResponse(trace))
}

func (s *Server) dlnaDiagnosticsManager() *dlna.Manager {
	if s.dlna != nil {
		return s.dlna
	}
	return dlna.NewManager(s.settings, "")
}

func dlnaRequestFromTrace(body DLNAProfileMatchTraceRequest) dlna.RendererRequest {
	return dlna.RendererRequest{
		UserAgent:    traceStringValue(body.UserAgent),
		FriendlyName: traceStringValue(body.FriendlyName),
		ClientIP:     traceStringValue(body.DeviceIp),
		RendererUUID: traceStringValue(body.RendererUuid),
		Headers:      headerFromTrace(body.Headers),
	}
}

func dlnaRequestFromDeliveryTrace(body DLNADeliveryTraceRequest) dlna.RendererRequest {
	return dlna.RendererRequest{
		UserAgent:    traceStringValue(body.UserAgent),
		FriendlyName: traceStringValue(body.FriendlyName),
		ClientIP:     traceStringValue(body.DeviceIp),
		RendererUUID: traceStringValue(body.RendererUuid),
		Headers:      headerFromTrace(body.Headers),
	}
}

func dlnaDeliveryTraceInput(body DLNADeliveryTraceRequest) dlna.DeliveryTraceInput {
	return dlna.DeliveryTraceInput{
		Request:    dlnaRequestFromDeliveryTrace(body),
		ProfileID:  traceStringValue(body.ProfileId),
		MediaPath:  traceStringValue(body.MediaPath),
		ObjectID:   traceStringValue(body.ObjectId),
		ResourceID: traceStringValue(body.ResourceId),
		StreamMode: traceStringValue(body.StreamMode),
		Probe:      probeFromTrace(body.Probe),
	}
}

func headerFromTrace(input *map[string]string) http.Header {
	headers := http.Header{}
	if input == nil {
		return headers
	}
	for key, value := range *input {
		headers.Set(key, value)
	}
	return headers
}

func probeFromTrace(input *DLNAMediaProbeRequest) delivery.ProbeResult {
	if input == nil {
		return delivery.ProbeResult{}
	}
	tracks := []delivery.Track{}
	if input.VideoCodec != nil {
		tracks = append(tracks, delivery.Track{
			Type:   delivery.TrackVideo,
			Codec:  input.VideoCodec,
			Height: input.Height,
		})
	}
	if input.AudioCodec != nil {
		tracks = append(tracks, delivery.Track{Type: delivery.TrackAudio, Codec: input.AudioCodec})
	}
	return delivery.ProbeResult{
		Container: delivery.Container{FormatName: input.Container},
		Tracks:    tracks,
	}
}

func dlnaProfileMatchTraceResponse(trace dlna.RendererProfileTrace) DLNAProfileMatchTraceResponse {
	explanation := trace.Match.Explanation
	return DLNAProfileMatchTraceResponse{
		ProfileId:           trace.Match.Profile.ID,
		ProfileName:         trace.Match.Profile.Name,
		SourceProfileId:     explanation.SourceProfileID,
		MatchSource:         explanation.MatchSource,
		MatchReason:         dlna.RendererMatchReason(explanation),
		WinningRule:         explanation.WinningRule,
		FallbackPath:        explanation.FallbackPath,
		Score:               int32(explanation.Score),
		CandidateProfileIds: append([]string{}, explanation.CandidateProfileIDs...),
		HeadersSummary:      append([]string{}, trace.HeadersSummary...),
		RuleTrace:           dlnaProfileRuleTraceResponse(trace.Rules),
		Candidates:          dlnaProfileMatchCandidateTraceResponse(trace.Candidates),
	}
}

func dlnaProfileMatchCandidateTraceResponse(traces []dlna.RendererProfileMatchCandidate) []DLNAProfileMatchCandidate {
	values := make([]DLNAProfileMatchCandidate, 0, len(traces))
	for _, trace := range traces {
		values = append(values, DLNAProfileMatchCandidate{
			ProfileId:    trace.ProfileID,
			ProfileName:  trace.ProfileName,
			Score:        int32(trace.Score),
			MinimumScore: int32(trace.MinimumScore),
			Priority:     int32(trace.Priority),
			Qualified:    trace.Qualified,
			Selected:     trace.Selected,
			RuleTrace:    dlnaProfileRuleTraceResponse(trace.RuleTrace),
		})
	}
	return values
}

func dlnaProfileRuleTraceResponse(traces []dlna.RendererProfileRuleTrace) []DLNAProfileRuleTrace {
	values := make([]DLNAProfileRuleTrace, 0, len(traces))
	for _, trace := range traces {
		values = append(values, DLNAProfileRuleTrace{
			ProfileId:   trace.ProfileID,
			ProfileName: trace.ProfileName,
			Field:       trace.Field,
			Value:       trace.Value,
			Rule:        trace.Rule,
			Score:       int32(trace.Score),
			Result:      DLNAProfileRuleTraceResult(trace.Result),
		})
	}
	return values
}

func dlnaDeliveryTraceResponse(trace dlna.DeliveryDecisionTrace) DLNADeliveryTraceResponse {
	return DLNADeliveryTraceResponse{
		ProfileId:        trace.ProfileID,
		ProfileName:      trace.ProfileName,
		MediaFileName:    trace.MediaFileName,
		ObjectId:         trace.ObjectID,
		ResourceId:       trace.ResourceID,
		StreamMode:       trace.StreamMode,
		DeliveryProtocol: DLNADeliveryTraceResponseDeliveryProtocol(trace.Decision.DeliveryProtocol),
		Mode:             DLNADeliveryTraceResponseMode(trace.Decision.Mode),
		VideoCodec:       trace.Decision.Plan.VideoCodec,
		AudioCodec:       trace.Decision.Plan.AudioCodec,
		ReasonCodes:      append([]string{}, trace.ReasonCodes...),
		CapabilityTrace:  dlnaCapabilityTraceResponse(trace.Trace),
	}
}

func dlnaCapabilityTraceResponse(traces []dlna.RendererCapabilityTrace) []DLNACapabilityTrace {
	values := make([]DLNACapabilityTrace, 0, len(traces))
	for _, trace := range traces {
		values = append(values, DLNACapabilityTrace{
			Field: trace.Field, Value: trace.Value, Rule: trace.Rule,
			Result: DLNACapabilityTraceResult(trace.Result),
		})
	}
	return values
}

func traceStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
