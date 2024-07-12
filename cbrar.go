package main

import (
	"context"
	"encoding/csv"
	"flag"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type BulkRedirect struct {
	Source    string
	ApiToken  string
	AccountId string
}

type TargetLists map[string][]BulkRedirect

func main() {
	host, dataFileName := parseCmdFlags()
	dataFile, err := os.Open(dataFileName)
	if err != nil {
		log.Fatalf("failed to open data file: %v", err)
	}
	defer dataFile.Close()
	targetLists, err := ReadAllTargetList(dataFile)
	if err != nil {
		log.Fatalf("failed to read data file: %v", err)
	}
	for targetList, redirects := range targetLists {
		err := CreateRedirect(context.Background(), host, targetList, redirects)
		if err != nil {
			log.Fatalf("failed to create redirect: %v", err)
		} else {
			for _, redirect := range redirects {
				log.Printf("redirect successfully created: %v", redirect.Source)
			}
		}
	}
}

func parseCmdFlags() (host, dataFlieName string) {
	flag.StringVar(&host, "host", "", "host to redirect at (Eg: https://example.com)")
	flag.StringVar(&dataFlieName, "dfile", "data.csv", "data file name")
	flag.Parse()
	if host == "" {
		flag.Usage()
		os.Exit(1)
	}
	return
}

func ReadAllTargetList(dataFile *os.File) (TargetLists, error) {
	r := csv.NewReader(dataFile)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	targetLists := make(TargetLists)
	for _, record := range records {
		targetLists[record[1]] = append(targetLists[record[1]], BulkRedirect{
			Source:    record[0],
			ApiToken:  record[2],
			AccountId: record[3],
		})
	}
	return targetLists, nil
}

func CreateRedirect(ctx context.Context, target_host, target_list string, redirects []BulkRedirect) error {
	api, err := cloudflare.NewWithAPIToken(redirects[0].ApiToken)
	if err != nil {
		return err
	}
	if len(redirects) < 1 {
		return nil
	}
	redirectCode := 302
	trueVar := true
	listItems := make([]cloudflare.ListItemCreateRequest, len(redirects))
	for i, redirect := range redirects {
		listItems[i] = cloudflare.ListItemCreateRequest{
			Redirect: &cloudflare.Redirect{
				SourceUrl:  redirect.Source,
				TargetUrl:  target_host,
				StatusCode: &redirectCode,
				PreserveQueryString: &trueVar,
				SubpathMatching:   &trueVar,
				PreservePathSuffix: &trueVar,
			},
		}
	}
	_, err = api.ReplaceListItems(context.TODO(), cloudflare.AccountIdentifier(redirects[0].AccountId), cloudflare.ListReplaceItemsParams{
		ID:    target_list,
		Items: listItems,
	})
	return err
}
