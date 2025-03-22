# Clickstream Dashboard Example
This project simulates clickstream events and streams that data to OpenSearch using Aiven for Apache Kafka. It serves as a proof of concept (POC) to demonstrate some of Aiven's core capabilities for real-time data streaming, stream processing, and data visualization.

---
### Project Overview
The idea is to simulate clickstream events performed by users of a fictional news website. For this POC, I model these events as JSON objects and publish them to a Kafka topic. A consumer reads from this topic and indexes each event into OpenSearch, where the data can be explored in OpenSearch Dashboards.

---
### Design Choices
- **Domain:** A news site with rudimentary users, and articles with categories like politics, tech, and sports.
- **Message format:** Events are simple go structs, serialized as JSON for seamless OpenSearch indexing. This format makes the data human-readable and works well with OpenSearch's default JSON-based ingestion pipeline. A future enhancement could adopt a more efficient format like Avro or Protobuf to reduce memory usage, depending upon the customer's usage and needs.
- **Event Modeling:**
  - All users must login before performing other actions.
  - Other event types, such as like and comment, represent these users interacting with articles.
  - Additional user/admin event types (e.g. register and logout) are left out to avoid spending time writing sequencing logic that doesn't add much to this POC, but could be added if they help demonstrate value for a particular customer's needs.
  - Currently, article IDs are simply randomly generated GUIDs. This could be enhanced by randomly generating new article IDs and tracking them for reuse, similarly to how user IDs are tracked and reused.

---
### Secrets Configuration
Secrets are stored locally in the ```conf/``` directory, the contents of which are included in the project's ```.gitignore``` file. The `opensearch.json` file contains keys like `host`, `username`, `password`, and `index`. To productionalize this code, these secrets should be moved to a secure store, such as HashiCorp Vault. Aiven's platform could be utilized by integrating HashiCorp Vault with [Aiven for Kafka Connect](https://aiven.io/docs/products/kafka/kafka-connect/howto/configure-secret-providers).

---
### Directory Structure

```
cmd/
├── producer/      # Generates and sends randomized clickstream events to Kafka
├── consumer/      # Reads Kafka events and indexes them into OpenSearch
conf/              # SSL certificates and OpenSearch credentials (included in .gitignore)
models/            # Event schema and event generation logic
```

---
### How to Run
1. Set up your Aiven Kafka and OpenSearch services in [Aiven Console](https://console.aiven.io).
2. Place your SSL certs and OpenSearch credentials in ```conf/```.
3. Navigate to the directory where you cloned this repo.
4. In one terminal:
    ```shell
    go run cmd/producer/main.go
    ```
5. In another terminal
    ```shell
    go run cmd/consumer/main.go
    ```
You'll see structured logs of events being produced and indexed in each terminal after the processes start. View the data and create visualizations in OpenSearch Dashboards.

---
### Suggested Visualizations

- **Pie Chart:** Distribution of `event_type` (e.g., login vs. like)
- **Line Graph:** Event volume over time (using the `timestamp` field)
- **Bar Chart:** Top `user_id` values by activity count
- **Heatmap or Table:** Device type by article category
