package settings

type Config struct {
	Server         serverSetting
	ServiceSetting serviceSetting
}

type serviceSetting struct {
	KafkaSetting kafkaSetting `mapstructure:"kafka"`
	S3Setting    s3Setting    `mapstructure:"s3"`
	SQSSeting    sqsSetting   `mapstructure:"sqs"`
}

type serverSetting struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	SQSUrl     string `mapstructure:"SQS_QUEUE_URL"`
	AWSRegion  string `mapstructure:"AWS_REGION"`
}

type kafkaSetting struct {
	KafkaBroker_1 string `mapstructure:"kafka_broker_1"`
	KafkaBroker_2 string `mapstructure:"kafka_broker_2"`
	KafkaBroker_3 string `mapstructure:"kafka_broker_3"`
}

type s3Setting struct {
	AwsAccessKeyId       string `mapstructure:"aws_access_key_id"`
	AwsSercetAccessKeyId string `mapstructure:"aws_sercet_access_key_id"`
	AwsRegion            string `mapstructure:"aws_region"`
	FilmBucketName       string `mapstructure:"film_bucket_name"`
}

type sqsSetting struct {
	QueueURL string `mapstructure:"queue_url"`
}
