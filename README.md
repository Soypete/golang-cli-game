# Web Server Game: 20 Questions

This is a game build for the Go WebServices in 3 weeks [course](https://github.com/Soypete/WebServices-in-3-weeks). This is a version of the game [20 questions](https://en.wikipedia.org/wiki/Twenty_questions).

[![Actions Status](https://github.com/soypete/{}/workflows/build/badge.svg)](https://github.com/soypete/{}/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/soypete/{}/branch/master/graph/badge.svg)](https://codecov.io/gh/soypete/{})

## How To play

1.  Add your user
2.  Start a game
3.  Invite players
4.  finish Your game

## Middleware:

This is introductory example of using middleware for metrics and auth. We are using [ExpVars](https://pkg.go.dev/expvar#section-documentation), [prometheus](https://github.com/prometheus/client_golang), and [basic auth](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#basic_authentication_scheme). This example is a starting point for software engineers to exand upon in their own services.

_NOTE_: while this is an example showing possible methods for implmenenting certain middleware technicques. It should not be considered a reference for best practices of security or monitoring of a production service.
