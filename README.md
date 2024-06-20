# FCBH Data Preparation for Deep Learning

## Introduction

> "And this gospel of the kingdom will be proclaimed throughout the whole world as a testimony to all nations, and then the end will come." Matt 24:14.
This seems to describe a contingency on our Lord's coming, that the gospel has been proclaimed to all nations.  Certainly, the work that SIL is doing in translating the scriptures, and the work FCBH is doing are critical parts.  But is it possible that the word "proclaimed" implies more than the scriptures themselves, but also the preaching of those scriptures?  Recent advances in artificial intelligence (AI) make it conceivable that technology might someday provide the ability to translate spoken words into any language.  And even though this technology might be developed by a company such as OpenAI, Google, or Microsoft; FCBH houses data that is critical to this development by having audio Bibles in a large number of languages.  Because each language is a translation of the same document, these audio translations will be especially useful to this AI task.  So, this project hopes to be a help by preparing data.

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

+ dataset_id
+ bible_id
+ audio_OT_id
+ audio_NT_id
+ text_OT_id
+ text_NT_id
+ text_source
+ language_iso
+ version_code
+ language_id
+ rolv_id
+ alphabet
+ language_name
+ version_name

In this sqlite implementation, a file contains only one identifying record.  In a design where languages are combined, bible_id would need to included in the script record as a foreign key to the Identifying record.

### Audio Script Record

+ script_id
+ dataset_id
+ book_id
+ chapter_num
+ chapter_end
+ audio_file
+ script_num
+ usfm_style
+ person
+ actor
+ verse_str
+ verse_end
+ script_text
+ script_begin_ts
+ script_end_ts
+ mfcc_json

The script_id is a surrogate auto generated integer primary key.  The logical primary key, which should be implemented as a unique index, consists of book_id, chapter_num, script_num.

### Audio Word Record

+ word_id
+ script_id
+ word_seq
+ verse_num
+ ttype
+ word
+ word_begin_ts
+ word_end_ts
+ word_enc
+ word_multi_enc
+ mfcc_json

The word_id is a surrogate auto generated integer primary key.  The logical primary key, which should be implemented as a unique index, consists of script_id, word_seq.  Script_id is a foreign key to the script table.

## Data Dictionary

### Identifying Attributes

**dataset_id** - A unique integer identifier for a dataset.  In this sqlite implementation, it is always 1.  But, in a central database implementation it would be essential.

**bible_id** - The FCBH bible_id, often 3 char of ISO + 3 char of Version.  It is the unique identifier of a bible, and is the common identifier given to all text, audio, and video filesets.

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

### Script Attributes

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

### Word Attributes

**word_id** - A surrogate primary key.  It is an integer that begins with 1 for the first record, and increments for each record inserted.  It is present primarily to make table updates efficient, and make it easy to update word_seq.

**script_id** - A foreign key to the script attributes

**word_seq** - An integer that defines the position of a word in the specific script line that it belongs to.  The columns (language_id, book_id, chapter_num, script_num, word_seq) are a unique index.

**verse_num** - This is typically a number, but can be a value like 2a. This column will be null when the word is part of a heading, reference, note, or other non-verse text. When a script segment crosses a verse boundary, the Context system can provide information about which word marks the beginning of a new verse.  Vessel and Excel cannot.  This data item belongs to a Word Attribute only if when a script crosses a verse boundary, we are able to identify the word where the new verse starts.

**ttype** - A code that identifies the type of data in word. It values are (W, S, P) meaning (Word, Space, Punctuation)

**word** - The word in UTF-8 format.  This could be more than one word if needed to correctly correspond to a word in the source language.

**word_begin_ts** - The timestamp for the start of a word in an audio script segment.

**word_end_ts** - The timestamp for the end of a word in an audio script segment.

**word_enc** - A numeric encoding that uniquely represents a word, and also carries some semantic meaning.

**word_multi_enc** - An encoding of the word which uses an encoding that is common to both languages.

**mfcc_json** - Mel-Frequency Cepstral Coefficients of the audio as produced by the python library librosa, and broken into word segments using the timestamps.

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

### Read Chapter of Scripts

Read each chapter of script capturing text and related metadata.  Parsing the text into words and punctuation.  Punctuation that trails words, such as (.,;:!?) are stored separately so that they can be included or excluded as needed.

### Checking Text 

Because changes could be made during the recording of the audio, and these textual changes might not have gotten back into the textual script that we have obtained.  It is best to attempt checking the correctness of the text, by using a speech to text tool, and comparing the generated text with the textual script.

At this time, this task is a black box.  I don’t know if we will be able to identify a speech to text tool that works well for little known languages.  Or, what will be required to do a comparison of its output text to the original text.

### Locating Each Word in Audio

*Aeneas* is used to find the beginning and ending of each word in the audio, and to output a file of timestamps that mark the beginning and ending of those words.  This example processes one chapter at a time.  It is an open question whether it might be more accurate to process each script segment one at a time.

### Generate Mel-Frequency Cepstral Coefficients (MFCC)

**librosa** is used to generate MFCC data of the audio.  It is the MFCC data that will be loaded into the neural network as a proxy for the audio file.  There is also a step here to normalize the MFCC output data.  This process is performed separately on each audio file.

### Final Preparation of MFCC data before loading into a tensor

After the MFCC's have been generated and stored for the entire document some additional processing is done to prepare it from use by the neural net.  The data should be normalized and the data should be padded to be of equal length.

### Prepare word encoding for the text

To prepare text code for loading into a neural net it will need to be encoded.  Possible methods include:

**Word2Vec** - from Gensim

**FastText** - from Meta

### Prepare multilingual encoding

The next step is to use a glossary file that describes the related words of two languages to generate a common encoding for the two languages.  Possible tools for doing this would include:

**mBERT** - from Google

**MUSE** - Multilingual Unsupervised or Supervised Word Embedding - from Meta


