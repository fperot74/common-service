package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	cs "github.com/cloudtrust/common-service"
	"github.com/gbrlsnchs/jwt"
)

// MakeHTTPBasicAuthenticationMW retrieve the token from the HTTP header 'Basic' and
// check if the password value match the allowed one.
// If there is no such header, the request is not allowed.
// If the password is correct, the username is added into the context:
//   - username: username extracted from the token
func MakeHTTPBasicAuthenticationMW(passwordToMatch string, logger cs.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var authorizationHeader = req.Header.Get("Authorization")

			if authorizationHeader == "" {
				logger.Log("Authorization Error", "Missing Authorization header")
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Missing Authorization header"), w)
				return
			}

			var matched, _ = regexp.MatchString(`^[Bb]asic *`, authorizationHeader)

			if !matched {
				logger.Log("Authorization Error", "Missing basic token")
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Missing basic token"), w)
				return
			}

			var splitToken = strings.Split(authorizationHeader, "Basic ")
			if len(splitToken) < 2 {
				splitToken = strings.Split(authorizationHeader, "basic ")
			}

			var accessToken = splitToken[1]

			// Decode base 64
			decodedToken, err := base64.StdEncoding.DecodeString(accessToken)

			if err != nil {
				logger.Log("Authorization Error", "Invalid base64 token")
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Invalid token"), w)
				return
			}

			// Extract username & password values
			var tokenSubparts = strings.Split(string(decodedToken), ":")

			if len(tokenSubparts) != 2 {
				logger.Log("Authorization Error", "Invalid token format (username:password)")
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Invalid token"), w)
				return
			}

			var username = tokenSubparts[0]
			var password = tokenSubparts[1]

			// Check password match
			if password != passwordToMatch {
				logger.Log("Authorization Error", "Invalid password value")
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Invalid token"), w)
				return
			}

			var ctx = context.WithValue(req.Context(), cs.CtContextUsername, username)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// KeycloakClient is the interface of the keycloak client.
type KeycloakClient interface {
	VerifyToken(realmName string, accessToken string) error
}

// MakeHTTPOIDCTokenValidationMW retrieve the oidc token from the HTTP header 'Bearer' and
// check its validity for the Keycloak instance binded to the component.
// If there is no such header, the request is not allowed.
// If the token is validated, the following informations are added into the context:
//   - access_token: the recieved access token in raw format
//   - realm: realm name extracted from the Issuer information of the token
//   - username: username extracted from the token
func MakeHTTPOIDCTokenValidationMW(keycloakClient KeycloakClient, audienceRequired string, logger cs.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var authorizationHeader = req.Header.Get("Authorization")

			if authorizationHeader == "" {
				logger.Log("Authorization Error", "Missing Authorization header")
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Missing Authorization header"), w)
				return
			}

			var matched, _ = regexp.MatchString(`^[Bb]earer *`, authorizationHeader)

			if !matched {
				logger.Log("Authorization Error", "Missing bearer token")
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Missing bearer token"), w)
				return
			}

			var splitToken = strings.Split(authorizationHeader, "Bearer ")

			if len(splitToken) < 2 {
				splitToken = strings.Split(authorizationHeader, "bearer ")
			}

			var accessToken = splitToken[1]

			payload, _, err := jwt.Parse(accessToken)
			if err != nil {
				logger.Log("Authorization Error", err)
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Invalid token"), w)
				return
			}

			var userID, username, issuer, realm string
			var groups []string

			// The audience in JWT may be a string array or a string.
			// First we try with a string array, if a failure occurs we try with a string
			{
				var jot TokenAudienceStringArray
				if err = jwt.Unmarshal(payload, &jot); err == nil {
					userID = jot.Subject
					username = jot.Username
					issuer = jot.Issuer
					var splitIssuer = strings.Split(issuer, "/auth/realms/")
					realm = splitIssuer[1]
					groups = extractGroups(jot.Groups)

					if !assertMatchingAudience(jot.Audience, audienceRequired) {
						logger.Log("Authorization Error", "Incorrect audience")
						httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Invalid token"), w)
						return
					}
				}
			}

			if err != nil {
				var jot TokenAudienceString
				if err = jwt.Unmarshal(payload, &jot); err == nil {
					userID = jot.Subject
					username = jot.Username
					issuer = jot.Issuer
					var splitIssuer = strings.Split(issuer, "/auth/realms/")
					realm = splitIssuer[1]
					groups = extractGroups(jot.Groups)

					if jot.Audience != audienceRequired {
						logger.Log("Authorization Error", "Incorrect audience")
						httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Invalid token"), w)
						return
					}
				} else {
					logger.Log("Authorization Error", err)
					httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Invalid token"), w)
					return
				}
			}

			if err = keycloakClient.VerifyToken(realm, accessToken); err != nil {
				logger.Log("Authorization Error", err)
				httpErrorHandler(context.TODO(), http.StatusForbidden, fmt.Errorf("Invalid token"), w)
				return
			}

			var ctx = context.WithValue(req.Context(), cs.CtContextAccessToken, accessToken)
			ctx = context.WithValue(ctx, cs.CtContextRealm, realm)
			ctx = context.WithValue(ctx, cs.CtContextUserID, userID)
			ctx = context.WithValue(ctx, cs.CtContextUsername, username)
			ctx = context.WithValue(ctx, cs.CtContextGroups, groups)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func assertMatchingAudience(jwtAudiences []string, requiredAudience string) bool {
	for _, jwtAudience := range jwtAudiences {
		if requiredAudience == jwtAudience {
			return true
		}
	}

	return false
}

func extractGroups(kcGroups []string) []string {
	var groups = []string{}

	for _, kcGroup := range kcGroups {
		groups = append(groups, strings.TrimPrefix(kcGroup, "/"))
	}

	return groups
}

// TokenAudienceStringArray is JWT token and the custom fields present in OIDC Token provided by Keycloak.
// Audience can be a string or a string array according the specification.
// The libraries are not supporting tit at this time (Fix in progress), meanwhile we circumvent it with a quick fix.
type TokenAudienceStringArray struct {
	hdr            *header
	Issuer         string   `json:"iss,omitempty"`
	Subject        string   `json:"sub,omitempty"`
	Audience       []string `json:"aud,omitempty"`
	ExpirationTime int64    `json:"exp,omitempty"`
	NotBefore      int64    `json:"nbf,omitempty"`
	IssuedAt       int64    `json:"iat,omitempty"`
	ID             string   `json:"jti,omitempty"`
	Username       string   `json:"preferred_username,omitempty"`
	Groups         []string `json:"groups,omitempty"`
}

// TokenAudienceString is JWT token with an Audience field represented as a string
type TokenAudienceString struct {
	hdr            *header
	Issuer         string   `json:"iss,omitempty"`
	Subject        string   `json:"sub,omitempty"`
	Audience       string   `json:"aud,omitempty"`
	ExpirationTime int64    `json:"exp,omitempty"`
	NotBefore      int64    `json:"nbf,omitempty"`
	IssuedAt       int64    `json:"iat,omitempty"`
	ID             string   `json:"jti,omitempty"`
	Username       string   `json:"preferred_username,omitempty"`
	Groups         []string `json:"groups,omitempty"`
}

type header struct {
	Algorithm   string `json:"alg,omitempty"`
	KeyID       string `json:"kid,omitempty"`
	Type        string `json:"typ,omitempty"`
	ContentType string `json:"cty,omitempty"`
}
