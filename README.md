# chess.com PGN Collator

I wondered what openings I should work on so I wrote this command line
application that analyses a collection of chess.com PGNs and counts the wins,
losses and draws for both colours. Obviously there's more to chess than the
opening, but these statistics give some insight into what to improve in your
reportoire.

## Install

This program is written in Go and can be installed using:

`go get -u github.com/bspaans/chess-archive-collator`

## Usage

Currently this only works for chess.com, as far as I know, as they add an "ECO"
code to their tags, which I've used to classify the opening. 

Get a month's worth of PGN data:

`curl https://api.chess.com/pub/player/bartspaans/games/2019/10/pgn > 2019_10.pgn`

Run the program:

`chess-archive-collator 2019_10.pgn`

Resulting in something like this (but hopefully with more wins):

![Example result](https://raw.githubusercontent.com/bspaans/chess-archive-collator/master/screenshot.png)

(classifying openings myself is work in progress/might never happen.
Incidentally if anyone knows of an open source opening database let me know).
