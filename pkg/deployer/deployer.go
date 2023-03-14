package deployer

import (
    "fmt"
    "io/ioutil"
	"os/exec"

    "gopkg.in/yaml.v3"
)

type Release struct {
    ReleaseName string `yaml:"releaseName"`
    RepoURL     string `yaml:"repoURL"`
    Chart       string `yaml:"chart"`
    Version     string `yaml:"version"`
    Namespace   string `yaml:"namespace"`
    Values      map[string]interface{} `yaml:"values,omitempty"`
}

func DeployReleases(inputFile string, kubeconfig string, releaseFilter string)  error {
	// Read the YAML file
    data := readFile(inputFile)

    // Parse the YAML data into a slice of Release structs
    releases := parseReleases(data)

    // Process each release
    for _, release := range releases {
		// Skip releases with a different name if a filter is specified
		if releaseFilter != "" && release.ReleaseName != releaseFilter {
			continue
		}

        // Generate a temporary file to store the values
        valuesFile, err := ioutil.TempFile("", "values-*.yaml")
        if err != nil {
            panic(err)
        }

        // Write the values to the temporary file
        if release.Values != nil {
			valuesData, err := yaml.Marshal(release.Values)
			if err != nil {
				panic(err)
			}
			_, err = valuesFile.Write(valuesData)
			if err != nil {
				panic(err)
			}
		}

		fmt.Printf("Deploying release %s from %s\n",
			release.ReleaseName, release.RepoURL)

		repoName := "tmprepo"

		// Add the repo
		addRepo(repoName, release.RepoURL)

        // Run the Helm command to upgrade or install the release
        cmd := fmt.Sprintf("helm upgrade --install %s %s/%s -f %s",
            release.ReleaseName, repoName, release.Chart, valuesFile.Name())
		if kubeconfig != "" {
			cmd = fmt.Sprintf("%s --kubeconfig=%s", cmd, kubeconfig)
		}
		if release.Namespace != "" {
			cmd = fmt.Sprintf("%s --namespace=%s", cmd, release.Namespace)
		}
		if release.Version != "" {
			cmd = fmt.Sprintf("%s --version=%s", cmd, release.Version)
		}
		
		runCommand(cmd)

		statusCmd := fmt.Sprintf("helm status %s", release.ReleaseName)
		if kubeconfig != "" {
			statusCmd = fmt.Sprintf("%s --kubeconfig=%s", statusCmd, kubeconfig)
		}
		if release.Namespace != "" {
			statusCmd = fmt.Sprintf("%s --namespace=%s", statusCmd, release.Namespace)
		}

		fmt.Printf("Done. Run\n\t%s\nto see the chart's deployment.\n\n", statusCmd)

		// Remove the repo
		removeRepo(repoName)
	}

	return nil
}

func UninstallReleases(inputFile string, kubeconfig string, releaseFilter string) error {
	// Read the YAML file
    data := readFile(inputFile)

    // Parse the YAML data into a slice of Release structs
    releases := parseReleases(data)

    // Process each release
    for _, release := range releases {
		// Skip releases with a different name if a filter is specified
		if releaseFilter != "" && release.ReleaseName != releaseFilter {
			continue
		}

		fmt.Printf("Uninstalling release %s\n", release.ReleaseName)

		// Run the Helm command to uninstall the release
		cmd := fmt.Sprintf("helm uninstall %s", release.ReleaseName)
		if kubeconfig != "" {
			cmd = fmt.Sprintf("%s --kubeconfig=%s", cmd, kubeconfig)
		}
		if release.Namespace != "" {
			cmd = fmt.Sprintf("%s --namespace=%s", cmd, release.Namespace)
		}
		runCommand(cmd)

		fmt.Printf("Done.\n\n")
	}

	return nil
}

func readFile(inputFile string) []byte {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	return data
}

func parseReleases(data []byte) []Release {
	var releases []Release
    err := yaml.Unmarshal(data, &releases)
    if err != nil {
        panic(err)
    }
	return releases
}

func runCommand(cmd string) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		fmt.Printf("Command output: %s\n", out)
		panic(err)
	}
	//fmt.Printf("Command output: %s\n", string(out))
}

func addRepo(repoName, repoURL string) {
	cmd := fmt.Sprintf("helm repo add %s %s", repoName, repoURL)
	runCommand(cmd)
}

func removeRepo(repoName string) {
	cmd := fmt.Sprintf("helm repo remove %s", repoName)
	runCommand(cmd)
}
