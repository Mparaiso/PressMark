package pressmark

import (
	"time"
    "fmt"
)

// User is a user
type User struct {
	ID      int64
	Name    string
	Email   string
	Created time.Time
	Updated time.Time
}
func (user User) String() string {
	return fmt.Sprintf("{ID:%d,Name:%s}", user.ID, user.Name)
}
