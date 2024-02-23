package database_package

import (
    "github.com/gorilla/sessions"
)

func CookieStore() *sessions.CookieStore {
    SecretKey := []byte("your-secret-key")
    cookieStore := sessions.NewCookieStore(SecretKey)
 
    return cookieStore
}