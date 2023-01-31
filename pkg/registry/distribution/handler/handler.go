package handler

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	distributionClient "tkestack.io/tke/pkg/registry/distribution/client"
	"tkestack.io/tke/pkg/util/log"
)

func DeleteRepo(ctx context.Context, client *distributionClient.Repository, userName, tenantID, repoName string) (err error) {
	tags, err := client.ListTag(repoName, userName, tenantID)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var errsMap sync.Map
	for _, tag := range tags {
		wg.Add(1)
		go func(repoName, tag, user, tenantID string, rc *distributionClient.Repository) {
			defer wg.Done()
			image := fmt.Sprintf("%s:%s", repoName, tag)
			log.Infof("delete label of image at database: %s", image)
			if err := rc.DeleteTag(repoName, tag, user, tenantID); err != nil {
				if regErr, ok := err.(*distributionClient.Error); ok {
					if regErr.Code == http.StatusNotFound {
						return
					}
				}
				errsMap.Store(tag, err.Error())
				log.Errorf("failed to delete %s: %v", image, err)
				return
			}
			log.Infof("delete tag at registry: %s", image)

		}(repoName, tag, userName, tenantID, client)
	}
	wg.Wait()
	errMessages := []string{}
	for _, tag := range tags {
		message, ok := errsMap.Load(tag)
		if !ok {
			continue
		}
		errMessages = append(errMessages, fmt.Sprintf("failed to delete %s tag %s: %s", repoName, tag, message.(string)))
	}
	if len(errMessages) > 0 {
		return fmt.Errorf("failed to delete repo: %v", errMessages)
	}
	return nil

}
