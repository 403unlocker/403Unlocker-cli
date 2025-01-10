package dns

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/salehborhani/403Unlocker-cli/internal/check"
	"github.com/urfave/cli/v2"
)

func URLValidator(URL string) bool {
	// Parse the URL
	u, err := url.Parse(URL)
	if err != nil {
		return false
	}
	// Check if the scheme is either "http" or "https"
	if u.Scheme != "http" && u.Scheme != "https" {
		u = AppendScheme(URL)
		if u == nil {
			return false
		}
	}
	// Check if the host is present
	if u.Host == "" {
		return false
	}
	return true
}

func AppendScheme(Url string) *url.URL {
	// Default to "http" when scheme is missing
	u, err := url.Parse("http://" + Url)
	if err != nil {
		return nil
	}
	return u
}

func CheckWithURL(c *cli.Context) error {
	fileToDownload := c.Args().First()
	dnsList, err := check.ReadDNSFromFile("config/dns.conf")
	if err != nil {
		fmt.Println("Error reading DNS list:", err)
		return err
	}

	// Map to store the total size downloaded by each DNS
	dnsSizeMap := make(map[string]int64)
	// Create a context with a timeout of 15 seconds

	for _, dns := range dnsList {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Create a custom HTTP client with the specified DNS
		clientWithCustomDNS := check.ChangeDNS(dns)
		client := grab.NewClient()
		client.HTTPClient = clientWithCustomDNS

		// Create a new download request
		req, err := grab.NewRequest(".", fileToDownload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for DNS %s: %v\n", dns, err)
		}
		req = req.WithContext(ctx)

		// Start the download
		resp := client.Do(req)
		// if err := resp.Err(); err != nil {

		// 	// fmt.Fprintf(os.Stderr, "Download failed for DNS %s: %v\n", dns, err)
		// }

		// Update the total size downloaded by this DNS
		dnsSizeMap[dns] += resp.BytesComplete() // Use BytesComplete() for partial downloads

		fmt.Printf("Downloaded %d KB using DNS %s\n", resp.BytesComplete()/8_000, dns)

	}

	// Determine which DNS downloaded the most data
	var maxDNS string
	var maxSize int64
	for dns, size := range dnsSizeMap {
		if size > maxSize {
			maxDNS = dns
			maxSize = size
		}
	}
	if maxDNS != "" {
		fmt.Printf("DNS %s downloaded the most data: %d KB\n", maxDNS, maxSize/8_000)
	} else {
		fmt.Println("No DNS server was able to download any data.")
	}

	return nil
}
