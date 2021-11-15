package vaultengine

import (
	"fmt"
	"strings"
)

// Folder defines a level of the Vault structure
type Folder map[string]interface{}

// FolderExport will export all subfolders and secrets from a specified location
// if readSecrets is false, the secret values would be replaced with an empty string
func (client *Client) FolderExport(path string, readSecrets bool) (Folder, error) {
	baseFolder := make(Folder)
	subFolders := make(Folder)

	err := client.PathReader(&subFolders, path, readSecrets)
	if err != nil {
		return nil, err
	}

	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")

	buildFolderStructure(&baseFolder, parts, subFolders)

	return baseFolder, nil
}

// buildFolderStructure creates the base tree structure
func buildFolderStructure(parentFolder *Folder, parts []string, subFolders Folder) error {
	nextPart := parts[0]
	parts = parts[1:]
	newSubFolder := make(Folder)

	if len(parts) == 0 {
		// If we are at the root level we overwrite the rootfolder with it's subfolder
		// so that we don't get empty keys in our export
		if nextPart == "" {
			*parentFolder = subFolders
		} else {
			(*parentFolder)[nextPart] = subFolders
		}

	} else {
		buildFolderStructure(&newSubFolder, parts, subFolders)
		(*parentFolder)[nextPart] = newSubFolder
	}

	return nil
}

//PathReader recursively reads the provided path and all subpaths
func (client *Client) PathReader(parentFolder *Folder, path string, readSecrets bool) error {
	folder, err := client.FolderRead(path)
	if err != nil {
		return err
	}

	for _, key := range folder {
		strKey := fmt.Sprintf("%v", key)
		newPath := path + strKey

		if IsFolder(strKey) {
			subFolder := make(Folder)
			keyName := strings.Replace(strKey, "/", "", -1)

			err = client.PathReader(&subFolder, newPath, readSecrets)
			if err != nil {
				return err
			}

			if (*parentFolder)[keyName] != nil {
				for key, elem := range (*parentFolder)[keyName].(map[string]interface{}) {
					subFolder[key] = elem
				}
			}
			(*parentFolder)[keyName] = subFolder
		} else {
			s := client.SecretRead(newPath, readSecrets)
			(*parentFolder)[strKey] = s
		}
	}

	return nil
}
