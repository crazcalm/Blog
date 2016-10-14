+++
Tags = [
  "Development",
  "Golang",
  "Editor",
  "Micro"
]
Categories = [
  "Development",
  "GoLang",
  "Editor",
  "Micro"
]
menu = "main"
date = "2016-10-07T19:03:20-04:00"
title = "Micro: My New Editor?"
Description = ""
Draft = "False"
+++

Why Learn Micro?
===============
I have no good reason as to why someone should learn micro over some other text editor. With that being said, here are some reasons as to why I want to learn how to use this editor.

1. It is written in Go.
    - I am trying to learn Go, so a text editor written in Go appeals to me.

2. Writing Go code in it is easy!
    - Micro has builtin features that make writting Go easier such as running go fmt after a file has been saved and showing the lines where compilation errors have occured.

3. I can run it in my Terminal.
    - A much have! I have a old netbook that I do a lot of programming on and that thing is not powerful enough to run a GUI application.

4. It is still in development.
    - I like the idea of this because it gives me the option of contributing to the project.

The Basics
==========
Command mode:
------------
`Ctrl + e`: Gets you into command mode. This mode is needed for almost everything!

The Tutorial
------------
While in command mode, enter `help tutorial`. This document will give you an overview of how micro works and how to customize it.

Getting to the Help Docs
------------------------
While in command mode, enter `help`. This is my go to document for finding more details about a micro feature.

Customizing the Editor Colors
-----------------------------
While in command mode, enter `help colors`. This will bring you to a document that will instruct you on how to change the programming syntax colors, editor theme colors, etc.

Keybidings
----------
I have been using nano for the past 2 years, so vim and emacs keybindings are something I know nothing about.

To see Micro's default keybindings, type ctlr + e to enter command mode and then type "help keybiddings".
This will result in a new horizontal window appearing with the default keybindings.

Here are a list of some of the key bindings that I find useful.

- `Ctrl + q` = Quit
- `Ctrl + o` = Open
- `Ctrl + s` = Save
- `Ctrl + f` = Find
- `Ctrl + n` = Find next
- `Ctrl + p` = Find previous
- `Ctrl + z` = Undo
- `Ctrl + y` = Redo
- `Ctrl + c` = Copy
- `Ctrl + x` = Cut
- `Ctrl + k` = Cut line
- `Ctrl + v` = Paste
- `Ctrl + a` = Select all
- `Ctrl + l` = Jump line
- `PageUp` = Cusor moves up a full page
- `PageDown` = Cursor moves down a full page
- `Ctrl + w` = Next split
- `Ctrl + t` = Add a Tab
- `Ctrl + /` = Next Tab
- `Ctrl + b` = Shell mode
- `Alt + left/right` = move the cursor a word in that direction
- `Ctrl + u` = Toggle Macro (recoding keystrokes?)
- `Ctrl + j` = Play Macro (play recorded steps)

You may create a `~/.config/micro/bindings.json` file where you can customize your keybindings.
Examples of this can be found in the keybindings help doc (in command mode run: `help keybindings`)