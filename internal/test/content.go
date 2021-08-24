/*******************************************************************************
*
* Copyright 2020 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/opencontainers/go-digest"
	"github.com/sapcc/go-bits/assert"
)

//Bytes groups a bytestring with its digest.
type Bytes struct {
	Contents  []byte
	Digest    digest.Digest
	MediaType string
}

//NewBytes makes a new Bytes instance.
func NewBytes(contents []byte) Bytes {
	return newBytesWithMediaType(contents, "application/octet-stream")
}

func newBytesWithMediaType(contents []byte, mediaType string) Bytes {
	return Bytes{contents, digest.Canonical.FromBytes(contents), mediaType}
}

//NewBytesFromFile creates a Bytes instance with the contents of the given file.
func NewBytesFromFile(path string) (Bytes, error) {
	buf, err := ioutil.ReadFile(path)
	return NewBytes(buf), err
}

//GenerateExampleLayer generates a blob of 1 MiB that can be used like an image
//layer when constructing image manifests for unit tests. The contents are
//generated deterministically from the given seed.
func GenerateExampleLayer(seed int64) Bytes {
	r := rand.New(rand.NewSource(seed))
	buf := make([]byte, 1<<20)
	r.Read(buf[:])
	return newBytesWithMediaType(buf, "application/vnd.docker.image.rootfs.diff.tar.gzip")
}

//Image contains all the pieces of a Docker image. The Layers and Config must
//be uploaded to the registry as blobs.
type Image struct {
	Layers   []Bytes
	Config   Bytes
	Manifest Bytes
}

var baseImageConfig = map[string]interface{}{
	"architecture": "amd64",
	"config": map[string]interface{}{
		"Hostname":     "",
		"Domainname":   "",
		"User":         "",
		"AttachStdin":  false,
		"AttachStdout": false,
		"AttachStderr": false,
		"Tty":          false,
		"OpenStdin":    false,
		"StdinOnce":    false,
		"Env": []string{
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"test_for=keppel",
		},
		"Cmd":        nil,
		"Image":      "",
		"Volumes":    nil,
		"WorkingDir": "",
		"Entrypoint": nil,
		"OnBuild":    nil,
		"Labels":     nil,
	},
	"container": "efd768c7229cf5030d391fb572f60cf4e22d5d85828fafb3aa5ff37997523c60",
	"container_config": map[string]interface{}{
		"Hostname":     "efd768c7229c",
		"Domainname":   "",
		"User":         "",
		"AttachStdin":  false,
		"AttachStdout": false,
		"AttachStderr": false,
		"Tty":          false,
		"OpenStdin":    false,
		"StdinOnce":    false,
		"Env": []string{
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"test_for=keppel",
		},
		"Cmd": []string{
			"/bin/sh",
			"-c",
			"#(nop) ",
			"ENV test_for=keppel",
		},
		"Image":      "",
		"Volumes":    nil,
		"WorkingDir": "",
		"Entrypoint": nil,
		"OnBuild":    nil,
		"Labels":     nil,
	},
	"created":        makeTimestamp(86400),
	"docker_version": "19.03.1-ce",
	"history": []map[string]interface{}{
		{
			"created":     makeTimestamp(0),
			"created_by":  "/bin/sh -c #(nop)  ENV test_for=keppel",
			"empty_layer": true,
		},
	},
	"os": "linux",
	"rootfs": map[string]interface{}{
		"type": "layers",
	},
}

//GenerateImage makes an Image from the given bytes in a deterministic manner.
func GenerateImage(layers ...Bytes) Image {
	//build image config referencing the given layers
	imageConfig := make(map[string]interface{})
	for k, v := range baseImageConfig {
		imageConfig[k] = v
	}
	history := []map[string]interface{}{imageConfig["history"].([]map[string]interface{})[0]}
	for idx, layer := range layers {
		history = append(history, map[string]interface{}{
			"created":    makeTimestamp(idx),
			"created_by": fmt.Sprintf("/bin/sh -c #(nop) ADD file:%s in / ", layer.Digest.String()),
		})
	}
	imageConfig["history"] = history
	imageConfigBytes, err := json.Marshal(imageConfig)
	if err != nil {
		panic(err.Error())
	}
	imageConfigBytesObj := newBytesWithMediaType(imageConfigBytes, "application/vnd.docker.container.image.v1+json")

	//build a manifest
	layerDescs := []map[string]interface{}{}
	for _, layer := range layers {
		layerDescs = append(layerDescs, map[string]interface{}{
			"mediaType": layer.MediaType,
			"size":      len(layer.Contents),
			"digest":    layer.Digest.String(),
		})
	}
	manifestData := map[string]interface{}{
		"schemaVersion": 2,
		"mediaType":     "application/vnd.docker.distribution.manifest.v2+json",
		"config": assert.JSONObject{
			"mediaType": imageConfigBytesObj.MediaType,
			"size":      len(imageConfigBytes),
			"digest":    imageConfigBytesObj.Digest.String(),
		},
		"layers": layerDescs,
	}
	manifestBytes, err := json.Marshal(manifestData)
	if err != nil {
		panic(err.Error())
	}

	return Image{
		Layers:   layers,
		Config:   imageConfigBytesObj,
		Manifest: newBytesWithMediaType(manifestBytes, "application/vnd.docker.distribution.manifest.v2+json"),
	}
}

//SizeBytes returns the value that we expect in the DB column
//`manifests.size_bytes` for this image.
func (i Image) SizeBytes() uint64 {
	imageSize := len(i.Manifest.Contents) + len(i.Config.Contents)
	for _, layer := range i.Layers {
		imageSize += len(layer.Contents)
	}
	return uint64(imageSize)
}

//ImageList contains all the pieces of a multi-architecture Docker image. This
//type is used for testing the behavior of Keppel with manifests that reference
//other manifests.
type ImageList struct {
	ImageManifests []Bytes
	Manifest       Bytes
}

//GenerateImageList makes an ImageList containing the given images in a
//deterministic manner.
func GenerateImageList(imageManifests ...Bytes) ImageList {
	manifestDescs := []map[string]interface{}{}
	testArchStrings := []string{"amd64", "arm", "arm64", "386", "ppc64le", "s390x"}
	for idx, m := range imageManifests {
		manifestDescs = append(manifestDescs, map[string]interface{}{
			"mediaType": m.MediaType,
			"size":      len(m.Contents),
			"digest":    m.Digest.String(),
			"platform": map[string]string{
				"os":           "linux",
				"architecture": testArchStrings[idx],
			},
		})
	}

	manifestListBytes, err := json.Marshal(map[string]interface{}{
		"schemaVersion": 2,
		"mediaType":     "application/vnd.docker.distribution.manifest.list.v2+json",
		"manifests":     manifestDescs,
	})
	if err != nil {
		panic(err.Error())
	}

	return ImageList{
		ImageManifests: imageManifests,
		Manifest:       newBytesWithMediaType(manifestListBytes, "application/vnd.docker.distribution.manifest.list.v2+json"),
	}
}

//SizeBytes returns the value that we expect in the DB column
//`manifests.size_bytes` for this image.
func (i ImageList) SizeBytes() uint64 {
	imageSize := len(i.Manifest.Contents)
	for _, m := range i.ImageManifests {
		imageSize += len(m.Contents)
	}
	return uint64(imageSize)
}

func makeTimestamp(seconds int) string {
	return time.Unix(int64(seconds), 0).UTC().Format(time.RFC3339Nano)
}
