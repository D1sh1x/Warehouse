package models

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type History struct {
	ID        int    `json:"id"`
	ItemID    int    `json:"item_id"`
	Action    string `json:"action"`
	ChangedBy string `json:"changed_by"`
	Timestamp string `json:"timestamp"`
	OldData   string `json:"old_data"`
	NewData   string `json:"new_data"`
}
