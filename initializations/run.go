package initializations

func Run() {
	loadConfigs()

	initSQSQueue()
}
