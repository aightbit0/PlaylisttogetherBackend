package memories

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func getFilesPaths(destinations []string) []byte {

	var files []ImagesResponse

	if len(destinations) == 0 {
		end, err := json.Marshal([]ImagesResponse{})
		if err != nil {
			fmt.Println(("Failed Marshal"))
		}

		return end
	}

	for i := 0; i < len(destinations); i++ {
		base64response, err := readFile(destinations[i])
		if err != nil {
			fmt.Println("Failed to read File")
		}
		p := ImagesResponse{
			Base64Code: base64response,
		}
		files = append(files, p)
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

	return base64Encoding, nil
}
