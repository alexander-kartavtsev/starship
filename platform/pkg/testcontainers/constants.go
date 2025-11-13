package testcontainers

// MongoDB constants
const (
	// MongoDB container constants
	MongoContainerName = "mongo-test"

	// MongoDB environment variables
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoDatabaseKey  = "MONGO_INITDB_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
)
