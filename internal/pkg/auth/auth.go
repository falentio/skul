package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/falentio/raid-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"

	"github.com/falentio/skul/internal/pkg/response"
)

var (
	ErrUnauthorized = response.NewUnauthorized(nil, "session was not found")
	ErrInvalidToken = response.NewBadRequest(nil, "token is invalid")
)

var (
	jwtIDFactory = raid.NewRaid().WithPrefix("jwt")
)

type Auth struct {
	Name          string
	Secure        bool
	SigningMethod jwt.SigningMethod
	Secret        []byte
	Logger        zerolog.Logger

	ctxKey string
}

func (a *Auth) ctx() string {
	if a.ctxKey == "" {
		a.ctxKey = raid.NewRaid().WithRandom().WithTimestampNow().String()
	}
	return a.ctxKey
}

func (a *Auth) Sign(c jwt.RegisteredClaims) (*http.Cookie, error) {
	c.ID = jwtIDFactory.WithRandom().WithTimestampNow().String()
	// c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))

	token, err := jwt.
		NewWithClaims(a.SigningMethod, c).
		SignedString(a.Secret)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     a.Name,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Secure:   a.Secure,
		// MaxAge:   86400 * 7,
	}
	return cookie, nil
}

func (a *Auth) Verify(token string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	t, err := jwt.
		NewParser().
		ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return a.Secret, nil
		})
	if err != nil {
		a.Logger.Debug().Err(err).Msg("error while verifying jwt")
		return nil, ErrInvalidToken
	}

	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (a *Auth) VerifyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(a.Name)
		if err != nil {
			response.HandleError(w, r, ErrUnauthorized)
			return
		}

		claims, err := a.Verify(token.Value)
		if err != nil {
			response.HandleError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), a.ctx(), claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (a *Auth) GetClaims(ctx context.Context) (*jwt.RegisteredClaims, error) {
	c, ok := ctx.Value(a.ctx()).(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrUnauthorized
	}
	return c, nil
}

func (a *Auth) GetSubjectRaid(ctx context.Context, prefix string) (raid.Raid, error) {
	c, err := a.GetClaims(ctx)
	if err != nil {
		return raid.NilRaid, err
	}

	id, err := raid.RaidFromString(c.Subject)
	if err != nil {
		return raid.NilRaid, ErrInvalidToken
	}

	if prefix != "" && id.Prefix() != prefix {
		return raid.NilRaid, ErrInvalidToken
	}

	return id, nil
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    a.Name,
		MaxAge:  0,
		Expires: time.Unix(0, 0).UTC(),
	})

	w.WriteHeader(204)
}
