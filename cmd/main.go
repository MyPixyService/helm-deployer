package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MyPixyService/helm-deployer/pkg/deployer"
)

func main() {
	// Parse command-line flags
	var inputFile string
	flag.StringVar(&inputFile, "f", "", "Input YAML file containing Helm release definitions")

	var kubeconfig string
	flag.StringVar(&kubeconfig, "k", "~/.kube/config", "Path to kubeconfig file")

	var releaseFilter string
	flag.StringVar(&releaseFilter, "r", "", "Only deploy/uninstall releases with the specified name")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Verbose output")

	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "Show help message")

	var uninstall bool
	flag.BoolVar(&uninstall, "uninstall", false, "Uninstall Helm releases")

	// Set usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Show help message if -f flag is present
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Check if filename is empty
	if inputFile == "" {
		fmt.Println("Error: input file is required (use -f flag)")
		os.Exit(1)
	}

	if uninstall {
		// Uninstall Helm releases
		if err := deployer.UninstallReleases(inputFile, kubeconfig, releaseFilter); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	} else {
		// Deploy Helm releases
		if err := deployer.DeployReleases(inputFile, kubeconfig, releaseFilter, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	}
}
