package models

type Account struct {
    ID          uint     `json:"-" gorm:"primary_key"`
    Username    string   `json:"username,omitempty"`
    Password    string   `json:"-"`
    EMail       string   `json:"enail"`
    Firstname   string   `json:"firstname"`
    Lastname    string   `json:"lastname"`
    Bookmarks []Bookmark `json:"bookmarks,omitempty"`
}
