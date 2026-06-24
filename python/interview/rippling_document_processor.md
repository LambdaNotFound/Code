### Overview
Design and build a Document Processor, which can processe a document and return result based on different critera. 

### Data Model
The document will be passed in as a string

The process logic will be consolidated into a Handler which will composite to the Documenet Processor

### API
Create a Documenet Processor class. A DocumentProcessor can have mulitple Handlers, and each of them process the documents independently. 

Create various Handler, listed as below:

1. LengthHandler: The first statistic we want to understand is the length of document (letters only, exluding whitespace or punctuation)

2. WordCountHandler: it takes the contents of a document as a string parameter and returns the number of words in the document.

3. CommonWordHandler: it takes the document as a string parameter and returns the most common word in the document.
- follow up: Let’s change that to instead return the top 3 words in the document. If there are < 3 words, return whatever is avaialble

### Low Level Design
Apply OOP best practices, the Handler class follows the strategy pattern with a process() behavior

For the CommonWordHandler, use a hashmap to track the word freqency and use a fixed size heap to track the top 3 most frequent words. 

Word tokenization: strip punctuation
Case sensitivity: insensitive
Tie-breaking in CommonWordHandler: use lexical order 



### Test Plan
1. test the LengthHandler with expected input string length
2. test the WordCountHandler that returns the correct number for words in the input string. 