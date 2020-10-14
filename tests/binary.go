package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func isTTY(file os.File) bool {
	if fileInfo, _ := file.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

func main() {

	showTTY := flag.Bool("showTTY", false, "Show TTY Status")
	sleep := flag.Uint("sleep", 0, "Sleep for duration")
	flag.Parse()

	if *showTTY {
		fmt.Printf("TTY Status -> Stdin : %v, Stdout : %v, Stderr : %v\n", isTTY(*os.Stdin), isTTY(*os.Stdout), isTTY(*os.Stderr))
	}

	if e := os.Getenv("SEXYENV"); e != "" {
		fmt.Printf("SEXYENV is set to %v\n", e)
	}

	fmt.Print("Random string which you should see ")

	if *sleep > 0 {
		fmt.Printf("Sleeping for %v seconds\n", *sleep)
		time.Sleep(time.Duration(*sleep) * time.Second)
	}

	fmt.Println("before seeing this string")

	// var input string
	// fmt.Print("Enter input: ")
	// fmt.Scanln(&input)
	// fmt.Println(input)
}
