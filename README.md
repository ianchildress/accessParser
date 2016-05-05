# accessParser

accessParser is meant to be used with a specific access log format and will currently break (easily) if the format is different. 

## Usage
Use command line arguments to specify the file to parse and the json out file path.
-in should point to the access log you want to parse
-out is optional and should point to where you want the output file to go

## Example
./accessParser -in access.log -out myFile.json
