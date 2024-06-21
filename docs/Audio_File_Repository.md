# Audio File Repository
## A General Design
### Gary Griswold
### March 6, 2024

## Requirement

The audio files that are produced by FCBH are first created in text files which are the script of the audio to be produced.  While these text files, correspond to the produced audio better than any other text files in FCBH repositories, it is now clear that the time that it takes to retrieve these documents makes them unusable as a resource for uses beyond the audio production process.  If these documents are to be used as sources for AI or other projects, an organized repository of these files will need to be deployed first.

## General Design

Considering the global nature of the audio group, a central system would be best. And AWS S3 would be a simple and familiar data store.

It would significantly add to the complexity of such a system to have a DBP like database that records what is in S3.  Instead, the program should simply read the files in S3 to know what is present.

If additional information is needed, beyond what is contained in the file path names, then it would be simplest to store that information in .json files or README files along with the data.

If a SQL query facility is really needed, but I suspect it is not needed, a database could be created on demand as needed using the filenames in the repositry.  It would take only a few seconds processing to create a database by reading stored filenames. 

The versioning capability of S3 should be useful for the audio group.

The user interface for uploading files could simply be an HTML form with for adding meta data, and a upload button.  The user interface for locating data should be a simple form that is based on the directory structure used similar to the navigation of a directory tree.

The design of the actual directory structure should be decided upon only after detailed discussion with the audio group.  But, I offer the following thoughts in the interest of clarity.

	/{iso}/{bibleId}/{mediaId}/{stage}/{type}/{actual file}

	{stage} refers to the stage of development, such as post-production.
	{type} is not the file type, but types like: audio parts, audio chapter, timestamps, text script, etc.

A system fitting this design should be only about 2K lines of code.  However, the time to load the data might take longer than development.

The language go is one of the best languages for developing servers, and is well suited to this application.




