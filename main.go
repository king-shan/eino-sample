package main

import "github.com/joho/godotenv"

func main() {
	err := godotenv.Load("example.env")
	if err != nil {
		panic(err)
	}
}
