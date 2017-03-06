package tenkft

import (
	"fmt"
	"os"
	"testing"
)

var c, _ = NewClient(os.Getenv("TEN_K_DEV"), Staging)
var projects = &Projects{}

func TestGetProjects(t *testing.T) {
	var err error
	projects, _, err = c.GetProjects(map[string]string{})
	if err != nil {
		t.Fatal("could not get projects", err)
	}
}

func TestGetProjectUsers(t *testing.T) {
	if len(projects.Data) == 0 {
		fmt.Println("There are no projects to test GetProjectUsers against")
		t.SkipNow()
	}

	p := projects.Data[0]
	_, _, err := c.GetProjectUsers(p.ID, map[string]string{})
	if err != nil {
		t.Errorf("could not get project users: %v", err.Error())
	}
}

func TestGetAllProjects(t *testing.T) {
	projects, _, err := c.GetAllProjects(map[string]string{})
	if err != nil {
		t.Errorf("could not get all projects: %v", err.Error())
	}

	if len(projects.Data) == 0 {
		fmt.Println("all projects returned an empty slice")
	}
}

func TestGetAllUsers(t *testing.T) {
	users, _, err := c.GetAllUsers(map[string]string{})
	if err != nil {
		t.Errorf("Could not get all users: %v", err.Error())
	}

	if len(users.Data) == 0 {
		fmt.Println("all projects returned an empty slice")
	}
}
