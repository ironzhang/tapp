#!/bin/bash

go_build() {
	local version=$1
	local git_commit=$2
	local build_time=$3
	local package=$4
	local target=$5

	if [ "$target" == "" ]; then
		go build -ldflags "-X $package.Version=$version -X $package.GitCommit=$git_commit -X $package.BuildTime=$build_time"
	else
		go build -o $target -ldflags "-X $package.Version=$version -X $package.GitCommit=$git_commit -X $package.BuildTime=$build_time"
	fi
}

version_build() {
	local version=$1
	local target=$2
	local package=main
	local git_commit=`git rev-parse HEAD`
	local build_time=`date '+%Y-%m-%d_%H:%M:%S'`

	go_build $version $git_commit $build_time $package $target
}

version_build v0.0.1 tapp
