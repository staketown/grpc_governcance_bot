package cosmos_governance_bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetIpfsData(ipfs string) *CascadiaDetails {
	if !strings.Contains(ipfs, "ipfs") {
		return nil
	}

	hash := strings.Split(ipfs, "://")[1]
	var details *CascadiaDetails
	response, err := http.Get("https://ipfs.io/ipfs/" + hash)

	if err != nil {
		fmt.Println(err)
		return &CascadiaDetails{
			Details: "",
			Title:   "",
		}
	}

	// close response body
	defer response.Body.Close()

	// read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &details)
	if err != nil {
		fmt.Println(err)
	}

	return details
}

type CascadiaDetails struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}
