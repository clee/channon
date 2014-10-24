# channon

Channon is a continuous integration system that doesn't try to do
a bunch of things half-assed. Instead, it tries to do only a few things
full-assed.

You define plans to describe how you want things to be run, and each
plan can contain multiple steps. Each step is a simple text file which
gets dumped onto the disk, made executable, and then executed. This means
that whatever kind of script you want to write, you just make sure to put
the correct hashbang in the first line and as long as your system has
the binary you need, it'll run your script with that binary.

Channon is completely agnostic about the kinds of scripts you write
the steps for your plan in. Want to use Python? Ruby? Bash? Scala? Go
right ahead. The environment variables set in each step persist into
the following steps, so you can set env vars to communicate status,
file locations, set up rbenv or virtualenv/pip, or whatever else you
need. If any step returns a non-zero value, Channon marks the run as
a failure, and optionally will run a “notify” step with all the
environment vars from the previous steps so you can email, jabber,
IRC, whatever, to tell people that things aren't working.

The same “shoot yourself with a bowel disruptor” attitude is applied to
version control. Make the first step in a plan a script to check the code
out fresh every time, or get a little more fancy and make it check out
the first time and then do updates afterwards. Or forget version control
entirely if you want, and just run simple scripts repeatedly.

Jobs can be triggered via HTTP request or scheduled using cron-style
syntax.

Channon's API embodies my ideas about what a modern server API should
look like, which means that it may be completely insane to you. Most
interactions will be over JSON, although there are also EventStream
endpoints so that if you want to drink from the firehose, you can. No
GUI of any kind is built in, but it should be straightforward to
implement one using a native toolkit or HTML5.
