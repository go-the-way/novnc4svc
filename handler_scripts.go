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
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	scriptsCache = map[string]string{}
	mu           = &sync.RWMutex{}
)

func scripts(prefix, scriptsRoute string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		scriptName := ctx.Param("script")
		fileName := prefix + scriptName
		mu.RLock()
		js, ok := scriptsCache[fileName]
		mu.RUnlock()
		if !ok {
			buf, err := scriptsFs.ReadFile(fileName)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, errors.New("file not found"))
				return
			}
			js = string(buf)
			mu.Lock()
			scriptsCache[fileName] = js
			mu.Unlock()
		}
		ctx.Header("Content-Type", "text/javascript; charset=UTF-8")
		js0 := strings.ReplaceAll(js, "@scriptsRoute@", scriptsRoute)
		ctx.String(http.StatusOK, js0)
	}
}
