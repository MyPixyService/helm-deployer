package deployer

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/MyPixyService/helm-deployer/pkg/utils"
)

func DeployReleases(inputFile string, kubeconfig string, releaseFilter string, verbose bool) error {
	// Read the YAML file
	data := utils.ReadFile(inputFile)

	// Parse the YAML data into a slice of Release structs
	releases := utils.ParseReleases(data)

	// Process each release
	for _, release := range releases {
		// Skip releases with a different name if a filter is specified
		if releaseFilter != "" && release.ReleaseName != releaseFilter {
			continue
		}

		// Skip releases that are disabled except if the filter is specified
		if !release.Enabled && releaseFilter == "" {
			continue
		}

		// Generate a temporary file to store the values
		tmpValuesFile, err := ioutil.TempFile("", "values-*.yaml")
		if err != nil {
			panic(err)
		}

		if release.ValuesFile != "" {
			// Read the values from the specified file
			valuesFileData := utils.ReadFile(release.ValuesFile)
			var values map[string]interface{}
			err = yaml.Unmarshal(valuesFileData, &values)
			if err != nil {
				panic(err)
			}

			utils.MergeMaps(release.Values, values)
		}

		// Write the values to the temporary file
		if release.Values != nil {
			valuesData, err := yaml.Marshal(release.Values)
			if err != nil {
				panic(err)
			}

			if verbose {
				fmt.Printf("Values for release %s:\n%s\n", release.ReleaseName, valuesData)
			}

			_, err = tmpValuesFile.Write(valuesData)
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
			release.ReleaseName, repoName, release.Chart, tmpValuesFile.Name())
		if kubeconfig != "" {
			cmd = fmt.Sprintf("%s --kubeconfig=%s", cmd, kubeconfig)
		}
		if release.Namespace != "" {
			cmd = fmt.Sprintf("%s --create-namespace --namespace=%s", cmd, release.Namespace)
		}
		if release.Version != "" {
			cmd = fmt.Sprintf("%s --version=%s", cmd, release.Version)
		}

		utils.RunCommand(cmd)

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
	data := utils.ReadFile(inputFile)

	// Parse the YAML data into a slice of Release structs
	releases := utils.ParseReleases(data)

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
		utils.RunCommand(cmd)

		fmt.Printf("Done.\n\n")
	}

	return nil
}

func addRepo(repoName, repoURL string) {
	cmd := fmt.Sprintf("helm repo add %s %s", repoName, repoURL)
	utils.RunCommand(cmd)
}

func removeRepo(repoName string) {
	cmd := fmt.Sprintf("helm repo remove %s", repoName)
	utils.RunCommand(cmd)
}
