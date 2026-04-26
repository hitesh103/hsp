package cmd

const (
	PriorityHigh = iota
	PriorityMedium
	PriorityLow
)

var highPriorityKeys = []string{
	"id", "name", "email", "token", "key", "secret", "password",
	"status", "message", "error", "code", "url", "type",
	"data", "result", "user", "userId", "access_token",
	"refresh_token", "auth", "authorization",
}

var lowPriorityKeys = []string{
	"_id", "created_at", "updated_at", "timestamp",
	"v", "__v", "index", "seq", "offset", "cursor",
	"page", "per_page", "total", "limit",
}

func GetPriority(key string) int {
	keyLower := toLower(key)
	if IsImportant(keyLower) {
		return PriorityHigh
	}
	if IsLowPriority(keyLower) {
		return PriorityLow
	}
	return PriorityMedium
}

func IsImportant(key string) bool {
	keyLower := toLower(key)
	for _, k := range highPriorityKeys {
		if keyLower == k || contains(keyLower, k) {
			return true
		}
	}
	return false
}

func IsLowPriority(key string) bool {
	keyLower := toLower(key)
	for _, k := range lowPriorityKeys {
		if keyLower == k || contains(keyLower, k) {
			return true
		}
	}
	return false
}

func DetectPriority(fieldName string) int {
	return GetPriority(fieldName)
}

func toLower(s string) string {
	if s == "" {
		return s
	}
	result := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}