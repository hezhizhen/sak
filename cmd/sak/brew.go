package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/hezhizhen/sak/internal/color"
	"github.com/hezhizhen/sak/internal/log"
	"github.com/hezhizhen/sak/internal/types"
	"github.com/spf13/cobra"
)

func brewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brew <keyword>",
		Short: "Search Homebrew packages and display their info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBrewSearch(args[0])
		},
	}
	return cmd
}

func runBrewSearch(keyword string) error {
	log.Debug("Searching for: %s", keyword)

	cmd := exec.Command("brew", "search", keyword)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("brew search failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		log.Info("No results found")
		return nil
	}

	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "==>") {
			continue
		}
		packages = append(packages, line)
	}

	log.Debug("Found %d packages", len(packages))

	var results []types.PackageInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 10)

	for _, name := range packages {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			info := getPackageInfo(name)
			mu.Lock()
			results = append(results, info)
			mu.Unlock()
		}(name)
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	var formulaResults, caskResults, unknownResults []types.PackageInfo
	for _, r := range results {
		switch r.Type {
		case "formula":
			formulaResults = append(formulaResults, r)
		case "cask":
			caskResults = append(caskResults, r)
		default:
			unknownResults = append(unknownResults, r)
		}
	}

	if len(formulaResults) > 0 {
		fmt.Println(color.Green("==> Formulae"))
		printBrewResults(formulaResults)
	}

	if len(caskResults) > 0 {
		if len(formulaResults) > 0 {
			fmt.Println()
		}
		fmt.Println(color.Blue("==> Casks"))
		printBrewResults(caskResults)
	}

	if len(unknownResults) > 0 {
		if len(formulaResults) > 0 || len(caskResults) > 0 {
			fmt.Println()
		}
		fmt.Println(color.Gray("==> Unknown"))
		printBrewResults(unknownResults)
	}

	return nil
}

func printBrewResults(results []types.PackageInfo) {
	maxNameLen := 0
	maxVersionLen := 0
	for _, r := range results {
		if len(r.Name) > maxNameLen {
			maxNameLen = len(r.Name)
		}
		if len(r.Version) > maxVersionLen {
			maxVersionLen = len(r.Version)
		}
	}

	for _, r := range results {
		padding := maxVersionLen - len(r.Version)
		marker := "  "
		if r.Installed {
			marker = "\u2713 "
		}
		fmt.Printf("%s%-*s  %s%*s  %s\n",
			marker,
			maxNameLen, r.Name,
			color.Yellow(r.Version),
			padding, "",
			r.URL)
	}
}

func getPackageInfo(name string) types.PackageInfo {
	log.Debug("Getting info for: %s", name)

	cmd := exec.Command("brew", "info", "--json=v2", name)
	output, err := cmd.Output()
	if err != nil {
		log.Debug("brew info failed for %s: %v", name, err)
		return types.PackageInfo{Name: name, Version: "unknown", Type: "unknown", Failed: true}
	}

	var data struct {
		Formulae []struct {
			Homepage  string `json:"homepage"`
			Installed []struct {
				Version string `json:"version"`
			} `json:"installed"`
			Versions struct {
				Stable string `json:"stable"`
			} `json:"versions"`
		} `json:"formulae"`
		Casks []struct {
			Homepage  string `json:"homepage"`
			Version   string `json:"version"`
			Installed string `json:"installed"`
		} `json:"casks"`
	}

	if err := json.Unmarshal(output, &data); err != nil {
		log.Debug("parse JSON failed for %s: %v", name, err)
		return types.PackageInfo{Name: name, Version: "unknown", Type: "unknown", Failed: true}
	}

	if len(data.Formulae) > 0 && data.Formulae[0].Versions.Stable != "" {
		return types.PackageInfo{
			Name:      name,
			Version:   data.Formulae[0].Versions.Stable,
			URL:       data.Formulae[0].Homepage,
			Type:      "formula",
			Installed: len(data.Formulae[0].Installed) > 0,
		}
	}
	if len(data.Casks) > 0 && data.Casks[0].Version != "" {
		return types.PackageInfo{
			Name:      name,
			Version:   data.Casks[0].Version,
			URL:       data.Casks[0].Homepage,
			Type:      "cask",
			Installed: data.Casks[0].Installed != "",
		}
	}

	log.Debug("no version info found for %s", name)
	return types.PackageInfo{Name: name, Version: "unknown", Type: "unknown", Failed: true}
}
