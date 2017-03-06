// Package tenkft provides a wrapper around the awesome https://www.10000ft.com API.
// All interactions with the tenkft API are done through the *Client struct.
// Usage:
//  import "github.com/workco/go-tenkft"
//
//  c, err := tenkft.NewClient("insert-your-token-here", tenkft.Staging) // or you can use tenkft.Production
//  handleErr(err)
//
//  projects, _, err := c.GetProjects(map[string]string{"fields": "tags,summmary"})
//  handleErr(err)
//
//  for _, project := range projects.Data {
//    fmt.Println(project.Name)
//  }
//
//  if projects.Paging.HasNext() {
//    nextPage := strconv.Itoa(projects.Paging.GetNextPage())
//    nextProjects, _, err := c.GetProjects(map[string]string{"page": nextPage})
//    ...
//  }
//
// You can also use MaxRetries to automatically retry a request when the tenkft API
// returns an error.
package tenkft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/workco/go-tenkft/utils"
)

const (
	// Production environment URL
	Production = "https://api.10000ft.com/api/v1"
	// Staging environment URL
	Staging = "https://vnext.10000ft.com/api/v1"
)

// Client use NewClient to return this instance type.
type Client struct {
	token      string
	env        string
	MaxRetries int
}

// NewClient takes credentials and returns client to perform API operations on
func NewClient(token, env string) (*Client, error) {
	if env != Production && env != Staging {
		return &Client{}, fmt.Errorf("env must be either %v, or %v", Production, Staging)
	}

	c := &Client{token: token, env: env}

	return c, nil
}

func queryfy(opts map[string]string) string {
	querySlice := []string{}
	for k, val := range opts {
		querySlice = append(querySlice, k+"="+val)
	}

	return strings.Join(querySlice, "&")
}

// GetAllProjects returns all projects - automatically paginates and returns accumulated projects.
// resp and err correspond to the latest one in the loop.
func (c *Client) GetAllProjects(opts map[string]string) (projects *Projects, resp *http.Response, err error) {
	projects = &Projects{Paging: &Paging{}}
	opts["per_page"] = "201"
	projects, resp, err = c.GetProjects(opts)
	if err != nil {
		return
	}

	for loop := projects.Paging.HasNext(); loop == true; loop = projects.Paging.HasNext() {
		opts["page"] = strconv.Itoa(projects.Paging.GetNextPage())
		newProjects, newResp, newErr := c.GetProjects(opts)
		resp = newResp
		if err != nil {
			err = newErr
			break
		}

		projects.Paging = newProjects.Paging
		projects.Data = append(projects.Data, newProjects.Data...)
	}

	return
}

// GetProjects returns all projects with default pagination
func (c *Client) GetProjects(opts map[string]string) (projects *Projects, resp *http.Response, err error) {
	projects = &Projects{Paging: &Paging{}}
	query := queryfy(opts)
	url, method, headers := c.env+"/projects?"+query, http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, projects)
	if err != nil {
		return
	}

	return
}

// GetUsers returns all users - manual pagination per opts paramater
// URL https://github.com/10Kft/10kft-api/blob/master/sections/users.md#endpoint-apiv1users
func (c *Client) GetUsers(opts map[string]string) (users *Users, resp *http.Response, err error) {
	users = &Users{Paging: &Paging{}}
	query := queryfy(opts)
	url, method, headers := c.env+"/users?"+query, http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, users)
	if err != nil {
		return
	}

	return
}

// GetUser returns a user based on a user object's ID
func (c *Client) GetUser(u *User, opts map[string]string) (resp *http.Response, err error) {
	query := queryfy(opts)
	url := c.env + "/users/" + strconv.Itoa(u.ID) + "?" + query
	method, headers := http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, u)
	if err != nil {
		return
	}

	return
}

// GetAllUsers returns all users - automatically paginates and returns the accumulated collection.
// resp and err correspond to the latest one in the loop.
// URL https://github.com/10Kft/10kft-api/blob/master/sections/users.md#endpoint-apiv1users
func (c *Client) GetAllUsers(opts map[string]string) (users *Users, resp *http.Response, err error) {
	users = &Users{Paging: &Paging{}}
	opts["per_page"] = "201"
	users, resp, err = c.GetUsers(opts)
	if err != nil {
		return
	}

	for loop := users.Paging.HasNext(); loop == true; loop = users.Paging.HasNext() {
		opts["page"] = strconv.Itoa(users.Paging.GetNextPage())
		newUsers, newResp, newErr := c.GetUsers(opts)
		resp = newResp
		if err != nil {
			err = newErr
			break
		}

		users.Paging = newUsers.Paging
		users.Data = append(users.Data, newUsers.Data...)
	}

	return
}

// CreateUser abstraction to POST /users
func (c *Client) CreateUser(u *User) (resp *http.Response, err error) {
	url, method, headers := c.env+"/users", http.MethodPost, map[string]string{"auth": c.token}

	body, err := json.Marshal(u.baseUser)
	if err != nil {
		return
	}

	fetcher, err := utils.NewFetchOpts(url, method, string(body), headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, u)
	if err != nil {
		return
	}

	return
}

// DeleteUser archives user by updating it with archived set to true
func (c *Client) DeleteUser(u *User) (*http.Response, error) {
	u.Archived = true
	return c.UpdateUser(u)
}

// UpdateUser abstraction to PUT /users/<id>
func (c *Client) UpdateUser(u *User) (resp *http.Response, err error) {
	url, method, headers := c.env+"/users/"+strconv.Itoa(u.ID), http.MethodPut, map[string]string{"auth": c.token}

	body, err := json.Marshal(u.baseUser)
	if err != nil {
		return
	}

	fetcher, err := utils.NewFetchOpts(url, method, string(body), headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, u)
	return
}

// CreateProject abstraction to POST /projects
func (c *Client) CreateProject(p *Project) (resp *http.Response, err error) {
	url, method, headers := c.env+"/projects", http.MethodPost, map[string]string{"auth": c.token}
	body, err := json.Marshal(p.baseProject)
	if err != nil {
		return
	}

	fetcher, err := utils.NewFetchOpts(url, method, string(body), headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, p)
	if err != nil {
		return
	}

	return
}

// DeleteProject calls UpdateProject with archive set to true
func (c *Client) DeleteProject(p *Project) (*http.Response, error) {
	p.baseProject = &baseProject{Archived: true}

	return c.UpdateProject(p)
}

// UpdateProject abstraction to PUT /projects/<id>
func (c *Client) UpdateProject(p *Project) (resp *http.Response, err error) {
	url := c.env + "/projects/" + strconv.Itoa(p.ID)
	method, headers := http.MethodPut, map[string]string{"auth": c.token}

	body, err := json.Marshal(p.baseProject)
	if err != nil {
		return
	}

	fetcher, err := utils.NewFetchOpts(url, method, string(body), headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, p)

	return
}

// GetAllUserAssignments - paginates through all assinments
func (c *Client) GetAllUserAssignments(u *User, opts map[string]string) (assignments *Assignments, resp *http.Response, err error) {
	opts["per_page"] = "250"
	assignments, resp, err = c.GetUserAssignments(u, opts)
	if err != nil {
		return
	}

	for loop := assignments.Paging.HasNext(); loop == true; loop = assignments.Paging.HasNext() {
		opts["page"] = strconv.Itoa(assignments.Paging.GetNextPage())
		newAssignments, newResp, newErr := c.GetUserAssignments(u, opts)
		resp = newResp
		if err != nil {
			err = newErr
			break
		}

		assignments.Paging = newAssignments.Paging
		assignments.Data = append(assignments.Data, newAssignments.Data...)
	}

	return
}

// GetUserAssignments retrieves all assignments for a user
// https://github.com/10Kft/10kft-api/blob/master/sections/assignments.md#endpoint-apiv1usersuser_idassignments
func (c *Client) GetUserAssignments(u *User, opts map[string]string) (assignments *Assignments, resp *http.Response, err error) {
	assignments = &Assignments{}
	query := queryfy(opts)
	url := c.env + "/users/" + strconv.Itoa(u.ID) + "/assignments?" + query
	method := http.MethodGet
	headers := map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, assignments)

	return
}

// GetProjectAssignments retrieves all assignments for a project
func (c *Client) GetProjectAssignments(p *Project, opts map[string]string) (assignments Assignments, resp *http.Response, err error) {
	query := queryfy(opts)
	url := c.env + "/projects/" + strconv.Itoa(p.ID) + "/assignments?" + query
	method := http.MethodGet
	headers := map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, assignments)

	return
}

// CreateUserAssignment abstraction to POST /users/<id>/assignments
func (c *Client) CreateUserAssignment(a *Assignment) (resp *http.Response, err error) {
	url := c.env + "/users/" + strconv.Itoa(a.UserID) + "/assignments"
	method, headers := http.MethodPost, map[string]string{"auth": c.token}

	body, err := json.Marshal(a.baseAssignment)
	if err != nil {
		return
	}

	fetcher, err := utils.NewFetchOpts(url, method, string(body), headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, a)

	return
}

// GetProjectPhases abstraction to GET /projects/<id>/phases
func (c *Client) GetProjectPhases(p *Project, opts map[string]string) (phases *Phases, resp *http.Response, err error) {
	phases = &Phases{}
	query := queryfy(opts)
	url := c.env + "/projects/" + strconv.Itoa(p.ID) + "/phases?" + query
	method, headers := http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, phases)
	if err != nil {
		return
	}

	return
}

// GetProjectByID abstraction to GET /projects/<id>
func (c *Client) GetProjectByID(ID int, opts map[string]string) (p *Project, resp *http.Response, err error) {
	p = &Project{}
	query := queryfy(opts)
	url := c.env + "/projects/" + strconv.Itoa(ID) + "?" + query
	method, headers := http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, p)
	return
}

// CreateProjectPhase abstraction to POST /projects/<id>/phases
func (c *Client) CreateProjectPhase(pID int, ph *Phase) (resp *http.Response, err error) {
	url := c.env + "/projects/" + strconv.Itoa(pID) + "/phases"
	method, headers := http.MethodPost, map[string]string{"auth": c.token}
	body, err := json.Marshal(ph.basePhase)
	if err != nil {
		return
	}

	fetcher, err := utils.NewFetchOpts(url, method, string(body), headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, ph)

	return
}

// CreateUserTags abstraction to POST /useres/<id>/tags
func (c *Client) CreateUserTags(u *User) (resp *http.Response, err error) {
	url := c.env + "/users/" + strconv.Itoa(u.ID) + "/tags"
	method := http.MethodPost
	headers := map[string]string{"auth": c.token}

	for _, t := range u.Tags.Data {
		body, err := json.Marshal(t.baseTag)
		if err != nil {
			return resp, err
		}

		fetcher, err := utils.NewFetchOpts(url, method, string(body), headers, c.MaxRetries)
		if err != nil {
			return resp, err
		}

		resp, err = fetcher.Fetch()
		if err != nil {
			return resp, err
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		err = json.Unmarshal(b, t)
		if err != nil {
			return resp, err
		}
	}

	return
}

// CreateProjectTags abstraction to POST /projects/<id>/tags for each project tag.
func (c *Client) CreateProjectTags(p *Project) (resp *http.Response, err error) {
	url := c.env + "/projects/" + strconv.Itoa(p.ID) + "/tags"
	method := http.MethodPost
	headers := map[string]string{"auth": c.token}

	for _, t := range p.Tags.Data {
		body, err := json.Marshal(t.baseTag)
		if err != nil {
			return resp, err
		}

		fetcher, err := utils.NewFetchOpts(url, method, string(body), headers, c.MaxRetries)
		if err != nil {
			return resp, err
		}

		resp, err = fetcher.Fetch()
		if err != nil {
			return resp, err
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		err = json.Unmarshal(b, t)
		if err != nil {
			return resp, err
		}
	}

	return
}

// GetLeaveTypes abstraction to GET /leave_types
func (c *Client) GetLeaveTypes(opts map[string]string) (leaveTypes *LeaveTypes, resp *http.Response, err error) {
	leaveTypes = &LeaveTypes{}
	query := queryfy(opts)
	url, method, headers := c.env+"/leave_types?"+query, http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, leaveTypes)
	if err != nil {
		return
	}

	return
}

// GetAllLeaveTypes returns all leave types - automatically paginates and returns accumulated leave types.
// resp and err correspond to the latest one in the loop.
func (c *Client) GetAllLeaveTypes(opts map[string]string) (leaveTypes *LeaveTypes, resp *http.Response, err error) {
	opts["per_page"] = "50"
	leaveTypes, resp, err = c.GetLeaveTypes(opts)
	if err != nil {
		return
	}

	for loop := leaveTypes.Paging.HasNext(); loop == true; loop = leaveTypes.Paging.HasNext() {
		opts["page"] = strconv.Itoa(leaveTypes.Paging.GetNextPage())
		newLeaveTypes, newResp, newErr := c.GetLeaveTypes(opts)
		resp = newResp
		if err != nil {
			err = newErr
			break
		}

		leaveTypes.Paging = newLeaveTypes.Paging
		leaveTypes.Data = append(leaveTypes.Data, newLeaveTypes.Data...)
	}

	return
}

// GetRoles returns all Role types for an account.
func (c *Client) GetRoles(opts map[string]string) (roles *Roles, resp *http.Response, err error) {
	roles = &Roles{}
	query := queryfy(opts)
	url, method, headers := c.env+"/roles?"+query, http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, roles)

	return
}

// GetAllRoles returns all role types - automatically paginates and returns accumulated roles
// resp and err correspond to the latest one in the loop.
func (c *Client) GetAllRoles(opts map[string]string) (roles *Roles, resp *http.Response, err error) {
	opts["per_page"] = "50"
	roles, resp, err = c.GetRoles(opts)
	if err != nil {
		return
	}

	for loop := roles.Paging.HasNext(); loop == true; loop = roles.Paging.HasNext() {
		opts["page"] = strconv.Itoa(roles.Paging.GetNextPage())
		newRoles, newResp, newErr := c.GetRoles(opts)
		resp = newResp
		if err != nil {
			err = newErr
			break
		}

		roles.Paging = newRoles.Paging
		roles.Data = append(roles.Data, newRoles.Data...)
	}

	return
}

// GetProjectBillRates returns all bill rates for a project.
func (c *Client) GetProjectBillRates(pID int, opts map[string]string) (billRates *BillRates, resp *http.Response, err error) {
	billRates = &BillRates{}
	query := queryfy(opts)
	url := c.env + "/projects/" + strconv.Itoa(pID) + "/bill_rates?" + query
	method, headers := http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, billRates)

	return
}

// GetAllProjectBillRates returns all project bill rates - automatically paginates and returns accumulated response
// resp and err correspond to the latest one in the loop.
func (c *Client) GetAllProjectBillRates(pID int, opts map[string]string) (billRates *BillRates, resp *http.Response, err error) {
	opts["per_page"] = "50"
	billRates, resp, err = c.GetProjectBillRates(pID, opts)
	if err != nil {
		return
	}

	for loop := billRates.Paging.HasNext(); loop == true; loop = billRates.Paging.HasNext() {
		opts["page"] = strconv.Itoa(billRates.Paging.GetNextPage())
		newBillRates, newResp, newErr := c.GetProjectBillRates(pID, opts)
		resp = newResp
		if err != nil {
			err = newErr
			break
		}

		billRates.Paging = newBillRates.Paging
		billRates.Data = append(billRates.Data, newBillRates.Data...)
	}

	return
}

// GetProjectUsers returns a project's users /projects/<id>/users
func (c *Client) GetProjectUsers(pID int, opts map[string]string) (users *Users, resp *http.Response, err error) {
	users = &Users{}
	query := queryfy(opts)
	url := c.env + "/projects/" + strconv.Itoa(pID) + "/users?" + query
	method, headers := http.MethodGet, map[string]string{"auth": c.token}

	fetcher, err := utils.NewFetchOpts(url, method, "", headers, c.MaxRetries)
	if err != nil {
		return
	}

	resp, err = fetcher.Fetch()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, users)

	return
}
