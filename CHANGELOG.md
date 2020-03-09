# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.3.0] - 2020-02-22

* We now have post and pre hooks for termination: _Terminating_ (pre-hook) and _Terminated_ (post-hook).  This allows things to be registered to be run once ALL of the `OnShutdown` hooks are run.

* BREAKING CHANGES:

    * Renamed: `IsDown` => `IsTerminating` (non-blocking boolean, as before)
    * Renamed: `Done` => `Terminating` (returns a channel, as before)
    * Added: `IsTerminated` (non-blocking boolean)
    * Added: `Terminated` (returns a channel)
