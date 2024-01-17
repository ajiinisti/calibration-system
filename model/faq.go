package model

type Faq struct {
	BaseModel
	Title  string
	Order  int
	Answer string
	Active bool
}
