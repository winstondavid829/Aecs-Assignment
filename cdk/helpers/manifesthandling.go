package helpers

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadKubernetesManifests(manifestsFolder string) ([]map[string]interface{}, error) {
	// Read all the manifest files from the folder.
	files, err := ioutil.ReadDir(manifestsFolder)
	if err != nil {
		return nil, err
	}

	// Initialize an empty array to store the manifest content.
	manifests := make([]map[string]interface{}, 0)

	// Loop through the files and read their content.
	for _, file := range files {
		if file.IsDir() {
			// Skip directories.
			continue
		}

		// Read the content of the manifest file.
		filePath := filepath.Join(manifestsFolder, file.Name())
		manifestContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		// Parse the manifest content into a map.
		var manifest map[string]interface{}
		if err := yaml.Unmarshal(manifestContent, &manifest); err != nil {
			return nil, err
		}

		// Append the parsed manifest to the array.
		manifests = append(manifests, manifest)
	}

	return manifests, nil
}
