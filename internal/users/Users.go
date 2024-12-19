package users

type Pengguna struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	RoleID    int    `json:"role_id"`
	ClientID  int    `json:"client_id"`
	CreatedAt string `json:"created_at"`
}
