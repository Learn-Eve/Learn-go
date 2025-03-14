package main

import "fast-learn/cmd"

func main() {
	defer cmd.Clean()
	cmd.Start()
}
