================
Embedding Probes
================

Overview
--------

``embedding-probes`` is a Go-based framework for evaluating the performance of text embedding models. It compares models like ``granite-embedding:latest`` and ``nomic-embed-text`` across various tasks, such as analogy detection, cross-language similarity, and semantic metric evidence analysis. Each task is implemented as a plugin, allowing easy extension with new probes. The framework generates a dynamically aligned results table, summarizing metrics and declaring winners based on task-specific criteria.

The project uses Go v1.24.0 and relies on an external embedding API (e.g., Ollama at ``http://localhost:11434/api/embeddings``) to generate embeddings. Results are displayed in a formatted table with columns sized to the widest entry, ensuring clean output.

Features
--------

- **Plugin-Based Architecture**: Tasks implement a ``Task`` interface with ``Name()``, ``MetricName()``, and ``Run()`` methods.
- **Dynamic Table Output**: Columns in the results table adjust to the widest entry, with proper separator alignment.
- **Extensible**: Add new probes by implementing the ``Task`` interface and registering them in the ``probes`` package.
- **Metrics**: Tasks use diverse metrics (e.g., Euclidean Distance, Accuracy, Weighted Similarity), with winners determined per task.
- **Cross-Language Support**: Tasks evaluate embeddings across languages like English, French, Mandarin, Russian, and Spanish.

Plugin Protocol
---------------

The plugin protocol is defined in ``probes/types.go``. Each task must implement the ``Task`` interface:

.. code-block:: go

    type Task interface {
        Name() string
        MetricName() string
        Run(models []string) (map[string]TaskResult, error)
    }

- **Name() string**: Returns the task’s display name (e.g., "Semantic Metric Evidence Task").
- **MetricName() string**: Returns the task’s metric name (e.g., "Weighted Similarity", "Accuracy").
- **Run(models []string) (map[string]TaskResult, error)**: Executes the task for the given models, returning a map of model names to results (``TaskResult`` contains ``Metric`` and ``Winner``).

Tasks register themselves using ``probes.RegisterTask()`` in their ``init()`` functions. See ``probes/analogy.go`` for an example.

Current Tasks
-------------

The framework includes nine tasks, each evaluating different aspects of embedding quality:

1. **Analogy Task**: Measures Euclidean Distance for analogy completion (e.g., "Paris is to France as London is to England").
2. **Cross-Language Capability Task**: Computes Cross-Language Similarity between Russian and French phrases.
3. **French Cross-Language Metric Evidence Task**: Evaluates Accuracy in identifying relevant French evidence for a wellness program metric.
4. **Mandarin Cross-Language Metric Evidence Task**: Measures Accuracy for Mandarin evidence relevance.
5. **Metric Evidence Task**: Assesses Accuracy in English evidence relevance (three evidence chunks).
6. **Russian Cross-Language Metric Evidence Task**: Evaluates Accuracy for Russian evidence relevance.
7. **Semantic Metric Evidence Task**: Computes Weighted Similarity for semantic relevance of English evidence.
8. **Semantic Similarity Task**: Measures Semantic Similarity between synonymous Russian phrases.
9. **Spanish Cross-Language Metric Evidence Task**: Evaluates Accuracy for mixed English/Spanish evidence relevance.

Directory Structure
-------------------

- ``main.go``: Entry point; loads config, runs tasks, and prints the results table.
- ``probes/``: Contains task implementations and shared utilities.
  - ``types.go``: Defines the ``Task`` interface, ``TaskResult``, and utility functions (e.g., ``getEmbedding``, ``cosineSimilarity``).
  - ``analogy.go``, ``cross_language.go``, etc.: Individual task implementations.
- ``config.json``: Specifies models to evaluate (e.g., ``["granite-embedding:latest", "nomic-embed-text"]``).
- ``go.mod``: Go module dependencies.

Setup
-----

1. **Install Go**: Ensure Go v1.24.0 or later is installed (``go version``).
2. **Clone the Repository**:

   .. code-block:: bash

       git clone <repository-url>
       cd embedding-probes

3. **Set Up Embedding API**: Run an embedding service (e.g., Ollama) at ``http://localhost:11434/api/embeddings`` supporting the configured models.
4. **Configure Models**: Edit ``config.json`` to list models, e.g.:

   .. code-block:: json

       {
           "models": ["granite-embedding:latest", "nomic-embed-text"]
       }

5. **Install Dependencies**:

   .. code-block:: bash

       go mod tidy

Usage
-----

Run the program to evaluate all tasks:

.. code-block:: bash

    go run ./main.go

**Output**:
- Lists registered tasks (e.g., "Registered tasks: 9").
- Displays per-task results (e.g., similarities, accuracies).
- Prints a ``Final Results Table`` with columns for Task, Task Name, Metric, model scores, and Winner.
- Summarizes overall reliability (e.g., "nomic-embed-text is more reliable (7 vs. 2 wins)").

Example table (hypothetical values):

.. code-block:: text

    Final Results Table:
    | Task | Task Name                                    | Metric                    | granite-embedding:latest | nomic-embed-text | Winner                   |
    |------|----------------------------------------------|---------------------------|--------------------------|------------------|--------------------------|
    | 1    | Analogy Task                                 | Euclidean Distance        | 40.9934                  | 19.3172          | nomic-embed-text         |
    | 2    | Cross-Language Capability Task               | Cross-Language Similarity | 0.6103                   | 0.4200           | granite-embedding:latest |
    ...

Extending the Framework
----------------------

To add a new task:

1. Create a new file in ``probes/`` (e.g., ``new_task.go``).
2. Implement the ``Task`` interface:

   .. code-block:: go

       package probes

       import "fmt"

       type newTask struct{}

       func (t *newTask) Name() string {
           return "New Task"
       }

       func (t *newTask) MetricName() string {
           return "New Metric"
       }

       func (t *newTask) Run(models []string) (map[string]TaskResult, error) {
           results := make(map[string]TaskResult)
           // Task logic here
           return results, nil
       }

       func init() {
           fmt.Println("Registering New Task")
           RegisterTask(&newTask{})
       }

3. Run ``go mod tidy`` and ``go run ./main.go`` to include the new task.

Contributing
------------

Contributions are welcome! Please:
- Submit pull requests with new tasks or improvements.
- Report issues via the repository’s issue tracker.
- Ensure code follows Go conventions and includes tests where applicable.

License
-------

MIT License. See ``LICENSE`` file for details.

Contact
-------

For questions, contact the maintainers via the repository’s issue tracker.