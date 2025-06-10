package common

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

const (
	// Color
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"

	// DNS config
	DNS_CONFIG_FILE         = ".config/403unlocker/dns.yml"
	CHECKED_DNS_CONFIG_FILE = ".config/403unlocker/checked_dns.yml"
	DOCKER_CONFIG_FILE      = ".config/403unlocker/dockerRegistry.yml"
	DNS_CONFIG_URL          = "https://raw.githubusercontent.com/403unlocker/403Unlocker-cli/refs/heads/main/config/dns.yml"
	DOCKER_CONFIG_URL       = "https://raw.githubusercontent.com/403unlocker/403Unlocker-cli/refs/heads/main/config/dockerRegistry.yml"

	// OS names
	WINDOWS_OS_NAME = "windows"
)

// FormatDataSize converts the size in bytes to a human-readable string in KB, MB, or GB.
func FormatDataSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d Bytes", bytes)
	}
}
func DownloadConfigFile(url, path string) error {
	homeDir := GetHomeDir()
	if homeDir == "" {
		fmt.Println("HOME environment variable not set")
		os.Exit(1)
	}
	filePath := AddPathToDir(homeDir, path)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return err
	}

	out, err := os.Create(filePath)

	if err != nil {
		fmt.Println(err)

		return err
	}
	defer out.Close()

	if err != nil {
		fmt.Println("Could not download config file.")
		return err
	}

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Could not get the response: ", err)
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		fmt.Println("Could not copy content file")
		return err
	}

	return nil
}

func WriteDNSToFile(filename string, dnsList []string) error {
	homeDir := GetHomeDir()
	if homeDir == "" {
		fmt.Println("HOME environment variable not set")
		os.Exit(1)
	}
	filename = AddPathToDir(homeDir, filename)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filename, err)
			return err
		}
		file.Close()
	}
	data := map[string][]string{
		"dnsServers": dnsList,
	}
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", filename, err)
		return err
	}

	return nil
}

func ReadDNSFromFile(filename string) ([]string, error) {
	homeDir := GetHomeDir()
	if homeDir == "" {
		fmt.Println("HOME environment variable not set")
		os.Exit(1)
	}
	filename = AddPathToDir(homeDir, filename)
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	// Get DNS list
	dnsServers := viper.GetStringSlice("dnsServers")
	if len(dnsServers) == 0 {
		return nil, fmt.Errorf("no DNS servers found in config")
	}
	return dnsServers, nil
}

func ReadDockerromFile(filename string) ([]string, error) {
	homeDir := GetHomeDir()
	if homeDir == "" {
		fmt.Println("HOME environment variable not set")
		os.Exit(1)
	}
	filename = AddPathToDir(homeDir, filename)
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	// Get DNS list
	registryList := viper.GetStringSlice("registryList")
	if len(registryList) == 0 {
		return nil, fmt.Errorf("no DNS servers found in config")
	}
	return registryList, nil
}

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
func GetHomeDir() string {
	if runtime.GOOS == WINDOWS_OS_NAME {
		return os.Getenv("USERPROFILE")
	} else {
		return os.Getenv("HOME")
	}
}
func GetTempDir() string {
	if runtime.GOOS == WINDOWS_OS_NAME {
		return os.Getenv("TEMP")
	} else {
		return "/tmp"
	}
}
func AddPathToDir(baseDir, extraDir string) string {
	if runtime.GOOS == WINDOWS_OS_NAME {
		return baseDir + "\\" + extraDir
	} else {
		return baseDir + "/" + extraDir
	}
}
