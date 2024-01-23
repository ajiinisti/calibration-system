package model

type Announcement struct {
	BaseModel
	Title    string
	Order    int
	FileName string
	File     []byte
	FileLink string
	Active   bool
}
