package util

import "fmt"

func DoneOrDieWithMesg(err error, mesg string) {
	if err != nil {
		fmt.Println(mesg)
		panic(err)
	}
}
