/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package app

import (
	log "github.com/Sirupsen/logrus"
	"errors"
	"fmt"
	"encoding/json"
)

var (
	ErrInvalidRequesttype = errors.New("invalid requesttype")
)

type InvalidProcessError struct {
	Msg string
}

func (e InvalidProcessError) Error() string {
	return fmt.Sprintf("process error: %s", e.Msg)
}

// Request represents the job for an unit in megam.
type Request string

func (r Request) String() string {
	return string(r)
}

func ParseRequest(req string) (Request, error) {
	switch req {
	case "build":
		return ReqBuild, nil
	case "building":
		return ReqBuilding, nil
	case "built":
		return ReqBuilt, nil
	case "create":
		return ReqCreate, nil
	case "creating":
		return ReqCreating, nil
	case "stateup":
		return ReqStateup, nil
	case "statedown":
		return ReqStatedown, nil
	case "created":
		return ReqCreated, nil
	case "delete":
		return ReqDelete, nil
	case "deleting":
		return ReqDeleting, nil
	case "deleted":
		return ReqDeleted, nil
	case "error":
		return ReqError, nil
	case "start":
		return ReqStart, nil
	case "starting":
		return ReqStarting, nil
	case "started":
		return ReqStarted, nil
	case "stop":
		return ReqStop, nil
	case "stoping":
		return ReqStoping, nil
	case "stopped":
		return ReqStopped, nil
	}
	return Request(""), ErrInvalidRequesttype
}

const (

    ReqBuild = Request("build")
    ReqBuilt = Request("built")
    ReqCreate = Request("create")
    ReqStateup = Request("stateup")
    ReqStatedown = Request("statedown")
    ReqCreated = Request("created")
    ReqDelete = Request("delete")
    ReqDeleting = Request("deleting")
    ReqDeleted = Request("deleted")
    ReqStart = Request("start")
    ReqStop = Request("stop")
    ReqStoping = Request("stoping")

	// ReqCreating is the initial status of a unit in the database,
	// it should transition shortly to a more specific status
	ReqCreating = Request("create")

	// ReqBuilding is the status for units being provisioned by the
	// provisioner, like in the deployment.
	ReqBuilding = Request("building")

	// ReqError is the status for units that failed to start, because of
	// an application error.
	ReqError = Request("error")

	// StatusStarting is set when the container is started in docker.
	ReqStarting = Request("starting")

	// StatusStarted is for cases where the unit is up and running, and bound
	// to the proper status, it's set by RegisterUnit and SetUnitStatus.
	ReqStarted = Request("started")

	// StatusStopped is for cases where the unit has been stopped.
	ReqStopped = Request("stopped")
)

type Requests struct {
	Id	             string     `json:"id"`
	AppId            string     `json:"cat_id"`
	AppName          string     `json:"name"`
	Action           string     `json:"cattype"`
    Category         string     `json:"category"`
	CreatedAt        string     `json:"created_at"`
}

/**
**fetch the request json from riak and parse the json to struct
**/
func (p *Payload) Convert() (*Requests, error) {
	log.Info("Get request %s", p.Id)
	r := &Requests{}
	riakUrl := "192.168.1.247:8087"

   conn, cerr := RiakConnection(riakUrl, "requests")
	if cerr != nil {
		return r, cerr
	}
	err := conn.FetchStruct(p.Id, r)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return r, nil
}

type Payload struct {
	Id string `json:"id"`
}

type PayloadConvertor interface {
	Convert(Payload) (*Requests, error)
}

func NewPayload(b []byte) (*Payload, error) {
	p := &Payload{}
	err := json.Unmarshal(b, &p)
	if err != nil {
		log.Error("Failed to parse the payload message:\n%s.", err)
		return nil, err
	}
	return p, err
}