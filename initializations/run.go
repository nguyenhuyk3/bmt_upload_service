package initializations

import messagebroker "bmt_upload_service/initializations/message_broker"

func Run() {
	loadConfigs()

	messagebroker.InitReaders()

	select {}
}
