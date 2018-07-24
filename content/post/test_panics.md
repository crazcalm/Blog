+++
draft = false
date = "2018-07-24"
tags = ["golang", "testing", "panic"]
description = "How to test panics tutorial"
title = "How to Test Panics"
highlight = true
css = []
scripts = []
+++

## Setup

Before we can test a panic, we need to have code that panics.

	package panic

	func willPanic() {
		panic("Told you I would panic")
	}

## Testing

Given that we know the code will panic, we can do the following:

	package panic

	import (
		"testing"
	)

	func TestWillPanic(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		//Call func
		willPanic()
	}
	
In the above, we try to recover and, if the recovery is equal to nil, we marked the test as failed.

## Thoughts

This is a contrived example. That being said, in most scenarios where you have code that can panic, the route to panicking is usually one outcome out of many possible outcomes. Given that this case is unique, I believe that it should be tested uniquely. That is, if you have code that can panic, write a test function to covers the "more expected" outcomes, and write a separate test function for the cases in which the code panics.
