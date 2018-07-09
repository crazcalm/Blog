+++
date = 2018-06-19
draft = true
tags = ["golang", "christmas", "term-quiz" ]
description = "A recap of how I used lazy programming techniques to build a Christmas gift form my fiancee"
title = "Lazy programming saves Christmas!"
highlight = true
css = []
scripts = []
+++

## Let's clear the air

You may be thinking? "How can you be lazy AND still build a Christmas gift?"

I know. To truly be lazy around Christmas time is to put in as little effort as possible into buy a gift. The one issue with this "True Lazy" definition is that you have to have money to buy a gift...

## Paint the picture

I was in China working for the Peace Corps (AKA NO SALARY!) and my fiancee was in New Jersey, which is otherwise known as the land far, far away.

In my mind, the my situation equated to something like this:

#### "Little time spent together" + "Shitty or no Christmas gift" = "No more relationship"

## Back to the point

To answer the original question ("How can you be lazy and build a Christmas gift?")

### Answer: You can't. This post is not about me being lazy. It is about how I program in a lazy way...

So, what is the programming problem that I am trying to solve in a lazy way?

### [Terminal Quizzes!](https://github.com/crazcalm/term-quiz)

I want to create quizzes that contain questions that come in 3 forms.

1. ABCD
2. True or False
3. Fill in the blank

#### So, the question is "How do you write questions that satisfy these 3 forms? (ABCD, True of False, Fill in the blank)"?

Short answer: You don't. That sounds hard. Lazy me does not like hard. Plus, this does not have to be perfect. If I write a quiz app that only uses one form, it's STILL A QUIZ APPLICATION! Now lets pick which form to use.

ummm..... ABCD

#### Okay, new question. "How do you write questions that satisfy a ABCD form?"

That's easy: You have one question with 4 possible answers!

    type Question struct {
      Question string
      Answers  []Answer
    }

#### New question: How do you know if one of the answers are correct?

That's easy: Mark it! As long as the questions knows which answer is correct, then we can check that with the user answer.

    type Answer struct {
      Answer string
      Correct bool
    }

    func (q *Question) CorrectAnswer () (*Answer, error) {
      for _, a := range q.Answers {
       if a.Correct {
         return a, nil
        }
      }
      return nil, errors.New("No Answer found") 
    }

    type UserAnswer struct {
      Question *Question
      Answer   *Answer
    } 

    func (u UserAnswer) Correct() (result bool, err error) {
      correctAnswer, err := u.Question.CorrectAnswer()
      if err != nil {
        return result, err
      }

      if currectAnswer == u.Answer { //This works because they are pointers to the same answer
        return true, nil
      }
      return false, nil
    }

(Lazy me writes out the rest of the application my downtime, but don't tell anyone!)

<image src="/img/secret.jpg">

Okay, now that lazy me has the ABCD case more or less solved. It is time to try solving one of the other cases.

## Cases left:
- True or False
- Fill in the Blank

True and False sounds similar enough to ABCD. Lets do that one!

### True or False. May the odds forever be in your favor!

*Takes a look at old code snippet*

    type Answer struct {
      Answer string
      Correct bool
    }

    func (q *Question) CorrectAnswer () (*Answer, error) {
      for _, a := range q.Answers {
       if a.Correct {
         return a, nil
        }
      }
      return nil, errors.New("No Answer found") 
    }

    type UserAnswer struct {
      Question *Question
      Answer   *Answer
    } 

    func (u UserAnswer) Correct() (result bool, err error) {
      correctAnswer, err := u.Question.CorrectAnswer()
      if err != nil {
        return result, err
      }

      if currectAnswer == u.Answer { //This works because they are pointers to the same answer
        return true, nil
      }
      return false, nil
    }

I don't think I need to change anything... The only thing that is kind of weird is that Answer.Answer is either True or False now, but I can live with that.

Two cases are now supported and I do not have to change a thing!

Next (and last) Case!

## Fill in the Blank

Well... This seems different. 

1. There are no answers to select from.
2. The user must input their own answer...

### Things we do know:

1. A questions has to have an answer
2. Answers have to be able to be checked/compared

##### --> This suggests that we can still use our Question struct. Though it will only have one answer... And that answer will be the correct answer.

I guess this means that we can still use our Answer struct too.

Q: "How do a check an answer?"

The below cannot work because your user answer cannot be the same pointer as your original answer.

    func (u UserAnswer) Correct() (result bool, err error) {
          correctAnswer, err := u.Question.CorrectAnswer()
          if err != nil {
            return result, err
          }

          if currectAnswer == u.Answer { //This works because they are pointers to the same answer
            return true, nil
          }
          return false, nil
        }


#### Well. Here is what I know.

1. The pointers cannot be the same.
2. The only information left to compare is the Answer string.

Let try this:

    func (u UserAnswer) Correct() (result bool, err error) {
          correctAnswer, err := u.Question.CorrectAnswer()
          if err != nil {
            return result, err
          }

          if currectAnswer == u.Answer { //This works because they are pointers to the same answer
            return true, nil
          }

          if strings.EqualFold(currectAnswer.Answer, u.Answer.Answer) {
            return true, nil
          }
          return false, nil
        }

Fill in the blank case is now done!

## To Summarize

We started out with the problem of "How do you write questions that satisfy these 3 forms? (ABCD, True of False, Fill in the blank)"?

I did not know how to solve that problem. So, lazy me created an easier problem to solve (satisfy the ABCD form). Once I solved the easy problem, I went back and added one level more of complexity to the problem (satisfy both the ABCD and True and False form). And I repeated this process until the original problem that I could not solve was solved!

Important note: Even if I was unable to satisfy the True and False or Fill in the Blank form, my fiancee would still have a gift!

## New Equation:

### "Little time spent together" + "Awesome Christmas gift" = "Happy future Wife"

<image src="/img/yes.jpg">

