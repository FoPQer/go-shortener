package main

import (
	"log"
	"os"
)

// helper is NOT the main function, so os.Exit and log.Fatal here must be flagged.
func helper() {
	os.Exit(1)              // want `os.Exit called outside main function of main package`
	log.Fatal("err")        // want `log.Fatal called outside main function of main package`
	log.Fatalf("%v", "err") // want `log.Fatalf called outside main function of main package`
	log.Fatalln("err")      // want `log.Fatalln called outside main function of main package`
}

// init is also outside main — must be flagged.
func init() {
	os.Exit(2) // want `os.Exit called outside main function of main package`
}

// main is the exempt function — nothing here should be flagged.
func main() {
	os.Exit(0)
	log.Fatal("fatal from main")
	log.Fatalf("fatal %s", "from main")
	log.Fatalln("fatal from main")
}
