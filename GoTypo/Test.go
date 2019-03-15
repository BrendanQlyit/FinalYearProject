package main

import (
				"fmt"
				"github.com/likexian/whois-go"
)

func main() {
				result, err := whois.Whois("hoogle.ie")

				if err != nil {
								fmt.Println(err)
				}

				fmt.Println(result)
}
