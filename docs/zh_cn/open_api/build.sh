#!/bin/bash
# Copyright 2022 The BFE Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

export PATH=/home/work/.jumbo/bin/:$PATH 

npm install --registry http://registry.npm.baidu-int.com
if command gitbook > /dev/null; then
    echo 'gitbook existed, skip install'
else 
    npm install gitbook-cli
    npm install gitbook-plugin-splitter
    npm install gitbook-plugin-page-toc-button
    npm install gitbook-plugin-search-pro
    npm install gitbook-plugin-hide-element
fi

gitbook build

rm -rf output
mv _book output
rm output/*/*.md
rm -rf output/.git*
