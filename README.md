# Pipeline Processing System

This project is a demonstration of a pipeline processing system built with Go. It consists of several components that work together to process integer inputs from the user in real-time. The system filters and manipulates data through a series of stages, utilizing concurrency mechanisms in Go such as goroutines and channels.

## Key Features

- **Data Source Input**: Accepts integer inputs from users via standard input (stdin). The special input "exit" terminates the program.
- **Negative Filter**: Filters out negative integers, allowing only positive integers to pass through.
- **Special Filter**: Further filters integers, passing through only those that are not zero and are multiples of 3.
- **Buffering System**: Temporalily stores processed data using a Ring Buffer mechanism, ensuring efficient data management and preventing overflow.
- **Dynamic Pipeline Construction**: The pipeline is built dynamically with the ability to easily insert or modify processing stages.

## System Components

### RingBuffer
- A custom, thread-safe implementation of a cyclic buffer for storing integers efficiently. Essential for throttling and managing the flow of data through the system.

### PipelineInt
- Orchestrates the flow of data through various processing stages. It is built flexibly to allow easy addition or rearrangement of stages.

### Stages
- Composable processing functions designed to perform specific tasks in the pipeline. This project includes:
  - negativeFilterStageInt: Filters out negative numbers.
  - specialFilterStageInt: Filters out numbers that are not multiples of 3 and not equal to zero.
  - bufferStageInt: Implements the buffering mechanism.


## Getting Started

### Prerequisites
- Go programming language installed on your system.

### Running the Program
1. Clone or download the repository to your local machine.
2. Open a terminal and navigate to the directory containing main.go.
3. Run the command go run main.go.
4. Start typing integers followed by the Enter key, observing the system's output.
5. Type exit when you wish to stop the program.

## How It Works
The program initiates by reading input from the user. Each input is sent through a series of stages designed to filter and process the data. Inter-stage communication is achieved via channels, with goroutines handling concurrent processing. The system showcases effective use of Go’s concurrency model to implement a non-trivial processing pipeline.
