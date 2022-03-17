package mydict

import "errors"

type Dictionary map[string]string

var (
	errNotFound   = errors.New("Not Found")
	errWordExists = errors.New("That aleady exists")
	errCantUpdate = errors.New("Cant update non-existing word")
	errCantDelete = errors.New("Cant Delete bcz not exists definition")
)

func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}

// Add a word to the dictionary
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists
	}
	return nil
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = definition
	case errNotFound:
		return errCantUpdate
	}
	return nil
}

func (d Dictionary) CustomDelete(word, definition string) error {
	value, err := d.Search(word)
	if err == errNotFound {
		return errNotFound

	} else if err == nil {
		if value == definition {
			d[word] = ""
		} else {
			return errCantDelete
		}
	}
	return nil
}
func (d Dictionary) Delete(word string) {
	delete(d, word)
}
