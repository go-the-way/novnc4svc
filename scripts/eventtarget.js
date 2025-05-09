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

export default class EventTargetMixin {
    constructor() {
        this._listeners = null;
    }

   addEventListener(type, callback) {
      if (!this._listeners) {
         this._listeners = new Map();
      }
      if (!this._listeners.has(type)) {
         this._listeners.set(type, new Set());
      }
      this._listeners.get(type).add(callback);
   }

   removeEventListener(type, callback) {
      if (!this._listeners || !this._listeners.has(type)) {
         return;
      }
      this._listeners.get(type).delete(callback);
   }

   dispatchEvent(event) {
      if (!this._listeners || !this._listeners.has(event.type)) {
         return true;
      }
      this._listeners.get(event.type).forEach((callback) => {
         callback.call(this, event);
      }, this);
      return !event.defaultPrevented;
   }
}
