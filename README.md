# clickstream-dashboard-example
This project generates events that feed data to OpenSearch to demo some of Aiven's basic capabilities.

The idea is to simulate clickstream events performed by users of a website. I've chosen to pretend this is a news site that publishes articles for users to read and interact with.

For this POC, I am starting with messages formatted in JSON because it's familiar and easy to work with--both in code and Opensearch. Moving to a lower-overhead message format, like avro or protobuf, could be a future enhancement. The tradeoff is more efficient serialization and lower memory usage, for lower ease of use.

User-related event types are sparse. A true MVP will probably include new registrations, deletions, logout events, among others. A real site probably also allows unregistered users to view some or all of the site's articles as well. I've omitted these to avoid building out logic around when these events occur in relation to other events, as this doesn't provide much value for demoing this initial POC.

Secrets are stored locally in the conf directory, but not checked into version control. This is not ideal, but works for now. A future enhancement would be to move these to HashiCorp Vault using Aiven for Apache Kafka Connect.

I'm sure there's a way to easily pipe the data from the Kafka topic to OpenSearch, but I didn't find it. I was able to very easily pass log data from my Kafka topic to OpenSearch, but I had trouble subscribing and reading events. I instead wrote a consumer in this repo to move the data. That code can be found in cmd/consumer/consumer.go.