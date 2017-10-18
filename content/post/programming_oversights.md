+++
highlight = true
date = "2017-10-18T13:44:46+08:00"
title = "Programming Oversight..."
description = ""
tags = ["windows", "terminal", "golang", "frustrated"]
draft = false
scripts = []
css = []

+++
## The Problem

I have a horrible memory. If I meet 10 people at a party, I might remember 2 of their names. Recently, I started teaching English in China, and I now have 300+ students. In the spirit of facing my memory problem head on, wrote a programmatic solution.

## Flash Cards App

<image src="/img/flashcard_app.png">
 
I built my own flashcard app! You can find the project on Github at [https://github.com/crazcalm/flash-cards](https://github.com/crazcalm/flash-cards). In summary, I used Golang to build a terminal application that reads a csv file and turns that content into flashcards.

After building this tool, I thought "Wow, this is great! Maybe I can do something similar to solve other classroom problems!"

## A Repo of Class Related Tools

I then create a repo called [student-csv](https://github.com/crazcalm/students-csv) to house random tools that I would use in class.

- I have tool that selects a random student.
- I have a tool to break the class into groups.
- I have a tool to take attendance.

All these tools are awesome on my Linux box, but this school only uses Window's 7 machines in class. For the past month, I have been telling myself that porting these apps to Windows should be easy. It is just a terminal app. How hard can it be?

Then I actually tried porting it last night....

## The Oversight

I was right, porting these apps to Windows was easy. The only thing that I had to modify were the system command calls. In each of these apps, I send to clear command to the terminal to clear the screen. On Windows, that command is slightly different.

**The oversight was that the Window's command prompt does not support utf-8...**

I am in China! All my students are Chinese people with Chinese names! Thus, even though all of these apps work fine, I cannot use them in class because the Window's command prompt cannot print out Chinese Characters!