# GURPS Character Sheet
![GitHub all releases](https://img.shields.io/github/downloads/richardwilkes/gcs/total?style=plastic)

GURPS[^1] Character Sheet (GCS) is a stand-alone, interactive, character sheet editor that allows you to build
characters for the [GURPS 4<sup>th</sup> Edition](http://www.sjgames.com/gurps) roleplaying game system.

GCS relies on another project of mine, [Unison](https://github.com/richardwilkes/unison),
for the UI and OS integration. The [prerequisites](https://github.com/richardwilkes/unison/blob/main/README.md) are
therefore the same as for that project. Once you have the prerequistes, you can build GCS by running the build script:
`./build.sh`. Add a `-h` to see available options.

[^1]: GURPS is a trademark of Steve Jackson Games, and its rules and art are copyrighted by Steve Jackson Games. All
rights are reserved by Steve Jackson Games. This game aid is the original creation of Richard A. Wilkes and is
released for free distribution, and not for resale, under the permissions granted in the
<a href="http://www.sjgames.com/general/online_policy.html">Steve Jackson Games Online Policy</a>.

## Changes

- Configuration & library files now follow more closely to the XDG folder spec (see [cd48601](https://github.com/DirectXMan12/gcs/commit/cd4860125859f4e7cc93e328f8d5d4e2838b4a94)):
  
  * Master Library: ~/.local/share/GCS/Master Library 
  * User Library: ~/.local/share/GCS/User Library
  * Config: ~/.config/GCS

- Building: build scripts have been modified slightly to build more nicely under Nix.

## Building

[default.nix](default.nix) pulls source from this GitHub repository.  May be used as a normal derivation (e.g. `nix-env -i -f default.nix`).
