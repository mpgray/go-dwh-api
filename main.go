package main

import (
	"go-dwh-api/app"
	u "go-dwh-api/utils"
)

func main() {

	app.Serve()

	defer u.Log.Infof("**Golang Backend API for Driveway Home Started Successfully**")
}
