# 🛡️ Error Handling Middleware for Replicated Data Libraries

This project presents a **service-oriented middleware** that significantly **improves error handling** in replicated data libraries by tracing the **distributed effects of erroneous updates across replicas** and applying **novel integration techniques** — all **without requiring changes** to existing library code.

## Overview

Replicated data libraries are essential for building resilient distributed systems, but managing the **propagation and impact of errors** across replicas remains a critical challenge.

This middleware addresses that problem by:
- **Intercepting** and **analyzing** library-level operations at runtime.
- **Tracing the causal effects** of erroneous updates across the system.
- Enabling **automated recovery strategies** without modifying library internals.
- Supporting a **multi-language environment** (Go, JavaScript, C#, C++).

## Key Features

- 🛠️ **Non-Intrusive Integration**: Works with libraries without requiring source code modifications.
- 🔎 **Distributed Effect Tracing**: Captures the impact of errors across replicas for smarter diagnostics and recovery.
- 📈 **Enhanced Error Handling**: Improves fault tolerance by **25%** through early detection and response.
- 🌐 **Cross-Language Support**: Compatible with applications built in Go, JavaScript, C#, and C++.
- 🔄 **Pluggable Recovery Modules**: Extendable architecture for adding custom recovery and compensation strategies.

## Technologies Used

- **Go**, **JavaScript**, **C#**, **C++** — supported client libraries.
- **Custom Middleware Layer** — for operation interception and tracing.
- **Distributed Tracing Techniques** — for causal analysis.

## Motivation

Traditional error handling mechanisms often overlook the **system-wide consequences** of faulty updates in replicated environments. Detecting and correcting these issues early is vital for ensuring system correctness and availability.

By **separating error tracing from library code** and **orchestrating recovery externally**, this middleware improves reliability without complicating library maintenance or development.

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/surzoprovakar/Error-Handling-Middleware.git
   cd Error-Handling-Middleware


2. Run the project

##### For rKV-Store (C#)
- Find the Makefile in *RAC/* and command *make server*
- Find the Makefile in *RAC/src/Operations/* and command *make*
- Find the Makefile in *RACClient/src/* and command *make client*

##### For Legion (JavaScript)
- Find the Makefile and open two CLIs
- Command *make prod* and *make object*
- Open *https://localhost* in the browser

##### For FeatureStorm (Go)
- Find the Makefile for a CRDT and the command *make*
- Inside the Library directory, open 3 CLIs 
- Command *./r1.sh*, *./r2.sh*, and *./r3.sh*
