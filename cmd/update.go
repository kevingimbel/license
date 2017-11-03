package cmd

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"github.com/kevingimbel/license/lib"
	"github.com/spf13/cobra"
)


// getCmd represents the get command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Fetch the latest license data",
	Long: fmt.Sprintf(`Download the latest license data from http://osl.kevin.codes/licenses/

	The downlaoded data is stored inside the %s file.
	`, lib.GetOutputFilePath()),
	Run: func(cmd *cobra.Command, args []string) {
		update()
	},
}

func update() (int, []lib.OSLJSONFormat) {
	// Create http client and prepare the request
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", lib.GetUpdateURL(), nil)

	if err != nil {
		fmt.Println("Unable to make requet.\n", err)
		return 1, nil
	}

	req.Header.Set("User-Agent", "osl-cli")
	response, err := httpClient.Do(req)

	defer response.Body.Close()

	data := []lib.OSLJSONFormat{}
	json.NewDecoder(response.Body).Decode(&data)

	file, err := os.Create(lib.GetOutputFilePath())

	if err != nil {
		fmt.Println("Error creating file.\n", err)
		return 1, nil
	}

	err = json.NewEncoder(file).Encode(data)

	if err != nil {
		fmt.Println("Error writing to file.\n", err)
		return 1, nil
	}

	fmt.Println("Update done.")
	return 0, data
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
