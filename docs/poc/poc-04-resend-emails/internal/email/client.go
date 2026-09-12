package email

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v4"
)

func NewClient() (*resend.Client, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Missing RESEND_API_KEY variable.")
	}

	return resend.NewClient(apiKey), nil
}
