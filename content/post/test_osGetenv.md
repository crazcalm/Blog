+++
draft = false
date = 2018-07-31
tags = ["golang", "testing"]
description = ""
title = "How to Test Environment Variables"
highlight = true
css = []
scripts = []
+++

## Background

I recently came across a piece of code that checks for the existence of an environment variable. If the environmental variable exists, the code will use the value of that variable. If it does not exist, the code will return a default value. 

Equivalently, this function could of set the environmental value to it's default value and then returned the value of the environmental variable.

	func envValue(env, defaultValue string) (result string, err error) {
		result = os.Getenv(env)
		if strings.EqualFold(result, "") {
			err = os.Setenv(env, defaultValue)
			if err != nil {
				return
			}
			result = os.Getenv(env)
		}
		return
	}

Code like this is fairly common (especially when your code needs to write documents to a specific location).

## How do you test that?

In order to create a test for this, you want to make sure that you cover both cases, which are when the environmental variable does and does not have a value set.

### The Implicit Way

To implicitly test this example, you can note what state your machine has before writing your test. That is, does your machine have this variable set? If so, you can use your machine's state to test the case where the environmental variable has a value and then 
set the environmental variable to the empty string to test the case where the environmental variable does not exist.

On my machine, the environmental variable NAME does not exist.
 Here is an example of how to implicitly test envValue using NAME:

	func TestEnvValueImplicit(t *testing.T) {
		tests := []struct {
			Key     string
			Default string
			Answer  string
		}{
			{"NAME", "default1", "default1"},
			{"NAME", "default2", "default1"},
		}

		for i, test := range tests {
			result, err := envValue(test.Key, test.Default)
			if err != nil {
				t.Errorf("envValue returned an err: %s", err.Error())
			}

			if !strings.EqualFold(result, test.Answer) {
				t.Errorf("Case (%d): Expected %s to be equal to %s, but got %s", i, test.Key, test.Answer, result)
			}
		}
	} 

During the first test case, since NAME did not exist, NAME got set to the test.Default value that was passed into the function, which was "default1".

During the second test case, NAME does exist and does not get set to the test.Default value, which was "default2". Instead, the value of NAME is "default1", which was set during the first test case.

In summary, we have leveraged the state of our machine to test the needed cases for this function.

#### Drawbacks

Given that the implicit way of testing this function utilizes the state of our machine in the test. If follows that a machine with a different initial state may get differing results.

For example: If someone else were to run this test, and they had the environmental variable NAME set to a value other than the empty string, then the test would not cover the "NAME does not exist" test case because, unlike my machine, the environmental variable NAME does exist on their machine.

### The Explicit Way

To avoid the drawbacks that come with implicitly testing this function, I would suggest that you explicitly test this function by setting NAME to the needed value before passing it into the function you want to test.

	func TestEnvValueExplicit(t *testing.T) {
		tests := []struct {
			Key     string
			Value   string
			Default string
			Answer  string
		}{
			{"NAME", "", "default1", "default1"},
			{"NAME", "default1", "default2", "default1"},
		}

		for i, test := range tests {
			//Need to keep original value so that I can clean up after the test
			oldValue := os.Getenv(test.Key)

			//Ensuring that the environmental variable exists
			err := os.Setenv(test.Key, test.Value)
			if err != nil {
				t.Fatalf("Failed to set environment variable: %s", err.Error())
			}

			result, err := envValue(test.Key, test.Default)
			if err != nil {
				t.Errorf("Case (%d): envValue returned an err: %s", i, err.Error())
			}

			if !strings.EqualFold(result, test.Answer) {
				t.Errorf("Case (%d): Expected %s to be equal to %s, but got %s", i, test.Key, test.Answer, result)
			}

			err = os.Setenv(test.Key, oldValue)
			if err != nil {
				t.Logf("Case (%d): Failed to set environment variable back to %s: %s", i, oldValue, err.Error())
			}
		}
	}

In the above code, we set NAME to the empty string to test that "NAME does not exist" case, and we set NAME to "default1" to test the "NAME does exist case". By explicitly setting NAME during each test case, the test no longer relies on the state of the machine. As a result, this test will work the same on machines that have and do not have the NAME variable set. 