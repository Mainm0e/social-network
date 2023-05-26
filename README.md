# SOCIAL-NETWORK, grit:lab javascript project

## USE
  
Frontend
start app "npm start" (if!! is not work try 
"   rm -rf node_modules
    npm install         "
)  

## PROJECT OVERVIEW

This project consists of creating a Facebook-like social network platform containing the following
features:  

> 1. Followers.
> 
> 2. Profile.
> 
> 3. Posts.
> 
> 4. Groups.
> 
> 5. Notifications
> 
> 6. Chats  

## CODING STANDARDS
  
The following is a draft set of guidelines for how code should be formatted etc. to increase readability and cleanliness of the resulting codebase. All of this is with a view to fascilitating collaboration and communication between members of the coding team, to allow for a better end-product to be produced in a shorter period of time. 

### NAMING SYNTAX  
  
Naming of variables and functions should be as semantically descriptive as possible (within reason, so as to not result in names which are too long / cumbersome).  

> __File names:__ *lowercase* (eg. thisishowtonamefiles.js)
>   
> __Local variables / functions:__ *camelcase* (e.g. variable *topSpeedCurrentCar*, function *findLetterInString* )  
> 
> __Global variables / functions:__ *TitleCase* (e.g. variable *RoomWidth*, function *FindAllGroups* )  
> 
> __Constants:__ *UPPERCASE WITH UNDERSCORE SEPARATORS* (e.g. constant variable POPULATION_CAPACITY)
  
**USE OF CONSTANTS, GLOBAL VARIABLES AND FUNCTIONS SHOULD BE KEPT TO A MINIMUM. WHENEVER THEY ARE ADDED TO THE CODEBASE, THE WHOLE TEAM SHOULD BE NOTIFIED AS TO THEIR NAME AND PURPOSE**
  
### DESCRIPTIONS / COMMENTING 
  
All functions should have a **descriptive body of text above them** which describe the function's purpose, inputs, outputs and plausible error cases. This body of text is then automatically referenced when the functions are used elsewhere in the code, and assists others in understanding a function's use when found elsewhere in the codebase. **Commenting within a function is recommended** where it may help in aiding readability or understanding, however this should be done **in moderation** (i.e. aim to structure code to be as *logical, clean and self-explanatory as possible so as to reduce the need for in-line commenting*).  
  
A function's description should consist of the following:  
  
> 1. Begin with the functions name with matching case (*i.e. if in camelcase, this first word should also be in camelcase*).  
> 
> 2. Continue by describing the inputs and what they represent in layman's terms, as well as what the function returns (outputs, if any).  
> 
> 3. Describe in general terms what the function does, as well as possible scenarios in which a non-nil error is returned (if any).  
  
### TESTING  
  
It is strongly encouraged that a *"Test Driven Development"* (TDD) methodology be followed by all members of the team. In a nutshell, this involves coding a basic unit test prior to developing any functional code and forces the programmer to think about important design decisions before adding to the code base. This also results in the codebase being built up in small, testable blocks of code, automatically guiding the programmer to work in a systematic way (with simple debugging opportunities) to produce robust code with sufficient error-handling.  
  
The **TDD methodology** can be summarised as follows *(of course check online if a more detailed run-through is required)*:  
  
> 1. Have a mirror test-file in place whenever possible *(e.g. systemfunc.go and systemfunc_test.go)*. The test-file will always be under the same package as the file with functions to be tested, and should also import the standard "testing" package (in addition to any other required packages). 
> 
> 2. As soon as you wish to code a function, begin the basic line of go code starting with: 1) Starting comment section to contain the function's description (leave empty for the meantime) as well as declaring the function (no need to specify inputs , outputs etc. at this stage), e.g. func whatIsTheName() {}.  
> 
> 3. Now move to the test file and begin coding a test, which always begins with "Test<function's name>" all in Title case, followed by the standard testing call, e.g. func TestWhatIsTheName(t *testing.T) {}.  
> 
> 4. Now code the test, considering 1 or 2 standard cases, and then any number of edge cases which you think may occur. Establish the necessary variables and call the function to be tested within the test function. Standard conditionals can be used to compare the results with the desired outputs, and t.Errorf("example error message text") used in the case of a fail.  
> 
> 5. Now code just enough code in the actual file (e.g. systemfunc.go) so as to be able to call the function. You will discover that you have already made many design decisions whilst formulating the test :) ... and run the test in the terminal with the command <go test -v>. **DO NOT WORRY IF YOUR FUNCTION FAILS THE TEST, THEY USUALLY DO THE FIRST COUPLE OF TIMES**.  
> 
> 6. Iterate between the actual function and test function until you have produced code that has passed the test and can handle all likely edge cases.  
  
### GIT OPERATIONS  
  
All members of the coding team must work on their own branches of code. The **master branch should only contain code that has been tested or approved via git pull requests (reviewed by an administrator)**. Thus it should be safe to pull from the master at any time.  

When a member of the team is satisfied with a block of code, and considers it ready for merging with the master branch, a **PULL REQUEST** should be initiated on Gitea (the administrator can also be notified on Discord of a pending pull request). The administrator will then merge with the master once having approved the request. For the meantime, all team members can be considered administrators although it is important that **code is reviewed by someone other than the person who has produced the code**.
 