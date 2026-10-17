This prototype addresses the critical need for robust agent-to-agent communication and collaboration in multi-agent systems, specifically focusing on enabling one agent to invoke another agent's capabilities as a tool. Current approaches often rely on simple message passing or custom, ad-hoc mechanisms, which lack discoverability, type safety, and a standardized invocation pattern.

**The Problem:** Building complex agentic AI systems requires modularity. Agents should be able to specialize in certain tasks and expose their functionalities to other agents. This goes beyond just sending messages; it's about one agent *calling* a function or capability provided by another agent, much like a function call in traditional programming. This is crucial for:
*   **Modularity:** Breaking down complex problems into smaller, manageable agent roles.
*   **Reusability:** Agents can expose generic capabilities that multiple other agents can leverage.
*   **Orchestration:** Enabling sophisticated workflows where agents coordinate by invoking each other's specific actions.
*   **Maintainability:** Easier to understand, debug, and evolve systems where agent interactions are clearly defined.

**Why this project shape (Go service with Protobuf and gRPC)?**
*   **Strong Typing and Schema Enforcement (Protobuf):** Defines agent capabilities (tools) and their parameters with strict types, ensuring clear contracts between agents. This eliminates ambiguity and reduces runtime errors.
*   **Efficient and Reliable Communication (gRPC):** Provides a high-performance, language-agnostic framework for inter-process communication. It's built for service-to-service interaction, which perfectly maps to agent-to-agent interaction in a distributed system.
*   **Language Agnostic:** While the prototype is in Go, gRPC generates client and server stubs for many languages, allowing agents written in different languages to seamlessly interact.
*   **Scalability:** gRPC is designed for microservices architectures, making it suitable for scaling multi-agent systems.
*   **Clear Tool Invocation Pattern:** gRPC methods directly map to agent "tools" or "capabilities," making the invocation explicit and structured.

**Setup and Usage (Zero API Keys):**

This prototype demonstrates a "Tool Registry" service where agents can register their capabilities, and other agents can discover and invoke them. We'll simulate agents interacting with this registry.

1.  **Install Go:** If you don't have Go installed, follow the instructions at [https://golang.org/doc/install](https://golang.org/doc/install).
2.  **Install Protobuf Compiler:**
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
    Ensure `protoc` itself is installed. For macOS with Homebrew: `brew install protobuf`. For other OS, check [https://grpc.io/docs/protoc-installation/](https://grpc.io/docs/protoc-installation/).
3.  **Clone the repository (or create the files manually):**
    ```bash
    # Assuming you've created a directory like 'agent-tool-registry'
    # Copy the files provided into this directory
    ```
4.  **Generate Protobuf Go code:**
    ```bash
    cd <your-project-directory>
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/registry.proto
    ```
5.  **Run the Tool Registry Service:**
    ```bash
    go run server/main.go
    ```
    You should see `Tool Registry Service started on :50051`.
6.  **Run the Agent Client Simulation:**
    In a *separate terminal*, run:
    ```bash
    go run client/main.go
    ```
    Observe the output. The client will simulate:
    *   `Agent A` registering a "CalculateSum" tool.
    *   `Agent B` registering a "FetchWeatherData" tool.
    *   `Agent C` discovering and invoking "CalculateSum".
    *   `Agent C` discovering and invoking "FetchWeatherData".

This demonstrates agents registering their capabilities and other agents invoking them through a central, typed registry. No LLM API keys are involved, as this is a pure orchestration and protocol layer.

**Optional Real-LLM Adapter:**
This prototype focuses on the *mechanism* of agent-to-agent tool invocation, not the *reasoning* about which tool to call or *how* to generate the parameters. LLMs would typically be used by an agent to decide *what* tool to call and *with what arguments* based on a prompt. This prototype provides the infrastructure that agent would then use to *execute* that decision. Therefore, no LLM adapter is provided, as the problem is purely about the communication protocol and not LLM integration itself.

---
