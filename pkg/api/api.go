/*******************************************************************************
*
* Copyright 2018 SAP SE
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

package api

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sapcc/go-bits/gopherpolicy"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/respondwith"
	"github.com/sapcc/keppel/pkg/database"
	"github.com/sapcc/keppel/pkg/openstack"
	gorp "gopkg.in/gorp.v2"
)

//KeppelV1 implements the /keppel/v1/ API endpoints.
type KeppelV1 struct {
	db *gorp.DbMap
	su *openstack.ServiceUser
	tv gopherpolicy.Validator
}

//NewKeppelV1 prepares a new KeppelV1 instance.
func NewKeppelV1(db *gorp.DbMap, su *openstack.ServiceUser) (*KeppelV1, error) {
	tv := gopherpolicy.TokenValidator{
		IdentityV3: su.IdentityV3,
	}
	policyPath := os.Getenv("KEPPEL_POLICY_PATH")
	if policyPath == "" {
		return nil, errors.New("missing env variable: KEPPEL_POLICY_PATH")
	}
	err := tv.LoadPolicyFile(policyPath)
	if err != nil {
		return nil, err
	}

	return &KeppelV1{
		db: db,
		su: su,
		tv: &tv,
	}, nil
}

//Router prepares a http.Handler
func (api *KeppelV1) Router() http.Handler {
	r := mux.NewRouter()

	//NOTE: Keppel account names appear in Swift container names, so they may not
	//contain any slashes.
	r.Methods("GET").Path("/keppel/v1/accounts").HandlerFunc(api.handleGetAccounts)
	r.Methods("GET").Path("/keppel/v1/accounts/{account:[^/]+}").HandlerFunc(api.handleGetAccount)
	r.Methods("PUT").Path("/keppel/v1/accounts/{account:[^/]+}").HandlerFunc(api.handlePutAccount)

	return r
}

func (api *KeppelV1) checkToken(r *http.Request) *gopherpolicy.Token {
	token := api.tv.CheckToken(r)
	token.Context.Logger = logg.Debug
	token.Context.Request = mux.Vars(r)
	return token
}

func (api *KeppelV1) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	token := api.checkToken(r)
	if !token.Require(w, "account:list") {
		return
	}

	var accounts []database.Account
	_, err := api.db.Select(&accounts, "SELECT * FROM accounts ORDER BY name")
	if respondwith.ErrorText(w, err) {
		return
	}

	//restrict accounts to those visible in the current scope
	var accountsFiltered []database.Account
	for _, account := range accounts {
		token.Context.Request["account_project_id"] = account.ProjectUUID
		if token.Check("account:show") {
			accountsFiltered = append(accountsFiltered, account)
		}
	}
	//ensure that this serializes as a list, not as null
	if len(accountsFiltered) == 0 {
		accountsFiltered = []database.Account{}
	}

	respondwith.JSON(w, http.StatusOK, map[string]interface{}{"accounts": accountsFiltered})
}

func (api *KeppelV1) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	token := api.checkToken(r)

	//first very permissive check: can this user GET any accounts AT ALL?
	token.Context.Request["account_project_id"] = token.Context.Auth["project_id"]
	if !token.Require(w, "account:show") {
		return
	}

	//get account from DB to find its project ID
	accountName := mux.Vars(r)["account"]
	account, err := api.findAccount(accountName)
	if respondwith.ErrorText(w, err) {
		return
	}

	//perform final authorization with that project ID
	if account != nil {
		token.Context.Request["account_project_id"] = account.ProjectUUID
		if !token.Check("account:show") {
			account = nil
		}
	}

	if account == nil {
		http.Error(w, "no such account", 404)
		return
	}

	respondwith.JSON(w, http.StatusOK, map[string]interface{}{"account": account})
}

func (api *KeppelV1) handlePutAccount(w http.ResponseWriter, r *http.Request) {
	//TODO
	w.Write([]byte("put account " + mux.Vars(r)["account"]))
}

func (api *KeppelV1) findAccount(name string) (*database.Account, error) {
	var account database.Account
	err := api.db.SelectOne(&account,
		"SELECT * FROM accounts WHERE name = $1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &account, err
}
