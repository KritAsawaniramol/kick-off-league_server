package util

import (
	"fmt"
	"strings"
)

type CreateTeamError struct {
	RequiredData []string
}

func (e *CreateTeamError) Error() string {
	return fmt.Sprintf("Please fill in the required information to create a team." + strings.Join(e.RequiredData, ","))
}
