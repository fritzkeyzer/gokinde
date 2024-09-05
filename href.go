package gokinde

import "fmt"

// UserLoginHref is the link you should present to your users to direct
// them to your Kinde login page.
func (cl *Client) UserLoginHref(redirectURI string) string {
	authURL := cl.cfg.KindeDomain + OAuth2AuthPath +
		"?response_type=code" +
		"&client_id=" + cl.cfg.ClientID +
		"&redirect_uri=" + redirectURI +
		"&scope=openid+profile+email" +
		"&state=abcabcabcabcabcabcabcabcabcabcabcabcabcabc" // TODO: wtf is this? and why does it work with and not work without?

	return authURL
}

// UserLogoutHref is the link you should present to your users, for them to log our from their session.
func (cl *Client) UserLogoutHref(redirectURI string) string {
	return fmt.Sprintf("%s/logout?redirect=%s",
		cl.cfg.KindeDomain,
		redirectURI,
	)
}
