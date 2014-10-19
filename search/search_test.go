/*
 * Copyright (c) 2013-2014, Jeremy Bingham (<jbingham@gmail.com>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package search

import (
	"encoding/gob"
	"fmt"
	"github.com/ctdk/goiardi/client"
	"github.com/ctdk/goiardi/databag"
	"github.com/ctdk/goiardi/environment"
	"github.com/ctdk/goiardi/node"
	"github.com/ctdk/goiardi/role"
	"github.com/ctdk/goiardi/organization"
	"testing"
	//"time"
)

// Most search testing can be handled fine with chef-pedant, but that's no
// reason to not have some go tests for it too.

var node1 *node.Node
var node2 *node.Node
var node3 *node.Node
var node4 *node.Node
var role1 *role.Role
var role2 *role.Role
var role3 *role.Role
var role4 *role.Role
var env1 *environment.ChefEnvironment
var env2 *environment.ChefEnvironment
var env3 *environment.ChefEnvironment
var env4 *environment.ChefEnvironment
var client1 *client.Client
var client2 *client.Client
var client3 *client.Client
var client4 *client.Client
var dbag1 *databag.DataBag
var dbag2 *databag.DataBag
var dbag3 *databag.DataBag
var dbag4 *databag.DataBag

var orgName = "default"
var org *organization.Organization

func orgInit() int {
	gob.Register(new(organization.Organization))
	var err error
	org, err = organization.New("default", "defaultyboo")
	if err != nil {
		panic(err)
	}
	err = org.Save()
	if err != nil {
		panic(err)
	}
	return 1
}

func makeSearchItems() int {
	/* Gotta populate the search index */
	nodes := make([]*node.Node, 4)
	roles := make([]*role.Role, 4)
	envs := make([]*environment.ChefEnvironment, 4)
	clients := make([]*client.Client, 4)
	dbags := make([]*databag.DataBag, 4)
	gob.Register(new(node.Node))
	gob.Register(new(role.Role))
	gob.Register(new(environment.ChefEnvironment))
	gob.Register(new(client.Client))
	gob.Register(new(databag.DataBag))

	for i := 0; i < 4; i++ {
		nodes[i], _ = node.New(org, fmt.Sprintf("node%d", i))
		nodes[i].Save()
		roles[i], _ = role.New(org, fmt.Sprintf("role%d", i))
		roles[i].Save()
		envs[i], _ = environment.New(org, fmt.Sprintf("env%d", i))
		envs[i].Save()
		clients[i], _ = client.New(org, fmt.Sprintf("client%d", i))
		clients[i].Save()
		dbags[i], _ = databag.New(org, fmt.Sprintf("databag%d", i))
		dbags[i].Save()
		dbi := make(map[string]interface{})
		dbi["id"] = fmt.Sprintf("dbi%d", i)
		dbi["foo"] = fmt.Sprintf("dbag_item_%d", i)
		dbags[i].NewDBItem(dbi)
	}
	node1 = nodes[0]
	node2 = nodes[1]
	node3 = nodes[2]
	node4 = nodes[3]
	role1 = roles[0]
	role2 = roles[1]
	role3 = roles[2]
	role4 = roles[3]
	env1 = envs[0]
	env2 = envs[1]
	env3 = envs[2]
	env4 = envs[3]
	client1 = clients[0]
	client2 = clients[1]
	client3 = clients[2]
	client4 = clients[3]
	dbag1 = dbags[0]
	dbag2 = dbags[1]
	dbag3 = dbags[2]
	dbag4 = dbags[3]

	/* Make this function return something so the compiler's happy building
	 * the tests. */
	return 1
}

var zz = orgInit()
var v = makeSearchItems()

func TestFoo(t *testing.T) {
	return
}

/* Only basic search tests are here. The stronger tests are handled in
 * chef-pedant, but these tests are meant to check basic search functionality.
 */

func TestSearchNode(t *testing.T) {
	n, _ := Search(org, "node", "name:node1")
	if n[0].(*node.Node).Name != "node1" {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchNodeAll(t *testing.T) {
	n, _ := Search(org, "node", "*:*")
	if len(n) != 4 {
		t.Errorf("Incorrect number of items returned, expected 4, got %d", len(n))
	}
}

func TestSearchRole(t *testing.T) {
	r, _ := Search(org, "role", "name:role1")
	if r[0].(*role.Role).Name != "role1" {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchRoleAll(t *testing.T) {
	n, _ := Search(org, "role", "*:*")
	if len(n) != 4 {
		t.Errorf("Incorrect number of items returned, expected 4, got %d", len(n))
	}
}

func TestSearchEnv(t *testing.T) {
	e, _ := Search(org, "environment", "name:env1")
	if e[0].(*environment.ChefEnvironment).Name != "env1" {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchEnvAll(t *testing.T) {
	n, _ := Search(org, "environment", "*:*")
	if len(n) != 4 {
		t.Errorf("Incorrect number of items returned, expected 4, got %d", len(n))
	}
}

func TestSearchClient(t *testing.T) {
	c, _ := Search(org, "client", "name:client1")
	if c[0].(*client.Client).Name != "client1" {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchClientAll(t *testing.T) {
	n, _ := Search(org, "client", "*:*")
	if len(n) != 4 {
		t.Errorf("Incorrect number of items returned, expected 4, got %d", len(n))
	}
}

func TestSearchDbag(t *testing.T) {
	d, _ := Search(org, "databag1", "foo:dbag_item_1")
	if len(d) == 0 {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchDbagAll(t *testing.T) {
	d, _ := Search(org, "databag1", "*:*")
	if len(d) != 1 {
		t.Errorf("Incorrect number of items returned, expected 1, got %d", len(d))
	}
}

func TestSecondOrg(t *testing.T) {
	sorgName := "boo"
	sorg, err := organization.New(sorgName, "booboo")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = sorg.Save()
	if err != nil {
		t.Errorf(err.Error())
	}
	snode, _ := node.New(sorg, "snode1")
	snode.Save()
	n, _ := Search(sorg, "node", "*:*")
	if len(n) != 1 {
		t.Errorf("Incorrect number of items returned, expected 1, got %d", len(n))
	}
	n, _ = Search(sorg, "node", "name:snode1")
	if len(n) != 1 {
		t.Errorf("Incorrect number of items returned with search by name, expected 1, got %d", len(n))
	}
	n, _ = Search(org, "node", "name:snode1")
	if len(n) != 0 {
		t.Errorf("searching the main test org for snode1 unexpectedly returned a result")
	}
	n, _ = Search(sorg, "node", "name:node1")
	if len(n) != 0 {
		t.Errorf("searching the second test org for node1 unexpectedly returned a result")
	}
}

// Probably don't want this as an always test, but it's handy to have available.
/*
func TestEmbiggenSearch(t *testing.T) {
	for i := 4; i < 35000; i++ {
		n, _ := node.New(org, fmt.Sprintf("node%d", i))
		n.Save()
		r, _ := role.New(org, fmt.Sprintf("role%d", i))
		r.Save()
		e, _ := environment.New(org, fmt.Sprintf("env%d", i))
		e.Save()
		c, _ := client.New(org, fmt.Sprintf("client%d", i))
		c.Save()
		d, _ := databag.New(org, fmt.Sprintf("databag%d", i))
		d.Save()
		dbi := make(map[string]interface{})
		dbi["id"] = fmt.Sprintf("dbi%d", i)
		dbi["foo"] = fmt.Sprintf("dbag_item_%d", i)
		d.NewDBItem(dbi)
	}
	time.Sleep(1 * time.Second)
	n, _ := Search(orgName, "client", "*:*")
	if len(n) != 35000 {
		t.Errorf("Incorrect number of items returned, expected 500, got %d", len(n))
	}
	c, _ := Search(orgName, "node", "*:*")
	if len(c) != 35000 {
		t.Errorf("Incorrect number of nodes returned, expected 500, got %d", len(n))
	}
	e, _ := Search(orgName, "environment", "name:env11666")
	if e[0].(*environment.ChefEnvironment).Name != "env11666" {
		t.Errorf("nothing returned from search")
	}
}
*/