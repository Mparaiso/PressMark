package pressmark

import (
	"fmt"
	"regexp"
	"strings"
)

// UserValidator validates a user
type UserValidator struct {
}

// Validate validates a user
func (UserValidator) Validate(model interface{}) (errors map[string][]string, err error) {
	user := model.(*User)
	errors =map[string][]string{}
	if strings.Trim(user.Name, "") == "" {
		errors["Name"] = append(errors["Name"], "Name must be present.")
	}
	if strings.Trim(user.Email, "") == "" {
		errors["Email"] = append(errors["Email"], "Email must be present.")
	}
	if match, err := regexp.MatchString("^(\\w+?)@(\\w+)\\.(\\w+)$", user.Email); err != nil || !match {
		errors["Email"] = append(errors["Email"], "Email must be a valid email.")
	}
	if len(errors) > 0 {
		err = fmt.Errorf("User %v has validation errors", user)
	}
	
	return
}
