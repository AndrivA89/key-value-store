package entity

type Store struct {
	Values map[string]string
}

var values map[string]string

func init() {
	// TODO: load from DB or file
	values = make(map[string]string)
}

func GetStore() *Store {
	return &Store{Values: values}
}

func (s *Store) Put(key, value string) error {
	s.Values[key] = value

	return nil
}

func (s *Store) Get(key string) (string, error) {
	value, ok := s.Values[key]
	if !ok {
		return value, NotFoundError
	}
	return value, nil
}

func (s *Store) Delete(key string) error {
	delete(s.Values, key)

	return nil
}
