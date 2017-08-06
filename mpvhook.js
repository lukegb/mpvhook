// Copyright 2017 Luke Granger-Brown. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

"use strict";

var matchRE = /^https:\/\/drive.google.com\/open\?(?:[^&]*&)*id=([^&]*)(?:&|$)/;

mp.add_hook("on_load", 5, function() {
  var url = mp.get_property("stream-open-filename");
  var m = url.match(matchRE);
  if (!m) {
    print("mpvhook skipping " + url);
    return;
  }
  var id = m[1];
  print("mpvhook running for " + url + " - id: " + id);

  var newURL = 'https://www.googleapis.com/drive/v3/files/' + id + '?alt=media';
  var tokenSP = mp.utils.subprocess({args: ["/home/lukegb/go/src/github.com/lukegb/mpvhook/mpvhook"]});
  if (tokenSP.status !== 0 || tokenSP.error != null) {
    print("mpvhook failed: " + tokenSP.status + " - " + tokenSP.error);
    return;
  }
  var token = tokenSP.stdout.replace(/\n*/g, '');

  mp.set_property("stream-open-filename", newURL);
  mp.set_property("file-local-options/http-header-fields", "Authorization: Bearer " + token + "");
});
