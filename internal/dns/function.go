package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"403unlocker-cli/internal/common"

	"github.com/cavaliergopher/grab/v3"
	"github.com/spf13/viper"
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

func CheckAndCacheDNS(url string) error {
	// Initialize Viper for DNS configuration
	dnsViper := viper.New()
	dnsViper.SetConfigFile(common.DNS_CONFIG_FILE)

	// Try to read the DNS config file
	if err := dnsViper.ReadInConfig(); err != nil {
		// If file doesn't exist, download it
		if err := common.DownloadConfigFile(common.DNS_CONFIG_URL, common.DNS_CONFIG_FILE); err != nil {
			return fmt.Errorf("error downloading DNS config: %w", err)
		}

		// Try reading again after download
		if err := dnsViper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading DNS config: %w", err)
		}
	}

	// Get DNS list from config
	dnsList := dnsViper.GetStringSlice("dnsServers")
	if len(dnsList) == 0 {
		return errors.New("no DNS servers found in config")
	}

	fmt.Println("\n+--------------------+------------+")
	fmt.Printf("| %-18s | %-10s |\n", "DNS Server", "Status")
	fmt.Println("+--------------------+------------+")

	var validDNSList []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, dns := range dnsList {
		wg.Add(1)
		go func(dns string) {
			defer wg.Done()

			client := common.ChangeDNS(dns)
			resp, err := client.Get(url)
			if err != nil {
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Red, "Error", common.Reset)
				return
			}
			defer resp.Body.Close()

			statusParts := strings.Split(resp.Status, " ")
			if len(statusParts) < 2 {
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Red, "Invalid", common.Reset)
				return
			}

			statusCode, err := strconv.Atoi(statusParts[0])
			if err != nil {
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Red, "Error", common.Reset)
				return
			}

			statusText := statusParts[1]
			if statusCode == http.StatusOK {
				mu.Lock()
				validDNSList = append(validDNSList, dns)
				mu.Unlock()
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Green, statusText, common.Reset)
			} else {
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Red, statusText, common.Reset)
			}
		}(dns)
	}

	wg.Wait()
	fmt.Println("+--------------------+------------+")

	// Initialize Viper for cache
	cacheViper := viper.New()
	cacheViper.SetConfigFile(common.CHECKED_DNS_CONFIG_FILE)

	// Save valid DNS list if any found
	if len(validDNSList) > 0 {
		cacheViper.Set("validDNSServers", validDNSList)
		if err := cacheViper.WriteConfig(); err != nil {
			if err := cacheViper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("error writing cached DNS config: %w", err)
			}
		}
		fmt.Printf("Cached %d valid DNS servers to %s\n", len(validDNSList), common.CHECKED_DNS_CONFIG_FILE)
	} else {
		fmt.Println("No valid DNS servers found to cache.")
	}

	return nil
}

func CheckWithURL(commandLintFirstArg string, check bool, timeout int) error {
	fileToDownload := commandLintFirstArg

	var dnsFile string
	if check {
		err := CheckAndCacheDNS(fileToDownload)
		if err != nil {
			return err
		}
		dnsFile = common.CHECKED_DNS_CONFIG_FILE
	} else {
		dnsFile = common.DNS_CONFIG_FILE
	}

	// Read the DNS list from the determined file
	dnsList, err := common.ReadDNSFromFile(dnsFile)
	if err != nil {
		// Fallback to download and read from the original DNS file
		err = common.DownloadConfigFile(common.DNS_CONFIG_URL, common.DNS_CONFIG_FILE)
		if err != nil {
			return fmt.Errorf("error downloading DNS config file: %w", err)
		}
		dnsList, err = common.ReadDNSFromFile(common.DNS_CONFIG_FILE)
		if err != nil {
			return fmt.Errorf("error reading DNS list from file: %w", err)
		}
	}

	dnsSizeMap := make(map[string]int64)

	fmt.Printf("\nTimeout: %d seconds\n", timeout)
	fmt.Printf("URL: %s\n\n", fileToDownload)

	// Print table header
	fmt.Println("+--------------------+----------------+")
	fmt.Printf("| %-18s | %-14s |\n", "DNS Server", "Download Speed")
	fmt.Println("+--------------------+----------------+")

	rand := time.Now().UnixMilli()
	tempDir := common.AddPathToDir(common.GetTempDir(), strconv.Itoa(int(rand)))

	var wg sync.WaitGroup
	for _, dns := range dnsList {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		clientWithCustomDNS := common.ChangeDNS(dns)
		client := grab.NewClient()
		client.HTTPClient = clientWithCustomDNS
		req, err := grab.NewRequest(fmt.Sprintf("%v", tempDir), fileToDownload)

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

	wg.Wait()
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

	os.RemoveAll(fmt.Sprintf("%v", tempDir))
	return nil
}
