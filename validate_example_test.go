package gokinde_test

import (
	"context"
	"log"
	"net/http"

	"github.com/fritzkeyzer/gokinde"
)

func ExampleClient_ValidateJWT() {
	// ⌄⌄⌄ server startup code: ⌄⌄⌄
	ctx := context.Background()
	kindeCl := gokinde.MustNewClient(ctx, gokinde.Cfg{
		ClientID:     "kinde-secret",
		ClientSecret: "kinde-secret",
		KindeDomain:  "https://your-app.kinde.com",
		ErrorLog: func(err error) {
			log.Println("ERR: kinde:", err)
		},
	})
	// ^^^ end ^^^

	// ⌄⌄⌄ auth handler code: ⌄⌄⌄
	validateRequest := func(r *http.Request) *gokinde.ValidatedJWT {
		cookie, _ := r.Cookie("kinde_access_token")
		if cookie == nil {
			return nil
		}

		valid, _ := kindeCl.ValidateJWT(cookie.Value)
		if valid == nil {
			return nil
		}

		log.Println(valid.UserID)
		log.Println(valid.Roles)
		log.Println(valid.Claims["custom_claim_key"])

		return valid
	}
	_ = validateRequest
	// ^^^ end ^^^
}
