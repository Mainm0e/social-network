package utils

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
createDirectory creates a directory with the given path if error occurs returns an error.
*/
func createDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory:", err)
		return errors.New("Error creating directory: " + err.Error())
	}
	return nil
}

/*
InitiateImagesPath creates the images directory and the subdirectories for the avatar, posts and comments images (if it does not exist)
using createDirectory function if error occurs returns an error.
*/
func InitiateImagesPath() error {
	directories := []string{
		"./images/avatars/",
		"./images/posts/",
		"./images/comments/",
	}

	for _, dir := range directories {
		err := createDirectory(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
ProcessImage decodes the Base64-encoded image data based on the data type prefix
and saves it to the desired location if error occurs returns an error
*/
func ProcessImage(data string, url string) (string, error) {
	// Remove the data type prefix if present
	parts := strings.SplitN(data, ",", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid avatar image data")
	}

	// Extract the image data after the comma
	imageData := parts[1]

	// Determine the image type based on the data type prefix
	imageType := ""
	if strings.HasPrefix(data, "data:image/jpeg") {
		imageType = "jpeg"
	} else if strings.HasPrefix(data, "data:image/png") {
		imageType = "png"
	} else if strings.HasPrefix(data, "data:image/gif") {
		imageType = "gif"
	} else {
		return "", errors.New("unsupported image type")
	}
	// Decode the Base64-encoded image data
	decodedData, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", errors.New("Error decoding avatar image  " + err.Error())
	}
	// Save the image to the desired location
	err = ioutil.WriteFile(url+"."+imageType, decodedData, 0644)
	if err != nil {
		return "", errors.New("Error saving avatar image  " + err.Error())
	}
	log.Println("Avatar image saved successfully")
	return url + "." + imageType, nil
}

/*
RetrieveImage retrieves the avatar image as Base64-encoded string from the local system
based on the provided file path.
If an error occurs, it returns an empty string and the error.
*/
func RetrieveImage(filePath string) (string, error) {
	// Read the image file as bytes
	imageData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading avatar image:", err)
		return "", errors.New("Error reading avatar image: " + err.Error())
	}

	// Determine the image type based on the file extension
	imageType := ""
	fileExtension := getFileExtension(filePath)
	if fileExtension == "jpg" || fileExtension == "jpeg" {
		imageType = "jpeg"
	} else if fileExtension == "png" {
		imageType = "png"
	} else if fileExtension == "gif" {
		imageType = "gif"
	} else {
		return "", errors.New("unsupported image type")
	}

	// Encode the image data to Base64
	encodedData := base64.StdEncoding.EncodeToString(imageData)

	// Create the data URI for the encoded image
	dataURI := "data:image/" + imageType + ";base64," + encodedData

	return dataURI, nil
}

// getFileExtension retrieves the file extension from a given file path
func getFileExtension(filePath string) string {
	parts := strings.Split(filePath, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}
