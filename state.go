package main

import "msg_repo_cli/ssw"

var (
	devices  = []string{}
	messages = []ssw.Message{}
)

var (
	currentRecipient string
	currentDevice    string
)
