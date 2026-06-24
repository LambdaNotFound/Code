### Overview
Design and build a log utility (Logger) to process logs. This log utility can support multiple different handlers, and the design will be extended to new logger handlers in the future. 

### API
Logger class, with log(s str) takes a input string, and all the log handlers process this string independently. 

LogHandler class, with a process(s str) behavior method for its sub-classes. Below are the concrete log handlers:

1. Filter log handler - Removes occurrences of a string (configurable) from the log message before printing it.

2. Truncation log handler - Truncates to N characters before printing.

3. Uppercase log handler - Capitalizes messages before printing.

4. Array log handler - Stores log messages in an array without printing.

### Low Level Design
Apply OOP principal in this design, using the strategy pattern for the LogHandler. 

A log handler processes the received log message individually (fan-out) and is independent of other log handlers. Each handler owns its output, there is no chaining.

The Logger class takes a list of LogHandlers during the initialization. 

### Test Plan
Calling logger.log("hello world") on the main logger will log the message to the console multiple times: as-is, then truncated, then capitalized. The message will also be stored in an array. Note: the message gets handled in a fan-out way, not chaining. 