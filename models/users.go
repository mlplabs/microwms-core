package models

type UserService struct {
	Storage *Storage
	Reference
}

// User пользователь
type User struct {
	RefItem
}

// GetUsers возвращает список пользователей
func (ps *UserService) GetUsers(offset int, limit int) ([]User, int, error) {
	items, count, err := ps.getItems(ps.Storage.Db, offset, limit)
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
func (ps *UserService) FindUserById(usrId int64) (*User, error) {
	item, err := ps.findItemById(ps.Storage.Db, usrId)
	u := new(User)
	u.RefItem = *item
	return u, err
}

// FindUserByName возвращает пользователя по наименованию
func (ps *UserService) FindUserByName(valName string) ([]User, error) {
	items, err := ps.findItemByName(ps.Storage.Db, valName)
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

func (ps *UserService) GetSuggestionUser(text string, limit int) ([]string, error) {
	return ps.getSuggestion(ps.Storage.Db, text, limit)
}

// CreateUser создает нового пользователя
func (ps *UserService) CreateUser(u *User) (int64, error) {
	return ps.createItem(ps.Storage.Db, u)
}

// UpdateUser обновляет пользователя
func (ps *UserService) UpdateUser(u *User) (int64, error) {
	return ps.updateItem(ps.Storage.Db, u)
}

// DeleteUser удаляет пользователя
func (ps *UserService) DeleteUser(u *User) (int64, error) {
	return ps.deleteItem(ps.Storage.Db, u.Id)
}
