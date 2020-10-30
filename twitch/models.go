package twitch

type Followers struct {
	Total      int         `json:"total"`
	Followers  []Followers `json:"data"`
	Pagination `json:"pagination"`
}

type Follower struct {
	FromID     string `json:"from_id"`
	FromName   string `json:"from_name"`
	ToID       string `json:"to_id"`
	ToName     string `json:"to_name"`
	FollowedAt string `json:"followed_at"`
}

type UsersList struct {
	Users []User `json:"data"`
}

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
}
type Pagination struct {
	Cursor string `json:"cursor"`
}

type TwitchErrorMessage struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
