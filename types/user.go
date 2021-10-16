package types

type User struct {
	entry

	username []string
}

func (u *User) Parent() Node {
	return nil
}

func (m *Library) GetUser(username string) *User {
	m.Lock()
	defer m.Unlock()

	return m.getUser(username)
}

func (m *Library) getUser(username string) *User {
	if v := m.getNode(UserNode, username); v != nil {
		return v.(*User)
	}
	return nil
}
