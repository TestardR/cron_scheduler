package main

/* We have a set of tasks, each running at least daily, which are scheduled with a simplified cron.
We want to find when each of them will next run.

The scheduler config looks like this:

30 1 /bin/run_me_daily
45 * /bin/run_me_hourly
* * /bin/run_me_every_minute
* 19 /bin/run_me_sixty_times

first field is the minutes past the hour
second field is the hour of the day
third is the command to run

Objectives:
We want you to write a command line program that when fed this config to stdin
and the simulated 'current time' in the format HH:MM as command line argument
outputs the soonest time at which each of the commands will fire and whether it is today or tomorrow.

Output:

For example given the above examples as input and
the simulated 'current time' command-line argument 16:10 the output should be

1:30 tomorrow - /bin/run_me_daily
16:45 today - /bin/run_me_hourly
16:10 today - /bin/run_me_every_minute
19:00 today - /bin/run_me_sixty_times

To run:
go run cmd/main.go -schedule="$(cat input.txt)" -timer=16:10
*/
