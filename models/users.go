package models

type UserService struct {
	Storage *Storage
	Reference
}

// User пользователь
type User struct {
	RefItem
}

type ReferenceUsers struct {
	Reference
}

// GetUsers возвращает список пользователей
func (ref *ReferenceUsers) GetUsers(offset int, limit int) ([]User, int, error) {
	items, count, err := ref.getItems(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	retVal := make([]User, len(items))
	for _, item := range items {
		u := new(User)
		u.RefItem = item
		retVal = append(retVal, *u)
	}
	return retVal, count, nil
}

// FindUserById возвращает пользователя по внутреннему идентификатору
func (ref *ReferenceUsers) FindUserById(usrId int64) (*User, error) {
	item, err := ref.findItemById(usrId)
	u := new(User)
	u.RefItem = *item
	return u, err
}

// FindUserByName возвращает пользователя по наименованию
func (ref *ReferenceUsers) FindUserByName(valName string) ([]User, error) {
	items, err := ref.findItemByName(valName)
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

func (ref *ReferenceUsers) GetSuggestionUser(text string, limit int) ([]string, error) {
	return ref.getSuggestion(text, limit)
}

// CreateUser создает нового пользователя
func (ref *ReferenceUsers) CreateUser(u *User) (int64, error) {
	return ref.createItem(u)
}

// UpdateUser обновляет пользователя
func (ref *ReferenceUsers) UpdateUser(u *User) (int64, error) {
	return ref.updateItem(u)
}

// DeleteUser удаляет пользователя
func (ref *ReferenceUsers) DeleteUser(u *User) (int64, error) {
	return ref.deleteItem(u.Id)
}
