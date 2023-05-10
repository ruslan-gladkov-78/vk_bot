package main

func StartBot() {
	sd := GetLongPollSessionData(groupID)
	AccessLongPoll(sd)
}

func main() {
	StartBot()
}
