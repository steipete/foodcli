package cli

import (
	"context"

	"github.com/steipete/ordercli/internal/browserauth"
	"github.com/steipete/ordercli/internal/chromecookies"
	"github.com/steipete/ordercli/internal/foodora"
)

var chromeLoadCookieHeader = chromecookies.LoadCookieHeader

var browserOAuthTokenPassword = func(ctx context.Context, req foodora.OAuthPasswordRequest, opts browserauth.PasswordOptions) (foodora.AuthToken, *foodora.MfaChallenge, browserauth.Session, error) {
	return browserauth.OAuthTokenPassword(ctx, req, opts)
}
