/******************************************************************************
*
*  Copyright 2020 SAP SE
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
*
******************************************************************************/

package tasks

import (
	"net/http"
	"net/url"
	"time"

	"github.com/sapcc/keppel/internal/keppel"
	"github.com/sapcc/keppel/internal/processor"
)

//janitorDummyRequest can be put in the Request field of type keppel.AuditContext.
var janitorDummyRequest = &http.Request{URL: &url.URL{
	Scheme: "http",
	Host:   "localhost",
	Path:   "keppel-janitor",
}}

//Janitor contains the toolbox of the keppel-janitor process.
type Janitor struct {
	cfg     keppel.Configuration
	fd      keppel.FederationDriver
	sd      keppel.StorageDriver
	icd     keppel.InboundCacheDriver
	db      *keppel.DB
	auditor keppel.Auditor

	//non-pure functions that can be replaced by deterministic doubles for unit tests
	timeNow           func() time.Time
	generateStorageID func() string
}

//NewJanitor creates a new Janitor.
func NewJanitor(cfg keppel.Configuration, fd keppel.FederationDriver, sd keppel.StorageDriver, icd keppel.InboundCacheDriver, db *keppel.DB, auditor keppel.Auditor) *Janitor {
	j := &Janitor{cfg, fd, sd, icd, db, auditor, time.Now, keppel.GenerateStorageID}
	j.initializeCounters()
	return j
}

//OverrideTimeNow replaces time.Now with a test double.
func (j *Janitor) OverrideTimeNow(timeNow func() time.Time) *Janitor {
	j.timeNow = timeNow
	return j
}

//OverrideGenerateStorageID replaces keppel.GenerateStorageID with a test double.
func (j *Janitor) OverrideGenerateStorageID(generateStorageID func() string) *Janitor {
	j.generateStorageID = generateStorageID
	return j
}

func (j *Janitor) processor() *processor.Processor {
	return processor.New(j.cfg, j.db, j.sd, j.icd, j.auditor).OverrideTimeNow(j.timeNow).OverrideGenerateStorageID(j.generateStorageID)
}
