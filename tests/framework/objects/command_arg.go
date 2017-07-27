/*
Copyright 2016 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package objects

// CommandArgs is a warpper for cmd args
type CommandArgs struct {
	Command             string
	SubCommand          string
	CmdArgs             []string
	OptionalArgs        []string
	PipeToStdIn         string
	EnvironmentVariable []string
}

//CommandOut is a wrapper for cmd out returned after executing command args
type CommandOut struct {
	StdOut   string
	StdErr   string
	ExitCode int
	Err      error
}
