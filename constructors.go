package tenkft

// all constructors are here.

// NewProjects - initializes a Projects struct with non nil fields.
func NewProjects() *Projects {
	return &Projects{Paging: &Paging{}, Data: []*Project{}}
}

// NewProject - initializes a Project struct with non nil fields.
func NewProject() *Project {
	return &Project{baseProject: &baseProject{}}
}

// NewUsers - initializes a Users struct with non nil fields.
func NewUsers() *Users {
	return &Users{Paging: &Paging{}, Data: []*User{}}
}

// NewUser - initializes a User struct with non nil fields.
func NewUser() *User {
	return &User{baseUser: &baseUser{}}
}
