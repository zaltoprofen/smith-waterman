# smith-waterman

Command line tools for string comparison with Smith-Waterman algorithm.

## Installation

```
$ go get -u github.com/zaltoprofen/smith-waterman
```

or Download binary from https://github.com/zaltoprofen/smith-waterman/releases


## Usage
```
$ smith-waterman string1 string2
```

e.g.
```
$ smith-waterman "アーチャー" "アーニャ"
アーチャ
｜｜　｜
アーニャ
```
