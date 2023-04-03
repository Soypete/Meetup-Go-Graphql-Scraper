// auth is a module used by this bot to connect to the meetup api and
// perform any needed oauth procedures or token management.
package auth

// AuthManager contains logic needed to get all api auth and tokens.
type AuthManager interface {
	GetBearerToken() (string, error)
}
