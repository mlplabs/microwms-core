package whs

// User warehouse user
type User struct {
	RefItem
}

// GetUsersItems returns a list of users
func (s *Storage) GetUsersItems(offset int, limit int, parentId int64) ([]User, int, error) {
	items, count, err := s.GetReference(tableRefUsers).getItems(offset, limit, parentId)
	if err != nil {
		return nil, 0, err
	}
	retVal := make([]User, len(items))
	for idx, item := range items {
		u := new(User)
		u.RefItem = item
		retVal[idx] = *u
	}
	return retVal, count, nil
}

// FindUserById finds a user by internal ID
func (s *Storage) FindUserById(usrId int64) (*User, error) {
	item, err := s.GetReference(tableRefUsers).findItemById(usrId)
	u := new(User)
	u.RefItem = *item
	return u, err
}

// FindUsersByName find users by name
func (s *Storage) FindUsersByName(valName string) ([]User, error) {
	items, err := s.GetReference(tableRefUsers).findItemByName(valName)
	if err != nil {
		return nil, err
	}
	retVal := make([]User, len(items))
	for _, item := range items {
		u := new(User)
		u.RefItem = item
		retVal = append(retVal, *u)
	}
	return retVal, err
}

func (s *Storage) GetUsersSuggestion(text string, limit int) ([]Suggestion, error) {
	return s.GetReference(tableRefUsers).getSuggestion(text, limit)
}

// CreateUser creates a new user
func (s *Storage) CreateUser(u *User) (int64, error) {
	return s.GetReference(tableRefUsers).createItem(u)
}

// UpdateUser updates the user
func (s *Storage) UpdateUser(u *User) (int64, error) {
	return s.GetReference(tableRefUsers).updateItem(u)
}

// DeleteUser deletes a user
func (s *Storage) DeleteUser(u *User) (int64, error) {
	return s.GetReference(tableRefUsers).deleteItem(u.Id)
}
