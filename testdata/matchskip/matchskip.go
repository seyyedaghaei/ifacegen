package matchskip

type User struct{}

type UserService struct{}

func (s *UserService) CreateUser(name string) error {
	return nil
}

func (s *UserService) GetUser(id int) (*User, error) {
	return nil, nil
}

type UserRepository struct{}

func (r *UserRepository) Store(user *User) error {
	return nil
}

// ifacegen:skip
type SkipService struct{}

func (s *SkipService) Hidden() error {
	return nil
}

// ifacegen:generate
type Generated struct{}

func (s *Generated) DoSomething() error {
	return nil
}

type MethodSkipService struct{}

func (s *MethodSkipService) Visible() error {
	return nil
}

// ifacegen:skip
func (s *MethodSkipService) Hidden() error {
	return nil
}
