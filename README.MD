## Description :
Ascii-art-web-export consists on making sure that it is possible to export the output of the web application. 
Our project implements the .txt and the .doc file formats.

## Authors :
Helena -    username : AthenaHTA <br>
Nathan -    username : Nathan<br>
Kingsley -  username : Kingyam <br>

## Usage :
- step 1 : Goto directory where program is stored in the console
- step 2 : Type `go run main.go`
- step 3 : In your browser go to `127.0.0.1:8080` or type localhost:8080
- step 4 : Select a banner and enter the text you want converted you can 
           also press enter or insert a \n to get multiline output
- step 5 : Click submit and check out the result
- step 6 : Select the radio button next to your choice of file format, e.g. doc or txt


## Implementation details :
We used the same algorithm as ascii art with a few tweaks mainly to handle the carriage return as well
which was a little different from the newline and we used a pre html tag to display the ascii art output 
because it deals with pre formatted text which made the ascii art easier to show and made the final 
product borderless which we thought looked better. <br>
The project transfers the file through the following Headers: Content-Type, Content-Length and Content-Disposition.