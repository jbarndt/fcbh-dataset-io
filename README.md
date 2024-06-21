# FCBH Data Preparation for Deep Learning

## Introduction

> "And this gospel of the kingdom will be proclaimed throughout the whole world as a testimony to all nations, and then the end will come." Matt 24:14.
This seems to describe a contingency on our Lord's coming, that the gospel has been proclaimed to all nations.  Certainly, the work that SIL is doing in translating the scriptures, and the work FCBH is doing are critical parts.  But is it possible that the word "proclaimed" implies more than the scriptures themselves, but also the preaching of those scriptures?  Recent advances in artificial intelligence (AI) make it conceivable that technology might someday provide the ability to translate spoken words into any language.  And even though this technology might be developed by a company such as OpenAI, Google, or Microsoft; FCBH houses data that is critical to this development by having audio Bibles in a large number of languages.  Because each language is a translation of the same document, these audio translations will be especially useful to this AI task.  So, this project hopes to be a help by preparing data.

## Methodology

The FCBH audio production process breaks text Bible chapters into script 
segments (called lines) that are often a sentence long, but always include one speaker.
The audio is recorded from these scripts into chapter files with timestamps 
that mark the beginning and end of each script segment.

Using the python module Aeneas, we process each chapter to obtain a list of
timestamps that mark the beginning and end of each script line, verse, or word.

Using a speech to text module for the language being processed, 
the generated text is used to test the correctness of the audio.

The audio data is converted into Mel-Frequency Cepstral Coefficients (MFCCs) 
using the python module librosa.  This output is then broken up into script line length,
or verse length, or word length segments using the timestamps found by Aeneas.  

FastText from Meta is used to create a word encoding of all of the available 
text in each language to be processed.  BERT and Word2Vect are two other tools 
that could also be used.  
Using a lexicon that provides equivalent meanings in each language, 
and other languages, these language encodings are used to create a 
single multilingual encoding that will be used for both languages.
There are a few possible solutions to Facebook's MUSE, Google's mBERT, 
Google's Universal Sentence Encoder (USE), or Byte Pair Encoding (BPE) 
and Sentence Piece.

The MFCC data for each word, and the corresponding multilingual word encoding 
of both the language and the source language are used to create a tensor 
as a timeseries with the corresponding MFCC, target language encoding.

Then the MFCC encoded word data is normalized and padded to be the same length in the
time dimension to prepare it for use by a neural net.

The tensor is loaded into a LLM (Large Language Model) or 
Recurrent Neural Net (RNN), or Natural Language Processor (NLP). 
It is likely the model would be designed to predict the next audio word.  

Once the utility of this process has been proven, the data could be loaded 
into DBP to simplify access.  AI researchers will expect all of the data 
required for the analysis of a language to be in one "dataset", often a 
zip file.

## Database Structure

The data is organized into three tables.  An Ident table, which has only a single row in a Sqlite3 database.
This is a denormalized table that contains important identifying information about the data collected for one Bible.

The Script is a normalized table that contains one record for each script line of an audio recording,
or one verse of a plain text Bible.

The Word is a normalized table that contains one record for each word of an audio recording.

### Ident Record

**dataset_id** - A unique integer identifier for a dataset.  In this sqlite implementation, it is always 1.  But, in a central database implementation it would uniquely identify each dataset.

**bible_id** - The FCBH bible_id, often 3 char of ISO + 3 char of Version.  It is the unique identifier of a Bible, and is the common identifier given to all text, audio, and video filesets.

**audio_OT_id** - The audio OT media_id if needed

**audio_NT_id** - The audio NT media_id if needed

**text_OT_id** - The text OT media_id if needed

**text_NT_id** - The text NT media_id if needed

**text_source** - This code defines the source of the text data collected.  Possible values include: script, text_plain, text_plain_edit, usx_edit

**language_iso** - The ISO language code using the ISO 639-3 standard.

**version_code** - The 3 character version code.  This is almost always the same as the last 3 digits of the bible_id.

**language_id** - The FCBH language id, which takes into account oral dialect

**rolv_id** - To be written

**alphabet** - The 4 digit code of the ISO 15924 standard.  It is also called script code.

**language_name** - The ISO 639-3 name of the language.

**version_name** - The name associated with the version_code.

### Script Record

**script_id** - A surrogate primary key.  It is an integer that begins with 1 for the first record, and increments for each record inserted.  It is present primarily to make table updates efficient.

**dataset_id** - A foreign key to the Ident table

**book_id** - The USFM 3 character book code.

**chapter_num** - An integer that defines the chapter number.

**chapter_end** - The end chapter of a piece text, almost always the same as chapter_num.

**script_num** - An integer that defines the script line that this word is part of when the chapter has been parsed into script segments.  The three fields (book_id, chapter_num, script_num) together uniquely identify a script in *any language*.

**usfm_style** - The USFM style code of the text. It is essential for identifying non-verse text, such as headings, titles, footnotes, cross references, and more.  Also, some AI researchers might consider the style information to be a useful source for their AI model.  Note: I think that the style codes being collected are solely the USFM paragraph codes, and not the USFM character codes, but this has not been verified.  Is it really certain there is only one of these for a script segment?

**person** - This is the person or character who is speaking in a script segment. Narrator is the most frequent person.  This data item is an attribute of a script segment, and is the same for each language.  It is included here because some AI researchers might find this information useful for the analysis of text language, since different people have different grammars and styles of speech.

**actor** - This is a number that identifies the actor who is speaking this script segment.  Since the Bible has more persons speaking than the number of actors available to record a Bible, actors will need to play many parts.  This data item is included because some AI researchers might find this information useful for the analysis of audio data.

**verse_str** - The starting verse number (string) of this piece of text.

**verse_end** - The ending verse number (string) of this piece of text.

**script_text** - This is the text of the script.

**script_begin_ts** -  The timestamp that marks the beginning of the text in the audio chapter file.

**script_end_ts** - The timestamp that marks the end of the script in the audio chapter file.

**mfcc_json** - Mel-Frequency Cepstral Coefficients of the audio as produced by the python library librosa, and broken into word segments using the timestamps.

### Word Record

**word_id** - A surrogate primary key.  It is an integer that begins with 1 for the first record, and increments for each record inserted.  It is present primarily to make table updates efficient, and make it easy to update.

**script_id** - A foreign key to the script attributes

**word_seq** - An integer that defines the position of a word in the specific script line that it belongs to.  The columns (script_id, word_seq) are a unique index.

**verse_num** - This is typically a number, but can be a value like 2a. This column will be empty when the word is part of a heading, reference, note, or other non-verse text. When a script segment crosses a verse boundary, the Context system can provide information about which word marks the beginning of a new verse.  Vessel and Excel cannot.  This data item belongs to a Word Attribute only if when a script crosses a verse boundary, we are able to identify the word where the new verse starts.

**ttype** - A code that identifies the type of data in word. It values are (W, S, P) meaning (Word, Space, Punctuation)

**word** - The word in UTF-8 format.  This could be more than one word if needed to correctly correspond to a word in the source language.

**word_begin_ts** - The timestamp for the start of a word in an audio script segment.

**word_end_ts** - The timestamp for the end of a word in an audio script segment.

**word_enc** - A numeric encoding that uniquely represents a word, and also carries some semantic meaning.

**word_multi_enc** - An encoding of the word which uses an encoding that is common to both languages.

**mfcc_json** - Mel-Frequency Cepstral Coefficients of the audio as produced by the python library librosa, and broken into word segments using the timestamps.

## Non-Canonical Text

The USFM and USX text sources include footnotes, notes, cross references,
and section headings that are not part of the canon of scripture.
The plain text copies of scripture contain none of this non-canonical text,
but contain only the verses, and none of the headings.
When audio scripts are produced only the canonical text of scripture is
included, except that book headings and chapter headings are included.

When comparing USFM, USX, plain text, and scripts it is essential that
these differences are taken into account.  While the Bible Brain system
recognizes Bible text type "text_usx", this system has a text type of "text_usx_edit",
which is USX data with all of the non-canonical text removed.  And, while
the Bible Brain system has Bible text type "text_plain", this system has
a text type of "text_plain_edit", which is text_plain with book and chapter
headings added.

Both text_plain_edit, and text_usx_edit formats will compare well with a script.

## Data Preparation Pipeline

There are a number of steps required to prepare the data for loading into
a LLM (Large Language Model) or other Neural Net model.

### Starting a project

**Sqlite** is used as the data store to hold the text, audio, and encodings.
It is much higher performance than any server based relational database
as long a multiple concurrent writes are not needed.  And they are not
because each project has its own database.

### Read Chapter of Scripts

The text is read into the data store using either a script, plain text, 
or USX as a data source.  Optionally, it is parsed into individual words
such that audio cut into word pieces can be labeled with the associated text.

### Checking Audio

Speech to text of the audio and comparison of that text with text
from another source is an important means to checking the correctness
of the audio.  Possible speech to text tools include:
**Whisper** - from OpenAI
**MMS** - from Meta

### Locating Each Word in Audio

**Aeneas** is used to find the beginning and ending of each script line, verse, 
or word of the audio.

### Generate Mel-Frequency Cepstral Coefficients (MFCC)

**librosa** is used to generate MFCC data of the audio that is to be used.
as the encoded input to a neural net.

### Prepare word encoding for the text

To prepare text code for loading into a neural net it will need to be encoded.  
Possible methods include:

**Word2Vec** - from Gensim

**FastText** - from Meta

### Prepare multilingual encoding

The next step is to use a glossary file that describes the related words of 
multiple languages to generate a common encoding for the two languages.  
Possible tools for doing this would include:

**mBERT** - from Google

**MUSE** - Multilingual Unsupervised or Supervised Word Embedding - from Meta

### Final Preparation of encoded data

After the MFCC's have been generated and stored for the entire document
some additional processing is done to prepare it from use by the neural net.  
The data should be normalized and the data should be padded to be of equal
length.

## Tensor Design

The tensor will include normalized and padded MFCC's of each audio script line or word, 
and multilingual encoding of each text word.  There are other attributes, such as USFM, person, and actor that an AI researcher might also find useful.

**label** - The associated text line or word

**mfccs_norm** - Each audio of a script line or word is replaced with a 
MFCC (Mel-Frequency Cepstral Coefficients) of the word, that has been normalized and padded so that all MFCC's are the same length.

**word_multi_enc** - The script line or word of text after being multilingual encoded. 
Each word of text is replaced by a numeric encoding that was created using Facebook MUSE or some other method.


