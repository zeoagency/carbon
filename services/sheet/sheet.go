package sheet

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// ImportFileToGoogleSheets imports the given data to google sheets.
func ImportFileToGoogleSheets(file io.Reader) (string, error) {
	ctx := context.Background()

	// Get google-credentials from the env file.
	config, err := google.ConfigFromJSON([]byte(
		os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")),
		drive.DriveFileScope,
	)
	if err != nil {
		return "", fmt.Errorf("Unable to parse credential secret key to config.")
	}

	// Create client by using google-credentials and token(defined at the env file).
	tok := &oauth2.Token{}
	err = json.NewDecoder(strings.NewReader(os.Getenv("GOOGLE_DRIVE_TOKEN_JSON"))).Decode(tok)
	if err != nil {
		return "", fmt.Errorf("Unable to parse client secret key to config.")
	}
	client := config.Client(ctx, tok)

	// Create the service.
	srv, err := drive.New(client)
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve Drive client.")
	}

	// Import the excel file to Google Sheets.
	gFile, err := srv.Files.Create(&drive.File{
		MimeType:        "application/vnd.google-apps.spreadsheet",
		Name:            "result.xlsx",
		WritersCanShare: true,
	}).Media(file, googleapi.ContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")).Do()
	if err != nil {
		return "", fmt.Errorf("Error occur while creating file on Google Sheets.")
	}

	// Set permissions to view-only.
	_, err = srv.Permissions.Create(gFile.Id, &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}).Do()
	if err != nil {
		return "", fmt.Errorf("Error occur while creating file permission on Google Sheets.")
	}

	// Importing was completed.
	url := "https://docs.google.com/spreadsheets/d/" + gFile.Id
	return url, nil
}
