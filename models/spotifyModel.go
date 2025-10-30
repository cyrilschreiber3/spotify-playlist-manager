package models

type Song struct {
	ID       string
	Title    string
	Artist   string
	Album    string
	CoverUrl string
}

type Playlist struct {
	ID          string
	Name        string
	Description string
	CoverUrl    string
	Length      int
	Songs       []Song
}
