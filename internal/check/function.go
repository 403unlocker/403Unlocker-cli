package check

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/salehborhani/403Unlocker-cli/internal/common"
	"github.com/urfave/cli/v2"
)

func ChangeDNS(dns string) *http.Client {
	dialer := &net.Dialer{}
	customResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dnsServer := fmt.Sprintf("%s:53", dns)
			return dialer.DialContext(ctx, "udp", dnsServer)
		},
	}
	customDialer := &net.Dialer{
		Resolver: customResolver,
	}
	transport := &http.Transport{
		DialContext: customDialer.DialContext,
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}

func CheckWithDNS(c *cli.Context) error {
	url := c.Args().First()
	url = ensureHTTPS(url)

	// Print header
	fmt.Println("\n+--------------------+------------+")
	fmt.Printf("| %-18s | %-10s |\n", "DNS Server", "Status")
	fmt.Println("+--------------------+------------+")

	dnsList, err := ReadDNSFromFile(common.DNS_CONFIG_FILE)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var wg sync.WaitGroup
	for _, dns := range dnsList {
		wg.Add(1)
		go func(dns string) {
			defer wg.Done()
			client := ChangeDNS(dns)
			resp, err := client.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			code := strings.Split(resp.Status, " ")
			statusCodeInt, err := strconv.Atoi(code[0])
			if err != nil {
				fmt.Println("Error converting status code:", err)
				return
			}

			// Format table row with colored status
			if statusCodeInt != http.StatusOK {
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Red, code[1], common.Reset)
			} else {
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Green, code[1], common.Reset)
			}

		}(dns)
	}
	wg.Wait()

	// Print footer
	fmt.Println("+--------------------+------------+")
	return nil
}

func ReadDNSFromFile(filename string) ([]string, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		fmt.Println("HOME environment variable not set")
		os.Exit(1)
	}
	filename = homeDir + "/" + filename
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	dnsServers := strings.Fields(string(data))
	return dnsServers, nil
}
func DomainValidator(domain string) bool {
	domainRegex := `^(http[s]?:\/\/)?([a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}).*?$`
	match, _ := regexp.MatchString(domainRegex, domain)
	if !match {
		return false
	}
	// Additional checks:
	// 1. The total length of the domain should not exceed 253 characters.
	if len(domain) > 253 {
		return false
	}
	// 2. Each segment between dots should be between 1 and 63 characters long.
	segments := strings.Split(domain, ".")
	for _, segment := range segments {
		if len(segment) < 1 || len(segment) > 63 {
			return false
		}
	}
	return true
}

func ensureHTTPS(URL string) string {
	// Regex to check if the URL starts with https://
	regexHTTPS := `^(https)://`
	reHTTPS, err := regexp.Compile(regexHTTPS)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return URL
	}
	regexHTTP := `^(http)://`
	reHTTP, err := regexp.Compile(regexHTTP)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return URL
	}
	if reHTTP.MatchString(URL) {
		URL = strings.TrimPrefix(URL, "http://")
	}
	if reHTTPS.MatchString(URL) {
		URL = strings.TrimPrefix(URL, "https://")
	}
	URL = "https://" + URL
	// Parse the URL to extract the host
	parsedURL, err := url.Parse(URL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}
	// Return only the scheme and host (e.g., https://example.com)
	return "https://" + parsedURL.Host + "/"
}
