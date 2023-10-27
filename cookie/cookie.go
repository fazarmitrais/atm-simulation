package cookie

import "github.com/gorilla/sessions"

type Cookie struct {
	Store *sessions.CookieStore
}

func New() *Cookie {
	key := []byte("super-secret-key")
	store := sessions.NewCookieStore(key)
	return &Cookie{Store: store}
}
