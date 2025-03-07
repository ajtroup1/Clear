# Project Draft

**Goal:** To create a programming language that adapts to and solves modern programming needs. Combines the:
- Readability and ease-of-learning of Python
- Power and versatility of  C++
- Design principles of Go


**Important Considerations:**

- Avoid reinventing the wheel: study and understand other languages to know what they do right and wrong. Make a solution that is unique to other languages and cannot be directly substituted with another due to ease-of-use, performance, design, etc.

- Balancing readability with performance and versatility: have a syntax that is readable, intuitive, “out-of-the-box”, and that minimizes boilerplate. This will be tough to compromise while being powerful, performant, and extremely dynamic/configurable.

- Developer experience, learning, and safety: make a solution that prioritizes:

  - Smooth experience: coding in this language should be fast, simple, intuitive, and dynamic/powerful.

  - Easily learnable: as easy as Python or Go to learn and develop powerful, real-world applications. This includes having ecosystem support to create frameworks that get projects up and running quickly.

  - Code safety: code should not be unpredictable and should crash at runtime as little as possible. I want this to be as bulletproof as possible when deployed.

  - Highlighting errors: the language should best describe to users what is wrong with their code. Preferably, the language should help the user debug to the point where they don’t even need the docs.

  - Fostering an ecosystem and scalability: a big key to any ‘big’ programming language’s success is contribution. Allowing other programmers to publish pre-built solutions and easily give new/experienced users a faster development time and ease-of-development is huge for scalability.


**Value / Uniqueness:**

- Safe, learnable, and descriptive: make it easy for programmers to adapt this language and feel like it is competent enough to be their choice. Also make the language descriptive and quick to pick up`


**Product Necessities:** 

- An intuitive but powerful syntax: The code should be understandable but also ideally handle any case imaginable. I really want to take from Python’s ease of readability and learning and combine that with the raw power and versatility of something like C++.

- A safe and descriptive solution: A big (but understandable) pitfall of C is its lack of safety. It is very powerful but can fail quite easily and annoyingly. This language needs to be as bulletproof as possible and be descriptive of what the user does wrong while using it.

- Documentation: Just the best possible documentation for how to use the language, its features, examples, whatever.

- Scalability and ecosystem: What makes big languages become big in the first place is the fact that independent programmers can publish their code and others can use it easily. Python and Go are both great examples of making a language ecosystem easy to use. I need to foster an easy and powerful way for users to develop an ecosystem within the language.


**Minimum Viable Product:**

- The core language with necessary language features.
Arithmetic operations, functions, operator precedence, etc.


**Tech Stack:**

- Go - core host language
- LLVM - backend to generate optimized machine code
- Look into GoCC, Participle, and GoYacc


**Architectural Considerations:**

- Compiled

- Object Oriented

- Rich prebuilt code library