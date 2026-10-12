
## Overview:

The goal of this project was to implement a key-value store that uses **linearizability** as its consistency model. A key-value store allows clients to retrieve, insert, and modify data through a set of predefined operations. The behavior and guarantees of these operations depend on the consistency model enforced by the system.

The client can invoke one of the following operations:

- *GET* (key): Retrieve the value associated with a key.
- *PUT* (key, value): Insert a new key-value pair or replace the existing value of a key.
- *APPEND* (key, value): Append the provided value to the existing server-side value associated with a key.

## Consistency Models:

A consistency model defines the guarantees a system provides regarding the visibility and ordering of data operations. As data is added, modified, and retrieved within a distributed system, a consistency model establishes rules that determine when updates become visible and how concurrent operations are observed by clients.

There are numerous consistency models, each offering different guarantees and trade-offs depending on the system's requirements and use cases.

### Linearizability (Strong Consistency):

Linearizability is widely regarded as one of the strongest consistency models because it guarantees that every operation appears to execute atomically at a single point between its invocation and completion.

Additionally, linearizability preserves real-time ordering. If one operation completes before another begins, the second operation must observe a system state consistent with the first operation having already occurred.

It can be summarized as: **"Every operation appears to take effect instantaneously, and all operations respect real-time ordering."**

### Sequential Consistency:

TODO:

### Causal Consistency:

TODO:

### Eventual Consistency:

TODO:

## Requirements:

The goal of this lab was to implement a key-value store that provides strong consistency while operating under strict memory constraints. Additionally, client requests were expected to execute correctly in the presence of failures, including network interruptions, lost responses, and server crashes.

A critical requirement was ensuring that operations were executed **at most once**, even when clients retried requests due to network failures or missing responses. This was particularly important for non-idempotent operations such as *APPEND*, where executing the same request multiple times could result in an incorrect system state.

## Design 1:

My initial approach consisted of storing client metadata on the server side to identify and deduplicate incoming requests. This metadata consisted of a combination of the `clientID` and `operationID`.

The `operationID` was represented as a 64-bit integer that would be incremented atomically for each new client-side request. Using the operation name as an identifier was not viable because a client issuing two or more concurrent requests of the same type could introduce collisions and break the deduplication logic.

Additionally, I selected a 64-bit integer to significantly reduce the likelihood of integer overflow, even for clients executing a large number of operations over an extended period.

Example:

…

Although this solution worked, it required the server to maintain a record of every operation executed by each client. As the number of clients and requests increased, the amount of metadata stored on the server grew indefinitely, making this solution unsustainable under the lab's memory constraints.

## Design 2:

To improve upon the previous implementation, I introduced a server-side background job responsible for periodically removing stale client request results.

Each cached response was assigned a TTL (time to live), which was configured based on the maximum duration a client could retry a request. This duration could vary depending on the retry strategy implemented by the client.

Once a cached response exceeded its TTL, the background job would remove it from memory, reducing the amount of metadata retained by the server.

Although this solution improved memory utilization by periodically removing stale request responses, it was still not the most optimal implementation. It also introduced additional complexity, as the correctness of the deduplication mechanism depended on ensuring that clients would never retry a request after its cached response had expired.

## Design 3:

The optimal solution was to have the server store only the most recent request response for each client. This reduced the number of cached responses to at most *N*, where *N* represents the number of clients.

However, this design required an additional constraint: **each client could have only one outstanding request at a time.**

To enforce this constraint, client-side requests were serialized using synchronization mechanisms. This ensured that a client could not issue a new request until its previous request had completed, including any necessary retries.

By serializing requests, the server was guaranteed to observe each client's operations in a predictable order. This significantly simplified the deduplication logic because the server only needed to distinguish between a new request and a retry of the client's most recently completed request.

Furthermore, the server only needed to cache responses for the *APPEND* operation. This was because *GET* operations do not modify the system's state, while *PUT* operations are idempotent when applied repeatedly with the same key and value.

However, *APPEND* is not idempotent. Retrying an already executed *APPEND* operation could result in the same value being appended multiple times. Therefore, caching the most recent *APPEND* response ensured that duplicate requests returned the original result without modifying the stored value again.

This design significantly reduced the amount of client metadata and cached responses maintained on the server, improving its memory complexity from **O(R)**, where *R* represents the total number of historical requests, to **O(N)**, where *N* represents the number of clients.

By combining client-side request serialization with server-side deduplication, the final implementation provided a more memory-efficient solution while preserving the correctness guarantees required by the lab.
