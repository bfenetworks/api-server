#!/bin/bash




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
