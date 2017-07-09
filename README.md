## tune.sh

download your OS specific binary from [releases](https://github.com/abhishekkr/tune.cli/releases)

---

#### usage

```
## to play searched songs on confirmation
tune.cli-linux-amd64 -search sout -type movie -song 1 -out play

## search for lucifer tv show soundtracks
tune.cli-linux-amd64 -search lucifer -type tv

## search for soundtracks for movie with 'sout' in title
tune.cli-linux-amd64 -search sout -type movie

## search for soundtracks for movie with 'sout' in title, song#3
tune.cli-linux-amd64 -search sout -type movie -song 3

## search for tv with lucifer in title, for season 1
tune.cli-linux-amd64 -search lucifer -type tv -season 1

## search for tv with lucifer in title, for season 1 episode 2
tune.cli-linux-amd64 -search lucifer -type tv -season 1 -episode 2

## search for tv with lucifer in title, for season 1 episode 2 song listing#4
tune.cli-linux-amd64 -search lucifer -type tv -season 1 -episode 1 -song 4

## search for soundtracks for artist with 'eminem' in title
tune.cli-linux-amd64 -search sout -type artist

```

* sample run

```
± % ./go-tasks.sh gorun -search south -type movie -song 1                                                                                              !11957
[ /movie/southpaw-2015 ]
[*] Beast (feat. Busta Rhymes, KXNG Crooked & Tech N9ne)
    [url](https://www.tunefind.com/song/85250/108074/)
    by [Eminem](https://www.tunefind.com/artist/eminem)
    listen at [youtube](https://www.tunefind.com/forward/song/85256?store=youtube&referrer=&x=osshgm&y=hyz50jpazt44sw4kc0kg0gs40s4gw4)
[ /movie/southside-with-you-2016 ]
[*] Start
    [url](https://www.tunefind.com/song/106149/136815/)
    by [Raphael Saadiq](https://www.tunefind.com/artist/raphael-saadiq)
    listen at [youtube](https://www.tunefind.com/forward/song/20032?store=youtube&referrer=&x=osshyo&y=813e3b1hg6o8o48sgg4gc8kg48c84ow)

```

---

[dev wiki](./wiki.developer.md)

---

