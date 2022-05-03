package memories

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"playlisttogether/backend/utils"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func getFilesPaths(destination string, foldername string) []byte {
	var files []ImagesResponse
	md5folder := utils.GetMD5Hash(foldername)
	folderInfo, err1 := os.Stat(destination + md5folder)

	if os.IsNotExist(err1) {
		fmt.Println("folder not exists")
		fmt.Println(folderInfo)
		end, err := json.Marshal(files)
		if err != nil {
			fmt.Println(("Failed Marshal"))

		}
		return end
	}

	err := filepath.Walk(destination+md5folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			base64response, err := readFile(path)
			if err != nil {
				fmt.Println("Failed to read File")
				return nil
			}

			p := ImagesResponse{
				Base64Code: base64response,
			}
			files = append(files, p)
		}

		return nil
	})

	if err != nil {
		fmt.Println("fail")
	}

	end, err := json.Marshal(files)
	if err != nil {
		fmt.Println(("Failed Marshal"))
	}

	return end
}

func readFile(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error Reading File")
		return "", err
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += toBase64(bytes)

	//fmt.Println(base64Encoding)

	return base64Encoding, nil
}
