package httpapi

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const streamTokenTTL = 6 * time.Hour

func newStreamTokenSecret() []byte {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err == nil {
		return secret
	}
	return []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
}

func (s *Server) newStreamToken(mediaID uuid.UUID, filePath string, expires int64) string {
	mac := hmac.New(sha256.New, s.streamSecret)
	_, _ = mac.Write(streamTokenMessage(mediaID, filePath, expires))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (s *Server) validStreamToken(mediaID uuid.UUID, filePath string, expires *int64, token *string) bool {
	if expires == nil || token == nil || *expires <= s.now().Unix() {
		return false
	}
	expected := s.newStreamToken(mediaID, filePath, *expires)
	return hmac.Equal([]byte(expected), []byte(*token))
}

func streamTokenMessage(mediaID uuid.UUID, filePath string, expires int64) []byte {
	return []byte(mediaID.String() + "\n" + filePath + "\n" + strconv.FormatInt(expires, 10))
}
