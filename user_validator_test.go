package pressmark_test

import (
	"github.com/interactiv/expect"
	"github.com/mparaiso/PressMark"
	"strings"
	"testing"
)

func TestUserValidator(t *testing.T) {
	e := expect.New(t)

	user := &pressmark.User{}
	userValidator := pressmark.UserValidator{}
	errors, err := userValidator.Validate(user)
	t.Log(errors)

	e.Expect(err).Not().ToBeNil()
	e.Expect(len(errors)).ToBe(2)

	user.Email = "john@acme.com"
	errors, err = userValidator.Validate(user)
	t.Log(errors)
	e.Expect(len(errors)).ToBe(1)
	user.Email = "john.com"
	errors, err = userValidator.Validate(user)
	e.Expect(len(errors)).ToBe(2)
	user.Email = "john@bar"
	errors, err = userValidator.Validate(user)
	e.Expect(len(errors)).ToBe(2)
	user = &pressmark.User{Name: "John Doe", Email: "john.doe@acme.com"}
	errors, err = userValidator.Validate(user)
	t.Log(errors)
	e.Expect(len(errors)).ToBe(0)
	user.Name = strings.Repeat("a", 60)
	errors, err = userValidator.Validate(user)
	e.Expect(len(errors)).ToBe(1)
	user = &pressmark.User{Name: "John Doe", Email: strings.Repeat("a", 300) + "@acme.com"}
	errors, err = userValidator.Validate(user)
	e.Expect(len(errors)).ToBe(1)
    t.Log(errors)
    	user = &pressmark.User{Name: "John Doe", Email: "john.doe@acme.co.uk"}
	errors, err = userValidator.Validate(user)
	e.Expect(len(errors)).ToBe(0)
}
