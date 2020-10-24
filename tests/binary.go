package main

import (
	"bufio"
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

	fmt.Printf("Args: %v\n", os.Args)

	showTTY := flag.Bool("tty", false, "Show TTY Status")
	showEnv := flag.Bool("env", false, "Show Environs")
	sleep := flag.Uint("sleep", 0, "Sleep for duration")
	takeInput := flag.Bool("input", false, "Prompt for input")

	flag.Parse()

	if *showTTY {
		fmt.Printf("TTY Status -> Stdin : %v, Stdout : %v, Stderr : %v\n", isTTY(*os.Stdin), isTTY(*os.Stdout), isTTY(*os.Stderr))
	}

	if *showEnv {
		fmt.Println(os.Environ())
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

	if *takeInput {
		fmt.Print("Give me an input and I will repeat it back> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fmt.Printf("Your input is %v", text)
	}

	// var input string
	// fmt.Print("Enter input: ")
	// fmt.Scanln(&input)
	// fmt.Println(input)
}
