# Spotgen

## What?

Spotgen is a minimal CLI tool for quickly generating spotify playlists.

**Note:** You will need access to a browser in order to register your app and initiate user verification

## Installation

To install, clone this repository and run `go build`

## Configuration

In order to use Spotgen you must first create an app in the Spotify developer portal. Once created go to settings and set your redirect URI to `https://localhost:8888/` and save. Then, using the client id and client secret provided by the developer portal, create a top-level `.env` file and set the `CLIENT_ID` & `CLIENT_SECRET` variables.

## Usage

There are two types of playlists that can be generated with Spotgen, featured or recommended. Featured playlists generate a playlist based on tracks in your featured playlist library, while recommended playlists generate based on given seed items.

```
spotgen
spotgen is a minimal CLI tool for quickly generating spotify playlists.

USAGE: 
    ./spotgen <COMMANDS> --name <Rest of flags>

COMMANDS:
    feat    Generate a featured playlist
        --name      Name of the playlist to be generated

        --len       Length of the playlist to be generated

        --desc      Description of the playlist to be generated

        --pub       Publicity of the playlist to be generated

        --collab    Collaboration capabilities of the playlist to be generated

    rec     Generate a recommended playlsit
        --name      Name of the playlist to be generated

        --len       Name of the playlist to be generated

        --desc      Name of the playlist to be generated

        --art       Comma-separated string of artists to use for playlist generation

        --gen       Comma-separated string of genres to use for playlist generation

        --pub       Publicity of the playlist to be generated

        --collab    Collaboration capabilities of the playlist to be generated

```
**Note:** You must use quotes when using the `--desc`, `--art`, & `--gen` flags if the value contains spaces or else it will parse incorrectly 

### Examples

Featured playlist with all fields:
 `./spotgen feat --name Generate1 --len 25 --desc "Spotgenerated" --pub false --collab true`

Featured playlist with the only required flag:
 `./spotgen feat --name Generate2`

Recommended playlist with all fields:
`./spotgen rec --name Generate3 --len 40 --desc "Spotgenerated2" --art "Drake,The Internet" --pub true --collab false`

Recommended playlist with default length:
`./spotgen rec --name Generate4 --desc "Spotgenerated3" --art "Luke Bryan, Babytron" --pub true --collab false`

## To-Do

 - Automate verification step
 