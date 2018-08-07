# Word Frequency Mini Project

A simple test project

## Problem

Given the attached text file as an argument, your program will read the file, and output the 20 most frequently used words in the file in order, along with their frequency. The output should be similar to that of the following bash program:

```
#!/usr/bin/env bash
cat $1 | tr -cs 'a-zA-Z' '[\n*]' | grep -v "^$" | tr '[:upper:]' '[:lower:]'| sort | uniq -c | sort -nr | head -20
```

## Constraint
Standard container libraries are not permitted (maps are not allowed, but slices are allowed). The use of I/O streams is permitted, the use of C++/Java/Go strings is discouraged.
