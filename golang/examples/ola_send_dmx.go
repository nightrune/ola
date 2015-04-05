/**
 *
 *
 */

package main

import "fmt"
import "ola"

func main() {
	fmt.Printf("Hello World!\n")
	var client = new(ola.Client)
	_, _ = client.FetchPlugins()
}
