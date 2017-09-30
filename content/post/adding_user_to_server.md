+++
scripts = []
css = []
date = "2017-09-30T15:30:48+08:00"
highlight = true
title = "Adding a user to your server"
description = ""
tags = ["ubuntu", "ssh", "server"]
draft = true

+++
# Why I am writing this blog?

I am writing this blog because I a tired of forgetting how to do this. I do not add users to my server very often. As such, every time I do so, I do it incorrectly and the new user is either not in the sudoer group, the user has no home directory, or that user cannot be logged into via ssh.

# Requirements:
- Create a new user
- Set a password
- Have a home directory with skeleton files
- Be able to use sudo
- Can use ssh to log in

# Use the "adduser" command

My linux box, running ubuntu, has two commands that can be used to created users; `useradd` and `adduser`. However, useradd is a low level interface that should never be used by the average person. Thus, we should use the high level interface that is adduser.

## Perks of adduser

By default, adduser command uses the `/etc/adduser.conf`, which creates a home directory for the new user, forces that user to set a password, copies skeleton files into the new directory, and some other things that I do not care too much about.

Command: `adduser <<username>>`

### Note:
- Skeleton files include files such as .bashrc and .profile.
	
- The skeleton file directory is located: /etc/skel/
	
- The files in that directory are copied over to the home directory of every new user that is created.
	

## Add to sudo group

The adduser command also allows you to add users to the sudo group! The one catch is that the user must already exist. See the below command.

Command: `adduser <<username>> sudo`

## Login via ssh

In the home directory of the new user, create a `.ssh` directory. Within that directory, create a file called `authorized_keys` and add your public ssh rsa key to that file. Once done, test logging into that user:

`ssh <<username>>@ip_address`