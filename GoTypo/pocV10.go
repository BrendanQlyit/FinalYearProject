package main

import (
		 "fmt"
		 "net"
		 "regexp"
		 "github.com/oschwald/geoip2-golang"
		 "github.com/oschwald/maxminddb-golang"
		 "github.com/fatih/color"
		 "log"
		 "strings"
		 "unicode"
		 "golang.org/x/net/publicsuffix"
		 "os"
		 "sync"
		 "text/tabwriter"

)

var(
 		domain string = ""
 		domainIp string
		option int
		w = new(tabwriter.Writer)
		wg = &sync.WaitGroup{}
		geolocate  bool = true
		total int
		g = color.New(color.FgHiGreen)
		r = color.New(color.FgHiRed)
		y = color.New(color.FgHiYellow)
		logo = `

 _____        _____                   _       _   _            _____        _
|   __|___   |  _  |___ ___ _____ _ _| |_ ___| |_|_|___ ___   |_   _|__ ___| |
|  |  | . |  |   __| -_|  _|     | | |  _| .'|  _| | . |   |    | || . | . | |
|_____|___|  |__|  |___|_| |_|_|_|___|_| |__,|_| |_|___|_|_|    |_||___|___|_| `


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
	total = total + len(results)
	return results
}

func extraCharAttack(domain string) []string {
	results := []string{}
	for i := 97; i < 123; i++ {
			results = append(results, fmt.Sprintf("%s%c", domain, i))
	}
	total = total + len(results)
	return results
}

// performs a bitsquat permutation attack
func bitsquattingAttack(domain string) []string {

	results := []string{}
	masks := []int32{1, 2, 4, 8, 16, 32, 64, 128}

	for i, c := range domain {
		for m := range masks {
			b := rune(int(c) ^ m)
			o := int(b)
			if (o >= 48 && o <= 57) || (o >= 97 && o <= 122) || o == 45 {
				results = append(results, fmt.Sprintf("%s%c%s", domain[:i], b, domain[i+1:]))
			}
		}
	}
	total = total + len(results)
	return results
}

// returns a count of characters in a word
func countChar(word string) map[rune]int {
	count := make(map[rune]int)
	for _, r := range []rune(word) {
		count[r]++
	}
	return count
}

func homographAttack(domain string) []string {
	// set local variables
	glyphs := map[rune][]rune{
		'a': {'à', 'á', 'â', 'ã', 'ä', 'å', 'ɑ', 'а', 'ạ', 'ǎ', 'ă', 'ȧ', 'α', 'ａ'},
		'b': {'d', 'ʙ', 'Ь', 'ɓ', 'Б', 'ß', 'β', 'ᛒ', '\u1E05', '\u1E03', '\u1D6C'}, // 'lb', 'ib'
		'c': {'ϲ', 'с', 'ƈ', 'ċ', 'ć', 'ç', 'ｃ'},
		'd': {'b', 'ԁ', 'ժ', 'ɗ', 'đ'}, // 'cl', 'dl', 'di'
		'e': {'é', 'ê', 'ë', 'ē', 'ĕ', 'ě', 'ė', 'е', 'ẹ', 'ę', 'є', 'ϵ', 'ҽ'},
		'f': {'Ϝ', 'ƒ', 'Ғ'},
		'g': {'q', 'ɢ', 'ɡ', 'Ԍ', 'Ԍ', 'ġ', 'ğ', 'ց', 'ǵ', 'ģ'},
		'h': {'һ', 'հ', '\u13C2', 'н'}, // 'lh', 'ih'
		'i': {'1', 'l', '\u13A5', 'í', 'ï', 'ı', 'ɩ', 'ι', 'ꙇ', 'ǐ', 'ĭ'},
		'j': {'ј', 'ʝ', 'ϳ', 'ɉ'},
		'k': {'κ', 'κ'}, // 'lk', 'ik', 'lc'
		'l': {'1', 'i', 'ɫ', 'ł'},
		'm': {'n', 'ṃ', 'ᴍ', 'м', 'ɱ'}, // 'nn', 'rn', 'rr'
		'n': {'m', 'r', 'ń'},
		'o': {'0', 'Ο', 'ο', 'О', 'о', 'Օ', 'ȯ', 'ọ', 'ỏ', 'ơ', 'ó', 'ö', 'ӧ', 'ｏ'},
		'p': {'ρ', 'р', 'ƿ', 'Ϸ', 'Þ'},
		'q': {'g', 'զ', 'ԛ', 'գ', 'ʠ'},
		'r': {'ʀ', 'Г', 'ᴦ', 'ɼ', 'ɽ'},
		's': {'Ⴝ', '\u13DA', 'ʂ', 'ś', 'ѕ'},
		't': {'τ', 'т', 'ţ'},
		'u': {'μ', 'υ', 'Ս', 'ս', 'ц', 'ᴜ', 'ǔ', 'ŭ'},
		'v': {'ѵ', 'ν', '\u1E7F', '\u1E7D'},      // 'v̇'
		'w': {'ѡ', 'ա', 'ԝ'}, // 'vv'
		'x': {'х', 'ҳ', '\u1E8B'},
		'y': {'ʏ', 'γ', 'у', 'Ү', 'ý'},
		'z': {'ʐ', 'ż', 'ź', 'ʐ', 'ᴢ'},
	}
	doneCount := make(map[rune]bool)
	results := []string{}
	runes := []rune(domain)
	count := countChar(domain)

	for i, char := range runes {
		// perform attack against single character
		for _, glyph := range glyphs[char] {
			results = append(results, fmt.Sprintf("%s%c%s", string(runes[:i]), glyph, string(runes[i+1:])))
		}
		// determine if character is a duplicate
		// and if the attack has already been performed
		// against all characters at the same time
		if count[char] > 1 && doneCount[char] != true {
			doneCount[char] = true
			for _, glyph := range glyphs[char] {
				result := strings.Replace(domain, string(char), string(glyph), -1)
				results = append(results, result)
			}
		}
	}
	total = total + len(results)
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
total = total + len(results)
return results
}

func swapAttack(domain string) []string{
	results := []string{}
	for i := 0; i < len(domain)-1; i++ {
		if domain[i+1] != domain[i] {
				results = append(results, fmt.Sprintf("%s%c%c%s", domain[:i], domain[i+1], domain[i], domain[i+2:]))
		}

	}
	total = total + len(results)
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
		total = total + len(results)
		return results
	}

func missingDot(domain string) []string {
		results := []string{}

		 var result string
				 result = "www" + domain
				 results = append(results, result)
				 total = total + len(results)
         return results

}

func tldAttack(domain string) []string {
		results := []string{}

		ccTLDs := []string{"com", "org", "net", "int", "edu", "gov", "mil",
											 "ac", "ad", "ae", "af", "ag", "ai", "al", "am", "an", "ao", "aq", "ar", "as", "at", "au", "aw", "ax", "az",
											 "ba", "bb", "bd", "be", "bf", "bg", "bh", "bi", "bj", "bl", "bm", "bn", "bo", "br", "bq", "bs", "bt", "bv", "bw", "by", "bz",
										 	 "ca", "cc", "cd", "cf", "cg", "ch", "ci", "ck", "cl", "cm", "cn", "co", "cr", "cs", "cu", "cv", "cw", "cy", "cz",
										 	 "dd", "de", "dj", "dk", "dm", "do", "dz",
											 "ec", "ee", "eg", "eh", "er", "es", "et", "eu",
										 	 "fi", "fj", "fk", "fm", "fo", "fr",
										   "ga", "gb", "gd", "ge", "gf", "gg", "gh", "gi", "gl", "gm", "gn", "gp", "gq", "gr", "gs", "gt", "gu", "gw", "gy",
										   "hk", "hm", "hn", "hr", "ht", "hu",
										 	 "id", "ie", "il", "im", "in", "io", "iq", "ir", "is", "it",
										 	 "je", "jm", "jo", "jp",
										 	 "ke", "kg", "kh", "ki", "km", "kn", "kp", "kr", "kw", "ky", "kz",
										 	 "la", "lb", "lc", "li", "lk", "lr", "ls", "lt", "lu", "lv", "ly",
										 	 "ma", "mc", "md", "me", "mf", "mg", "mh", "mk", "ml", "mm", "mn", "mo", "mp", "mq", "mr", "ms", "mt", "mu", "mv", "mw", "mx", "my", "mz",
										 	 "na", "nc", "ne", "nf", "ng", "ni", "nl", "no", "np", "nr", "nu", "nz",
										 	 "om",
										 	 "pa", "pe", "pf", "pg", "ph", "pk", "pl", "pm", "pn", "pr", "ps", "pt", "pw", "py",
										 	 "qa",
										 	 "re", "ro", "rs", "ru", "rw",
										 	 "sa", "sb", "sc", "sd", "se", "sg", "sh", "si", "sj", "sk", "sl", "sm", "sn", "so", "sr", "ss", "st", "su", "sv", "sx", "sz",
										 	 "tc", "td", "tf", "tg", "th", "tj", "tk", "tl", "tm", "tn", "to", "tp", "tr", "tt", "tv", "tw", "tz",
										 	 "ua", "ug", "uk", "um", "us", "uy", "uz",
										 	 "va", "vc", "ve", "vg", "vi", "vn", "vu",
										 	 "wf", "ws",
										 	 "ye", "yt", "yu",
										 	 "za", "zm", "zr", "zw"}

			var result string
						for i := range ccTLDs{
							result =  domain + "." + ccTLDs[i]
							results = append(results, result)
						}

			total = total + len(results)
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
			printReport("Additional Char", extraCharAttack(sanitizedDomain), tld)
			printReport("Bitsquatting", bitsquattingAttack(sanitizedDomain), tld)
			printReport("Homograph", homographAttack(sanitizedDomain), tld)
			printReport("TLD", tldAttack(sanitizedDomain), "")

	}
}


func getOption() int{
	fmt.Println("Select a search option")
	fmt.Println("1:	List perumutations")
	fmt.Println("2:	List perumutations + method")
	fmt.Println("3:	List perumutations + IP address")
	fmt.Println("4:	List perumutations + IP address + Geo Location")
	fmt.Println("5:	List Description of attack types")

	fmt.Print("Option: ")
	fmt.Scan(&option)

	if option == 1 ||  option == 2 || option == 3 || option == 4{
		return option
		} else if option == 5{
			printDescription()
			return option
		}else{
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

	if(strings.Contains(Domain, ".")){
		r.Domain = Domain + "" + tld
  }else{
		r.Domain = Domain + "." + tld
	}


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
}/*else if option == 5{
		printDescription()
}*/

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

func printLogo(){
	g.Printf(logo)
	fmt.Println()
	}

func printDescription(){

	r.Printf("Additional Char: ")
	y.Println("\tAdds an extra character ranging from a-z to the end of the domain.")
	r.Printf("Bitsquatting: ")
	fmt.Println("\t\tBitsquatting refers to the registration of a domain names one bit different\n\t\t\tthan a popular domain.")
	r.Printf("Character Swap: ")
	y.Println("\tSwaps the order of the characters in the domain.")
	r.Printf("Duplication: ")
	fmt.Println("\t\tDoubles each character within the domain one at a time.")
	r.Printf("Fat Finger: ")
	y.Println("\t\tReplaces each character wihin the domain with a neighboring character of\n\t\t\tthe keyboard.")
	r.Printf("Homograph: ")
	fmt.Println("\t\tHomograph attacks are phishing schemes in which the phisher takes advantage\n\t\t\tof the ability to register internationalized domain names (IDNs) using\n\t\t\tnon-Latin characters that look the same as Latin characters\n\t\t\t(such as some Cyrillic or Greek characters, for example).")
	r.Printf("Missing Dot: ")
	y.Println("\t\tRemoves the '.' after www")
	r.Printf("Omission: ")
	fmt.Println("\t\tRemoves one character form the domain")
	r.Printf("TLD: ")
	y.Println("\t\t\tReplaces the TLD of the domain with every possible TLD")








	}

func main (){
	printLogo()
	fmt.Println("\n\n")
	y.Println("Welcome to Go Permutation Tool")
	fmt.Print("Please enter a domain name: ")
	fmt.Scanf("%s\r", &domain)

	if !validateDomainName(domain) {
		 for !validateDomainName(domain){
					fmt.Printf("Domain Name %s is invalid", domain)
					println("")
					fmt.Print("Please enter a domain name: ")
					n, err := fmt.Scanf("%s\r", &domain)
					if err != nil{
						fmt.Println(n, err)
					}
					
				}
	}

	if validateDomainName(domain){

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
			//	tldAttack(sanitizedDomain)
				targets := []string{sanitizedDomain + "." + tld}
				runPermutations(targets)
				fmt.Println("Number of Permutations:- ",  total)



                 }

         }
}
}
