# System Test Plan

## General Guidelines

1. Establish a Test User so that generated files are in one directory.
2. Give each test a project name that is unique and descriptive
3. Add info on memory usage and CPU to the log
4. Define a testing package to hold the system tests
5. Put request.yaml files into a directory
    1. ? How do we initiate the execution of a specific test.  Do I create a test function for each test?  Should it submit to the server?  Auto submitting to a CLI would require the use of cmd.Exec
    2. ? Where are regression test results stored

These are all script level tests

When the system is used for textual compares, some of the runs will have no output, except the stored data.

How many languages will be loaded?
- [ ] Load Plain Text Script from BB-API
- [ ] Load Plain Text + Headings from BB-API
- [ ] Load USX Text edited from BB-API
- [ ] Load Script Using File
    - [ ] What is the Output that is desired

In each language loaded:
Compare Plain Text - Plain Text Edit
Plain Text Edit - USX Edit
USX Edit - Script
Plain Text Edit - Script
For each compare output, an html file
The html file should contain count of the number of character differences
(Purpose is primarily to test correctness of data)

- [ ] Load Audio from BB-API for n languages, do speech to text, output the text as JSON
- [ ] compare text to USXEdit, output difference (Mark Scott use case)
- [ ] Load Audio from POST, do speech to text, repeat with same request (OBT use case)

- [ ] Load Project Text with PlainText Edit, because it is verse aligned.  Load time stamps from BB-API
- [ ] Load same Text, and timestamp with aeneas, compare the timestamps using float to find in difference they are equal (float.SubTo, float.Abs, float.Max.  Or, possibly I need to look at Sandeep’s method
- [ ] Postpone refinement of Word level timestamps

- [ ] MFCC calculation, load n audios and compute timestamps and mfcc.  I don’t know how to validate this.  Test must output mfcc data in csv and json.
- [ ] Panda load test.  Load MFCC data with text into pandas using both CSV and JSON.

- [ ] Word Test, use word parse on each text type, and some languages, run unit test to compare results with script.
- [ ] Is there a need to rerun the text comparison tests?
- [ ] Test the timestamps at a word level (how) aeneas is the only way.
- [ ] Test MFCC, and FastText, output csv and json
- [ ] Load both output types into Pandas

