package mypackage

import (
    "github.com/gorilla/sessions"
)

func cookieStore() *sessions.CookieStore {
    SecretKey := []byte("your-secret-key")
    cookieStore := sessions.NewCookieStore(SecretKey)
 
    return cookieStore
}