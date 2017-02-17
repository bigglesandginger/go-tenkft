# go-tenkft

tenkft is a golang package that provides a wrapper around the awesome https://www.10000ft.com API.

All interactions with the tenkft API is done through the `*tenkft.Client` struct.

## Usage:

 ```go
import "github.com/workco/go-tenkft"

c, err := tenkft.NewClient("insert-your-token-here", tenkft.Staging) // or you can use tenkft.Production
handleErr(err)

projects, _, err := c.GetProjects(map[string]string{"fields": "tags,summmary"})
handleErr(err)

for _, project := range projects.Data {
  fmt.Println(project.Name)
}

if projects.Paging.HasNext() {
  nextPage := strconv.Itoa(projects.Paging.GetNextPage())
  nextProjects, \_, err := c.GetProjects(map[string]string{"page": nextPage})
  ...
}
```

 You can also use `MaxRetries` to automatically retry a request when the tenkft API
 returns an error.
