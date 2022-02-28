## GURPS Character Sheet (Nix Derivation)

GURPS Character Sheet (GCS) is a stand-alone, interactive, character sheet
editor that allows you to build characters for the
[GURPS 4<sup>th</sup> Edition](http://www.sjgames.com/gurps) roleplaying game
system.

![Build Status](https://github.com/richardwilkes/gcs/actions/workflows/build.yml/badge.svg?branch=master)

## Changes

- Configuration & library files now follow more closely to the XDG folder spec (see [6210acf](https://github.com/DirectXMan12/gcs/commit/6210acf1cf35af9bfa6378ea5c770ea6caf5a5a0)):
  
  * Master Library: ~/.local/share/GCS/Master Library 
  * User Library: ~/.local/share/GCS/User Library
  * Config: ~/.config/GCS

- Building: build scripts have been modified slightly to build more nicely under Nix.

## Building

[default.nix](default.nix) pulls source from this GitHub repository.  May be used as a normal derivation (e.g. `nix-env -i -f default.nix`).
