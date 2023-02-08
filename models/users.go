package models

// User пользователь
type User struct {
	RefItem
}

type ReferenceUsers struct {
	Reference
}

// GetItems возвращает список пользователей
func (ref *ReferenceUsers) GetItems(offset int, limit int) ([]User, int, error) {
	items, count, err := ref.getItems(offset, limit, 0)
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

// FindById возвращает пользователя по внутреннему идентификатору
func (ref *ReferenceUsers) FindById(usrId int64) (*User, error) {
	item, err := ref.findItemById(usrId)
	u := new(User)
	u.RefItem = *item
	return u, err
}

// FindByName возвращает пользователя по наименованию
func (ref *ReferenceUsers) FindByName(valName string) ([]User, error) {
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

func (ref *ReferenceUsers) GetSuggestion(text string, limit int) ([]string, error) {
	return ref.getSuggestion(text, limit)
}

// Create создает нового пользователя
func (ref *ReferenceUsers) Create(u *User) (int64, error) {
	return ref.createItem(u)
}

// Update обновляет пользователя
func (ref *ReferenceUsers) Update(u *User) (int64, error) {
	return ref.updateItem(u)
}

// Delete удаляет пользователя
func (ref *ReferenceUsers) Delete(u *User) (int64, error) {
	return ref.deleteItem(u.Id)
}
