# clickstream-dashboard-example
An example clickstream event dashboard using Aiven Platform

The idea here is to simulate clickstream events performed by users of a website. I've chosen to pretend this is a news site that publishes articles for users to read and interact with.

For this POC, I am starting with messages formatted in JSON because it's familiar and easy to work with--both in code and Opensearch. Moving to a lower-overhead message format, like avro or protobuf, could be a future enhancement. The tradeoff is more efficient serialization and lower memory usage, for lower ease of use.

User-related event types are sparse. A true MVP will probably include new registrations, deletions, logout events, among others. A real site probably also allows unregistered users to view some or all of the site's articles as well. I've omitted these to avoid building out logic around when these events occur in relation to other events, as this doesn't provide much value for demoing this inital POC.
