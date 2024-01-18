# Luzifer / accounting

## What's this?

`accounting` is a server based software to keep track of my budget on all the different accounts I do have. It ensures my money is where I need it and is available when I need it.

It is stronly inspired by [YNAB](https://www.ynab.com/) and uses some of the same principles but in some parts deviates from them.

## Why write an own software?

There is plenty of accounting software out there (including YNAB) which doesn't really fit my use-case. Also after being a happy customer of YNAB for nearly 9 years I've seen many changes implemented into YNAB which didn't really fit my use-case anymore and some of them even broke the ways I was using the software.

This piece of software if my way to ensure I can continue to work the way I want while ditching all the features I don't need. So this is not a clone and probably not a replacement for YNAB to you but for me it is.

## Should I use it?

Well, depends on whether it fits your needs. This software carries a strong API-first approach so my automation can do stuff for me. The interface is in a "first usable version" state and I wouldn't call it ready. Feature-wise this is not ready yet:

- [x] Manage accounts and categories
- [x] Manage a (monthly) budget based on the categories
- [x] Differentiate between budget and tracking accounts
- [x] Have an API to automatically enter transactions, transfer money between accounts & categories
- [x] Can properly calculate based on your transactions how many money you have in which category
- [ ] Have reports telling you about your net-worth, where the money went, …
- [ ] Have a fancy UI with keyboard-shortcuts and all the cool stuff

So back to your question: Get yourself a copy, start it (the default database is an in-memory SQLite database) and take it for a test-drive. You like it? Put a proper database under it, install it on a server and happy budgeting… You don't like it? No offence taken, just put it into the trash and use something else.

## Will you accept contributions?

In short: Maybe. Longer version: This software started as a system to do **my** budget and furthermost needs to fit my needs. I don't want feature-creep in it, I don't want features I'll never use. If you plan to hack on it: Just create a fork, hack on it, build something cool for yourself. If you think I'd like the changes you made: Lets talk about it but keep in mind this explanation and don't be angry with me when the response is "Nah, I don't think this will fit my version…"!

## How to run it?

For local testing just clone the repo, ensure you got a fairly new `nodejs` (including NPM) and `go` toolchain installed and

```console
# make run
```

then go to http://localhost:5000/ and try it.

To really host it, `make build` it, put it on your server and let it run. (It's a binary so probably use systemd to run it and an nginx as proxy in front of it doing auth-stuff…)
