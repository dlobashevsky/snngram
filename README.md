# Simple scoring of text fragments

## Idea

We have a collection of words with frequencies (corpus).
We need a fast estimator of frequncy for new words arrived. If arrived word is contained in the corpus we already lnow it frequency.
If the word cannot be found - we can do following:
1. Try all possible binary splits  `W = W1 . W2`
2. Estimate frequencies/probabilities of W1 and W2
3. Pray for independency of W1 and W2 and select best split for maximal probability as P(W)=P(W1)*P(W2) :)

So the process is recursive, and we can achieve a tree-like parsed structure.
Actualy this way is a pure semi-naive bayes approach, but it works. ))

For scoring we use notation minus-log10, so score 7.65 is equal to probability 10^-7.65.


## Test corpus
We can use Google N-Gram dataset for testing

wget -ci google.1gram && zcat google*.gz | grep '^[a-zA-Z0-9]*[[:blank:]]' | cut -f 1,3 | LC_ALL=C sort >big.corpus
This corpus have duplicated records, so you need to join same ngrams by awk or other preferred language.

## Compile

go mod tidy
go build


