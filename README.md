
## Quick Start

To run the project:

```bash
make run
```
This command will start the server and client containers using Docker Compose. There are 3 clients send commands concurrently for demonstrational purposes. 

## Help

```shell
Available commands:
  deps   - Install dependencies
  build  - Build the project (requires deps)
  clean  - Clean up containers and output directory
  run    - Run the project using docker-compose
  tidy   - Tidy up the go.mod file
  help   - Show this help message
```

## Overview
The system consists of a server that processes commands (add, delete, get, getAll) on an ordered map data structure. Clients send commands through RabbitMQ, and the server writes results to an output file.

## Main Trade-offs
As with many systems, the main performance bottleneck is likely to be disk I/O, as the server needs to write to the same output file.
Reading from the ordered map is concurrent, thanks to the use of a read-write mutex (`RWMutex`). This allows multiple read operations to occur simultaneously, improving performance.

### Parallelism vs. Consistency

Middleground consistency chosen over maximum parallelism. Full map lock used during write and delete operations, rather than fine-grained locking or atomic operations. 
This approach ensures the map's state remains consistent, particularly during operations like `getAllItems` which may take longer to execute than a `delete` operation.

While the current implementation is suitable for moderate workloads, scaling to handle very large datasets or high concurrency might require a distributed approach. For example:

- Sharding the data across multiple nodes
- Replicating the ordered map for read scalability

However, maintaining consistency across distributed nodes would introduce significant complexity. The current design prioritizes simplicity and correctness for moderate-scale use cases.

