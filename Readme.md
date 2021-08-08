wordcount
=========

## usage

```
wordcount - print the number of unique words in a stream

Usage: wordcount [file]

Example:

interactive mode: wordcount
read from file:   wordcount test.txt
read from pipe:   echo lorem ipsum lorem | wordcount

Output:

hello: 500
world:  50
how:    10
are:     5
you:     5
```

## interactive

```
Enter text to capture word counts.
The following keywords will not be counted (beginning with '-')
-stats: print statistics
-reset: reset statistics
-help: print help
-exit: quit
> Hello hello hello. 
> How are you?
> Are you doing well?
> -stats
hello: 3
are:   2
you:   2
doing: 1
how:   1
well:  1
> -exit
```

## other ideas

* add a flag to control whether we should ignore case or not
* print stats when killed with Ctrl+C
