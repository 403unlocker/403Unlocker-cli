package dns

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/salehborhani/403Unlocker-cli/internal/check"
	"github.com/salehborhani/403Unlocker-cli/internal/common"
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
		return false
	}
	// Check if the host is present
	if u.Host == "" {
		return false
	}
	return true
}
func CheckWithURL(c *cli.Context) error {
	fileToDownload := c.Args().First()
	timeout := c.Int("timeout")
	dnsList, err := check.ReadDNSFromFile(common.DNS_CONFIG_FILE)
	if err != nil {
		fmt.Println("Error reading DNS list:", err)
		return err
	}

	dnsSizeMap := make(map[string]int64)
	fmt.Printf("\nTimeout: %d seconds\n", timeout)
	fmt.Printf("URL: %s\n\n", fileToDownload)

	// Print table header
	fmt.Println("+--------------------+----------------+")
	fmt.Printf("| %-18s | %-14s |\n", "DNS Server", "Download Speed")
	fmt.Println("+--------------------+----------------+")

	tempDir := time.Now().UnixMilli()
	for _, dns := range dnsList {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		clientWithCustomDNS := check.ChangeDNS(dns)
		client := grab.NewClient()
		client.HTTPClient = clientWithCustomDNS

		req, err := grab.NewRequest(fmt.Sprintf("/tmp/%v", tempDir), fileToDownload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for DNS %s: %v\n", dns, err)
		}
		req = req.WithContext(ctx)

		resp := client.Do(req)
		dnsSizeMap[dns] = resp.BytesComplete()

		speed := common.FormatDataSize(resp.BytesComplete() / int64(timeout))
		if resp.BytesComplete() == 0 {
			fmt.Printf("| %-18s | %s%-14s%s |\n", dns, common.Red, speed+"/s", common.Reset)
		} else {
			fmt.Printf("| %-18s | %-14s |\n", dns, speed+"/s")
		}
	}

	// Print table footer
	fmt.Println("+--------------------+----------------+")

	// Find and display the best DNS
	var maxDNS string
	var maxSize int64
	for dns, size := range dnsSizeMap {
		if size > maxSize {
			maxDNS = dns
			maxSize = size
		}
	}

	fmt.Println() // Add a blank line for separation
	if maxDNS != "" {
		bestSpeed := common.FormatDataSize(maxSize / int64(timeout))
		fmt.Printf("Best DNS: %s%s%s (%s%s/s%s)\n",
			common.Green, maxDNS, common.Reset,
			common.Green, bestSpeed, common.Reset)
	} else {
		fmt.Println("No DNS server was able to download any data.")
	}

	os.RemoveAll(fmt.Sprintf("/tmp/%v", tempDir))
	return nil
}
