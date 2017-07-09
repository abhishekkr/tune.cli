
### Wiki for running from code

```
## fetch dependencies
./go-tasks.sh deps

## to play searched songs on confirmation
./go-tasks.sh run -search sout -type movie -song 1 -out play

## search for lucifer tv show soundtracks
./go-tasks.sh run -search lucifer -type tv                                                                                                       !11843

## search for soundtracks for movie with 'sout' in title
./go-tasks.sh run -search sout -type movie

## search for soundtracks for movie with 'sout' in title, song#3
./go-tasks.sh run -search sout -type movie -song 3

## search for tv with lucifer in title, for season 1
./go-tasks.sh run -search lucifer -type tv -season 1

## search for tv with lucifer in title, for season 1 episode 2
./go-tasks.sh run -search lucifer -type tv -season 1 -episode 2

## search for tv with lucifer in title, for season 1 episode 2 song listing#4
./go-tasks.sh run -search lucifer -type tv -season 1 -episode 1 -song 4

## search for soundtracks for artist with 'eminem' in title
./go-tasks.sh run -search sout -type artist

```
