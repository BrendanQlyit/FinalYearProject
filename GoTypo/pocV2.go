package main

import ( "fmt"
		 "net"
		 "regexp"
		 "github.com/oschwald/geoip2-golang"
		 "log"
		 "strings"

)


func validateDomainName(domain string) bool {


         RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z
									  ]{2,3})$`)

         return RegExp.MatchString(domain)
 }



func missingDot(domain string) string {
		 
		 var removeDot string 
		// removeDot = domain
         if(strings.Contains(domain, "www")){
         removeDot = strings.Replace(domain, ".", "", 1)

         fmt.Println("Before : ", domain)
         fmt.Println("After : ", removeDot)
         }
         return removeDot
	
}

func main (){
	
	var domain string
	var domainIp string
	
	fmt.Println("Typosquatting POC!")
	fmt.Print("Please enter a domain name: ")
	fmt.Scanf("%s", &domain)

         if !validateDomainName(domain) {
                 fmt.Printf("Domain Name %s is invalid\n", domain)
         } else {
                 fmt.Printf("Domain Name %s is VALID\n", domain)
                 addr, err1 := net.ResolveIPAddr("ip4", domain)
                 if addr == nil{
	                 fmt.Println(err1)
                 } else if err1 == nil{ 
	                 domainIp = addr.String()
	                 fmt.Println("The IP address of " + domain + " = " + domainIp)
	                 
	                 
						db, err := geoip2.Open("GeoLite2-City.mmdb")
						if err != nil {
							log.Fatal(err)
						}
						defer db.Close()
						ip := net.ParseIP(domainIp)
						record, err := db.City(ip)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println(record.Country.IsoCode + 
							" " + record.City.Names["en"] )
						
						
						fmt.Println("Missing Dot " + missingDot(domain) )

							
                 }

         }
}
	
	



