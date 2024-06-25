# Timestamps

It is a partial intent of this project to produce accurate timestamps for
all FCBH audio files.  And to provide a facility that can be used for new
audio files as they are produced.

This page documents the study of timestamps that will hopefully lead to
a solution to this task.

### June 20, 2024

Errors were found in some database timestamps, when Whisper was given the
ability to process by verse.  The chapter files were chopped into verse
files using ffmpeg and Bible Brain timestamps.  Just two filesets were
done and in those two there were a few cases where words had moved to 
the next verse indicating that the timestamp had been incorrect.

ENGWEB
Titus 2:3 “which is good” move to vs 4
3John 1:4 “walking in truth” move to vs 5
3John 1:10 “beloved don’t imitate that which is evil” move to vs 11

APFCMU
Titus 1:7 “nang masaya amang patalus magagustuhan tumagaganda” move to vs 8

### June 25, 2024

Using the program cli_misc/timestamps_compare/timestamp_experiment.go
it was demonstrated that the timestamps in Bible Brain and the Aeneas 
timestamps in Sandeep's repository are identical with just an occasional
difference of 1e-17.

