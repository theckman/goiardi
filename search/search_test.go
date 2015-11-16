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
	"testing"
	"time"

	"github.com/theckman/goiardi/client"
	"github.com/theckman/goiardi/config"
	"github.com/theckman/goiardi/databag"
	"github.com/theckman/goiardi/environment"
	"github.com/theckman/goiardi/indexer"
	"github.com/theckman/goiardi/node"
	"github.com/theckman/goiardi/role"
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

func makeSearchItems() int {
	indexer.Initialize(config.Config)
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

	// circleci is sometimes weird about the index having everything. This
	// *never* comes up locally. ??? Could possibly be because the indexer
	// hasn't had a chance to finish indexing?
	reindexObjs := make([]indexer.Indexable, 0, 4*5)
	for i := 0; i < 4; i++ {
		nodes[i], _ = node.New(fmt.Sprintf("node%d", i))
		nodes[i].Default["baz"] = fmt.Sprintf("borb")
		nodes[i].Default["blurg"] = fmt.Sprintf("b%d", i)
		nodes[i].Save()
		roles[i], _ = role.New(fmt.Sprintf("role%d", i))
		roles[i].Save()
		envs[i], _ = environment.New(fmt.Sprintf("env%d", i))
		envs[i].Save()
		clients[i], _ = client.New(fmt.Sprintf("client%d", i))
		clients[i].Save()
		dbags[i], _ = databag.New(fmt.Sprintf("databag%d", i))
		dbags[i].Save()
		dbi := make(map[string]interface{})
		dbi["id"] = fmt.Sprintf("dbi%d", i)
		dbi["foo"] = fmt.Sprintf("dbag_item_%d", i)
		dbags[i].NewDBItem(dbi)
		reindexObjs = append(reindexObjs, nodes[i])
		reindexObjs = append(reindexObjs, roles[i])
		reindexObjs = append(reindexObjs, envs[i])
		reindexObjs = append(reindexObjs, clients[i])
		dbis, _ := dbags[i].AllDBItems()
		for _, d := range dbis {
			reindexObjs = append(reindexObjs, d)
		}
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

	// Let the indexing functions catch up. This has not been a problem in
	// The Real World™ (famous last words), but it's *definitely* a problem
	// when running go test with GOMAXPROCS > 1.
	time.Sleep(1 * time.Second)

	/* Make this function return something so the compiler's happy building
	 * the tests. */
	return 1
}

var v = makeSearchItems()

var searcher = &TrieSearch{}

func TestFoo(t *testing.T) {
	return
}

/* Only basic search tests are here. The stronger tests are handled in
 * chef-pedant, but these tests are meant to check basic search functionality.
 */

func TestSearchNode(t *testing.T) {
	n, _ := searcher.Search("node", "name:node1", 1000, "id ASC", 0, nil)
	if len(n) == 0 || n[0]["name"] != "node1" {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchNodeAll(t *testing.T) {
	n, _ := searcher.Search("node", "*:*", 1000, "id ASC", 0, nil)
	if len(n) != 4 {
		t.Errorf("Incorrect number of items returned, expected 4, got %d :: %v", len(n), n)
	}
}

func TestSearchNodeFalse(t *testing.T) {
	n, _ := searcher.Search("node", "foo:bar AND NOT foo:bar", 1000, "id ASC", 0, nil)
	if len(n) != 0 {
		t.Errorf("Incorrect number of items returned, expected 0, got %d", len(n))
	}
}

func TestSearchNodeAttr(t *testing.T) {
	n, _ := searcher.Search("node", "name:node1 AND NOT baz:urb", 1000, "id ASC", 0, nil)
	if len(n) != 1 {
		t.Errorf("Incorrect number of items returned, expected 1, got %d", len(n))
	}
}

func TestSearchNodeAttrExists(t *testing.T) {
	n, _ := searcher.Search("node", "name:node1 AND NOT baz:borb", 1000, "id ASC", 0, nil)
	if len(n) != 0 {
		t.Errorf("Incorrect number of items returned, expected 0, got %d", len(n))
	}
}

func TestSearchNodeAttrAndExists(t *testing.T) {
	n, _ := searcher.Search("node", "name:node1 AND baz:borb", 1000, "id ASC", 0, nil)
	if len(n) != 1 {
		t.Errorf("Incorrect number of items returned, expected 1, got %d", len(n))
	}
}

func TestSearchNodeAttrAndNotExists(t *testing.T) {
	n, _ := searcher.Search("node", "name:node1 AND baz:urb", 1000, "id ASC", 0, nil)
	if len(n) != 0 {
		t.Errorf("Incorrect number of items returned, expected 0, got %d", len(n))
	}
}

func TestSearchRole(t *testing.T) {
	r, _ := searcher.Search("role", "name:role1", 1000, "id ASC", 0, nil)
	if len(r) == 0 || r[0]["name"] != "role1" {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchRoleAll(t *testing.T) {
	n, _ := searcher.Search("role", "*:*", 1000, "id ASC", 0, nil)
	if len(n) != 4 {
		t.Errorf("Incorrect number of items returned, expected 4, got %d", len(n))
	}
}

func TestSearchEnv(t *testing.T) {
	e, _ := searcher.Search("environment", "name:env1", 1000, "id ASC", 0, nil)
	if len(e) == 0 || e[0]["name"] != "env1" {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchEnvAll(t *testing.T) {
	n, _ := searcher.Search("environment", "*:*", 1000, "id ASC", 0, nil)
	if len(n) != 4 {
		t.Errorf("Incorrect number of items returned, expected 4, got %d", len(n))
	}
}

func TestSearchClient(t *testing.T) {
	c, _ := searcher.Search("client", "name:client1", 1000, "id ASC", 0, nil)
	if len(c) == 0 || c[0]["name"] != "client1" {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchClientAll(t *testing.T) {
	n, _ := searcher.Search("client", "*:*", 1000, "id ASC", 0, nil)
	if len(n) != 4 {
		t.Errorf("Incorrect number of items returned, expected 4, got %d", len(n))
	}
}

func TestSearchDbag(t *testing.T) {
	d, _ := searcher.Search("databag1", "foo:dbag_item_1", 1000, "id ASC", 0, nil)
	if len(d) == 0 {
		t.Errorf("nothing returned from search")
	}
}

func TestSearchDbagAll(t *testing.T) {
	d, _ := searcher.Search("databag1", "*:*", 1000, "id ASC", 0, nil)
	if len(d) != 1 {
		t.Errorf("Incorrect number of items returned, expected 1, got %d", len(d))
	}
}

// Probably don't want this as an always test, but it's handy to have available.
/*
func TestEmbiggenSearch(t *testing.T) {
	for i := 4; i < 35000; i++ {
		n, _ := node.New(fmt.Sprintf("node%d", i))
		n.Save()
		r, _ := role.New(fmt.Sprintf("role%d", i))
		r.Save()
		e, _ := environment.New(fmt.Sprintf("env%d", i))
		e.Save()
		c, _ := client.New(fmt.Sprintf("client%d", i))
		c.Save()
		d, _ := databag.New(fmt.Sprintf("databag%d", i))
		d.Save()
		dbi := make(map[string]interface{})
		dbi["id"] = fmt.Sprintf("dbi%d", i)
		dbi["foo"] = fmt.Sprintf("dbag_item_%d", i)
		d.NewDBItem(dbi)
	}
	time.Sleep(1 * time.Second)
	n, _ := searcher.Search("client", "*:*", 1000, "id ASC", 0, nil)
	if len(n) != 35000 {
		t.Errorf("Incorrect number of items returned, expected 500, got %d", len(n))
	}
	c, _ := searcher.Search("node", "*:*", 1000, "id ASC", 0, nil)
	if len(c) != 35000 {
		t.Errorf("Incorrect number of nodes returned, expected 500, got %d", len(n))
	}
	e, _ := searcher.Search("environment", "name:env11666", 1000, "id ASC", 0, nil)
	if e[0].(*environment.ChefEnvironment).Name != "env11666" {
		t.Errorf("nothing returned from search")
	}
}
*/
