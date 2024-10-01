/*
Copyright The Rotate Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package google

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/internal/environment"
	"github.com/raniellyferreira/rotate-files/pkg/providers"
	"github.com/raniellyferreira/rotate-files/pkg/utils"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GoogleProvider struct {
	client *storage.Client
}

func NewGoogleProvider() (*GoogleProvider, error) {
	credsFile := environment.GetEnv("GOOGLE_APPLICATION_CREDENTIALS", "")
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(credsFile))
	if err != nil {
		return nil, err
	}
	return &GoogleProvider{client: client}, nil
}

func (g *GoogleProvider) Delete(fullPath string) error {
	bucket, path := utils.GetBucketAndKey(fullPath)
	obj := g.client.Bucket(bucket).Object(path)
	return obj.Delete(context.Background())
}

func (g *GoogleProvider) ListFiles(fullPath string) ([]*providers.BackupInfo, error) {
	bucket, prefix := utils.GetBucketAndKey(fullPath)
	it := g.client.Bucket(bucket).Objects(context.Background(), &storage.Query{Prefix: prefix})

	var files []*providers.BackupInfo
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		files = append(files, &providers.BackupInfo{
			Path:      fmt.Sprintf("gs://%s/%s", bucket, objAttrs.Name),
			Size:      objAttrs.Size,
			Timestamp: carbon.CreateFromTimestamp(objAttrs.Created.Unix()),
		})
	}

	return files, nil
}
