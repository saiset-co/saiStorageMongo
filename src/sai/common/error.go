package common

import (
	"time"
	"fmt"
	"saiStorageMongo/src/github.com/fatih/color"
)

type Error struct {
	What string `json:"error"`
	Code int    `json:"code"`
}

func (e *Error) Error() string {
	r := color.New(color.FgRed, color.Bold)
	r.Println(fmt.Sprintf("Error at %v: %s", time.Now(), e.What))

	return string(ConvertInterfaceToJson(e))
}
