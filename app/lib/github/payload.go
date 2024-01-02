package github

type Webhook struct {
	Ref string `json:"ref"`
	// Before     string     `json:"before"`
	// After      string     `json:"after"`
	Commits    []Commit   `json:"commits"`
	Pusher     Author     `json:"pusher"`
	Repository Repository `json:"repository"`
	Sender     *User      `json:"sender"`
}

type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Commit struct {
	Id        string   `json:"id"`
	Timestamp string   `json:"timestamp"`
	Message   string   `json:"message"`
	Author    Author   `json:"author"`
	Committer Author   `json:"committer"`
	Url       string   `json:"url"`
	Distinct  bool     `json:"distinct"`
	Added     []string `json:"added"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
}

type Repository struct {
	Id              int    `json:"id"`
	NodeID          string `json:"node_id"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	HtmlUrl         string `json:"html_url"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	Language        string `json:"language"`
	Private         bool   `json:"private"`
	Owner           User   `json:"owner"`
}

type User struct {
	Login   string `json:"login"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Id      int    `json:"id"`
	HtmlUrl string `json:"html_url"`
}
