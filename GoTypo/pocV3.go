package main

import ( "fmt"
		 "net"
		 "regexp"
		 "github.com/oschwald/geoip2-golang"
		 "log"
		 "strings"
		 "unicode"

)

var(
 		domain string
 		domainIp string
)

func geoCheck(ipAdd string) string {
	if ipAdd != "" {
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
		return record.Country.IsoCode +
			" " + record.City.Names["en"]
	}
	return ""
}


func validateDomainName(domain string) bool {


         RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z
									  ]{2,3})$`)

         return RegExp.MatchString(domain)
 }

func omissionAttack(domain string) []string {
	results := []string{}
	for i := range domain{
		results = append(results, fmt.Sprintf("%s%s", domain[:i], domain[i+1:]))
	}
	fmt.Print("Omission Attack: ")
	fmt.Println(strings.Join(results, "\nOmission Attack: "))
	return results
}

func duplicationAttack(domain string) []string {
	results := []string{}
	count := make(map[string]int)
	for i, c :=range domain {
		if unicode.IsLetter(c) {
		result := fmt.Sprintf("%s%c%c%s", domain[:i], domain[i], domain[i], domain[i+1:])
		//remove duplicates
		count[result]++
		if count[result] <2 {
			results = append(results, result)
		}
	}
}
fmt.Print("Duplication Attack: ")
fmt.Println(strings.Join(results, "\nDuplication Attack: "))

return results
}

func swapAttack(domain string) []string{
	results := []string{}
	for i := 0; i < len(domain)-1; i++ {
		if domain[i+1] != domain[i] {
				results = append(results, fmt.Sprintf("%s%c%c%s", domain[:i], domain[i+1], domain[i], domain[i+2:]))
		}

	}
	fmt.Print("Swap Attack: ")
	fmt.Println(strings.Join(results, "\nSwap Attack: "))
	return results
}

func FFAttack(domain string) []string {
	results := []string{}
		count := make(map[string]int)
		keyboard := map[rune]string{'1': "2q", '2': "3wq1", '3': "4ew2", '4': "5re3", '5': "6tr4", '6': "7yt5", '7': "8uy6", '8': "9iu7", '9': "0oi8", '0': "po9",
			'q': "12wa",'w': "3esaq2", 'e': "4rdsw3", 'r': "5tfde4", 't': "6ygfr5", 'y': "7uhgt6", 'u': "8ijhy7", 'i': "9okju8", 'o': "0plki9", 'p': "lo0",
			'a': "qwsz", 's': "edxzaw", 'd': "rfcxse", 'f': "tgvcdr", 'g': "yhbvft", 'h': "ujnbgy", 'j': "ikmnhu", 'k': "olmji", 'l': "kop",
			'z': "asx", 'x': "zsdc", 'c': "xdfv", 'v': "cfgb", 'b': "vghn", 'n': "bhjm", 'm': "njk"}
		for i, c := range domain {
				for _, char := range []rune(keyboard[c]) {
					result := fmt.Sprintf("%s%c%s", domain[:i], char, domain[i+1:])
					// remove duplicates
					count[result]++
					if count[result] < 2 {
						results = append(results, result)
					}

							fmt.Println("Fat Finger: ", result, "Address: ", performLookUp(result), "GEO: ", geoCheck(performLookUp(result)))
				}
		}
		return results
	}

func missingDot(domain string) string {

		 var removeDot string
		// removeDot = domain
         if(strings.Contains(domain, "www")){
         removeDot = strings.Replace(domain, ".", "", 1)
				 fmt.Println("Missing Dot Attack: ", removeDot, "Address: ", performLookUp(removeDot), "GEO:", geoCheck(performLookUp(removeDot)) )
         }

         return removeDot

}

func runPermutations(domain string){
	missingDot(domain)
	omissionAttack(domain)
	duplicationAttack(domain)
	FFAttack(domain)
	swapAttack(domain)
}

func performLookUp(domain string) string{

		 addr, err1 := net.ResolveIPAddr("ip4", domain)

				if err1 != nil{
					 return ""
				 }
					return addr.String()

}

func main (){


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


					//	fmt.Println("Missing Dot " + missingDot(domain) )
					//	fmt.Println(strings.Join(omissionAttack(domain), "\n"))
					//	fmt.Println(strings.Join(duplicationAttack(domain), "\n"))
				//		fmt.Println(strings.Join(FFAttack(domain), "\n"))
							runPermutations(domain)

					//		fmt.Println(strings.Join(swapAttack(domain), "\n"))


                 }

         }
}
