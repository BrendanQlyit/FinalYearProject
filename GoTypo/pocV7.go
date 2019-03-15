package main

import (
		 "fmt"
		 "net"
		 "regexp"
		 "github.com/oschwald/geoip2-golang"
		 "github.com/oschwald/maxminddb-golang"
		 "log"
		 "strings"
		 "unicode"
		 "golang.org/x/net/publicsuffix"
		 "os"
		 "sync"
		 "text/tabwriter"

)

var(
 		domain string
 		domainIp string
		option int
		w = new(tabwriter.Writer)
		wg = &sync.WaitGroup{}
		geolocate  bool = true
)

type GeoIPRecord struct {
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

type Record struct {
	Technique   string `json:"technique"`
	Domain      string `json:"domain"`
	A           string `json:"a_record"`
	Geolocation string `json:"geolocation"`
}

func geoCheck(ipAdd string) string {
	if ipAdd != "" {
		db, err := maxminddb.Open("GeoLite2-City.mmdb")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		ip := net.ParseIP(domainIp)
		var record GeoIPRecord
	  err = db.Lookup(ip, &record)
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

return results
}

func swapAttack(domain string) []string{
	results := []string{}
	for i := 0; i < len(domain)-1; i++ {
		if domain[i+1] != domain[i] {
				results = append(results, fmt.Sprintf("%s%c%c%s", domain[:i], domain[i+1], domain[i], domain[i+2:]))
		}

	}
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
				}
		}
		return results
	}

func missingDot(domain string) []string {
		results := []string{}

		 var result string
				 result = "www" + domain
				 results = append(results, result)
         return results

}


func runPermutations(targets []string) {

		for _, target := range targets {
			sanitizedDomain, tld := sepInput(target)
			printReport("missing dot", missingDot(sanitizedDomain), tld)
			printReport("Omission", omissionAttack(sanitizedDomain), tld)
			printReport("Duplication", duplicationAttack(sanitizedDomain), tld)
			printReport("Fat Finger", FFAttack(sanitizedDomain), tld)
			printReport("Character Swap", swapAttack(sanitizedDomain), tld)

	}
}


func getOption() int{
	fmt.Println("Select a search option")
	fmt.Println("1:	List perumutations")
	fmt.Println("2:	List perumutations + method")
	fmt.Println("3:	List perumutations + IP address")
	fmt.Println("4:	List perumutations + IP address + Geo Location")
	fmt.Print("Option: ")
	fmt.Scan(&option)

	if option == 1 ||  option == 2 || option == 3{
		return option
		} else {
			return 9
		}

}

//Seperates the domain name and TLD
func sepInput(domain string) (sepDomain, tld string){

			tld, _ = publicsuffix.PublicSuffix(domain)
			sepDomain = strings.Replace(domain, "."+tld, "", -1) //remove tld from domain
			sepDomain = strings.Replace(sepDomain, "www.", "", -1) //remove www.
			return sepDomain, tld


}

func performLookUp(domain string) string{

		 addr, err1 := net.ResolveIPAddr("ip4", domain)

				if err1 != nil{
					 return ""
				 }
					return addr.String()

}


// performs lookups on individual records
func doLookups(Technique, Domain, tld string, out chan<- Record, geolocate bool) {
	defer wg.Done()
	r := new(Record)
	r.Technique = Technique
	r.Domain = Domain + "." + tld

		r.A = performLookUp(r.Domain)

	if geolocate {
		r.Geolocation = geoCheck(performLookUp(r.Domain))
	}
	out <- *r
}

// runs bulk lookups on list of domains
func runLookups(technique string, results []string, tld string, out chan<- Record,  geolocate bool) {
	for _, r := range results {
		wg.Add(1)
		go doLookups(technique, r, tld, out, geolocate)
	}
}

func printReport(technique string, results []string, tld string) {
	out := make(chan Record)
	w.Init(os.Stdout, 18, 8, 0, '\t', 0)

	if option == 4{
		runLookups(technique, results, tld, out, true)
	}else if option == 1{
		for _, result := range results {
			fmt.Println(result + "." + tld)
			}
		}else if option == 2{
			for _, result := range results {
				printResults(w, technique, result, tld)
			}
	}else if option == 3{
		runLookups(technique, results, tld, out, false)
}

	go monitorWorker(wg, out)
	for r := range out {
		r.printRecordData(w)
	}
}

func printResults(writer *tabwriter.Writer, technique, result, tld string) {
		fmt.Fprintln(w, technique+"\t"+result+"."+tld+"\t")
		w.Flush()
}

func monitorWorker(wg *sync.WaitGroup, channel chan Record) {
	wg.Wait()
	close(channel)
}

func (r *Record) printRecordData(writer *tabwriter.Writer) {
			if option == 4{
			fmt.Fprintln(writer, r.Technique+"\t"+r.Domain+"\t"+"IP:"+r.A+"\t"+"GEO:"+r.Geolocation+"\t")
			writer.Flush()
		}else if option == 3{
			fmt.Fprintln(writer, r.Technique+"\t"+r.Domain+"\t"+"IP:"+r.A+"\t")
			writer.Flush()
		}

	}

func main (){

	fmt.Println("Typosquatting POC!")
	fmt.Print("Please enter a domain name: ")
	fmt.Scanf("%s", &domain)
	getOption()


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



				sanitizedDomain, tld := sepInput(domain)
				targets := []string{sanitizedDomain + "." + tld}
				runPermutations(targets)



                 }

         }
}
