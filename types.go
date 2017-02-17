package tenkft

// Projects a collection of project - emulates /projects
type Projects struct {
	Data   []*Project `json:"data"`
	Paging *Paging    `json:"paging"`
}

// GetByID get a project from collection by id
func (ps *Projects) GetByID(id int) (targetProject *Project) {
	for _, p := range ps.Data {
		if p.ID == id {
			targetProject = p
			return
		}
	}

	return
}

// Find finds a person based on a callback that returns a boolean
func (ps *Projects) Find(cb func(*Project) bool) (p *Project) {
	for _, project := range ps.Data {
		if cb(project) {
			p = project
			return
		}
	}

	return
}

type baseProject struct {
	Archived     bool   `json:"archived,omitempty"`
	Name         string `json:"name,omitempty"`
	EndsAt       string `json:"ends_at,omitempty"`
	StartsAt     string `json:"starts_at,omitempty"`
	Description  string `json:"description,omitempty"`
	Client       string `json:"client,omitempty"`
	ProjectState string `json:"project_state,omitempty"`
	PhaseName    string `json:"phase_name,omitempty,omitempty"`
	ProjectCode  string `json:"project_code,omitempty,omitempty"`
}

// Project abstraction to the /project schema
type Project struct {
	*baseProject
	ID                  int         `json:"id"`
	ArchivedAt          string      `json:"archived_at"`
	GUID                string      `json:"guid"`
	ParentID            int         `json:"parent_id"`
	SecureURL           string      `json:"secureurl"`
	SecureURLExpiration string      `json:"secureurl_expiration"`
	Settings            interface{} `json:"settings"`
	TimeentryLockout    interface{} `json:"timeentry_lockout"`
	DeletedAt           string      `json:"deleted_at"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	UseParentBillRates  bool        `json:"use_parent_bill_rates"`
	Thumbnail           string      `json:"thumbnail"`
	Type                string      `json:"type"`
	HasPendingUpdates   bool        `json:"has_pending_updates"`
	Tags                Tags        `json:"tags"`
	Assignments         Assignments `json:"assignments"`
	BoundingStartdate   string      `json:"bounding_startdate"`
	BoundingEnddate     string      `json:"bounding_enddate"`
	ConfirmedHours      float64     `json:"confirmed_hours"`
	ConfirmedDollars    float64     `json:"confirmed_dollars"`
	ApprovedHours       float64     `json:"approved_hours"`
	ApprovedDollars     float64     `json:"approved_dollars"`
	UnconfirmedHours    float64     `json:"unconfirmed_hours"`
	UnconfirmedDollars  float64     `json:"unconfirmed_dollars"`
	ScheduledHours      float64     `json:"scheduled_hours"`
	ScheduledDollars    float64     `json:"scheduled_dollars"`
	FutureHours         float64     `json:"future_hours"`
	FutureDollars       float64     `json:"future_dollars"`
}

type baseUser struct {
	Archived          bool        `json:"archived,omitempty"`
	Discipline        string      `json:"discipline"`
	Email             string      `json:"email"`
	FirstName         string      `json:"first_name"`
	HireDate          interface{} `json:"hire_date"`
	LastName          string      `json:"last_name"`
	Location          string      `json:"location"`
	MobilePhone       interface{} `json:"mobile_phone"`
	Role              string      `json:"role"`
	BillabilityTarget float64     `json:"billability_target"`
}

// User abstraction to the /user schema
type User struct {
	*baseUser
	AccountOwner      bool        `json:"account_owner"`
	ArchivedAt        string      `json:"archived_at"`
	Billable          bool        `json:"billable"`
	Billrate          float64     `json:"billrate"`
	CreatedAt         string      `json:"created_at"`
	Deleted           bool        `json:"deleted"`
	DeletedAt         string      `json:"deleted_at"`
	DisplayName       string      `json:"display_name"`
	EmployeeNumber    interface{} `json:"employee_number"`
	GUID              string      `json:"guid"`
	HasLogin          bool        `json:"has_login"`
	ID                int         `json:"id"`
	InvitationPending bool        `json:"invitation_pending"`
	LoginType         string      `json:"login_type"`
	OfficePhone       interface{} `json:"office_phone"`
	TerminationDate   string      `json:"termination_date"`
	Thumbnail         interface{} `json:"thumbnail"`
	Type              string      `json:"type"`
	UserSettings      float64     `json:"user_settings"`
	UserTypeID        int         `json:"user_type_id"`
	Tags              Tags        `json:"tags"`
	Assignments       Assignments `json:"assignments"`
}

// // MarshalJSON interface implementation
// func (u *User) MarshalJSON() (bytes []byte, err error) {
// 	newUser := &NewUser{}
// 	newUser.FirstName = u.FirstName
// 	newUser.LastName = u.LastName
// 	bytes, err = json.Marshal(newUser)
// 	return
// }

// Tags holds a collection of tags - only reachable from a user or project.
type Tags struct {
	Data   []*Tag  `json:"data"`
	Paging *Paging `json:"paging"`
}

type baseTag struct {
	Value string `json:"value"`
}

// Tag holds a tag - only reachable from a user or a project.
type Tag struct {
	*baseTag
	ID int `json:"id"`
}

// Users holds a collection of users and also indicates whether paginating is available.
type Users struct {
	Data   []*User `json:"data"`
	Paging *Paging `json:"paging"`
}

// GetNonOwnerCount returns the number of users who are not account owners
func (users *Users) GetNonOwnerCount() int {
	var count int
	for _, u := range users.Data {
		if u.AccountOwner == false {
			count++
		}
	}

	return count
}

// Paging abstracts paging parameters
type Paging struct {
	PerPage  int    `json:"per_page"`
	Page     int    `json:"page"`
	Previous string `json:"previous"`
	Self     string `json:"self"`
	Next     string `json:"next"`
}

// HasNext confirms whether there is a next pagination page.
func (p *Paging) HasNext() bool {
	return p.Next != "null" && p.Next != ""
}

// GetNextPage returns next page in pagination
func (p *Paging) GetNextPage() int {
	return p.Page + 1
}

// Assignments abstraction to /assignments schema
type Assignments struct {
	Data   []*Assignment `json:"data"`
	Paging *Paging       `json:"paging"`
}

type baseAssignment struct {
	AllocationMode string  `json:"allocation_mode"`
	AssignableID   int     `json:"assignable_id"`
	EndsAt         string  `json:"ends_at"`
	FixedHours     float64 `json:"fixed_hours,omitempty"`
	HoursPerDay    float64 `json:"hours_per_day,omitempty"`
	Percent        float64 `json:"percent,omitempty"`
	StartsAt       string  `json:"starts_at"`
}

// Assignment an abstraction to an assignment schema
type Assignment struct {
	*baseAssignment
	AllDayAssignment  bool    `json:"all_day_assignment"`
	BillRate          float64 `json:"bill_rate"`
	BillRateID        int     `json:"bill_rate_id"`
	CreatedAt         string  `json:"created_at"`
	ID                int     `json:"id"`
	RepetitionID      int     `json:"repetition_id"`
	ResourceRequestID int     `json:"resource_request_id"`
	Status            string  `json:"status"`
	UpdatedAt         string  `json:"updated_at"`
	UserID            int     `json:"user_id"`
}

// Phases abstraction to project phases schema
type Phases struct {
	Data   []*Phase `json:"data"`
	Paging *Paging  `json:"paging"`
}

type basePhase struct {
	Archived  bool   `json:"archived,omitempty"`
	PhaseName string `json:"phase_name"`
	EndsAt    string `json:"ends_at"`
	StartsAt  string `json:"starts_at"`
}

// Phase abstraction to a project phase object
type Phase struct {
	*basePhase
	ID                  int         `json:"id"`
	ArchivedAt          string      `json:"archived_at"`
	Description         string      `json:"description"`
	GUID                string      `json:"guid"`
	Name                string      `json:"name"`
	ParentID            int         `json:"parent_id"`
	ProjectCode         string      `json:"project_code"`
	SecureURL           string      `json:"secureurl"`
	SecureURLExpiration string      `json:"secureurl_expiration"`
	Settings            interface{} `json:"settings"`
	TimeentryLockout    interface{} `json:"timeentry_lockout"`
	DeletedAt           string      `json:"deleted_at"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	UseParentBillRates  bool        `json:"use_parent_bill_rates"`
	Thumbnail           string      `json:"thumbnail"`
	Type                string      `json:"type"`
	HasPendingUpdates   bool        `json:"has_pending_updates"`
	Client              string      `json:"client"`
	ProjectState        string      `json:"project_state"`
}

// PlaceholderResources abstraction to /placeholder_resources
type PlaceholderResources struct {
	Data   []*PlaceholderResource `json:"data"`
	Paging *Paging                `json:"paging"`
}

// PlaceholderResource abstraction to a PlaceholderResource object.
type PlaceholderResource struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	UserTypeID   int     `json:"user_type_id"`
	GUID         string  `json:"guid"`
	Role         string  `json:"role"`
	Discipline   string  `json:"discipline"`
	Location     string  `json:"location"`
	CreatedAt    string  `json:"created_at"`
	Billrate     float64 `json:"billrate"`
	DisplayName  string  `json:"displayName"`
	Type         string  `json:"type"`
	Thumbnail    string  `json:"thumbnail"`
	Abbreviation string  `json:"abbreviation"`
	Color        string  `json:"color"`
}

// LeaveTypes abstraction to /leave_types response collection
type LeaveTypes struct {
	Data   []*LeaveType `json:"data"`
	Paging *Paging      `json:"paging"`
}

// FindByName finds a *LeaveType by its name
func (lts *LeaveTypes) FindByName(name string) (lt *LeaveType) {
	for _, lt = range lts.Data {
		if lt.Name == name {
			return
		}
	}

	return
}

// LeaveType abstraction to LeaveType object
type LeaveType struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	GUID        string `json:"guid"`
	Name        string `json:"name"`
	DeletedAt   string `json:"deleted_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Type        string `json:"type"`
}

// Roles abstraction to /roles schema
type Roles struct {
	Data   []*Role `json:"data"`
	Paging *Paging `json:"paging"`
}

// Role abstraction to a role object
type Role struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// BillRates abstraction to /roles schema
type BillRates struct {
	Data   []*BillRate `json:"data"`
	Paging *Paging     `json:"paging"`
}

// BillRate abstraction to a role object
type BillRate struct {
	ID           int     `json:"id"`
	Rate         float64 `json:"rate"`
	AssignableID int     `json:"assignable_id"`
	DisciplineID int     `json:"discipline_id"`
	RoleID       int     `json:"role_id"`
	UserID       int     `json:"user_id"`
	StartsAt     string  `json:"starts_at"`
	EndsAt       string  `json:"ends_at"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Startdate    string  `json:"startdate"`
	Enddate      string  `json:"enddate"`
}

// // Time extension
// type Time struct {
// 	time.Time
// }
//
// // UnmarshalJSON interface
// func (t *Time) UnmarshalJSON(b []byte) (err error) {
// 	str := string(b)
// 	if str == "null" || str == "" {
// 		return
// 	}
//
// 	err = t.Time.UnmarshalJSON(b)
// 	return
// }
