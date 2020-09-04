package models

import (
	"fmt"
	"time"
)

//Project schema para a tabela projects
type Project struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Lot struct {
	LotID     int64     `json:"lotid"`
	Price     int64     `json:"price"`
	Quantity  int64     `json:"quantity"`
	Buydate   time.Time `json:"buydate"`
	ProjectID int64     `json:"projectID"`
}

func models() {
	t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fmt.Printf("Go launched at %s\n", t.Local())
}
