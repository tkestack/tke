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

package util

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/util/log"
)

func PushRepository(ctx context.Context, registryClient *registryinternalclient.RegistryClient, namespace *registry.Namespace, repository *registry.Repository, repoName, tag, digest string) error {
	needIncreaseRepoCount := false
	if repository == nil {
		needIncreaseRepoCount = true
		if _, err := registryClient.Repositories(namespace.ObjectMeta.Name).Create(ctx, &registry.Repository{
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
		}, metav1.CreateOptions{}); err != nil {
			log.Error("Failed to create repository while received notification",
				log.String("tenantID", namespace.Spec.TenantID),
				log.String("namespace", namespace.Spec.Name),
				log.String("repo", repoName),
				log.String("tag", tag),
				log.Err(err))
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
					if _, err := registryClient.Repositories(namespace.ObjectMeta.Name).UpdateStatus(ctx, repository, metav1.UpdateOptions{}); err != nil {
						log.Error("Failed to update repository tag while received notification",
							log.String("tenantID", namespace.Spec.TenantID),
							log.String("namespace", namespace.Spec.Name),
							log.String("repo", repoName),
							log.String("tag", tag),
							log.Err(err))
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
			if _, err := registryClient.Repositories(namespace.ObjectMeta.Name).UpdateStatus(ctx, repository, metav1.UpdateOptions{}); err != nil {
				log.Error("Failed to create repository tag while received notification",
					log.String("tenantID", namespace.Spec.TenantID),
					log.String("namespace", namespace.Spec.Name),
					log.String("repo", repoName),
					log.String("tag", tag),
					log.Err(err))
				return err
			}
		}
	}

	if needIncreaseRepoCount {
		// update namespace repo count
		namespace.Status.RepoCount = namespace.Status.RepoCount + 1
		if _, err := registryClient.Namespaces().UpdateStatus(ctx, namespace, metav1.UpdateOptions{}); err != nil {
			log.Error("Failed to update namespace repo count while received notification",
				log.String("tenantID", namespace.Spec.TenantID),
				log.String("namespace", namespace.Spec.Name),
				log.String("repo", repoName),
				log.String("tag", tag),
				log.Err(err))
			return err
		}
	}
	return nil
}

func PullRepository(ctx context.Context, registryClient *registryinternalclient.RegistryClient, namespace *registry.Namespace, repository *registry.Repository, repoName, tag string) error {
	if repository == nil {
		return fmt.Errorf("repository %s not exist", repoName)
	}
	repository.Status.PullCount = repository.Status.PullCount + 1
	if _, err := registryClient.Repositories(namespace.ObjectMeta.Name).UpdateStatus(ctx, repository, metav1.UpdateOptions{}); err != nil {
		log.Error("Failed to update repository pull count while received notification",
			log.String("tenantID", namespace.Spec.TenantID),
			log.String("namespace", namespace.Spec.Name),
			log.String("repo", repoName),
			log.String("tag", tag),
			log.Err(err))
		return err
	}
	return nil
}
