# FCBH Data Preparation for Deep Learning

## Technical Notes

This project is currently in development, and at a very experimental stage.  The files of python code mostly run as individual programs to perform a specific task, which is not the final intent of this project.  Almost all of these programs use Sqlite3 as a datastore, and will create a database file in your current directory, or expect one.  In order to locate database files elsewhere, set the environment variable FCBH_DATASET_DB to that location.

+ FCBHDownload.py - A program for downloading audio and text fileset from the Faith Comes by Hearing “Bible Brain” repository.

+ ScriptAdapter.py - This program reads excel files that contain the text in audio script form, and loads the text into a database.

+ WhisperAdapter.py - Whisper is an OpenAI speech to text tool.  This program uses Whisper to process an audio fileset, and loads the text into a database.

+ WordParser.py - This program reads the audio_script table, breaks scripts into words, punctuation, and whitespace, and stores the result in the audio_words table

+ There are many more python scripts that will be added to this list when they are given a CLI interface.

## Introduction

> "And this gospel of the kingdom will be proclaimed throughout the whole world as a testimony to all nations, and then the end will come." Matt 24:14.
This seems to describe a contingency on our Lord's coming, that the gospel has been proclaimed to all nations.  Certainly, the work that SIL is doing in translating the scriptures, and the work FCBH is doing are critical parts.  But is it possible that the word "proclaimed" implies more than the scriptures themselves, but also the preaching of those scriptures?  Recent advances in artificial intelligence (AI) make it conceivable that technology might someday provide the ability to translate spoken words into any language.  And even though this technology might be developed by a company such as OpenAI, Google, or Microsoft; FCBH houses data that is critical to this development by having audio Bibles in a large number of languages.  Because each language is a translation of the same document, these audio translations will be especially useful to this AI task.  So, this document describes a means to prepare this data for use by AI researchers.

## Methodology

The FCBH audio production process breaks text Bible chapters into script segments (called lines) that are often a sentence long, but always include one speaker.  For all languages, the chapters are broken into the same script segments, so that the content of any script segment is the same in all languages.  The audio is recorded from these scripts into chapter files with timestamps that mark the beginning and end of each script segment.

Provided we can identify a speech to text module for the language being processed.  That module will be used to process the audio, and compare the text generated to the input text to verify correctness.  Even though we are using a script prepared for the audio, it is possible that textual corrections were made that did not make it back into the text that was saved for the audio.

Using the python module Aeneas, we process each chapter to obtain a list of timestamps that mark the beginning and end of each word.  While it might be most efficient to use Aeneas on one chapter at a time.  It might be more accurate to do this on each script segment.

The audio data is converted into Mel-Frequency Cepstral Coefficients (MFCCs) using the python module librosa.  This output is then broken up into word length segments using the timestamps found by Aeneas.  Then the MFCC data is normalized and padded to be the same length in the time dimension to prepare it for use by a neural net.

FastText from Facebook research is used to create a word encoding of all of the available text in each language to be processed.  BERT and Word2Vect are two other tools that could also be used.  Using a lexicon that provides equivalent meanings in each language, and other languages, these language encodings are used to create a single multilingual encoding that will be used for both languages.  There are a few possible solutions to Facebook's MUSE, Google's mBERT, Google's Universal Sentence Encoder (USE), or Byte Pair Encoding (BPE) and Sentence Piece.

The MFCC data for each word, and the corresponding multilingual word encoding of both the language and the source language are used to create a tensor as a timeseries with the corresponding MFCC, target language encoding.

The tensor is loaded into a Recurrent Neural Net (RNN) and/or Natural Language Processor (NLP). I am assuming the model would be designed to predict the next audio word.  

Once the utility of this process has been proven, the data could be loaded into DBP to simplify access.  AI researchers will expect all of the data required for the analysis of a language to be in one "dataset", often a zip file.

## Data Sources

This design proposes loading an AI model with records of audio and text data one word at a time.  For an AI model to produce good results, it will be essential that the text words and the audio are perfectly aligned, that is, each record must contain the text and the audio of the same word.  

In order to produce an audio, the text is first parsed into sentence/voice segments, and then the audio of that text is recorded by an actor speaking.  These audio and text segments are later checked and rechecked by people to ensure the correctness of the audio, and then repaired or re-recorded when necessary.  This appears to be the best source, although at this moment, it is unclear what data will be gotten from Context, and what data will be gotten from the Vessel system, or the older Excel spreadsheet.

It is also essential that the input source includes usfm style codes, because these codes identify many different kinds of titles, section heading, and cross references which are not part of the scripture canon, and are not be included in all language translations.  It is not possible to correctly align script segments from one language to the next without identifying these non-verse text segments.

## Data Structures

The data is organized into three tables.  An Identifying table, which has only a single row in a Sqlite3 database.  This is a denormalized table that contains important identifying information about the data collected for one Bible.

An Audio Script Record is a normalized table that contains one record for each script line of an audio recording.
An Audio Word Record is a normalized table that contains one record for each word of an audio recording.

### Identifying Record

bible_id
language_iso
version_code
source_code
language_id
rolv_id
alphabet_code
language_name
version_name
alphabet_name

In this sqlite implementation, a file contains only one identifying record.  In a design where languages are combined, bible_id would need to included in the script record as a foreign key to the Identifying record.

### Audio Script Record

+ script_id
+ book_id
+ chapter_num
+ audio_file
+ script_num
+ usfm_style
+ person
+ actor
+ script_text
+ script_begin_ts
+ script_end_ts
+ script_mfcc

The script_id is a surrogate auto generated integer primary key.  The logical primary key, which should be implemented as a unique index, consists of language_id, book_id, chapter_num, script_num.

### Audio Word Record

+ word_id
+ script_id
+ word_seq
+ verse_num
+ word
+ punct
+ src_language (questionable)
+ src_word (questionable)

Additional data is produced by a data preparation pipeline.

+ word_begin_ts
+ word_end_ts
+ mfccs
+ mfccs_norm (questionable, possibly post process)
+ word_enc
+ src_word_enc
+ word_multi_enc
+ src_word_multi_enc

The word_id is a surrogate auto generated integer primary key.  The logical primary key, which should be implemented as a unique index, consists of script_id, word_seq.  Script_id is a foreign key to the script table.

## Data Dictionary

### Identifying Attributes

**bible_id** - The FCBH bible_id, often 3 char of ISO + 3 char of Version.  It is the unique identifier of a bible, and is the common identifier given to all text, audio, and video filesets.

**language_iso** - The ISO language code using the ISO 639-3 standard.

**version_code** - The 3 character version code.  This is almost always the same as the last 3 digits of the bible_id.

**source_code** - This code is unique to this system.  It defines the source of the text data collected.  Possible values include: usfm, usx, jsonBible, jsonChapter, excel, etc.

**language_id** - The FCBH language id, which takes into account oral dialect

**rolv_id** - To be written

**alphabet** - The 4 digit code of the ISO 15924 standard.  It is also called script code.  The need for this is uncertain, but if it is useful, it might be necessary to analyze the unicode values of the characters to obtain it.

**language_name** - The ISO 639-3 name of the language.

**version_name** - The name associated with the version_code.

### Book Attributes

**book_id** - The USFM 3 character book code.
    - Context:
    - Vessel: BookAbbreviation
    - Excel: Col B

### Chapter Attributes

**chapter_num** - An integer that defines the chapter number.
    - Context:
    - Vessel: StartChapter [[also, EndChapter can I ignore it?]]
    - Excel: Col C

**audio_file** - The filename or full pathname audio containing the chapter.
    - Context:
    - Vessel: 
    - Excel:

### Script Attributes

**script_id** - A surrogate primary key.  It is an integer that begins with 1 for the first record, and increments for each record inserted.  It is present primarily to make table updates efficient.

**script_num** - An integer that defines the script line that this word is part of when the chapter has been parsed into script segments.  The three fields (book_id, chapter_num, script_num) together uniquely identify a script in *any language*.
    - Context:
    - Vessel:
    - Excel: Col I (4 digit number that increments for an entire book, not for a chapter.)

**usfm_style** - The USFM style code of the text. It is essential for identifying non-verse text, such as headings, titles, footnotes, cross references, and more.  Also, some AI researchers might consider the style information to be a useful source for their AI model.  Note: I think that the style codes being collected are solely the USFM paragraph codes, and not the USFM character codes, but this has not been verified.  Is it really certain there is only one of these for a script segment?
    - Context:
    - Vessel:
    - Excel: only << in verse column to indicate non-verse heading

**person** - This is the person or character who is speaking in a script segment. Narrator is the most frequent person.  This data item is an attribute of a script segment, and is the same for each language.  It is included here because some AI researchers might find this information useful for the analysis of text language, since different people have different grammars and styles of speech.
    - Context:
    - Vessel: CharacterId
    - Excel: Col E (names like, Jesus, John the Baptist, Religious Leaders, Narr 01)

**actor** - This is a number that identifies the actor who is speaking this script segment.  Since the Bible has more persons speaking than the number of actors available to record a Bible, actors will need to play many parts.  This data item is included because some AI researchers might find this information useful for the analysis of audio data.
    - Context:
    - Vessel:
    - Excel: Col F (2 digit number)

**script_text** - This is the text of the script.  The keys book_id, chapter_num, script_num will address corresponding script text in any language.
    - Context:
    - Vessel:
    - Excel: Col M

**script_begin_ts** -  The timestamp that marks the beginning of the script in the audio chapter file.  This could be used to process Aeneas and librosa one script segment at a time, if that produces more accurate results.
    - Context:
    - Vessel:
    - Excel: Verse Timing File

**script_end_ts** - The timestamp that marks the end of the script in the audio chapter file.
    - Context:
    - Vessel:
    - Excel: Verse Timing File

**script_mfcc** - Mel-Frequency Cepstral Coefficients of the audio as produced by the python library librosa, and broken into word segments using the timestamps.
Librosa module

### Word Attributes

**word_id** - A surrogate primary key.  It is an integer that begins with 1 for the first record, and increments for each record inserted.  It is present primarily to make table updates efficient, and make it easy to update word_seq.

**word_seq** - An integer that defines the position of a word in the specific script line that it belongs to.  The columns (language_id, book_id, chapter_num, script_num, word_seq) are a unique index.
    - Source is implicit in script_text sequence

**verse_num** - This is typically a number, but can be a value like 2a. This column will be null when the word is part of a heading, reference, note, or other non-verse text. When a script segment crosses a verse boundary, the Context system can provide information about which word marks the beginning of a new verse.  Vessel and Excel cannot.  This data item belongs to a Word Attribute only if when a script crosses a verse boundary, we are able to identify the word where the new verse starts.
    - Context: 
    - Vessel: StartVerse [[EndVerse, StartSubVerse hopefully ignore]]
    - Excel: verse numbers are embedded in script_text in curly braces, such as {2}.

**word** - The word of an audio segment in UTF-8 format.  This could be more than one word if needed to correctly correspond to a word in the source language.
    - Context
    - Vessel: Text
    - Excel:

**punct** - Any punctuation that followed the word such as: (.,:;!?).  It is parsed out of the word so that it can be included or excluded for different parts of the process as deemed necessary by the AI researcher.
    - Context:
    - Vessel: Text
    - Excel:

**src_language** - The ISO 639-3 code of the source language that was translated.  This datum is placed here, because some passages might be translated from different languages, such as: Hebrew, Greek, Aramic, Latin, or English.  This data item is only needed if src_word is provided.

**src_word** - It is not clear that it will be possible to collect this data item.  But, if it is possible, this would be the word in Hebrew, Greek, Aramic, Latin, English, or some other language that was translated to the data item stored in *word*.  Both the word or the src_word could be multiple words if that was needed to correctly align their meaning.

**word_begin_ts** - The timestamp for the start of a word in an audio script segment.
    - Aeneas module

**word_end_ts** - The timestamp for the end of a word in an audio script segment.
    - Aeneas module

**mfccs** - Mel-Frequency Cepstral Coefficients of the audio as produced by the python library librosa, and broken into word segments using the timestamps.
    - librosa module

**mfccs_norm** - The MFCC data after it has been normalized and padded so the segments are of equal length. [[It is not clear that this belongs in the record.  It might be a post process.]]

**word_enc** - A numeric encoding that uniquely represents a word, and also carries some semantic meaning.
    - FastText, Word2Vec, or BERT

**src_word_enc** - An encoding of the src_word using the same mechanism as the word.
    - FastText, Word2Vec, or BERT

**word_multi_enc** - An encoding of the word which uses an encoding that is common to both languages.
    - MUSE, mBERT, or other procrustes

**src_word_multi_enc** - An encoding of the source word that uses an encoding that is common to both languages.
    - MUSE, mBERT, or other procrustes

## DBP Records

While the data extracted from Context, Vessel or Excel spreadsheets, could be stored in simple records as described above.  When they are to be made externally available, they could be stored in DBP as a means to make it externally available to those outside FCBH who should be given access.

The audio_script table could be integrated into DBP by removing the following columns:
+ book_id
+ chapter_num
+ audio_file

And replacing them with a bible_file_id, which would be a foreign key reference to the DBP bible_files table.

The audio_words table could be integrated into DBP as currently defined.  It contains a script_id column as a foreign key reference to the audio_script table.

When the data is stored in DBP, a query that joined the following tables would provide all the same data described here, and other current DBP data.
audio_words
audio_scripts
bible_files
books
bible_filesets
bibles
languages

## Non-Canonical Text

This document assumes that it is possible to align the words of text with corresponding words of the source text. If this is possible, one problem that must be handled is text that is outside of verses that might or might not come from the source text.

Some headings that are outside of verses are part of the scripture where they appear and others are not.  For example, the Psalm 119 is an acrostic poem whose sections begin with Hebrew letters, and these letters are part of the poem, and are part of the original Hebrew text.  On the other hand, there are descriptive headings, such as "Jesus at a Wedding in Cana of Galilee".  This is an example of a heading that is not part of the original text, but added in the publishing of some version of the Bible.

This matters because to prepare the text data in a tensor, we need to associate the words in the language with the related words in the source language.  But, when a heading is not part of the canon, there might not be a match for that text in the source version, if the source version is the original Greek or Hebrew.  So, code that is aligning text will need to try to distinguish the headings that are part of the original text from those that are not.  I am hoping the style code will make this possible.

## Tensor Design

The tensor will include normalized and padded MFCC's of each audio word, and multilingual encoding of each text word.  If it is possible, it would also contain a multilingual encoding of the source language word for each word in the language.  There are other attributes, such as USFM, person, and actor that an AI researcher might also find useful.

**mfccs_norm** - Each audio of a word is replaced with a MFCC (Mel-Frequency Cepstral Coefficients) of the word, that has been normalized and padded so that all MFCC's are the same length.

**word_multi_enc** - The word of text after being multilingual encoded. Each word of text is replaced by a numeric encoding that was created using Facebook MUSE or some other method.

**src_word_multi_enc** - The source text after being multilingual encoded.  Each word of the source language text is also replaced with a numeric encoding that was created using Facebook MUSE or some other method.

## Data Preparation Pipeline

There are a number of steps required to prepare the data for loading the tensor using some existing tools.  The following shows the essential code to use those tools.

### Starting a project

*Sqlite* could be used as a repository to hold the data while it is being created.  Using sqlite as a temporary data store will provide a laboratory where the data pipeline can be experimentally revised, and isolate revisions from a more formal system.

```
# Create an experimental database
db = DBAdapter(language_iso, language_id, language_name)
```

### Read Chapter of Scripts

Read each chapter of script capturing text and related metadata.  Parsing the text into words and punctuation.  Punctuation that trails words, such as (.,;:!?) are stored separately so that they can be included or excluded as needed.

```
# Load the words from some source
file = FileAdapter(db)
file.loadScriptsCSV("filepath.csv")
```
### Checking Text 

Because changes could be made during the recording of the audio, and these textual changes might not have gotten back into the textual script that we have obtained.  It is best to attempt checking the correctness of the text, by using a speech to text tool, and comparing the generated text with the textual script.

At this time, this task is a black box.  I don’t know if we will be able to identify a speech to text tool that works well for little known languages.  Or, what will be required to do a comparison of its output text to the original text.

### Locating Each Word in Audio

*Aeneas* is used to find the beginning and ending of each word in the audio, and to output a file of timestamps that mark the beginning and ending of those words.  This example processes one chapter at a time.  It is an open question whether it might be more accurate to process each script segment one at a time.

```
# install ffmpeg
# install espeak
# pip install numpy
# pip install aeneas

import os
import json
import subprocess
from DBAdapter import *

# Create a text file of the words using the parsed script in the Database
def createWordsFile(db, audioFile, outputFile):
    with open(outputFile, 'w') as file:
        resultSet = db.selectWordsByFile(audioFile)
        for (word_id, word, punct) in resultSet:
            file.write(word + '\n')

# Use Aeneas to produce timestamps for the beginning and ending of each word
def aeneas(language, audioFile, textFile, outputFile):
    command = [
        "python3", "-m", "aeneas.tools.execute_task",
        audioFile,
        textFile,
        f"task_language={language}|os_task_file_format=json|is_text_type=plain",
        outputFile,
        #"-example-words-multilevel --presets-word"
        "-example-words --presets-word"
    ]
    subprocess.run(command)

# Check that the generated output is consistent with the input,
# and store timestamps.
def storeAeneas(db, audioFile, outputFile):
    resultSet = db.selectWordsByFile(audioFile)
    with open(outputFile, 'r') as file:
        timestamps = json.load(file)
        segments = timestamps['fragments']
        if len(segments) != len(resultSet):
            print("ERROR: Num Text Words =", len(resultSet), 
                "Num Audio Words =", len(segments))
        for index, seg in enumerate(segments):
            (word_id, word, punct) = resultSet[index]
            if len(seg['children']) > 0:
                print("Error in segments there are children", seg)
            if len(seg['lines']) != 1:
                print("Error lines is not 1 word", seg)
            elif word != seg['lines'][0]:
                print("Error parsed word and aeneas do not match")
            db.addWordTimestamp(word_id, float(seg['begin']), float(seg['end']))
    db.updateWordTimestamps()
```

### Generate Mel-Frequency Cepstral Coefficients (MFCC)

*librosa* is used to generate MFCC data of the audio.  It is the MFCC data that will be loaded into the neural network as a proxy for the audio file.  There is also a step here to normalize the MFCC output data.  This process is performed separately on each audio file.

```
import librosa
import numpy as np
from DBAdapter import *

def genMFCC(db, audioDir, audio_file):
    audioPath = os.path.join(audioDir, audio_file)
    audioData, sampleRate = librosa.load(audioPath)
    mfccs = librosa.feature.mfcc(y=audioData, sr=sampleRate, n_mfcc=13)
    hopLength = 512 # librosa default
    frameRate = sampleRate / hopLength
    resultSet = db.selectWordTimestampsByFile(audio_file)
    for (word_id, word, word_begin_ts, word_end_ts) in resultSet:
        startIndex = int(word_begin_ts * frameRate)
        endIndex = int(word_end_ts * frameRate)
        # Slice the MFCC data
        segment = mfccs[:, startIndex:endIndex]
        db.addWordMFCC(word_id, segment)
    db.updateWordMFCCs()
```

### Final Preparation of MFCC data before loading into a tensor

After the MFCC's have been generated and stored for the entire document some additional processing is done to prepare it from use by the neural net.  The data should be normalized and the data should be padded to be of equal length.

```
import sqlite
import numpy as np

def normPadMFCC(db, normalize):
mfccTuples = db.selectWordMFCCs() 
    mfccList = []
    for (word_id, mfcc) in mfccTuples:
        mfccList.append(mfcc)
    joinedMFCCs = np.concatenate(mfccList, axis=1)
    mean = np.mean(joinedMFCCs, axis=1)
    stds = np.std(joinedMFCCs, axis=1)
    maxLen = max(array.shape[1] for array in mfccList)
    for (word_id, mfcc) in mfccTuples:
        if normalize:
            mfcc2 = (mfcc - mean[:, None]) / stds[:, None]
        else:
            mfcc2 = mfcc
        padded = np.pad(mfcc, ((0, 0), (0, maxLen - mfcc2.shape[1])), 'constant')
        db.addPadWordMFCC(word_id, padded)
    db.updatePadWordMFCCs()
```

### Prepare word encoding for the language and source language

*FastText* word encoding will be used for all of the translated languages, using the OT and NT text as a source.  Preexisting encodings should not be used, because they would have a very different size (number of words) than the translated languages.

```
import fasttext
import tempfile
from DBAdapter import *

def encodeWords(db, modelName):
    resultSet = db.selectWords()
    filename = os.path.join(tempfile.mkdtemp(), "scripture.text")
    print("Words to encode in", filename)
    file = open(filename, "w")
    first = 0
    for (word_id, word, punct, src_word) in resultSet:
        file.write(word)
        if punct != None:
            file.write(punct)      
        file.write(" ")
    file.close()
    model = fasttext.train_unsupervised(filename, "cbow")
    model.save_model(modelName)
    for (word_id, word, punct, src_word) in resultSet:
        word_enc = model.get_word_vector(word)
        #print(word, type(word_enc.dtype), word_enc.shape)
        db.addWordEncoding(word_id, word_enc)
    db.updateWordEncoding()

```

### Prepare multilingual encoding

The next step is to use a glossary file that describes the related words of two languages to generate a common encoding for the two languages.  The following is a solution using scipy.linalg.

```
# install MUSE
git clone https://github.com/facebookresearch/MUSE.git
cd MUSE
pip install -r requirements.txt

# prepare a bilingual dictionary that contains one line per word with the word in one language on the left, and the corresponding word in the other language separated by a space.

python supervised.py --src_lang en --tgt_lang fr --src_emb data/wiki.en.vec --tgt_emb data/wiki.fr.vec --dico_train data/en-fr.0-5000.txt --exp_path ./muse_output/

[[Note: this code has not been tested]]
```

## Markdown Notes

Markdown was used for formatting so this document could easily become a Github Readme file.

https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax

# H1
## H2
### H3
Blank line is new paragraph

*italics*
**bold**
***italics and bold***
~strikethrough~
<sub>subscript</sub>
<sup>superscript</sup>
[clickable](https://something.com)
> block quote
`unformatted`
```
code
```
+ unordered list item
1. ordered list item (always 1.)
Tab - nested list item
-[ ] Task item
-[x] Checked Task item
Text needing footnote [^1]
[^1] footnote text

