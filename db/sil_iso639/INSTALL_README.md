
Download
https://iso639-3.sil.org/code_tables/download_tables

sqlite3 iso_639_3.db

sqlite> .mode tabs

sqlite> .import iso-639-3.tab languages

sqlite> .import iso-639-3-macrolanguages.tab macro

sqlite> .import whisper.tab whisper

https://www.ethnologue.com/codes/

sqlite> .import LanguagesCodes.tag language_country

sqlite> .import CountryCodes.tab country

https://glottolog.org/meta/downloads



sqlite> .exit

mv iso_639_3.db $FCBH_DATASET_DB