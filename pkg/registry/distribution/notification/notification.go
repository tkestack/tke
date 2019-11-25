/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package notification

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"net/http"
	"regexp"
	"strings"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/registry/distribution/tenant"
	"tkestack.io/tke/pkg/util/log"
)

const Path = "/registry/notification"

const manifestPattern = `^application/vnd.docker.distribution.manifest.v\d\+(json|prettyjws)`

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type handler struct {
	manifestRegexp *regexp.Regexp
	registryClient *registryinternalclient.RegistryClient
}

func NewHandler(loopbackConfig *restclient.Config) (http.Handler, error) {
	re, err := regexp.Compile(manifestPattern)
	if err != nil {
		return nil, err
	}

	registryClient, err := registryinternalclient.NewForConfig(loopbackConfig)
	if err != nil {
		return nil, err
	}

	return &handler{
		manifestRegexp: re,
		registryClient: registryClient,
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer func() { _ = req.Body.Close() }()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Failed to read notification body from distribution", log.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var notification Notification
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Error("Failed to unmarshal notification body from distribution", log.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events, err := filterEvents(&notification, h.manifestRegexp)
	if err != nil {
		log.Errorf("Failed to filter events: %v", err)
		return
	}

	for _, event := range events {
		repository := event.Target.Repository
		tenantID, namespace, repoName := ParseRepository(repository)
		tag := event.Target.Tag
		digest := event.Target.Digest
		action := event.Action

		user := event.Actor.Name
		if len(user) == 0 {
			user = "anonymous"
		}

		if err := updateRepository(h.registryClient, tenantID, namespace, action, repoName, tag, digest); err != nil {
			log.Error("Failed to handler distribution notification event",
				log.String("tenantID", tenantID),
				log.String("namespace", namespace),
				log.String("repository", repository),
				log.String("tag", tag),
				log.String("action", action),
				log.String("user", user),
				log.Err(err))
		}
	}

	w.WriteHeader(http.StatusOK)
}

func updateRepository(registryClient *registryinternalclient.RegistryClient, tenantID, namespace, action, repoName, tag, digest string) error {
	namespaceList, err := registryClient.Namespaces().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", tenantID, namespace),
	})
	if err != nil {
		return err
	}
	if len(namespaceList.Items) == 0 {
		return fmt.Errorf("namespace %s in tenant %s not exist", namespace, tenantID)
	}
	namespaceObject := namespaceList.Items[0]
	repoList, err := registryClient.Repositories(namespaceObject.ObjectMeta.Name).List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s,spec.namespaceName=%s", tenantID, repoName, namespace),
	})
	if err != nil {
		return err
	}

	var repoObject *registry.Repository
	if len(repoList.Items) > 0 {
		repoObject = &repoList.Items[0]
	}

	switch action {
	case "push":
		return pushRepository(registryClient, &namespaceObject, repoObject, repoName, tag, digest)
	case "pull":
		return pullRepository(registryClient, &namespaceObject, repoObject, repoName, tag)
	}

	return fmt.Errorf("unknown action in distribution notification event handler")
}

func pushRepository(registryClient *registryinternalclient.RegistryClient, namespace *registry.Namespace, repository *registry.Repository, repoName, tag, digest string) error {
	needIncreaseRepoCount := false
	if repository == nil {
		needIncreaseRepoCount = true
		if _, err := registryClient.Repositories(namespace.ObjectMeta.Name).Create(&registry.Repository{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace.ObjectMeta.Name,
			},
			Spec: registry.RepositorySpec{
				Name:          repoName,
				TenantID:      namespace.Spec.TenantID,
				NamespaceName: namespace.Spec.Name,
				Visibility:    namespace.Spec.Visibility,
			},
			Status: registry.RepositoryStatus{
				PullCount: 0,
				Tags: []registry.RepositoryTag{
					{
						Name:        tag,
						Digest:      digest,
						TimeCreated: metav1.Now(),
					},
				},
			},
		}); err != nil {
			log.Error("Failed to create repository while received notification",
				log.String("tenantID", namespace.Spec.TenantID),
				log.String("namespace", namespace.Spec.Name),
				log.String("repo", repoName),
				log.String("tag", tag))
			return err
		}
	} else {
		existTag := false
		if len(repository.Status.Tags) == 0 {
			needIncreaseRepoCount = true
		} else {
			for k, v := range repository.Status.Tags {
				if v.Name == tag {
					existTag = true
					repository.Status.Tags[k] = registry.RepositoryTag{
						Name:        tag,
						Digest:      digest,
						TimeCreated: metav1.Now(),
					}
					if _, err := registryClient.Repositories(namespace.ObjectMeta.Name).UpdateStatus(repository); err != nil {
						log.Error("Failed to update repository tag while received notification",
							log.String("tenantID", namespace.Spec.TenantID),
							log.String("namespace", namespace.Spec.Name),
							log.String("repo", repoName),
							log.String("tag", tag))
						return err
					}
					break
				}
			}
		}

		if !existTag {
			repository.Status.Tags = append(repository.Status.Tags, registry.RepositoryTag{
				Name:        tag,
				Digest:      digest,
				TimeCreated: metav1.Now(),
			})
			if _, err := registryClient.Repositories(namespace.ObjectMeta.Name).UpdateStatus(repository); err != nil {
				log.Error("Failed to create repository tag while received notification",
					log.String("tenantID", namespace.Spec.TenantID),
					log.String("namespace", namespace.Spec.Name),
					log.String("repo", repoName),
					log.String("tag", tag))
				return err
			}
		}
	}

	if needIncreaseRepoCount {
		// update namespace repo count
		namespace.Status.RepoCount = namespace.Status.RepoCount + 1
		if _, err := registryClient.Namespaces().UpdateStatus(namespace); err != nil {
			log.Error("Failed to update namespace repo count while received notification",
				log.String("tenantID", namespace.Spec.TenantID),
				log.String("namespace", namespace.Spec.Name),
				log.String("repo", repoName),
				log.String("tag", tag))
			return err
		}
	}
	return nil
}

func pullRepository(registryClient *registryinternalclient.RegistryClient, namespace *registry.Namespace, repository *registry.Repository, repoName, tag string) error {
	if repository == nil {
		return fmt.Errorf("repository %s not exist", repoName)
	}
	repository.Status.PullCount = repository.Status.PullCount + 1
	if _, err := registryClient.Repositories(namespace.ObjectMeta.Name).UpdateStatus(repository); err != nil {
		log.Error("Failed to update repository pull count while received notification",
			log.String("tenantID", namespace.Spec.TenantID),
			log.String("namespace", namespace.Spec.Name),
			log.String("repo", repoName),
			log.String("tag", tag))
		return err
	}
	return nil
}

func filterEvents(notification *Notification, re *regexp.Regexp) ([]*Event, error) {
	events := make([]*Event, 0)

	for _, event := range notification.Events {
		log.Debug("Received a distribution event",
			log.String("id", event.ID),
			log.String("target", fmt.Sprintf("%s:%s", event.Target.Repository, event.Target.Tag)),
			log.String("digest", event.Target.Digest),
			log.String("action", event.Action),
			log.String("mediaType", event.Target.MediaType),
			log.String("userAgent", event.Request.UserAgent))

		if !re.MatchString(event.Target.MediaType) {
			continue
		}

		if strings.HasPrefix(event.Target.Repository, fmt.Sprintf("%s/", tenant.CrossTenantNamespace)) ||
			strings.HasPrefix(event.Target.Repository, fmt.Sprintf("/%s/", tenant.CrossTenantNamespace)) {
			log.Debugf("Ignore a library repo event", log.String("target", fmt.Sprintf("%s:%s", event.Target.Repository, event.Target.Tag)))
			continue
		}

		if checkEvent(&event) {
			events = append(events, &event)
			log.Debugf("Add event to collection: %s", event.ID)
			continue
		}
	}

	return events, nil
}

func checkEvent(event *Event) bool {
	// push action
	if event.Action == "push" {
		return true
	}
	// if it is pull action, check the user-agent
	userAgent := strings.ToLower(strings.TrimSpace(event.Request.UserAgent))
	return userAgent != "tke-registry-client"
}

// ParseRepository splits a repository into three parts: tenantID, namespace and rest
func ParseRepository(repository string) (tenantID, namespace, repo string) {
	repository = strings.TrimLeft(repository, "/")
	repository = strings.TrimRight(repository, "/")
	if !strings.ContainsRune(repository, '/') {
		repo = repository
		return
	}
	index := strings.Index(repository, "/")
	namespaceFull := repository[0:index]
	repo = repository[index+1:]
	namespaces := strings.SplitN(namespaceFull, "-", 2)
	if len(namespaces) == 2 {
		tenantID = namespaces[0]
		namespace = namespaces[1]
	} else {
		namespace = namespaceFull
	}
	return
}
