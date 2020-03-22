package main

import (
	"github.com/jorie1234/corona/corona"
)


func main() {
	c := corona.GetCoronaData()
	if c != nil {
		corona.SaveCoronaImage(c, "corona.png")
	}
}
