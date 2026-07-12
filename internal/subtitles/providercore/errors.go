package providercore

import "errors"

var (
	ErrCaptchaRequired             = errors.New("captcha_required")
	ErrAntiCaptchaRequired         = errors.New("anti_captcha_required")
	ErrPrivateMembershipRequired   = errors.New("private_membership_required")
	ErrReleaseProvenanceRequired   = errors.New("release_provenance_required")
	ErrProviderBrokenUpstream      = errors.New("provider_broken_upstream")
	ErrProviderPrerequisiteMissing = errors.New("provider_prerequisite_missing")
)
