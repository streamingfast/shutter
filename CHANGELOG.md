# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.4.1] - 2020-06-23

* Added missing license files

## [v1.4.0] - 2020-03-09

### Changed

* License changed to Apache 2.0

* We now have post and pre hooks for termination: _Terminating_
  (pre-hook) and _Terminated_ (post-hook).  This allows things to be
  registered to be run once ALL of the `OnTerminating` (the previous
  `OnShutdown`) hooks are run.

* BREAKING CHANGES:

    * Moved package from `github.com/eoscanada/shutter` to `github.com/dfuse-io/shutter`
    * Renamed: `OnShutdown` => `OnTerminating` (the pre-hook)
    * Renamed: `IsDown` => `IsTerminating` (non-blocking boolean, as before)
    * Renamed: `Done` => `Terminating` (returns a channel, as before)
    * Added: `OnTerminated(func(err error))` (callback called AFTER all the `OnTerminating` callbacks are done)
    * Added: `IsTerminated` (non-blocking boolean)
    * Added: `Terminated` (returns a channel)
    * Removed `NewWithCallback`, replaced by `New(RegisterOnTerminated(f1), RegisterOnTerminating(f2))` options-style configuration.

* `New(opts ...)` extended to accept options-style configurators (was previously `New()`)
