// Copyright 2025 novnc4svc Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package novnc4svc

import (
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/go-the-way/svc"
)

var (
	//go:embed templates/index.html
	indexFs string
	//go:embed scripts
	scriptsFs embed.FS
)

const (
	noVNCRoute   = "/no_vnc"
	wsRoute      = "/cloud_vnc"
	scriptsRoute = "/scripts"
)

type (
	Transformer func(id string) (wsUrl string)
	plugin      struct {
		noVNCRoute, wsRoute, scriptsRoute, idQueryName string

		transformer Transformer
	}
)

func DefaultPlugin(transformer Transformer) *plugin {
	return Plugin(noVNCRoute, wsRoute, scriptsRoute, "vnc_id", transformer)
}

func Plugin(noVNCRoute, wsRoute, scriptsRoute, idQueryName string, transformer Transformer) *plugin {
	return &plugin{
		noVNCRoute:   noVNCRoute,
		wsRoute:      wsRoute,
		scriptsRoute: scriptsRoute,
		idQueryName:  idQueryName,
		transformer:  transformer,
	}
}

func (p *plugin) Plug(engine *gin.Engine) {
	engine.StaticFS(p.scriptsRoute, svc.NewFSAdapter(scriptsRoute, scriptsFs))
	engine.GET(p.noVNCRoute, noVNC(p.wsRoute, p.scriptsRoute))
	engine.GET(p.wsRoute, wsVNC(p.idQueryName, p.transformer))
}
