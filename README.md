# FCBH Data Preparation for Deep Learning

## Introduction

> And this gospel of the kingdom will be proclaimed throughout the whole world as a testimony to all nations, and then the end will come." Matt 24:14.
This seems to describe a contingency on our Lord's coming, that the gospel has been proclaimed to all nations.  Certainly, the work that SIL is doing in translating the scriptures, and the work FCBH is doing are critical parts.  But is it possible that the word "proclaimed" implies more than the scriptures themselves, but also the preaching of those scriptures.  Recent advances in artificial intelligence (AI) make it conceivable that technology might someday provide the ability to translate spoken words into any language.  And even though this technology might be developed by a  company such as OpenAI, Google, or Microsoft; FCBH houses data that is critical to this development by having audio Bibles in a large number of languages.  and because each language is a translation of the same document.  These audio translations will be especially useful to this AI task.  So, this document describes a means to prepare this data for use by AI researchers.

## Methodology

We start with audio data that is in chapter form with timestamps marking the beginning and ending of each segment that was recorded.  These segments are typically no more than a sentence long, and have only one actor's voice.  Associated with each segment, we have the text in the same language for each audio.  [[Not sure: The timestamps will be used to break the chapter audio back into the original segments that were recorded.]]

Using the python module Aeneas, we process the small segments one at a time to obtain a list of timestamps for individual word breaks.  [[Questionable: This process is performed on the smaller files rather than the chapter files so that word alignment errors, if they occur, only affect a small amount of data.]]

The audio data is converted into Mel-Frequency Cepstral Coefficients (MFCCs) using the python module librosa.  This output is then broken up into word length segments using the timestamps found by Aeneas. And then the segments are normalized and padded to be of equal length.

The text is prepared for the tensor by using multilingual word embedding.  There are a few possible solutions to Facebook's MUSE, Google's mBERT, Google's Universal Sentence Encoder (USE), or Byte Pair Encoding (BPE) and Sentence Piece. In this document, MUSE is being used, but without any certainty that this is the best solution. 

FastText from Facebook research is used to create a word encoding of all of the available text in the language to be processed, and the source languages.  BERT and Word2Vect are two other tools.  These language encodings and their source language encoding and a translation dictionary are uaed to create a single multilingual encoding that will be used for both languages.

The MFCC data for each word, and the corresponding multilingual word encoding of both the language and the source language are entered into a tensor as a timeseries with the corresponding MFCC, source language encoding, and target language encoding in the same sample.

The tensor is loaded into a Recurrent Neural Net (RNN). The model should be designed to predict the next audio word.  It could also be used to predict the next text word, but if that were the task is there any advantage in having the audio data in the tensor?

Once the utility of this process has been proven, the data could be loaded into DBP for further access.  AI researchers will expect all of the data required for the analysis of a language to be in one "dataset", often a zip file.

## Data Sources

This design proposes loading an AI model with records of audio and text data one word at a time.  For an AI model to produce good results, it will be essential that the text words and the audio are perfectly aligned, that is, each record must contain the text and the audio of the same word.  

In order to produce an audio, the text is first parsed into sentence/voice segments, and then the audio of that text is recorded by an actor speaking.  These audio and text segments are later checked and rechecked by people to insure the correctness of the audio, and then repaired or rerecorded when necessary.  This appears to be the best source, although some of this data is in the newer Vessel system, and others are in the Excel system that has a history of versions.

A key problem with other possible sources of this data is getting a consistent handling of non-verse text, for example the section headings at the start of a passage of scripture.  USFM and USX 2.0 contain markers at the beginning of each verse, but they do not explicitly mark the end of a verse, instead the non verse text has a wide variety of format codes that identify the text's purpose.  Consistent handling of this is problematic.  But these are handled consistently in the audio production process, because each audio is an explicit recording made from the parsed text, regardless of how the non-verse text was parsed.

## Data Structures

At this time a single denormalized table is recommended to store the data.  If a normalized database design were employed, the time involved in reworking the database design as this experimental process evolved could be excessive.

If I were to normalize the tables, they would be language, books, script, and words.

### Audio Word Record

+ id
+ language_id
+ language_iso
+ language_name
+ alphabet (questionable)
+ book_id
+ chapter_num
+ script_num
+ word_seq
+ verse_num
+ usfm_style
+ person
+ actor
+ word
+ punct
+ src_language (questionable)
+ src_word (questionable)
+ audio_file

Additional data is produced by a data preparation pipeline.

+ word_begin_ts
+ word_end_ts
+ mfccs
+ mfccs_norm
+ word_enc
+ src_word_enc
+ word_multi_enc
+ src_word_multi_enc


The logical primary  key, which should be implemented as a unique index, consists of language_id, book_id, chapter_num, script_num, word_seq.

## Data Dictionary

*id* - A surrogate primary key.  It is an integer that begins with 1 for the first record, and increments for each record inserted.

*language_id* - The FCBH language id, which takes into account oral dialect
    - Vessel: 
    - Excel: 

*language_iso* - The ISO language code using the ISO 639-3 standard.
    - Vessel: 
    - Excel: 

*language_name* - The ISO 639-3 name of the language
    - Vessel: 
    - Excel: 

*alphabet* - The 4 digit code of the ISO 15924 standard.  Also called script code.  The need for this is uncertain, but if it is useful, it will be necessary to analyze the unicode values of the characters to obtain it.
    - Vessel: analyze unicode of text
    - Excel: analyze unicode of text

*book_id* - The USFM 3 character book code.
    - Vessel: BookAbbreviation
    - Excel: 

*chapter_num* - An integer that defines the chapter number.
    - Vessel: StartChapter [[EndChapter hopefully ignore]]
    - Excel:

*script_num* - An integer that defines the script line that this word is part of when the chapter has been parsed into script segments.  The three fields (book_id, chapter_num, script_num) together uniquely identify a script in any language.

*word_seq* - An integer that defines the position of a word in the specific script line that it belongs to.  
    - Vessel: implicit in JSON seq
    - Excel:

*verse_num* - This will be an integer. This column will be null when the the word is part of a heading, reference, note, or other non-verse text.
[[Excel has verse numbers like 2b-3a.  While Vessel has a start verse number and an end verse number field.]]
    - Vessel: StartVerse [[EndVerse, StartSubVerse hopefully ignore]]
    - Excel:

*usfm_style* - The USFM style code of the text. In the excel spreadsheet this item is in the verse field for non-verse audio segments, but for clarity it is separated. [[A style code can be a hint to whether text is part of the canon of scripture.  Can a style provide the answer with certainty in all cases?]] USFM document [Click here](https://ubsicap.github.io/usfm/)
    - Vessel: USFMMarker
    - Excel:

*person or character* - This is an integer that identifies the character speaking in a script_segment.  Each integer identifies a specific Bible person or the narrator.  The meaning of these integers is the same across all languages.
    - Vessel: CharacterId
    - Excel: 

*actor* - This is an integer that identifies an actor that performed this part.  One actor will often perform multiple parts.  There is no reason for this application to know the identity of these people, but there is a chance that knowing when the same actor is used will be useful information to a neural net. 

*word* - The word of the audio segment in UTF-8 format.  This could be more than one word if needed to correctly correspond to a word in the source language.
    - Vessel: Text
    - Excel:

*punct* - Any punctuation that followed the word such as: (.,:;!?).  It is parsed out of the word so that it can be included or excluded for different parts of the process.
    - Vessel: Text

*src_language* - The ISO 639-3 code of the source language that was translated.  This datum is placed here, because some passages might be translated from the source language, and other parts from languages, such as English.
    - Vessel: I have not found this
    - Excel:

*src_word* - The related word in the source language that the text was translated from.  This could be more than one word if needed to correctly correspond to a word in the translated language.
    - Vessel: LinkedReferenceText
    - Excel

*audio_file* - The filename or full pathname audio containing the chapter.
    - Vessel: 
    - Excel:

*word_begin_ts* - The timestamp for the start of the audio_seqment
    - Vessel: TBD
    - Excel: TBD

*word_end_ts* - The timestamp for the end of the audio_segment
    - Vessel: TBD
    - Excell: TBD

*mfccs* - Mel-Frequency Cepstral Coefficients of the audio as produced by the python library librosa, and broken into word segments using the timestamps.

*mfccs_norm* - The MFCC data after it has been normalized and padded so the segments are of equal length.

*word_enc* - A numeric encoding that uniquely represents a word, and also carries some semantic meaning.

*src_word_enc* - An encoding of the src_word using the same mechanism as the word.

*word_multi_enc* - An encoding of the word which uses an encoding that is common to both languages.

*src_word_multi_enc* - An encoding of the source word that uses an encoding that is common to both languages.

## DBP Records

While the data extracted from Vessel and Excel spreadsheets, could be stored in simple records as described above.  When they are to be made externally available, they should be stored in DBP as a means to make it externally available to those outside FCBH who should be given access.

```
CREATE TABLE audio_script (
    file_id INT NOT NULL
    script_num INT NOT NULL,
    verse_num INT NULL, -- non-verse text is null
    usfm_style CHAR(4) NULL, -- required for non-verse text
    person INT NOT NULL, 
    actor INT NOT NULL,
    word VARCHAR(1024) NOT NULL,
    punct CHAR 1 NOT NULL,
    src_language CHAR(3),
    src_word VARCHAR(2048),
    audio_begin_ts BIGINT NOT NULL,
    audio_end_ts BIGINT NOT NULL,
    PRIMARY KEY (file_id, script_num),
    FOREIGN KEY person REFERENCES person_id ?? (person_id) -- questionable
);
-- The above code is NOT tested.
```

When the data is stored in DBP, the following query would return all of the metadata that is being made available to the neural net processing.  One advantage of this approach is that additional data stored in the bible table, and some other tables could be added to the query result.

```
SELECT bibles.language_id, 
    language.iso AS language_code, 
    language.name AS language_name, 
    script as alphabet, /* there are multiple in alphabet_language */
    audio_words.script_seq, 
    bible_file.book_id, 
    book.protestant_order AS book_seq, 
    file.chapter_start AS chapter_num, 
    audio_words.verse_num, 
    audio_words.style_code,
    audio_words.person, 
    audio_words.actor, 
    audio_words.word,
    audio_words.punct,
    audio_words.source_language, 
    audio_words.source_word, 
    bible_file.file_name AS audio_file,
    audio_words.audio_start_ts, 
    audio_words.audio_end_ts
FROM audio_words
JOIN bible_file ON audio_words.file_id = bible_file.id
JOIN books ON bible_file.book_id = books.id
JOIN bible_filesets ON bible_file.hash_id = bible_filesets.hash_id
JOIN bibles ON file_filesets.id = bibles.id
JOIN languages ON bibles.language_id = audio_words.languages.id
ORDER BY book.protestant_order, file.chapter_start, audio_words.script_seq
-- The above code is NOT tested
```

## Non-Canonical Text

This document assumes that it is possible to align the words of text with corresponding words of the source text. If this is possible, one problem that must be handled is text that is outside of verses that might or might not come from the source text.

Some headings that are outside of verses are part of the scripture where they appear and others are not.  For example, the Psalm 119 is an acrostic poem whose sections begin with Hebrew letters, and these letters are part of the poem, and are part of the original Hebrew text.  On the other hand, there are descriptive headings, such as "Jesus at a Wedding in Cana of Galilee".  This is an example of a heading that is not part of the original text, but added in the publishing of some versions of the Bible.

This matters because to prepare the text data in a tensor, we need to associate the words in the language with the related words in the source language.  But, when a heading is not part of the canon, there might not be a match for that text in the source version, if the source version is the original Greek or Hebrew.  So, code that is aligning text will need to try to distinguish the headings that are part of the original text from those that are not.  I am hoping the style code will make this possible. [[This needs to be discussed with SIL]]

## Tensor Design

The tensor will contain one word per record with at least three columns.  One for the language text word, another for the source text word, and a final one for the audio of the same word.  But, each of these data items is given a proxy data item that is suitable for input to a Neural Net.

**word_multi_enc** - The word of text after being multilingual encoded. Each word of text is replaced by a numeric encoding that was created using Facebook MUSE or some other method.

**src_word_multi_enc** - The source text after being multilingual encoded.  Each word of the source language text is also replaced with a numeric encoding that was created using Facebook MUSE or some other method.

**mfccs_norm** - Each audio of a word is replaced with a MFCC (Mel-Frequency Cepstral Coefficients) of the word, that has been normalized and padded so that all MFCC's are the same length.

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
file.loadCSVv1("filepath.csv")
[[What is the correct function name????]]
```

### Locating Each Word in Audio

*Aeneas* is used to find the beginning and ending of each word in the audio, and to output a file of timestamps that mark the beginning and ending of those words.  [[Questionable: This task is performed on the smaller sentence size files so that when an alignment error occurs, the entire chapter is not affected.]]  This process is performed separately on each audio file.

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
        resultSet = db.selectWordsForFile(audioFile)
        for (id, word, src_word) in resultSet:
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
    resultSet = db.selectWordsForFile(audioFile)
    with open(outputFile, 'r') as file:
        timestamps = json.load(file)
        segments = timestamps['fragments']
        if len(segments) != len(resultSet):
            print("ERROR: Num Text Words =", len(resultSet), 
                "Num Audio Words =", len(segments))
        for index, seg in enumerate(segments):
            (id, word, src_word) = resultSet[index]
            if len(seg['children']) > 0:
                print("Error in segments there are children", seg)
            if len(seg['lines']) != 1:
                print("Error lines is not 1 word", seg)
            elif word != seg['lines'][0]:
                print("Error parsed word and aeneas do not match")
            db.updateTimestamps(id, float(seg['begin']), float(seg['end']))
    

# A more flexible alternative is to execute the aeneas library
# [aeneas lib](https://www.readbeyond.it/aeneas/docs/libtutorial.html#overview)
```

### Generate Mel-Frequency Cepstral Coefficients (MFCC)

*librosa* is used to generate MFCC data of the audio.  It is the MFCC data that will be loaded into the neural network as a proxy for the audio file.  There is also a step here to normalize the MFCC output data.  This process is performed separately on each audio file.

```
import librosa
import numpy as np
import json
from DBAdapter import *

def genMFCC(db, audio_file):
    audioData, sampleRate = librosa.load(audio_file)
    print("sampleRate", sampleRate)
    mfccs = librosa.feature.mfcc(y=audioData, sr=sampleRate, n_mfcc=13)
    print("mfccs shape", mfccs.shape)
    # Load your timestamps from the JSON file
    hopLength = 512 # librosa default
    frameRate = sampleRate / hopLength
    resultSet = db.selectTimestamps(audio_file)
    for (id, word, audio_begin_ts, audio_end_ts) in resultSet:
        print(id, word, audio_begin_ts, audio_end_ts)
        startIndex = int(audio_begin_ts * frameRate)
        endIndex = int(audio_end_ts * frameRate)
        # Slice the MFCC data
        segment = mfccs[:, startIndex:endIndex]
        print("start", startIndex, "end", endIndex, "shape", segment.shape)
        db.updateMFCC(id, segment)
```

### Final Preparation of MFCC data before loading into a tensor

After the MFCC's have been generated and stored for the entire document some additional processing is done to prepare it from use by the neural net.  The data should be normalized and the data should be padded to be of equal length.

```
import sqlite
import numpy as np

def prepareMFCCs(normalize):
    mfccsList = sqlite.selecMFCCs() # selects id, MFCCs in numpy
    sum = 0.0
    count = 0.0
    maxLen = 0
    for id, mfccs in mfccsList:
        for mfcc in mfccs:
            for value in mfcc[:]:
            count += 1
            sum += value
            if len(mfcc[0]) > maxLen:
                maxLen = len(mfcc[0])
    average = sum / count
    # iterate again to normalize and pad
            if normalize:
            # Pad MFCC segments to the maximum length
padded_mfcc_segments = []
for segment in mfcc_segments:
    padded_segment = np.pad(segment, ((0, 0), (0, max_length - len(segment[0]))), mode='constant')
    padded_mfcc_segments.append(padded_segment)
    db.updateNormalizedMFCC(id, padded_mfcc_segments)

```

### Prepare word encoding for the language and source language

*FastText* word encoding will be used for all of the translated languages, using the OT and NT text as a source.  Preexisting encodings should not be used, because they would have a very different size (number of words) than the translated languages.

```
import fasttext
from DBAdapter import *

def encodeWords(db):
    resultSet = db.selectWords()
    filename = "scripture.text"
    file = open(filename, "w")
    first = 0
    for (id, book_id, chapter_num, script_num, word_seq, 
            verse_num, usfm_style, person, word) in resultSet:
        if word not in {'.', ',', ';', ':', '!', '?'} and id > 1:
            file.write(" ")
        file.write(word)
    file.close()
    print("start model")
    model = fasttext.train_unsupervised(filename)#, "cbow")
    print("model finished")
    model.save("sonnet.model")
    print("model saved")
    for (id, book_id, chapter_num, script_num, word_seq, 
            verse_num, usfm_style, person, word) in resultSet:
        word_enc = model.get_word_vector(word)
        print(word, type(word.dtype), word.shape)
        db.updateEncoding(id, word_enc)
```

### Prepare multilingual encoding

[[This step requires glossaries for the languages being processed.  I need to obtain this in order to finish this code.]]
The following is a solution using scipy.linalg.

```
import numpy as np
from scipy.linalg import orthogonal_procrustes

# Assuming you have loaded FastText models for both languages
ft_model_lang1 = ...  # FastText model for language 1
ft_model_lang2 = ...  # FastText model for language 2

# Assuming you have a bilingual dictionary
bilingual_dict = {'word_in_lang1': 'word_in_lang2', ...}

# Extracting the word vectors for the bilingual dictionary
matrix_lang1 = np.array([ft_model_lang1[word] for word in bilingual_dict.keys()])
matrix_lang2 = np.array([ft_model_lang2[bilingual_dict[word]] for word in bilingual_dict.keys()])

# Learn the transformation matrix
R, _ = orthogonal_procrustes(matrix_lang1, matrix_lang2)

# Apply the transformation to the entire embedding matrix of one language
transformed_matrix_lang1 = np.dot(ft_model_lang1.get_embedding_matrix(), R)

# Now, transformed_matrix_lang1 is aligned with ft_model_lang2's space
```

The following is a solution using Facebook MUSE

```

# install MUSE
git clone https://github.com/facebookresearch/MUSE.git
cd MUSE
pip install -r requirements.txt

# prepare a bilingual dictionary that contains one line per word with the word in one language on the left, and the corresponding word in the other language separated by a space.

python supervised.py --src_lang en --tgt_lang fr --src_emb data/wiki.en.vec --tgt_emb data/wiki.fr.vec --dico_train data/en-fr.0-5000.txt --exp_path ./muse_output/

[[change the languages and filenames]]

```

### Prepare the text for the tensor

Using the multilingual encoding that was prepared in the step above, locate the correct encoding for each word of the language and each word of the source language in the audio files to be processed.  And create a tensor with the encoding aligned with the MFCC of the audio that spoke the  word.

```
resultSet = db.
```

## Markdown Notes

Markdown was used for formatting so that the document could easily become a Github Readme file.

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

## Notes on Audio Production by Vessel

How audios are produced is not truly a part of this project, except that it is from the audio production process that we obtain the data needed for this project.

Context is an internal tool that is processing USFM.  This is to be retired and incorporated into vessel.  Context is reading USFM.

Vessel makes algorithmic decisions in order to break text into sentences and voice.  It recognizes the following:
1. verse blocks
2. sentences by period
3. looks for other punctuation, like quote marks

In this way they are parsing the text into script_segments, and algorithmically getting it 90% right.
Then Vessel is used interactively to present issues to remote users for them to correct.

The data is stored in CouchBase using JSON format.  CouchBase provides a distributed advantage that a remote worker can sync to a local device without an internet connection for a long time, and then it will sync to the central system when that is possible.

The database stores the data in a hierarchy of book level, chapter level, paragraph level, script level (I don’t recall if he also said verse level)


# Questions
+ In order to generate a multilingual word encoding for a language, I need a lexicon that associates words in the language with words in the source language.  Can I obain such a file?
+ Can the style code of a bit of non-verse text be used to determine whether the text is part of the canon of scripture or translator added explanation?
+ I am creating a tensor that contains 1 word per record with one data item being the text word encoding, and another being the MFCC of the audio of that word.  Are there SIL tools that would make it possible to identify the source language word or words in order to put them in the record as well.




