package main

import (
	"fmt"
	"github.com/dkuntz2/fate"
	"net/http"
	"os"
)

func main() {
	f := fate.New()
	f.Run()

	listenSpec := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("Starting server on http://%s\n", listenSpec)
	http.ListenAndServe(listenSpec, f.Router())
}
