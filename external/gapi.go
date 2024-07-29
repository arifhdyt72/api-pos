package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

func serviceAccount(credentialFile string) (*jwt.Config, error) {
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		return nil, err
	}
	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &c)
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/drive",
		},
		TokenURL: google.JWTTokenURL,
	}
	config.Subject = "application@pcsindonesia.co.id"
	return config, nil
	// token, err := config.TokenSource(oauth2.NoContext).Token()
	// if err != nil {
	// 	return nil, err
	// }
	// return token, nil
}

func GAuth(filename string, installType int) (string, error) {
	conf, err := serviceAccount("cred.json")
	if err != nil {
		fmt.Printf("Unable to retrieve Drive client: %v", err)
		return "", err
	}

	client := conf.Client(context.TODO())
	srv, err := drive.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Drive client: %v", err)
		return "", err
	}
	// year, month, _ := time.Now().Date()
	label := time.Now().Format("2006-01")
	r, err := srv.Files.List().
		Q("mimeType='application/vnd.google-apps.folder' and title = '" + label + "' and '" + os.Getenv("DIR_PARENT_ID") + "' in parents").
		MaxResults(1).
		Fields("nextPageToken, items(id, title)").
		Do()
	if err != nil {
		fmt.Printf("Unable to retrieve files: %v", err)
		return "", err
	}

	var folder *drive.File
	if len(r.Items) == 0 {
		rp, err := srv.Files.List().
			Q("mimeType='application/vnd.google-apps.folder' and title = '" + os.Getenv("DRIVE_DIR") + "'").
			MaxResults(1).
			Fields("nextPageToken, items(id, title)").
			Do()
		if err != nil || len(rp.Items) == 0 {
			fmt.Printf("Unable to retrieve files: %v", err)
			return "", err
		}
		currentDir := drive.ParentReference{
			Id: rp.Items[0].Id,
		}
		folder, _ = srv.Files.Insert(&drive.File{
			Title:    label,
			Parents:  []*drive.ParentReference{&currentDir},
			MimeType: "application/vnd.google-apps.folder",
		}).Do()
		// file = service.files().create(body=file_metadata, fields='id').execute()
		// rp.Items[0]
	} else {
		folder = r.Items[0]
	}
	var title string
	if installType == 1 {
		title = "Installation - New SIK"
	} else if installType == 2 {
		title = "Installation - Replacement"
	} else if installType == 3 {
		title = "Corrective Maintenance"
	} else if installType == 4 {
		title = "Preventive Maintenance"
	} else if installType == 5 {
		title = "Dismantle"
	}
	fmt.Println("Title: " + title)

	rp, err := srv.Files.List().
		Q("mimeType='application/vnd.google-apps.folder' and title = '" + title + "' and '" + folder.Id + "' in parents").
		MaxResults(1).
		Fields("nextPageToken, items(id, title)").
		Do()

	if err != nil || len(rp.Items) == 0 {
		currentDir := drive.ParentReference{
			Id: folder.Id,
		}
		folder, _ = srv.Files.Insert(&drive.File{
			Title:    title,
			Parents:  []*drive.ParentReference{&currentDir},
			MimeType: "application/vnd.google-apps.folder",
		}).Do()
	} else {
		folder = rp.Items[0]
	}

	parent := drive.ParentReference{
		Id: folder.Id,
	}
	file, _ := os.Open(filename)
	result, err := srv.Files.Insert(&drive.File{
		Title:   filename,
		Parents: []*drive.ParentReference{&parent},
	}).Media(file).Do()

	if err != nil {
		fmt.Printf("Unable to retrieve files: %v", err)
		return "", err
	}

	fmt.Println("done")
	return result.AlternateLink, nil
}
