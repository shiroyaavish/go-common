package rabbit

type TopicI interface {
	GetTopicName() string
	GetBody() []byte
}
