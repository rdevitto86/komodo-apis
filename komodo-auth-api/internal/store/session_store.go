package store

type sessionStore struct {
	GetLocalSession    func(token string) (string, bool)
	SetLocalSession    func(token string, userID string)
	DeleteLocalSession func(token string)
	ClearLocalSessions func()
}

// Stores session tokens for the auth service
var localSessionData = make(map[string]string)

var LocalSessionStore = func() *sessionStore {
	return &sessionStore{
		GetLocalSession: func(token string) (string, bool) {
			value, exists := localSessionData[token]
			return value, exists
		},
		SetLocalSession: func(token string, userID string) {
			localSessionData[token] = userID
		},
		DeleteLocalSession: func(token string) {
			delete(localSessionData, token)
		},
		ClearLocalSessions: func() {
			for k := range localSessionData {
				delete(localSessionData, k)
			}
		},
	}
}
