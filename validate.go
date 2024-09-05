package gokinde

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type ValidatedJWT struct {
	// Kinde UserID
	UserID string

	// If the user has assigned roles in Kinde, the role keys are listed here.
	// Note that you need to configure your app:
	// 	> tokens > token customisation > access token > additional claims > enable roles
	Roles []string

	// All JWT claims, including custom ones
	Claims jwt.MapClaims
}

var InvalidJWTError = fmt.Errorf("token invalid")

// ValidateJWT returns a map of claims if the token is valid.
// If the token is invalid - InvalidJWTError is returned.
// Other errors could be returned if parsing or format errors occur.
func (cl *Client) ValidateJWT(jwtB64 string) (*ValidatedJWT, error) {
	token, err := jwt.Parse(jwtB64, cl.jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("jwt parse token: %w", err)
	}

	if !token.Valid {
		return nil, InvalidJWTError
	}

	claimsM, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("jwt claims format")
	}

	userId, ok := claimsM["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("jwt claims missing userID at claims['sub']")
	}

	rolesList, ok := claimsM["roles"].([]any)
	var roles []string
	if ok {
		for _, roleRaw := range rolesList {
			role := roleRaw.(map[string]any)
			roles = append(roles, role["key"].(string))
		}
	}

	return &ValidatedJWT{
		UserID: userId,
		Claims: claimsM,
		Roles:  roles,
	}, nil
}
