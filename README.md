# RobacBOT

RobacBOT is a discord bot built in Golang. The planned structure allows a package-based design for adding new commands to the bot, without modifying too much in the core code.
TEST

## Getting Started

The project should be able to run on any platform with Golang supported on it.

## Prerequisites

Project has been compiled with Golang 1.14 and Discordgo v0.21.1

## Installing

Download the repo anywhere and compile with go.

## Configuration

A file called 'config.json' contains settings for the bot, such as the bot's token. Certain command packages have their own .json configuration files.

### Adding new command packages

Go's structure of packages and modules permits some freedom regarding dividing commands in different files and packages, with the purpose to allow multiple people to work on different commands and to avoid having everything in just one file.

In order to add your own command package, you must make a new folder and at least one .go file in it, with the first line being 'package X', X being the name of your folder. You would write your code in that .go file.
Afterwards, you are supposed to import the package in commands.go, and then add a new case to the switch, with the necessary arguments.

### Troubleshooting

For issues, contact the developers at the Robac discord server.
