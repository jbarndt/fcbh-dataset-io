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

ENGWEB - Titus, Philemon, 3rd John results
Titus 2:3 “which is good” move to vs 4
3John 1:4 “walking in truth” move to vs 5
3John 1:10 “beloved don’t imitate that which is evil” move to vs 11

APFCMU - Titus, Philemon, 3rd John results
Titus 1:7 “nang masaya amang patalus magagustuhan tumagaganda” move to vs 8

### June 25, 2024

Using the program cli_misc/timestamps_compare/timestamp_experiment.go
it was demonstrated that the timestamps in Bible Brain and the Aeneas 
timestamps in Sandeep's repository are identical with just an occasional
difference of 1.0e-17.

### June 26, 2024

A program was written to compare BB timestamps to Aeneas timestamps
over all available filesets in TSdata.  Only 3rd John was run,
and the results are usually off by -0.3 seconds.  But, there were a
number of cases where the BB has a zero timestamp for the last record
throwing the statistics off considerably.  It might be better to do
the same test using Sandeep's Aeneas numbers.

### June 26, 2024

Using the sil database created from their public information there are no
filesets that I can access which have compatibility to Whisper, except
ENGWEB.  So, any testing involving Whisper is not possible.
