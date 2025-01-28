package youtube

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func Search(query string) ([][]string, error) {
	developerKey := os.Getenv("YOUTUBE_API_KEY")
	if developerKey == "" {
		return nil, fmt.Errorf("YouTube API key not set. Please set YOUTUBE_API_KEY environment variable")
	}

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
		Timeout:   10 * time.Second,
	}

	service, err := youtube.New(client)
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube client: %v", err)
	}

	call := service.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(25).
		Type("video")

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error performing search: %v", err)
	}

	var results [][]string
	for _, item := range response.Items {
		if item.Id.VideoId != "" {
			results = append(results, []string{
				item.Snippet.Title,
				item.Snippet.ChannelTitle,
				item.Id.VideoId,
			})
		}
	}

	return results, nil
}

func GetVideoURL(videoID string) string {
	return fmt.Sprintf("https://youtube.com/watch?v=%s", videoID)
}
