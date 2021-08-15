package selftest

import "errors"

var (
	ErrMissingArgs    = errors.New("missing arguments")
	ErrPersonNotFound = errors.New("person not found")
)

type Searcher interface {
	Search(people []*Person, firstName, lastName string) *Person
}

type Person struct {
	FirstName string
	LastName  string
	Phone     string
}

type Phonebook struct {
	People []*Person
}

func (p *Phonebook) Find(searcher Searcher, firstName, lastName string) (string, error) {
	if firstName == "" || lastName == "" {
		return "", ErrMissingArgs
	}

	s := searcher.Search(p.People, firstName, lastName)
	if s == nil {
		return "", ErrPersonNotFound
	}

	return s.Phone, nil
}
