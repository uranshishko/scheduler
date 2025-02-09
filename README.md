# Scheduler Package

## Overview
The `scheduler` package provides a lightweight scheduling system that allows users to run tasks at predefined or custom intervals. It supports both standard cron-like expressions (`@yearly`, `@monthly`, etc.) and custom durations (`@every 10s`, `@every 5m`, etc.).

## Features
- Supports predefined time expressions such as `@daily`, `@weekly`, and `@hourly`.
- Allows custom time intervals using the `@every` syntax.
- Supports duration expressions like `10h20m5s100ms1200ns` and so on.
- Runs scheduled tasks asynchronously.
- Provides a cancel function to stop scheduled tasks.

## Installation
To use this package in your Go project, install it using:

```sh
 go get github.com/uranshishko/scheduler
```

Then import it in your Go code:

```go
import "github.com/uranshishko/scheduler"
```

## Usage
### Creating a Scheduler
To start using the scheduler, initialize it with a start time:

```go
import (
    "fmt"
    "time"
    "github.com/uranshishko/scheduler"
)

func main() {
    s := scheduler.New(time.Now())
}
```

### Scheduling a Task
Use the `Schedule` method to set up a task:

```go
func task(event scheduler.Event) error {
    fmt.Println("Task executed at:", event.Time)
    return nil
}

func main() {
    s := scheduler.New(time.Now())
    cancel, err := s.Schedule("@every 10s", task)
    if err != nil {
        fmt.Println("Error scheduling task:", err)
        return
    }
    
    defer cancel() // Ensure cancellation when the program exits
    time.Sleep(30 * time.Second) // Let the scheduler run for a while
}
```

### Canceling a Scheduled Task
The `Schedule` method returns a `cancel` function that stops the task execution:

```go
cancel() // Stops the scheduled task
```

## Expression Syntax
The scheduler recognizes two types of expressions:

### Predefined Expressions
- `@yearly`   → Runs once a year
- `@monthly`  → Runs once a month
- `@weekly`   → Runs once a week
- `@daily`    → Runs once a day
- `@hourly`   → Runs once an hour

### Custom Intervals
- `@every 10s` → Runs every 10 seconds
- `@every 5m`  → Runs every 5 minutes
- `@every 1h`  → Runs every 1 hour
- Supports complex duration expressions like `10h20m5s100ms1200ns`.

## Error Handling
If an invalid expression is provided, the `Schedule` method returns an error:

```go
cancel, err := s.Schedule("invalid-expression", task)
if err != nil {
    fmt.Println("Error:", err)
}
```

If the handler function returns an error, the task stops execution.

## Contributing
Contributions are welcome! Feel free to submit a pull request or open an issue for bugs or feature requests.

## License
This project is licensed under the MIT License.

## Author
Developed by Uran. Reach out on GitHub for any questions or improvements!
