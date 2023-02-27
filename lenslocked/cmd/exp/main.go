package main

import (
	stdctx "context"
	"fmt"

	"github.com/sjadczak/webdev-go/lenslocked/context"
	"github.com/sjadczak/webdev-go/lenslocked/models"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := stdctx.Background()
	user := models.User{
		Email: "steve@jadczak.com",
	}

	ctx = context.WithUser(ctx, &user)
	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)
}
