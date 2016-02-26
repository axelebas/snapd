// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main_test

import (
	"fmt"
	"net/http"

	. "gopkg.in/check.v1"

	. "github.com/ubuntu-core/snappy/cmd/snap"
)

func (s *SnapSuite) TestConnectHelp(c *C) {
	msg := `Usage:
  snap.test [OPTIONS] connect <snap>:<plug> <snap>:<slot>

The connect command connects a plug to a slot.
It may be called in the following ways:

$ snap connect <snap>:<plug> <snap>:<slot>

Connects the specific plug to the specific slot.

$ snap connect <snap>:<plug> <snap>

Connects the specific plug to the only slot in the provided snap that matches
the connected interface. If more than one potential slot exists, the command
fails.

$ snap connect <plug> <snap>[:<slot>]

Without a name for the snap offering the plug, the plug name is looked at in
the gadget snap, the kernel snap, and then the os snap, in that order. The
first of these snaps that has a matching plug name is used and the command
proceeds as above.

Help Options:
  -h, --help               Show this help message

[connect command arguments]
  <snap>:<plug>
  <snap>:<slot>
`
	rest, err := Parser().ParseArgs([]string{"connect", "--help"})
	c.Assert(err.Error(), Equals, msg)
	c.Assert(rest, DeepEquals, []string{})
}

func (s *SnapSuite) TestConnectExplicitEverything(c *C) {
	s.RedirectClientToTestServer(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.Method, Equals, "POST")
		c.Check(r.URL.Path, Equals, "/2.0/skills")
		c.Check(DecodedRequestBody(c, r), DeepEquals, map[string]interface{}{
			"action": "grant",
			"skill": map[string]interface{}{
				"snap": "producer",
				"name": "plug",
			},
			"slot": map[string]interface{}{
				"snap": "consumer",
				"name": "slot",
			},
		})
		fmt.Fprintln(w, `{"type":"sync", "result":{}}`)
	})
	rest, err := Parser().ParseArgs([]string{"connect", "producer:plug", "consumer:slot"})
	c.Assert(err, IsNil)
	c.Assert(rest, DeepEquals, []string{})
}

func (s *SnapSuite) TestConnectExplicitPlugImplicitSlot(c *C) {
	s.RedirectClientToTestServer(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.Method, Equals, "POST")
		c.Check(r.URL.Path, Equals, "/2.0/skills")
		c.Check(DecodedRequestBody(c, r), DeepEquals, map[string]interface{}{
			"action": "grant",
			"skill": map[string]interface{}{
				"snap": "producer",
				"name": "plug",
			},
			"slot": map[string]interface{}{
				"snap": "consumer",
				"name": "",
			},
		})
		fmt.Fprintln(w, `{"type":"sync", "result":{}}`)
	})
	rest, err := Parser().ParseArgs([]string{"connect", "producer:plug", "consumer"})
	c.Assert(err, IsNil)
	c.Assert(rest, DeepEquals, []string{})
}

func (s *SnapSuite) TestConnectImplicitPlugExplicitSlot(c *C) {
	s.RedirectClientToTestServer(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.Method, Equals, "POST")
		c.Check(r.URL.Path, Equals, "/2.0/skills")
		c.Check(DecodedRequestBody(c, r), DeepEquals, map[string]interface{}{
			"action": "grant",
			"skill": map[string]interface{}{
				"snap": "",
				"name": "plug",
			},
			"slot": map[string]interface{}{
				"snap": "consumer",
				"name": "slot",
			},
		})
		fmt.Fprintln(w, `{"type":"sync", "result":{}}`)
	})
	rest, err := Parser().ParseArgs([]string{"connect", "plug", "consumer:slot"})
	c.Assert(err, IsNil)
	c.Assert(rest, DeepEquals, []string{})
}

func (s *SnapSuite) TestConnectImplicitPlugImplicitSlot(c *C) {
	s.RedirectClientToTestServer(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.Method, Equals, "POST")
		c.Check(r.URL.Path, Equals, "/2.0/skills")
		c.Check(DecodedRequestBody(c, r), DeepEquals, map[string]interface{}{
			"action": "grant",
			"skill": map[string]interface{}{
				"snap": "",
				"name": "plug",
			},
			"slot": map[string]interface{}{
				"snap": "consumer",
				"name": "",
			},
		})
		fmt.Fprintln(w, `{"type":"sync", "result":{}}`)
	})
	rest, err := Parser().ParseArgs([]string{"connect", "plug", "consumer"})
	c.Assert(err, IsNil)
	c.Assert(rest, DeepEquals, []string{})
}