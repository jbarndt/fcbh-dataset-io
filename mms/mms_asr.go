package mms

// This program

// checks the incoming lang go to see if it is unsupported, and tries to get a new one if it is not.
// receives a set a files as input that in an []Input slice that contain identifying information
// start up the mms_asr.py with cmd, passing lang in as a parameter
// set up the pipes.
// iterate over the files in []Input
// lookup verse timestamps for each file.
// there might be a case for script timestamps?
// create temp directory,
// using ffmpeg, split the audio files into verse segments
// iterate over segments,
// send each segment filename to mms_asr.py
// read output, and process by combining it with identifying, book, chapter, verse in some struct
// at the end of each chapter, or book, insert the data into scripts table
// when entirely done, send it an exit message, such as ctrl-D

// How would the py program be written.
// py program would have a main that accepts the lang as sys.argv and initializes the model
// it would then read waiting for input, and process when it has finished, and output to stdout
// when it receives a ctrl-D signal it will stop

// Can this design be used on the command line?
// Yes, one starts the program with a command line argument
// One preparse a list of input files into a data file.
// Pass the input file using <
// The output should display on stdout.
