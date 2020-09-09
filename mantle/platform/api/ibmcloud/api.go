// Copyright 2020 IBM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ibmcloud

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/coreos/mantle/auth"
	"github.com/coreos/mantle/platform"
	"github.com/coreos/pkg/capnslog"
	"github.com/coreos/pkg/multierror"

	"github.com/ibm/vpc-go-sdk/"
)

var plog = capnslog.NewPackageLogger("github.com/coreos/mantle", "platform/api/ibmcloud")

type Options struct {
	*platform.Options
	// The ibmcloud region regional api calls should use
	Region string

	// Config file. Defaults to ~/.ibmcloud/config.json
	ConfigPath string
	// The profile to use when resolving credentials, if applicable
	Profile string

	// AccessKeyID is the optional access key to use. It will override all other sources
	AccessKeyID string
	// SecretKey is the optional secret key to use. It will override all other sources
	SecretKey string
}

type API struct {
	ecs  *ecs.Client
	oss  *oss.Client
	opts *Options
}

// New creates a new ibmcloud API wrapper. It uses credentials from any of the
// standard credentials sources, including the environment and the profile
// configured in ~/.ibmcloud.
func New(opts *Options) (*API, error) {
	return nil, nil
}

// CopyImage replicates an image to a new region
func (a *API) CopyImage(source_id, dest_name, dest_region, dest_description, kms_key_id string, encrypted bool) (string, error) {
	return "1", nil
}

// ImportImage attempts to import an image from COS returning the image_id & error
//
// NOTE: this function will re-use existing images that share the same final name
// if the name is not unique then provide force to pre-remove any images with the
// specified name
func (a *API) ImportImage(format, bucket, object, image_size, device, name, description, architecture string, force bool) (string, error) {
	return "1", nil
}

// Wait for the import image task and return the image id. See also similar
// code in AWS' finishSnapshotTask.
func (a *API) finishImportImageTask(importImageResponse *ecs.ImportImageResponse) (string, error) {
	return "1", nil
}

// GetImages retrieves a list of images by ImageName
func (a *API) GetImages(name string) (*ecs.DescribeImagesResponse, error) {
	return nil
}

// GetImagesByID retrieves a list of images by ImageId
func (a *API) GetImagesByID(id string) (*ecs.DescribeImagesResponse, error) {
	return nil
}

// DeleteImage deletes an image and it's underlying snapshots
func (a *API) DeleteImage(id string, force bool) error {
	return nil
}

// DeleteSnapshot deletes a snapshot
func (a *API) DeleteSnapshot(id string, force bool) error {
	return nil
}

// UploadFile is a multipart upload, use for larger files
//
// NOTE: this function will return early if an object already exists
// at the specified path, if it might not be unique provide the force
// option to skip these checks
func (a *API) UploadFile(filepath, bucket, path string, force bool) error {
	return nil
}

// DeleteFile deletes a file from a COS bucket
func (a *API) DeleteFile(bucket, path string) error {
	return nil
}

// PutObject performs a singlepart upload into a COS bucket
func (a *API) PutObject(r io.Reader, bucket, path string, force bool) error {
	return nil
}

// ListRegions lists the enabled regions in ibmcloud implicitly
// by the Profile and Region options.
func (a *API) ListRegions() ([]string, error) {
	return nil, nil
}
