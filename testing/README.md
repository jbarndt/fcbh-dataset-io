# System Test Plan

## General Guidelines

- Create a single testing package to hold all system tests.
- Establish a Test User so that generated files are in one directory. 
- Give each test a project name that is unique and descriptive. 
- Add info on memory usage and duration to log.
- Put request.yaml files into a directory, name each file using RequestName + StepName.
- Make bibleId a variable in each test, and LanguageISO? 
- Give each request.yaml a func Test{RequestName}{StepName}.
- Each func can process the request with a variable list of BibleIds.
- There is a system wide list of BibleId's to test as well as a list in each test.
- There is a test harness that submits request to the server or cli.
- The test harness will always move output to a specific test output directory.
- Each test func will have some kind of verification, at least a record or line count.
- Each test function will have a comment that describes what manual verification is recommended.

## Last Minute Changes

+ Change all database columns to not null with go consistent default values.
+ Remove VersionCode from Request.
+ Can LanguageISO be removed, or is it not consistently three chars of BibleId

### Terminology

- Plain Text refers to the DBP plain text dataset
- Plain Text Edit refers to DBP plain text with headings extracted from USX
- USX Text Edit refers to USX text edited to contain only content in Script
- Script refers to Excel spreadsheets used in process of recording Audio Bible
- Whisper is a Speech to Text program produced by OpenAI
- FastText is a Word Encoding program produced by Meta
- Aeneas is a Timestamp generating program
- Librosa is an MFCC generating program

## Script level text tests

- [ ] Load Plain Text Script from BB-API, CSV output
- [ ] Load Plain Text Edit from BB-API, JSON output
- [ ] Load USX Text Edit from BB-API, CSV output
- [ ] Load Script Using File, CSV output

- [ ] Compare Plain Text to Plain Text Edit, HTML output
- [ ] Compare Plain Text Edit to USX Edit, HTML output
- [ ] Compare USX Edit to Script, HTML output (Mark Scott use case)
- [ ] Plain Text Edit to Script, HTML output
- Note: Part of the purpose of these comparison tests is to verify that each kind of data parsing is correct.

## Script level audio tests

- [ ] Load Audio from BB-API do Whisper speech to text, output the text as JSON
- [ ] compare text to USXEdit, output HTML (Mark Scott use case)
- [ ] Load Audio from POST, do speech to text, output HTML (OBT use case)
- [ ] Reuse request with second and third passage (OBT use case)

- [ ] Load Project Text with Plain Text Edit, because it is verse aligned.  Load time stamps from BB-API, Output CSV
- [ ] Load same Text, and timestamp with aeneas, compare the timestamps using float to find in difference they are equal (float.SubTo, float.Abs, float.Max.  Or, possibly I need to look at Sandeep’s method
- [ ] Load Script and script level timestamps, compare with above case where the script is verse aligned.
- Note: The comparison code to do this has not yet been written (5/1/24)

- [ ] Load audios and compute timestamps and mfcc.  Verification is limited to record counts. output CSV
- [ ] Repeat and output JSON
- [ ] Load CSV file into Pandas
- [ ] Load Json file into Pandas

## Word level text tests

- [ ] Using existing Plain Text Edit, generate words, validate word split, fast text words, generate CSV
- [ ] Using existing USX Text Edit, generate words, validate word split, fast text words, generate JSON
- [ ] Load CSV file into Pandas
- [ ] Load Json file into Pandas

## Word level audio tests

- [ ] Load audio into Word Plain Text Test, timestamp words using script timestamps, generate MFCC, output CSV
- [ ] Load audio into Word USX Text Edit, timestamp words using script timestamps, generate MFCC, output JSON
- [ ] Load CSV output into Pandas
- [ ] Load JSON output into Pandas

## Full Test

- [ ] Load USX Text Edit, Audio, timestamp, compute MFCC, word Split, timestamp words, compute MFCC of words, fast text encode words, output JSON.
- [ ] Rerun and output CSV

