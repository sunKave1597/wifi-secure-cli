package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Scan WiFi and analyze security levels",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📡 Scan WiFi...")

		out, err := exec.Command("nmcli", "-f", "SSID,SIGNAL,SECURITY", "device", "wifi", "list").Output()
		if err != nil {
			fmt.Println("ERROR !!! :", err)
			return
		}

		lines := strings.Split(string(out), "\n")

		fmt.Printf("\n%-25s %-10s %-20s %-20s\n", "SSID", "SIGNAL", "ENCRYPTION", "SECURITY LEVEL")
		fmt.Println(strings.Repeat("-", 80))

		for i, line := range lines {
			if i == 0 || strings.TrimSpace(line) == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}

			ssid := fields[0]
			signal := fields[1]
			security := strings.Join(fields[2:], " ")

			level := assessSecurityLevel(security)

			fmt.Printf("%-25s %-10s %-20s %-20s\n", ssid, signal, security, level)
		}
	},
}

func assessSecurityLevel(enc string) string {
	enc = strings.ToUpper(enc)

	switch {
	case strings.Contains(enc, "WEP"):
		return "Very High"
	case strings.Contains(enc, "OPEN"):
		return "Have no pass!"
	case strings.Contains(enc, "WPA3"):
		return "Very Good !!! (WPA3)"
	case strings.Contains(enc, "WPA2") && strings.Contains(enc, "AES"):
		return "Good"
	case strings.Contains(enc, "WPA2"):
		return "Good but not AES"
	case strings.Contains(enc, "WPA"):
		return "It's old"
	default:
		return "Unknown"
	}
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
