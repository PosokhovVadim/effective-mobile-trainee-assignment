package storage

type Storage interface {
	AddSong()
	DeleteSong()
	GetSong()
	UpdateSong()
}
