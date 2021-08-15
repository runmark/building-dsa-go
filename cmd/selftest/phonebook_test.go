package selftest

import (
	"testing"
)

type DummySearcher struct{}

func (ds DummySearcher) Search(people []*Person, firstName, lastName string) *Person {
	return &Person{}
}

func dummySearcherFind(t *testing.T) {
	p := Phonebook{}

	want := ErrMissingArgs
	_, got := p.Find(DummySearcher{}, "", "")

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

// ---------------------------------------------------

type StubSearcher struct {
	Phone string
}

func (ss StubSearcher) Search(people []*Person, firstName, lastName string) *Person {
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     ss.Phone,
	}
}

func stubSearcherFind(t *testing.T) {
	p := Phonebook{}

	stubPhone := "+86 123456"
	want := stubPhone
	got, err := p.Find(StubSearcher{stubPhone}, "Joe", "Done")
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}

}

// ---------------------------------------------------

type SpySearcher struct {
	Phone     string
	WasCalled bool
}

func (ss *SpySearcher) Search(people []*Person, firstName, lastName string) *Person {
	ss.WasCalled = true
	return &Person{
		firstName, lastName, ss.Phone,
	}
}

func spySearcherFind(t *testing.T) {
	p := Phonebook{}

	stubPhone := "+86 123456"
	want := stubPhone

	spy := SpySearcher{Phone: stubPhone}
	got, err := p.Find(&spy, "Joe", "Done")
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	if !spy.WasCalled {
		t.Errorf("expected to call Search in Find, but it wasn't!")
	}

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

// ---------------------------------------------------

type MockSearcher struct {
	phone     string
	methodsToCall map[string]bool
}

func (ms *MockSearcher) Search(people []*Person, firstName, lastName string) *Person {
	ms.methodsToCall["Search"] = true
	return &Person{
		firstName, lastName, ms.phone,
	}
}

func (ms *MockSearcher) ExpectToCall(methodName string) {
	if ms.methodsToCall == nil {
		ms.methodsToCall = make(map[string]bool)
	}
	ms.methodsToCall[methodName] = false
}

func (ms *MockSearcher) Verify(t *testing.T) {
	for methodName, called := range ms.methodsToCall {
		if !called {
			t.Errorf("Expected to call %v, but it wasn't.", methodName)
		}
	}
}

func mockSearcherFind(t *testing.T) {
	p := Phonebook{}

	stubPhone := "+86 123456"
	want := stubPhone

	mock := MockSearcher{
		phone: stubPhone,
	}

	mock.ExpectToCall("Search")

	got, _ := p.Find(&mock, "Joe", "Done")

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
	mock.Verify(t)
}

// ---------------------------------------------------

type FakeSearcher struct {}

func (fs FakeSearcher) Search(people []*Person, firstName, lastName string) *Person {
	if len(people) == 0 {
		return nil
	}
	return people[0]
}

func fakeSearcherFind(t *testing.T) {
	p := Phonebook{}
	fake := FakeSearcher{}

	got, _ := p.Find(fake, "Jaone", "Done")

	if got != "" {
		t.Errorf("wanted '', got '%s'", got)
	}
}

// ---------------------------------------------------

func TestFind(t *testing.T) {
	t.Run("dummySearcherFind", dummySearcherFind)
	t.Run("stubSearcherFind", stubSearcherFind)
	t.Run("spySearcherFind", spySearcherFind)
	t.Run("mockerSearcherFind", mockSearcherFind)
	t.Run("fakeSearcherFind", fakeSearcherFind)
}
