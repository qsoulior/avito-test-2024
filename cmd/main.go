package main

import (
	"os"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/app"
)

func main() {
	code := app.Run()
	os.Exit(code)
}
