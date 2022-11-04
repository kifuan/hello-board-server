package models

import (
	"os"
)

var (
	DSN          = os.Getenv("DSN")
	GRAVATAR     = os.Getenv("GRAVATAR")
	ADMIN_EMAIL  = os.Getenv("ADMIN_EMAIL")
	ADMIN_SECRET = os.Getenv("ADMIN_SECRET")
)
