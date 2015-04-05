/**
 *
 *
 */

package main

import (
	"fmt"
	"net"
	"ola"
)

func HandleError(err error) {
	fmt.Printf("%s\n", err.Error())
}

func main() {
	fmt.Printf("Hello World!\n")
	sock, err := net.Dial("tcp", "localhost:9010")
	if err != nil {
		HandleError(err)
		return
	}
	client := ola.NewClient(sock)
	fmt.Printf("Going to attempt to grab Plugins")
	plugins, err := client.FetchPlugins()
	fmt.Printf("We got a response or error!")
	if err != nil {
		HandleError(err)
		return
	}
	fmt.Printf("%s", plugins)
}
