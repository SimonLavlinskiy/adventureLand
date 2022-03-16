package repository

type Result struct {
	ID            uint   `gorm:"primaryKey"`
	Type          string `gorm:"embedded"`
	ItemId        *uint  `gorm:"embedded"`
	Item          *Item
	CountItem     *uint `gorm:"embedded"`
	SpecialItemId *uint `gorm:"embedded"`
	SpecialItem   *Item
	Experience    *int `gorm:"embedded"`
}

type Casual struct {
	Experience int
}

func (r *Result) Casual() *Casual {
	return &Casual{
		Experience: *r.Experience,
	}
}

type CasualPlus struct {
	Experience int
	Item       Item
	CountItem  uint
}

func (r *Result) CasualPlus() *CasualPlus {
	return &CasualPlus{
		Experience: *r.Experience,
		Item:       *r.Item,
		CountItem:  *r.CountItem,
	}
}

type SuperCasual struct {
	Experience  int
	Item        Item
	CountItem   uint
	SpecialItem Item
}

func (r *Result) SuperCasual() *SuperCasual {
	if r == nil {
		return nil
	}
	return &SuperCasual{
		Experience:  *r.Experience,
		Item:        *r.Item,
		CountItem:   *r.CountItem,
		SpecialItem: *r.SpecialItem,
	}
}

func (u User) UserGetResult(r Result) {
	//var UsResult struct
	switch r.Type {
	case "casual":
		//var UsResult = r.Casual
	case "casualPlus":
		//usResult := r.CasualPlus
	case "superCasual":
		//usResult := r.SuperCasual
	}

	//fmt.Println(UsResult)
}
