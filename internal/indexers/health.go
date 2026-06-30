package indexers

import "errors"

func StatusCode(err error) *int32 {
	var statusErr StatusError
	if !errors.As(err, &statusErr) {
		return nil
	}
	code := int32(statusErr.StatusCode)
	return &code
}

func IsPermanentFailure(statusCode *int32) bool {
	if statusCode == nil {
		return false
	}
	code := *statusCode
	if code == 429 || code == 408 {
		return false
	}
	return code >= 400 && code < 500
}

func StatusCodeFromDetails(details map[string]interface{}) *int32 {
	value, ok := details["statusCode"]
	if !ok {
		return nil
	}
	switch typed := value.(type) {
	case int:
		code := int32(typed)
		return &code
	case int32:
		code := typed
		return &code
	case int64:
		code := int32(typed)
		return &code
	case float64:
		code := int32(typed)
		return &code
	default:
		return nil
	}
}
