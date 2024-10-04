package api

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Validate() bool {
	if len(u.Username) == 0 {
		return false
	}
	if len(u.Password) == 0 {
		return false
	}

	return true
}
