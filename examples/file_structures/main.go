package main

import (
	"fmt"

	events "github.com/Basemint-Community/Confession/Events"
	utils "github.com/Basemint-Community/Confession/Utils"
)

func init() {
	utils.LoadEnv()
    events.ConfessionBot()
}

func main() {
    fmt.Println("Yes Saar")
}


 