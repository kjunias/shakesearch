# Change Log

## Setup
Go and node Docker images were updated

## Backend
* To make the queries case insensitivity, the text content was lowercased while building the index. 
* The underlying text content remained the same since suffixarray (very cool) returns the indices in the text of a particular lookup. This makes it possible to then just return a susbstring of the original text with the indices. 
* For consistency the indices returned by suffixarray are sorted.

* A basic mechanism of returning batches of 20 results at most is implemented. A batch number parameter was also added to the search handler to make this happen.

## Frontend
* The load button is hooked up to the controller on the frontend.
* A model of currently loaded results and batch number is added to the frontend and processed in the controller
* Search results are appended and batch number incremented when the Load More button is clicked
* Search results and batch number are reset to [] and 0 when the Search button is clicked

