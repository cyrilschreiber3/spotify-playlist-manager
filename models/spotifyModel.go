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
	Songs       []Song
}

func NewSong(id, title, artist, album, coverUrl string) Song {
	return Song{
		ID:       id,
		Title:    title,
		Artist:   artist,
		Album:    album,
		CoverUrl: coverUrl,
	}
}

func NewPlaylist(id, name, description, coverUrl string, songs []Song) Playlist {
	return Playlist{
		ID:          id,
		Name:        name,
		Description: description,
		CoverUrl:    coverUrl,
		Songs:       songs,
	}
}
